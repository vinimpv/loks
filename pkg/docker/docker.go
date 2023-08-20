package docker

import (
	"os"
	"os/exec"
)

func PullImage(image string) error {
	cmd := exec.Command("docker", "pull", image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func CheckImageExists(image string) bool {
	cmd := exec.Command("docker", "inspect", image)
	err := cmd.Run()
	return err == nil
}
