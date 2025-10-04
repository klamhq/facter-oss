package packages

import (
	"context"
	"testing"

	"github.com/klamhq/facter-oss/pkg/options"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	cfg := options.RunOptions{}
	res := New(logrus.New(), &cfg.Facter.Inventory.Packages)
	assert.NotNil(t, res)
}

func TestCollectPackages(t *testing.T) {
	cfg := options.RunOptions{}
	ctx := context.Background()
	p := New(logrus.New(), &cfg.Facter.Inventory.Packages)
	res, err := p.CollectPackages(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}
