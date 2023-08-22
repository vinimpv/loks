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
# Define the components, representing services running in the cluster.
# A component can have multiple deployments, each in its own pod with specific configurations.
components:
  - name: postgres
    image: postgres:latest # Docker image for deployment.
    # Environment variables at the component level, available to all deployments. They can be overridden at the deployment level.
    env:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: postgres
    deployments:
      - name: postgres
        ports:
          - port: 5432
            hostPort: 30432 # Exposed to the host machine, must be in the range of 30000-32767.
        # Kubernetes liveness probe to check if the pod is ready. Using exec probe to run 'pg_isready' to check database readiness.
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
    # Script run after deployment, e.g., creating an S3 bucket in localstack.
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
    skip_image_pull: true # Skip image pull when developing locally.
    # Script to build the development image (tagged ':dev'). Required for 'loks update'.
    build_dev: |
      docker build -t backend:dev .
    # Script run before deployment, e.g., database migrations.
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
    # List of dependent components. Ensures deployment after dependencies.
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
        # Path on the Docker image where the project folder will be mounted.
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
