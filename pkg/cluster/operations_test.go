package cluster

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDestroyCluster(t *testing.T) {
	DestroyCluster("test")
	err := CreateCluster("test", "../../example", []int{})

	// Assert
	assert.NoError(t, err)

	err = DestroyCluster("test")
	assert.NoError(t, err)

}
