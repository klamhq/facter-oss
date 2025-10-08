package sink

import (
	"github.com/klamhq/facter-oss/pkg/agent/store"
	"github.com/klamhq/facter-oss/pkg/options"
	"github.com/klamhq/facter-oss/pkg/utils"
	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
	"github.com/sirupsen/logrus"
)

// SinkInventory sinks the inventory to the configured output(s).
func SinkInventory(cfg *options.RunOptions, logger *logrus.Logger, store store.InventoryStore, inventory *schema.InventoryRequest, fullInventory *schema.HostInventory) error {
	var err error
	// Save the inventory to a file if configured or to remote if enabled
	hostname := utils.GetHostnameFromInventory(inventory)
	switch cfg.Facter.Sink.Output.Type {
	case "file":
		err = exportToFile(inventory, logger, cfg)
		if err != nil {
			logger.WithError(err).Error("Failed to export inventory to file")
		}
	case "remote":
		err = sendOverGrpc(&cfg.Facter.Sink.Output.FacterServer, inventory, logger)
		if err != nil {
			logger.WithError(err).Error("Failed to send inventory to remote server")
		}
	}
	if err != nil {
		errStore := store.Delete(hostname)
		if errStore != nil {
			logger.WithError(errStore).Error("Failed te delete inventory store")
		}
		logger.Warnf("Deleted snapshot for host %s due to failed send", hostname)
		return err
	}

	if err = store.Save(hostname, fullInventory); err != nil {
		logger.Error("Failed to save inventory:", err)
		return err
	}
	logger.Infof("Inventory for host %s saved to local store %s", hostname, cfg.Facter.Store.Path)
	if err := store.Close(); err != nil {
		logger.WithError(err).Error("Failed to close inventory store")
	}
	return nil
}
