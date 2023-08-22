# loks

## Overview

Loks is an attempt to simplify the process of bootstrapping a local Kubernetes environment for development. It aims to emulate an experience similar to Docker Compose, although it's still in early stages.

## How It Works

Under the hood, loks uses [kind] to create a Kubernetes cluster, share volumes with the host machine, and expose ports to the host machine. It employs `ytt` for the YAML templating and [kapp] for the deployment.

## How to Use It

### Configuration

Create a `loks.yaml` file in the root folder of your projects (assuming you want to run multiple projects in the same cluster). Here's an example of what the `loks.yaml` file should look like:

```yaml
name: loks
# Components are the different services that will be running in the cluster
# each component can have multiple deployments, each deployment will be
# running in a different pod and can have different configurations
components:
  - name: postgres
    # image is the docker image that will be used for the deployment of the one or more pods
    image: postgres:latest
    # you can specify env variables at the component level, these will be available to all deployments
    # these can be overriden at the deployment level
    env:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: postgres
    deployments:
      - name: postgres
        # ports is a list of ports that will be exposed by the deployment
        ports:
          - port: 5432
            # hostPort is the port that will be exposed to the host machine, these have to be in the range of 30000-32767
            hostPort: 30432
        # livinessProbe is a kubernetes probe that will be used to check if the pod is ready
        # in this case we are using a exec probe that will run the command pg_isready to check if the database is ready
        livenessProbe:
          exec:
            command:
              - pg_isready
              - -U
              - postgres
              - -d
              - postgres
          initialDelaySeconds: 3
          periodSeconds: 3

  - name: localstack
    image: localstack/localstack:latest
    # post_deploy_script is a script that will be run after the deployment of the component
    # this can be used to run commands like creating a bucket in s3 (in this case we are using localstack)
    post_deploy_script: |
      aws --endpoint-url=${AWS_ENDPOINT_URL} s3 mb s3://${AWS_S3_BUCKET_NAME}
    deployments:
      - name: localstack
        ports:
          - port: 4566
            hostPort: 30566
    env:
      SERVICES: s3
      AWS_DEFAULT_REGION: us-east-1
      AWS_ACCESS_KEY_ID: test
      AWS_SECRET_ACCESS_KEY: test
      AWS_S3_BUCKET_NAME: loks
      AWS_ENDPOINT_URL: http://localstack.default.svc.cluster.local:4566

  - name: backend
    # skip_image_pull is used to skip the image pull, this is useful when you are developing the image locally
    # and you don't want to pull the production image
    # another option here would be to specify the production image like in the other components, then you can
    # skip building the development image
    skip_image_pull: true
    # build_dev is a script that will be run to build the development image
    # its necessary to build the development image and for the `loks update` command to be able to build new images
    # important: the script should build the image with the tag `:dev`
    build_dev: |
      docker build -t backend:dev .
    # pre_deploy_script is a script that will be run before the deployment of the component
    # this can be used to run commands like migrations
    pre_deploy_script: |
      sleep 5
      echo "Migration completed"
    env:
      DATABASE_URL: postgresql://postgres:postgres@postgres:5432/postgres
    deployments:
      - name: backend
        ports:
          - port: 8000
            hostPort: 30800
        mount_path: /app
        livenessProbe:
          httpGet:
            path: /health
            port: 8000
          initialDelaySeconds: 3
          periodSeconds: 3
    # dependencies is a list of components that this component depends on
    # this will make sure that the component is deployed after the dependencies are deployed
    dependencies:
      - redis
      - postgres
      - localstack

  - name: frontend
    skip_image_pull: true
    build_dev: |
      docker build -t frontend:dev .
    deployments:
      - name: frontend
        ports:
          - port: 80
            hostPort: 30080
        # mount_path is a path on the docker image that we will mount the project folder to
        mount_path: /app
    dependencies:
      - backend
```

You can find an example project in the [`examples` folder](https://github.com/vinimpv/loks/tree/main/example).

## Installation

### Prerequisites

- [Docker](https://www.docker.com/get-started/)
- [kind](https://kind.sigs.k8s.io/docs/user/quick-start/)
- [ytt](https://carvel.dev/ytt/)
- [kapp](https://carvel.dev/kapp/)

### Build from Source

```bash
git clone git@github.com:vinimpv/loks.git
cd loks
make build
ln -s $(pwd)/build/loks /usr/local/bin/loks
```

### From Releases

Download the latest release from [here](https://github.com/vinimpv/loks/releases)

### Commands

#### `loks up`

This command will create the cluster and deploy all the components. It will also mount the project folder to the mount path specified in the component configuration.

#### `loks update`

This command will build the development images for the components that have the `build_dev` script specified. It will also deploy the components that have been updated.

#### `loks down`

This command will delete the cluster and all the components.
