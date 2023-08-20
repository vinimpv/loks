package docker

import (
	"fmt"
	"os"
	"os/exec"
)

func PullImage(image string) error {
	cmd := exec.Command("docker", "pull", image)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error pulling image %s, %v", image, err)
	}
	return nil
}

func CheckImageExists(image string) bool {
	cmd := exec.Command("docker", "inspect", image)
	err := cmd.Run()
	return err == nil
}

func BuildDev(path, command string) error {
	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error building %s, %v", path, err)
	}
	return nil
}

func Build(path, tag string) error {
	cmd := exec.Command("docker", "build", path, "-t", tag)
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error building %s, %v", path, err)
	}
	return nil
}

func Tag(source, target string) error {
	cmd := exec.Command("docker", "tag", source, target)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error tagging %s, %v", source, err)
	}
	return nil
}
