package system

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSystem(t *testing.T) {

	s := *GetSystem()
	assert.NotEmpty(t, s.CPU)
	assert.NotEmpty(t, s.Memory)
	assert.NotEmpty(t, s.Load)

}
