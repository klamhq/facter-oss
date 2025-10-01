package network

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetInterfaces(t *testing.T) {
	var iFaces, err = GetNetworkInterfaces()
	assert.NoError(t, err)
	assert.Nil(t, err)
	assert.NotNil(t, iFaces)
	assert.NotEmpty(t, iFaces)

	for _, iFace := range iFaces {
		assert.True(t, len(iFace.Name) > 0)
		assert.True(t, len(iFace.Flags) > 0)
	}
}
