package initSystem

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGatherSystemdInfoFail(t *testing.T) {
	logger := logrus.New()
	s, err := GatherSystemdInfo(logger)
	assert.Error(t, err)
	assert.Nil(t, s)
}

func TestListAllServicesFail(t *testing.T) {
	logger := logrus.New()
	s, err := listAllServices(logger)
	assert.Error(t, err)
	assert.Nil(t, s)
}

func TestGetServiceDetailsFail(t *testing.T) {
	s, err := getServiceDetails("nginx")
	assert.Error(t, err)
	assert.Nil(t, s)
}
