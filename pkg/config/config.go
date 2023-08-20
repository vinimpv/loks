package config

import (
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
	Name        string            `mapstructure:"name"`
	Image       string            `mapstructure:"image"`
	Env         map[string]string `mapstructure:"env,omitempty"`
	Deployments []Deployment      `mapstructure:"deployments"`
	BuildDev    string            `mapstructure:"build_dev,omitempty"`
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

func LoadConfig(path string) (*Config, error) {
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
