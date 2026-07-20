package networks

import (
	"context"
	"testing"

	"github.com/klamhq/facter-oss/pkg/options"
	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCollectNetworks_WithPublicIP(t *testing.T) {
	cfg := options.RunOptions{}
	cfg.Facter.Inventory.Networks.PublicIp.Enabled = true
	cfg.Facter.Inventory.Networks.PublicIp.PublicIpApiUrl = "https://ifconfig.me/"
	cfg.Facter.Inventory.Networks.PublicIp.Timeout = 2
	
	c := New(logrus.New(), &cfg.Facter.Inventory.Networks)
	ctx := context.Background()
	res, err := c.CollectNetworks(ctx)
	
	// May error due to network restrictions, but shouldn't panic
	if err != nil {
		t.Logf("Network test skipped due to: %v", err)
		return
	}
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.DnsInfo)
}

func TestCollectNetworks_GeoIPMissingConfig(t *testing.T) {
	cfg := options.RunOptions{}
	cfg.Facter.Inventory.Networks.GeoIp.Enabled = true
	// Don't set API key or URL
	
	c := New(logrus.New(), &cfg.Facter.Inventory.Networks)
	ctx := context.Background()
	res, err := c.CollectNetworks(ctx)
	
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.GeoipInfo)
	// Should have default values when config missing
	assert.Equal(t, float64(0), res.GeoipInfo.Latitude)
	assert.Equal(t, float64(0), res.GeoipInfo.Longitude)
}

func TestCollectNetworks_AllFeaturesDisabled(t *testing.T) {
	cfg := options.RunOptions{}
	cfg.Facter.Inventory.Networks.PublicIp.Enabled = false
	cfg.Facter.Inventory.Networks.GeoIp.Enabled = false
	cfg.Facter.Inventory.Networks.Connections.Enabled = false
	cfg.Facter.Inventory.Networks.Firewall.Enabled = false
	
	c := New(logrus.New(), &cfg.Facter.Inventory.Networks)
	ctx := context.Background()
	res, err := c.CollectNetworks(ctx)
	
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.DnsInfo)
	assert.NotEmpty(t, res.Interfaces)
}

func TestCraftConnections(t *testing.T) {
	cfg := options.RunOptions{}
	c := New(logrus.New(), &cfg.Facter.Inventory.Networks)
	
	network := &schema.Network{}
	err := c.craftConnections(network)
	
	assert.NoError(t, err)
	// Connections should be populated even if empty
	assert.NotNil(t, network.Connections)
}

func TestCraftConnections_MultipleCallsProd(t *testing.T) {
	cfg := options.RunOptions{}
	c := New(logrus.New(), &cfg.Facter.Inventory.Networks)
	
	network := &schema.Network{}
	
	// Call multiple times to ensure no issues
	err1 := c.craftConnections(network)
	err2 := c.craftConnections(network)
	
	assert.NoError(t, err1)
	assert.NoError(t, err2)
}

func TestNew_ValidConfig(t *testing.T) {
	logger := logrus.New()
	cfg := &options.NetworksOptions{
		Enabled: true,
		PublicIp: options.PublicIpOptions{
			Enabled: true,
		},
	}
	
	collector := New(logger, cfg)
	
	assert.NotNil(t, collector)
	assert.Equal(t, logger, collector.log)
	assert.Equal(t, cfg, collector.cfg)
}

func TestCollectNetworks_ContextCancelled(t *testing.T) {
	cfg := options.RunOptions{}
	c := New(logrus.New(), &cfg.Facter.Inventory.Networks)
	
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately
	
	// Should still work even with cancelled context
	// as the implementation doesn't check context
	res, err := c.CollectNetworks(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}
