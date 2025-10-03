package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsRoot(t *testing.T) {
	if os.Geteuid() == 0 {
		b := IsRoot()
		assert.True(t, b, "Should be run as root")
	}
	b := IsRoot()
	assert.False(t, b, "Should not be run as root in tests")
}
