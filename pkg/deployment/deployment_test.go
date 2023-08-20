package deployment

import (
	"testing"
	"vinimpv/loks/pkg/cluster"
	"vinimpv/loks/pkg/config"
	"vinimpv/loks/pkg/renderer"

	"github.com/stretchr/testify/assert"
)

func TestDeployToCluster(t *testing.T) {
	config, err := config.LoadConfig("../../example/loks.yaml")
	assert.NoError(t, err)
	renderedYaml, err := renderer.Render("../../example/loks.yaml")
	assert.NoError(t, err)
	err = cluster.CreateCluster(config.Name)
	assert.NoError(t, err)
	defer cluster.DestroyCluster(config.Name)
	err = DeployToCluster(config.Name, renderedYaml)
	assert.NoError(t, err)
}
