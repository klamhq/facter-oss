package sink

import (
	"testing"

	"github.com/klamhq/facter-oss/pkg/agent/store"
	"github.com/klamhq/facter-oss/pkg/options"
	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	mockHostname string
)

func TestSinkInventory_ExportToFileFails(t *testing.T) {
	mockHostname = "host1"

	logger := logrus.New()
	cfg := &options.RunOptions{}
	cfg.Facter.Store.Path = "/tmp/store"
	s := func() store.InventoryStore {
		s, err := store.NewBoltInventoryStore(cfg.Facter.Store.Path)
		if err != nil {
			logrus.WithError(err).Fatal("Unable to create inventory boltdb store")
		}
		return s
	}()

	inventory := &schema.InventoryRequest{}
	fullInventory := &schema.HostInventory{}
	fullInventory.Hostname = mockHostname

	err := SinkInventory(cfg, logger, s, inventory, fullInventory)
	assert.Error(t, err)
	assert.Nil(t, s.Save(mockHostname, fullInventory))
	assert.NoError(t, s.Delete(mockHostname))
}

func TestSinkInventory_SaveFails(t *testing.T) {
	mockHostname = "host2"

	logger := logrus.New()
	cfg := &options.RunOptions{}
	cfg.Facter.Store.Path = "/tmp/store2"
	s := func() store.InventoryStore {
		s, err := store.NewBoltInventoryStore(cfg.Facter.Store.Path)
		if err != nil {
			logrus.WithError(err).Fatal("Unable to create inventory boltdb store")
		}
		return s
	}()
	// generate an error by passing a nil inventory, key is hostname and cannot be empty
	inventory := &schema.InventoryRequest{}
	fullInventory := &schema.HostInventory{}
	fullInventory.Hostname = mockHostname

	err := SinkInventory(cfg, logger, s, inventory, fullInventory)
	assert.Error(t, err)
	assert.Nil(t, s.Save(mockHostname, fullInventory))
	assert.NoError(t, s.Delete(mockHostname))
}

func TestSinkInventory_Success(t *testing.T) {
	mockHostname = "host3"

	logger := logrus.New()
	cfg := &options.RunOptions{}
	cfg.Facter.Store.Path = "/tmp/store3"
	s := func() store.InventoryStore {
		s, err := store.NewBoltInventoryStore(cfg.Facter.Store.Path)
		if err != nil {
			logrus.WithError(err).Fatal("Unable to create inventory boltdb store")
		}
		return s
	}()
	inventory := &schema.InventoryRequest{}
	inventory.Content = &schema.InventoryRequest_Full{Full: &schema.HostInventory{Hostname: mockHostname}}
	fullInventory := &schema.HostInventory{}
	fullInventory.Hostname = mockHostname

	err := SinkInventory(cfg, logger, s, inventory, fullInventory)
	assert.NoError(t, err)
}
