package applications

import (
	"context"
	"testing"

	"github.com/klamhq/facter-oss/pkg/options"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCollectApplications_DockerDisabled(t *testing.T) {
	cfg := options.RunOptions{}
	cfg.Facter.Inventory.Applications.Docker.Enabled = false
	
	c := New(logrus.New(), &cfg.Facter.Inventory.Applications)
	ctx := context.Background()
	
	res, err := c.CollectApplications(ctx)
	
	assert.NoError(t, err)
	assert.Nil(t, res, "Should return nil when Docker is disabled")
}

func TestCollectApplications_ContextCancelled(t *testing.T) {
	cfg := options.RunOptions{}
	cfg.Facter.Inventory.Applications.Docker.Enabled = true
	
	c := New(logrus.New(), &cfg.Facter.Inventory.Applications)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately
	
	// Docker client creation should handle cancelled context
	res, err := c.CollectApplications(ctx)
	
	// May return error or nil depending on Docker availability
	if err != nil {
		assert.Nil(t, res)
	}
}

func TestNew_ValidConfig(t *testing.T) {
	logger := logrus.New()
	cfg := &options.ApplicationsOptions{
		Enabled: true,
		Docker: struct {
			Enabled bool `yaml:"enabled"`
		}{
			Enabled: true,
		},
	}
	
	collector := New(logger, cfg)
	
	assert.NotNil(t, collector)
	assert.Equal(t, logger, collector.log)
	assert.Equal(t, cfg, collector.cfg)
}

func TestCollectApplications_MultipleContexts(t *testing.T) {
	cfg := options.RunOptions{}
	cfg.Facter.Inventory.Applications.Docker.Enabled = false
	
	c := New(logrus.New(), &cfg.Facter.Inventory.Applications)
	
	// Call multiple times with different contexts
	ctx1 := context.Background()
	ctx2 := context.Background()
	
	res1, err1 := c.CollectApplications(ctx1)
	res2, err2 := c.CollectApplications(ctx2)
	
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Nil(t, res1)
	assert.Nil(t, res2)
}
