package agent

import (
	"testing"

	"github.com/klamhq/facter-oss/pkg/options"
	"github.com/stretchr/testify/assert"
)

func TestRunDebugAndProfiling(t *testing.T) {
	cfg := options.RunOptions{}
	cfg.Facter.Logs.DebugMode = true
	cfg.Facter.PerformanceProfiling.Enabled = true
	err := Run(&cfg)
	assert.Error(t, err)
}

func TestRunNoBotldFile(t *testing.T) {
	cfg := options.RunOptions{}
	err := Run(&cfg)
	assert.Error(t, err)
}

func TestRunOnlyPath(t *testing.T) {
	cfg := options.RunOptions{}
	dir := t.TempDir()
	cfg.Facter.Store.Path = dir + "file"
	err := Run(&cfg)
	assert.NoError(t, err)
}

func TestRunBoltDirFail(t *testing.T) {
	cfg := options.RunOptions{}
	cfg.Facter.Inventory.Packages.Enabled = true

	dir := t.TempDir()
	cfg.Facter.Store.Path = dir
	err := Run(&cfg)
	assert.Error(t, err)
}

func TestRunFullOpts(t *testing.T) {
	cfg := options.RunOptions{}
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
	dir := t.TempDir()
	cfg.Facter.Store.Path = dir + "store"
	err := Run(&cfg)
	assert.NoError(t, err)
}
