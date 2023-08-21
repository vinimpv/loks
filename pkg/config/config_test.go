package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	configPath := "../../example/loks.yaml"
	config, err := LoadConfigFromPath(configPath)
	assert.NoError(t, err, "Failed to load config")

	// Checking the number of components
	assert.Len(t, config.Components, 5, "Expected 5 components")

	// Check specific properties for each component

	// Redis component
	redis := config.Components[0]
	assert.Equal(t, "redis", redis.Name)
	assert.Equal(t, "redis:latest", redis.Image)
	assert.Len(t, redis.Deployments, 1)
	assert.Equal(t, "redis", redis.Deployments[0].Name)
	assert.Equal(t, 6379, redis.Deployments[0].Ports[0].Port)
	assert.Equal(t, 6379, redis.Deployments[0].HostPort)
	assert.Nil(t, redis.Env)

	// Postgresql component
	postgres := config.Components[1]
	assert.Equal(t, "postgresql", postgres.Name)
	assert.Equal(t, "postgres:latest", postgres.Image)
	assert.Equal(t, "postgres", postgres.Env["POSTGRESQL_PASSWORD"])
	assert.Len(t, postgres.Deployments, 1)
	assert.Equal(t, "postgres", postgres.Deployments[0].Name)
	assert.Equal(t, 5432, postgres.Deployments[0].Ports[0].Port)
	assert.Equal(t, 5432, postgres.Deployments[0].HostPort)

	// Localstack component
	localstack := config.Components[2]
	assert.Equal(t, "localstack", localstack.Name)
	assert.Equal(t, "localstack/localstack:latest", localstack.Image)
	assert.Len(t, localstack.Deployments, 1)
	assert.Equal(t, "localstack", localstack.Deployments[0].Name)
	assert.Equal(t, 4566, localstack.Deployments[0].Ports[0].Port)
	assert.Equal(t, 4566, localstack.Deployments[0].HostPort)
	assert.Equal(t, "s3", localstack.Deployments[0].Env["SERVICES"])

	// Backend-app component
	backend := config.Components[3]
	assert.Equal(t, "backend-app", backend.Name)
	assert.Equal(t, "docker build -t backend:dev .", backend.BuildDev)
	assert.Equal(t, "redis", backend.Env["REDIS_HOST"])
	assert.Equal(t, "postgresql", backend.Env["POSTGRESQL_HOST"])
	assert.Equal(t, "localstack:4566", backend.Env["S3_ENDPOINT"])
	assert.Len(t, backend.Deployments, 1)
	assert.Equal(t, "backend", backend.Deployments[0].Name)
	assert.Equal(t, 5000, backend.Deployments[0].Ports[0].Port)
	assert.ElementsMatch(t, []string{"redis", "postgres", "localstack"}, backend.Deployments[0].Dependencies)

	// Frontend-app component
	frontend := config.Components[4]
	assert.Equal(t, "frontend-app", frontend.Name)
	assert.Equal(t, "docker build -t frontend:dev .", frontend.BuildDev)
	assert.Equal(t, "backend-app:5000", frontend.Env["BACKEND_ENDPOINT"])
	assert.Len(t, frontend.Deployments, 1)
	assert.Equal(t, "frontend", frontend.Deployments[0].Name)
	assert.Equal(t, 80, frontend.Deployments[0].Ports[0].Port)
	assert.ElementsMatch(t, []string{"backend"}, frontend.Deployments[0].Dependencies)

}
