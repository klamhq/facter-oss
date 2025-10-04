package initSystem

import (
	"runtime"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGatherSystemdInfoFailDarwin(t *testing.T) {
	if runtime.GOOS == "darwin" {
		logger := logrus.New()
		s, err := GatherSystemdInfo(logger)
		assert.Error(t, err)
		assert.Nil(t, s)
	}
	logger := logrus.New()
	s, err := GatherSystemdInfo(logger)
	assert.Error(t, err)
	assert.Nil(t, s)
}
func TestGatherSystemdInfo(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("no systemd in darwin")
	}
	logger := logrus.New()
	s, err := GatherSystemdInfo(logger)
	assert.NoError(t, err)
	assert.Nil(t, s)
}

func TestListAllServicesFail(t *testing.T) {
	if runtime.GOOS == "darwin" {
		logger := logrus.New()
		s, err := listAllServices(logger)
		assert.Error(t, err)
		assert.Nil(t, s)
	}
}

func TestListAllServices(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("no systemd in darwin")
	}
	logger := logrus.New()
	s, err := listAllServices(logger)
	assert.NoError(t, err)
	assert.Nil(t, s)

}

func TestGetServiceDetails(t *testing.T) {
	if runtime.GOOS == "darwin" {
		s, err := getServiceDetails("nginx")
		assert.Error(t, err)
		assert.Nil(t, s)
	} else {
		s, err := getServiceDetails("nginx")
		assert.NoError(t, err)
		assert.Nil(t, s)
	}
}
