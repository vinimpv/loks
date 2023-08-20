package renderer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRrender(t *testing.T) {
	out, err := Render("../../example/loks.yaml")
	assert.NoError(t, err)
	assert.NotEmpty(t, out)

}
