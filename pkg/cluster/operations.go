package cluster

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"

	_ "embed"
)

var KIND_CONFIG_TEMPLATE = `
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraMounts:
    - hostPath: {{.HostPathToMount}}
      containerPath: /workspace
  kubeadmConfigPatches:
    - |
      kind: InitConfiguration
      nodeRegistration:
        kubeletExtraArgs:
          node-labels: "ingress-ready=true"
  extraPortMappings:
  {{range .Ports}}
  - containerPort: {{.}}
    hostPort: {{.}}
    protocol: TCP
  {{end}}
`

type CommandTemplateData struct {
	HostPathToMount string
	Ports           []int
}

func renderKindConfig(name, hostPathToMount string, ports []int) (string, error) {
	tmpl, err := template.New("kindConfig").Parse(KIND_CONFIG_TEMPLATE)
	if err != nil {
		return "", fmt.Errorf("error parsing kind template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, CommandTemplateData{
		HostPathToMount: hostPathToMount,
		Ports:           ports,
	})

	if err != nil {
		return "", fmt.Errorf("error executing kind template: %w", err)
	}

	return buf.String(), nil
}

func getCreateClusterCmd(name, hostPathToMount string, ports []int) (string, error) {
	kindConfig, err := renderKindConfig(name, hostPathToMount, ports)
	if err != nil {
		return "", fmt.Errorf("failed to render kind config: %w", err)
	}
	return fmt.Sprintf("kind create cluster --name %s --config - <<EOF\n%s\nEOF", name, kindConfig), nil
}

func CreateCluster(name, hostPathToMount string, ports []int) error {
	creationCmd, err := getCreateClusterCmd(name, hostPathToMount, ports)
	if err != nil {
		return fmt.Errorf("failed to create cluster: %w", err)
	}
	cmd := exec.Command("bash", "-c", creationCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
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
