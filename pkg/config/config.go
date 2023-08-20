package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Deployment struct {
	Name         string            `mapstructure:"name"`
	Port         int               `mapstructure:"port"`
	HostPort     int               `mapstructure:"hostPort"`
	Env          map[string]string `mapstructure:"env,omitempty"`
	Dependencies []string          `mapstructure:"dependencies,omitempty"`
}

type Component struct {
	Name          string            `mapstructure:"name"`
	Image         string            `mapstructure:"image,omitempty"`
	SkipPullImage bool              `mapstructure:"skip_pull_image,omitempty"`
	Env           map[string]string `mapstructure:"env,omitempty"`
	Deployments   []Deployment      `mapstructure:"deployments"`
	BuildDev      string            `mapstructure:"build_dev,omitempty"`
}

type Ingress struct {
	Host  string `mapstructure:"host"`
	Paths []struct {
		Path    string `mapstructure:"path"`
		Service string `mapstructure:"service"`
	} `mapstructure:"paths"`
}

type Config struct {
	Name       string      `mapstructure:"name"`
	Components []Component `mapstructure:"components"`
	Ingress    []Ingress   `mapstructure:"ingress"`
}

func GetCurrentContextRootPath() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %v", err)
	}
	configPathInCurrentDir := filepath.Join(currentDir, "loks.yaml")
	configPathInParentDir := filepath.Join(currentDir, "..", "loks.yaml")

	var configPath string
	if _, err := os.Stat(configPathInCurrentDir); err == nil {
		configPath = configPathInCurrentDir
	} else if _, err := os.Stat(configPathInParentDir); err == nil {
		configPath = configPathInParentDir
	} else {
		return "", fmt.Errorf("no config file found")
	}
	return filepath.Dir(configPath), nil
}

// LoadConfig loads the config from the current directory or the parent directory
// if the config file is not found in the current directory.
// It returns an error if the config file is not found in either directory.
// It returns the config if the config file is found.
// It returns an error if the config file is found but it cannot be parsed.
// The user must be in the root projects directory or inside one of the projects root folders.
func LoadUserConfig() (*Config, error) {
	configPath, err := GetCurrentContextRootPath()
	if err != nil {
		return nil, fmt.Errorf("error getting current context root path: %v", err)
	}
	return LoadConfigFromPath(filepath.Join(configPath, "loks.yaml"))

}

// LoadConfig loads the config from the given path.
func LoadConfigFromPath(path string) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
