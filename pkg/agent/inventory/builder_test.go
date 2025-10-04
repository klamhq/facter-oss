package inventory

import (
	"context"
	"testing"

	"github.com/klamhq/facter-oss/pkg/agent/collect/applications"
	"github.com/klamhq/facter-oss/pkg/agent/collect/networks"
	"github.com/klamhq/facter-oss/pkg/agent/collect/packages"
	"github.com/klamhq/facter-oss/pkg/agent/collect/platform"
	"github.com/klamhq/facter-oss/pkg/agent/collect/process"
	"github.com/klamhq/facter-oss/pkg/agent/collect/ssh"
	"github.com/klamhq/facter-oss/pkg/agent/collect/systemservices"
	"github.com/klamhq/facter-oss/pkg/agent/collect/users"
	"github.com/klamhq/facter-oss/pkg/models"
	"github.com/klamhq/facter-oss/pkg/options"
	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewBuilder_DisablesSSHWhenUserDisabled(t *testing.T) {
	cfg := options.RunOptions{}
	cfg.Facter.Store.Path = "/tmp/store"
	cfg.Facter.Inventory.User.Enabled = false
	cfg.Facter.Inventory.SSH.Enabled = true // Should be disabled by NewBuilder

	system := &models.System{}
	logger := logrus.New()

	b, err := NewBuilder(cfg, system, logger)
	assert.NoError(t, err)
	assert.False(t, b.Cfg.Facter.Inventory.SSH.Enabled, "SSH should be disabled when User collector is disabled")
	b.Store.Close()
}

func TestNewBuilder_DefaultValues(t *testing.T) {
	cfg := options.RunOptions{}
	cfg.Facter.Store.Path = "/tmp/store"
	system := &models.System{}
	system.Host.Hostname = "test"
	logger := logrus.New()

	b, err := NewBuilder(cfg, system, logger)
	assert.NoError(t, err)
	assert.NotNil(t, b.Log)
	assert.Equal(t, cfg, b.Cfg)
	assert.Equal(t, system, b.SystemGather)
	assert.NotNil(t, b.Now)
	assert.NotNil(t, b.WhoAmI)
	assert.NotNil(t, b.Store)
	assert.Nil(t, b.Platform)
	err = b.Store.Delete("test")
	assert.NoError(t, err)
	b.Store.Close()
}

func TestBuilder_Build_MetadataAndHostname(t *testing.T) {
	cfg := options.RunOptions{}
	cfg.Facter.Store.Path = "/tmp/store"
	system := &models.System{}
	system.Host.Hostname = "myhost"
	logger := logrus.New()

	b, err := NewBuilder(cfg, system, logger)
	assert.NoError(t, err)
	inv, err := b.Build(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "myhost", inv.Hostname)
	assert.NotNil(t, inv.Metadata)
	assert.Equal(t, "0.1.0", inv.Metadata.FacterVersion)
	assert.NotEmpty(t, inv.Metadata.RunningDate)
	err = b.Store.Delete("myhost")
	assert.NoError(t, err)
	b.Store.Close()
}

func TestBuilder_ManageDelta_FullSentWhenNoPrevious(t *testing.T) {
	cfg := options.RunOptions{}
	cfg.Facter.Store.Path = "/tmp/store"
	system := &models.System{}
	system.Host.Hostname = "host1"
	logger := logrus.New()
	b, err := NewBuilder(cfg, system, logger)
	assert.NoError(t, err)

	fullInv := &schema.HostInventory{Hostname: "host1"}
	req, returned := b.ManageDelta(fullInv)
	assert.NotNil(t, req)
	assert.Equal(t, fullInv, returned)
	assert.NotNil(t, req.GetFull())
	err = b.Store.Delete("host1")
	assert.NoError(t, err)
	b.Store.Close()

}

func TestBuilder_ManageDelta_DeltaSentWhenChanged(t *testing.T) {
	cfg := options.RunOptions{}
	cfg.Facter.Store.Path = "/tmp/store"
	system := &models.System{}
	system.Host.Hostname = "host2"
	logger := logrus.New()
	b, err := NewBuilder(cfg, system, logger)
	assert.NoError(t, err)

	// Inventaire initial
	pkg := []*schema.Package{{Name: "pkg1", Version: "1.0.0"}}
	fullInv := &schema.HostInventory{Hostname: "host2"}
	fullInv.Packages = pkg
	err = b.Store.Save("host2", fullInv)
	assert.NoError(t, err)

	// Inventaire modifié
	pkgA := []*schema.Package{{Name: "pkg3", Version: "1.0.3"}}
	fullInvA := &schema.HostInventory{Hostname: "host2"}
	fullInvA.Packages = pkgA

	req, returned := b.ManageDelta(fullInvA)
	assert.NotNil(t, req)
	assert.Equal(t, fullInvA, returned)
	assert.NotNil(t, req.GetDelta())

	err = b.Store.Delete("host2")
	assert.NoError(t, err)
	b.Store.Close()
}

func TestBuilder_ManageDelta_NoDeltaSentWhenNoChange(t *testing.T) {
	cfg := options.RunOptions{}
	cfg.Facter.Store.Path = "/tmp/store"
	system := &models.System{}
	system.Host.Hostname = "host3"
	logger := logrus.New()
	b, err := NewBuilder(cfg, system, logger)
	assert.NoError(t, err)

	pkg := []*schema.Package{{Name: "pkg1", Version: "1.0.0"}}
	fullInv := &schema.HostInventory{Hostname: "host3"}
	fullInv.Packages = pkg

	// Enregistre l'inventaire initial dans le store
	err = b.Store.Save("host3", fullInv)
	assert.NoError(t, err)

	// Appelle ManageDelta avec le même inventaire
	req, returned := b.ManageDelta(fullInv)
	assert.Nil(t, req)
	assert.Nil(t, returned)

	err = b.Store.Delete("host3")
	assert.NoError(t, err)
	b.Store.Close()
}

func TestBuilder_CollectAll(t *testing.T) {
	cfg := options.RunOptions{}
	cfg.Facter.Store.Path = "/tmp/store"
	cfg.Facter.Inventory.Packages.Enabled = true
	cfg.Facter.Inventory.Networks.Enabled = true
	cfg.Facter.Inventory.Networks.Ports.Enabled = true
	cfg.Facter.Inventory.Networks.Connections.Enabled = true
	cfg.Facter.Inventory.SSH.Enabled = true
	cfg.Facter.Inventory.User.Enabled = true
	cfg.Facter.Inventory.Platform.Enabled = true
	cfg.Facter.Inventory.Applications.Enabled = true
	cfg.Facter.Inventory.SystemdService.Enabled = true
	cfg.Facter.Inventory.Process.Enabled = true
	system := &models.System{}
	system.Host.Hostname = "host4"
	logger := logrus.New()
	b, err := NewBuilder(cfg, system, logger)
	assert.NoError(t, err)
	b.Packages = packages.New(logger, &cfg.Facter.Inventory.Packages)
	b.Networks = networks.New(logger, &cfg.Facter.Inventory.Networks)
	b.Users = users.New(logger, &cfg.Facter.Inventory.User)
	b.SSHInfos = ssh.New(logger, &cfg.Facter.Inventory.SSH)
	b.Platform = platform.New(logger, &cfg.Facter.Inventory.Platform, models.SystemPaths{}, system)
	b.Applications = applications.New(logger, &cfg.Facter.Inventory.Applications)
	b.SystemServices = systemservices.New(logger, &cfg.Facter.Inventory.SystemdService)
	b.Processes = process.New(logger, &cfg.Facter.Inventory.Process)
	ctx := context.Background()
	inv, err := b.Build(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, inv)
}
