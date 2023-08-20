package deployment

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func DeployToCluster(clusterName, renderedYaml string) error {

	cmd := exec.Command("kapp", "deploy", "-y", "-a", clusterName, "-f", "-")
	cmd.Stdin = bytes.NewBufferString(renderedYaml)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to deploy to cluster: %w", err)
	}
	return nil
}
