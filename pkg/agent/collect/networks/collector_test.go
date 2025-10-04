package networks

import (
	"context"
	"testing"

	"github.com/klamhq/facter-oss/pkg/options"
	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	cfg := options.RunOptions{}
	res := New(logrus.New(), &cfg.Facter.Inventory.Networks)
	assert.NotNil(t, res)
}

func TestCollectNetworks(t *testing.T) {
	cfg := options.RunOptions{}
	c := New(logrus.New(), &cfg.Facter.Inventory.Networks)
	ctx := context.Background()
	res, err := c.CollectNetworks(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.DnsInfo)
	assert.NotEmpty(t, res.Interfaces)
	assert.Empty(t, res.ExternalIp)
	assert.Empty(t, res.Firewall)
	assert.Empty(t, res.GeoipInfo)
	assert.Empty(t, res.Connections)
}

func TestCollectNetworksWithConnections(t *testing.T) {
	cfg := options.RunOptions{}
	cfg.Facter.Inventory.Networks.Connections.Enabled = true
	c := New(logrus.New(), &cfg.Facter.Inventory.Networks)
	ctx := context.Background()
	res, err := c.CollectNetworks(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.DnsInfo)
	assert.NotEmpty(t, res.Interfaces)
	assert.Empty(t, res.ExternalIp)
	assert.Empty(t, res.Firewall)
	assert.Empty(t, res.GeoipInfo)
	assert.NotEmpty(t, res.Connections)
}

func TestCraftFirewallFail(t *testing.T) {
	cfg := options.RunOptions{}
	cfg.Facter.Inventory.Networks.Connections.Enabled = true
	c := New(logrus.New(), &cfg.Facter.Inventory.Networks)
	network := &schema.Network{}
	err := c.craftFirewall(network)
	assert.Error(t, err)

}
