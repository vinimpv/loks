package renderer

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRrender(t *testing.T) {
	out, err := Render("../../example/loks.yaml")
	assert.NoError(t, err)
	assert.NotEmpty(t, out)

	//write output to file
	err = os.WriteFile("test.yaml", []byte(out), 0644)
	assert.NoError(t, err)

}
