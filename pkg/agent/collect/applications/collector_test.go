package applications

import (
	"context"
	"os"
	"testing"

	"github.com/klamhq/facter-oss/pkg/options"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	cfg := options.RunOptions{}
	res := New(logrus.New(), &cfg.Facter.Inventory.Applications)
	assert.NotNil(t, res)
}

func TestCollectionApplicationsEmpty(t *testing.T) {
	cfg := options.RunOptions{}
	c := New(logrus.New(), &cfg.Facter.Inventory.Applications)
	ctx := context.Background()
	res, err := c.CollectApplications(ctx)
	assert.NoError(t, err)
	assert.Nil(t, res)
}

func TestCollectionApplicationsDockerEnabled(t *testing.T) {
	cfg := options.RunOptions{}
	cfg.Facter.Inventory.Applications.Docker.Enabled = true
	c := New(logrus.New(), &cfg.Facter.Inventory.Applications)
	ctx := context.Background()
	res, err := c.CollectApplications(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestCollectionApplicationsDockerFail(t *testing.T) {
	cfg := options.RunOptions{}
	cfg.Facter.Inventory.Applications.Docker.Enabled = true
	os.Setenv("DOCKER_HOST", "http://fake")
	c := New(logrus.New(), &cfg.Facter.Inventory.Applications)
	ctx := context.Background()
	res, err := c.CollectApplications(ctx)
	assert.Error(t, err)
	assert.Nil(t, res)
}
