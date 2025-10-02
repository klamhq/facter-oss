package agent

import (
	"context"
	"time"

	"github.com/klamhq/facter-oss/pkg/agent/collect/applications"
	"github.com/klamhq/facter-oss/pkg/agent/collect/networks"
	"github.com/klamhq/facter-oss/pkg/agent/collect/packages"
	"github.com/klamhq/facter-oss/pkg/agent/collect/platform"
	"github.com/klamhq/facter-oss/pkg/agent/collect/process"
	"github.com/klamhq/facter-oss/pkg/agent/collect/systemservices"
	"github.com/klamhq/facter-oss/pkg/agent/collect/users"
	"github.com/klamhq/facter-oss/pkg/agent/collectors/system"
	"github.com/klamhq/facter-oss/pkg/agent/inventory"
	"github.com/klamhq/facter-oss/pkg/agent/sink"
	"github.com/klamhq/facter-oss/pkg/models"
	"github.com/klamhq/facter-oss/pkg/options"
	"github.com/klamhq/facter-oss/pkg/performance"
	"github.com/klamhq/facter-oss/pkg/utils"
	"github.com/sirupsen/logrus"
)

// RunAgent is the main function to run the agent.
// It collects system facts, crafts a protobuf message, and sends it to the configured output.
// It also handles performance profiling if enabled in the configuration.
func Run(cfg *options.RunOptions) {
	start := time.Now() // Used to mesure running duration
	// get default value
	options.DefaultNewRunOptions()
	defaultLogLevel := logrus.InfoLevel
	logrus.Info("[AGENT] Collecting system facts...")
	if cfg.Facter.Logs.DebugMode {
		defaultLogLevel = logrus.DebugLevel
	}
	var factory utils.LoggerFactory = &utils.DefaultLoggerFactory{}
	logger := factory.New(defaultLogLevel)

	logger.Debugf("Log verbosity set to %s", defaultLogLevel)

	if cfg.Facter.PerformanceProfiling.Enabled {
		performance.Profiling(logger)
	}
	systemGather := system.GetSystem()

	b := inventory.NewBuilder(*cfg, systemGather, logger)

	b.Platform = platform.New(b.Log, &b.Cfg.Facter.Inventory.Platform, models.SystemPaths{
		InitCheckPath: b.Cfg.Facter.Inventory.Platform.System.InitCheckPath,
		MachineID:     b.Cfg.Facter.Inventory.Platform.System.MachineID,
		MachineUUID:   b.Cfg.Facter.Inventory.Platform.System.MachineUUID,
	}, systemGather)

	b.Packages = packages.New(b.Log, &b.Cfg.Facter.Inventory.Packages)

	b.Applications = applications.New(b.Log, &b.Cfg.Facter.Inventory.Applications)

	b.SystemServices = systemservices.New(b.Log, &b.Cfg.Facter.Inventory.SystemdService)

	b.Networks = networks.New(b.Log, &b.Cfg.Facter.Inventory.Networks)

	b.Users = users.New(b.Log, &b.Cfg.Facter.Inventory.User)

	b.Processes = process.New(b.Log, &b.Cfg.Facter.Inventory.Process)

	inventory, err := b.Build(context.Background())
	if err != nil {
		logger.WithError(err).Error("Unable to build inventory")
	}
	inventoryMsg, fullInventory := b.ManageDelta(inventory)
	if inventoryMsg == nil {
		logger.Info("No inventory changes detected, nothing to do !")
	}

	err = sink.SinkInventory(cfg, logger, b.Store, inventoryMsg, fullInventory)
	if err != nil {
		logger.WithError(err).Error("Failed to sink inventory")
	}

	elapsed := time.Since(start).Round(time.Millisecond)
	logger.Infof("Runned in %s", elapsed)
}
