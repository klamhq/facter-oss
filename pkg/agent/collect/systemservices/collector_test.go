package systemservices

import (
	"context"
	"testing"

	"github.com/klamhq/facter-oss/pkg/options"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	cfg := options.RunOptions{}
	res := New(logrus.New(), &cfg.Facter.Inventory.SystemdService)
	assert.NotNil(t, res)
}

func TestCollectSystemServicesBadSysInit(t *testing.T) {
	cfg := options.RunOptions{}
	s := New(logrus.New(), &cfg.Facter.Inventory.SystemdService)
	ctx := context.Background()
	res, err := s.CollectSystemServices(ctx, "fake")
	assert.Empty(t, res)
	assert.Nil(t, err)
}
