package gefyra

import (
	"fmt"
	"os"
	"os/exec"
)

func Start(clusterName string) error {
	cmd := exec.Command("gefyra", "up", "--context", fmt.Sprintf("kind-%s", clusterName))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error starting gefyra: %v", err)
	}
	return nil
}

func RunContainer(image, containerName, namespace string, ports, volumes []string) error {
	args := []string{"run", "-i", image, "-N", containerName, "-n", namespace}
	for _, port := range ports {
		args = append(args, "--expose", port)
	}
	for _, volume := range volumes {
		args = append(args, "-v", volume)
	}
	cmd := exec.Command("gefyra", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error running gefyra container: %v", err)
	}
	return nil
}

func Bridge(image, containerName, deploymentName, namespace string, ports, volumes []string) error {
	args := []string{"bridge", "-N", containerName, "--target", "deploy/" + deploymentName, "-n", namespace, "--port", ports[0]}
	cmd := exec.Command("gefyra", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error running gefyra bridge: %v", err)
	}
	return nil
}

func Unbridge(name string) {
	cmd := exec.Command("gefyra", "unbridge", "-N", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("error running gefyra unbridge: %v", err)
	}
}
