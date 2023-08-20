package renderer

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"vinimpv/loks/pkg/utils"
)

//go:embed schema.yaml
var schema string

//go:embed template.yaml
var template string

func Render(configPath string) (string, error) {
	// Open the config file and get its contents as a string
	file, err := os.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("failed to open config file: %w", err)
	}
	fileStr := string(file)
	fileStr = fmt.Sprintf("#@data/values\n---\n%s", fileStr)

	tempConfigFilePath, err := utils.WriteTempFile(fileStr)
	if err != nil {
		return "", fmt.Errorf("failed to write config to temporary file: %w", err)
	}
	defer os.Remove(tempConfigFilePath)

	tempTemplateFilePath, err := utils.WriteTempFile(template)
	if err != nil {
		return "", fmt.Errorf("failed to write template to temporary file: %w", err)
	}
	defer os.Remove(tempTemplateFilePath)

	tempSchemaFilePath, err := utils.WriteTempFile(schema)
	if err != nil {
		return "", fmt.Errorf("failed to write schema to temporary file: %w", err)
	}
	defer os.Remove(tempSchemaFilePath)

	out, err := ytt(tempSchemaFilePath, tempTemplateFilePath, tempConfigFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to render: %w", err)
	}
	return out, nil

}

func ytt(schemaPath, templatePath, valuesPath string) (string, error) {
	// Assuming ytt is available as a command line tool
	cmd := exec.Command("ytt", "-f", schemaPath, "-f", templatePath, "-f", valuesPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute ytt: %v\nOutput: %s", err, output)
	}
	return string(output), nil
}
