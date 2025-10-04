package process

import (
	"context"
	"testing"

	"github.com/klamhq/facter-oss/pkg/options"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	cfg := options.RunOptions{}
	res := New(logrus.New(), &cfg.Facter.Inventory.Process)
	assert.NotNil(t, res)
}

func TestCollectProcess(t *testing.T) {
	cfg := options.RunOptions{}
	ctx := context.Background()
	p := New(logrus.New(), &cfg.Facter.Inventory.Process)
	res, err := p.CollectProcess(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}
