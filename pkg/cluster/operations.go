package cluster

import (
	"fmt"
	"os"
	"os/exec"
)

const KIND_CREATION_CMD = `cat <<EOF | kind create cluster --name %s --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
    - |
      kind: InitConfiguration
      nodeRegistration:
        kubeletExtraArgs:
          node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 31820
    hostPort: 31820
    protocol: udp
EOF
`

func CreateCluster(name string) error {
	creationCmd := fmt.Sprintf(KIND_CREATION_CMD, name)
	cmd := exec.Command("bash", "-c", creationCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to create cluster")
	}
	fmt.Println("cluster created successfully")
	return nil

}

func DestroyCluster(name string) error {
	cmd := exec.Command("kind", "delete", "clusters", name)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to destroy cluster: %w", err)
	}
	fmt.Printf("cluster %s destroyed successfully\n", name)
	return nil
}

func LoadImage(clusterName, image string) error {
	cmd := exec.Command("kind", "load", "docker-image", image, "--name", clusterName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to load image: %w", err)
	}
	fmt.Printf("image %s loaded successfully to cluster %s\n", image, clusterName)
	return nil
}
