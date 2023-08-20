package utils

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/google/uuid"
)

func WriteTempFile(contents string) (string, error) {
	// write a certain string to a temporary file in a certain folder, return the file path
	uuid := uuid.New()
	filePath := fmt.Sprintf("%s/loks-%s.yaml", os.TempDir(), uuid)
	err := os.WriteFile(filePath, []byte(contents), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write temporary file: %w", err)
	}
	return filePath, nil
}

func IsToolInstalled(tool string) bool {
	_, err := exec.LookPath(tool)
	return err == nil
}
