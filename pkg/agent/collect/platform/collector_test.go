package platform

import (
	"context"
	"testing"

	"github.com/klamhq/facter-oss/pkg/models"
	"github.com/klamhq/facter-oss/pkg/options"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	cfg := options.RunOptions{}
	path := models.SystemPaths{}
	sys := &models.System{}
	res := New(logrus.New(), &cfg.Facter.Inventory.Platform, path, sys)
	assert.NotNil(t, res)
}

func TestCollectPlatform(t *testing.T) {
	cfg := options.RunOptions{}
	ctx := context.Background()
	path := models.SystemPaths{}
	sys := &models.System{}
	p := New(logrus.New(), &cfg.Facter.Inventory.Platform, path, sys)
	res, err := p.CollectPlatform(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}
