package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseResourceString(t *testing.T) {
	resources, err := parseResourceString("")
	assert.NoError(t, err)
	assert.Empty(t, resources)
}
