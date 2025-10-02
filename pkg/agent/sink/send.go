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

	// Save the inventory to a file if configured
	hostname := utils.GetHostnameFromInventory(inventory)
	err := exportToFile(inventory, logger, cfg)
	if err != nil {
		logger.WithError(err).Fatal("Failed to export inventory to file")
		_ = store.Delete(hostname)
		logger.Warnf("Deleted snapshot for host %s due to failed send", hostname)
		return err
	}
	if err := store.Save(hostname, fullInventory); err != nil {
		logger.Error("Failed to save inventory:", err)
		return err
	}
	logger.Infof("Inventory for host %s saved to local store %s", hostname, cfg.Facter.Store.Path)
	if err := store.Close(); err != nil {
		logger.WithError(err).Error("Failed to close inventory store")
	}
	return nil
}
