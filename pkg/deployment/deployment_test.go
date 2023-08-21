package deployment

import (
	"testing"
	"vinimpv/loks/pkg/cluster"
	"vinimpv/loks/pkg/config"
	"vinimpv/loks/pkg/renderer"

	"github.com/stretchr/testify/assert"
)

func TestDeployToCluster(t *testing.T) {
	config, err := config.LoadConfigFromPath("../../example/loks.yaml")
	assert.NoError(t, err)
	renderedYaml, err := renderer.Render("../../example/loks.yaml")
	assert.NoError(t, err)
	err = cluster.CreateCluster(config.Name, "../../example", []int{8000, 8001, 8002, 8003, 8004, 8005, 8006, 8007})
	assert.NoError(t, err)
	defer cluster.DestroyCluster(config.Name)
	err = DeployToCluster(config.Name, renderedYaml)
	assert.NoError(t, err)
}
