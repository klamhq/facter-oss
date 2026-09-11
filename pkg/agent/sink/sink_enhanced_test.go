package sink

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/klamhq/facter-oss/pkg/agent/store"
	"github.com/klamhq/facter-oss/pkg/options"
	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSinkInventory_FileOutput(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "test-store.db")
	outputPath := filepath.Join(tmpDir, "output.json")

	logger := logrus.New()
	cfg := &options.RunOptions{}
	cfg.Facter.Store.Path = storePath
	cfg.Facter.Sink.Output.Type = "file"
	cfg.Facter.Sink.Output.Format = "json"
	cfg.Facter.Sink.Output.OutputDirectory = tmpDir
	cfg.Facter.Sink.Output.OutputFilename = "output.json"

	s, err := store.NewBoltInventoryStore(cfg.Facter.Store.Path)
	assert.NoError(t, err)

	fullInventory := &schema.HostInventory{
		Hostname: "test-host",
	}
	inventory := &schema.InventoryRequest{
		Content: &schema.InventoryRequest_Full{Full: fullInventory},
	}

	err = SinkInventory(cfg, logger, s, inventory, fullInventory)
	assert.NoError(t, err)

	// Verify output file was created
	_, err = os.Stat(outputPath)
	assert.NoError(t, err)

	// Verify store was saved
	_, err = os.Stat(storePath)
	assert.NoError(t, err)
}

func TestSinkInventory_FileOutputProtobuf(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "test-store.db")
	outputPath := filepath.Join(tmpDir, "output.pb")

	logger := logrus.New()
	cfg := &options.RunOptions{}
	cfg.Facter.Store.Path = storePath
	cfg.Facter.Sink.Output.Type = "file"
	cfg.Facter.Sink.Output.Format = "proto"
	cfg.Facter.Sink.Output.OutputDirectory = tmpDir
	cfg.Facter.Sink.Output.OutputFilename = "output.pb"

	s, err := store.NewBoltInventoryStore(cfg.Facter.Store.Path)
	assert.NoError(t, err)

	fullInventory := &schema.HostInventory{
		Hostname: "test-host-proto",
	}
	inventory := &schema.InventoryRequest{
		Content: &schema.InventoryRequest_Full{Full: fullInventory},
	}

	err = SinkInventory(cfg, logger, s, inventory, fullInventory)
	assert.NoError(t, err)

	// Verify output file was created
	_, err = os.Stat(outputPath)
	assert.NoError(t, err)
}

func TestSinkInventory_InvalidOutputType(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "test-store.db")

	logger := logrus.New()
	cfg := &options.RunOptions{}
	cfg.Facter.Store.Path = storePath
	cfg.Facter.Sink.Output.Type = "unknown"
	cfg.Facter.Sink.Output.Format = "json"
	cfg.Facter.Sink.Output.OutputDirectory = tmpDir
	cfg.Facter.Sink.Output.OutputFilename = "output.json"

	s, err := store.NewBoltInventoryStore(cfg.Facter.Store.Path)
	assert.NoError(t, err)

	fullInventory := &schema.HostInventory{
		Hostname: "test-host-unknown",
	}
	inventory := &schema.InventoryRequest{
		Content: &schema.InventoryRequest_Full{Full: fullInventory},
	}

	// With unknown type, it should still work (no error in switch default)
	err = SinkInventory(cfg, logger, s, inventory, fullInventory)
	assert.NoError(t, err)
}

func TestSinkInventory_FileOutputError(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "test-store.db")

	logger := logrus.New()
	cfg := &options.RunOptions{}
	cfg.Facter.Store.Path = storePath
	cfg.Facter.Sink.Output.Type = "file"
	cfg.Facter.Sink.Output.Format = "json"
	cfg.Facter.Sink.Output.OutputDirectory = "/invalid/path/does/not/exist"
	cfg.Facter.Sink.Output.OutputFilename = "output.json"

	s, err := store.NewBoltInventoryStore(cfg.Facter.Store.Path)
	assert.NoError(t, err)

	fullInventory := &schema.HostInventory{
		Hostname: "test-host-error",
	}
	inventory := &schema.InventoryRequest{
		Content: &schema.InventoryRequest_Full{Full: fullInventory},
	}

	err = SinkInventory(cfg, logger, s, inventory, fullInventory)
	assert.Error(t, err)
}

func TestSinkInventory_StoreDeleteOnError(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "test-store.db")

	logger := logrus.New()
	cfg := &options.RunOptions{}
	cfg.Facter.Store.Path = storePath
	cfg.Facter.Sink.Output.Type = "file"
	cfg.Facter.Sink.Output.Format = "json"
	cfg.Facter.Sink.Output.OutputDirectory = "/invalid/directory"
	cfg.Facter.Sink.Output.OutputFilename = "output.json"

	s, err := store.NewBoltInventoryStore(cfg.Facter.Store.Path)
	assert.NoError(t, err)

	fullInventory := &schema.HostInventory{
		Hostname: "test-delete-on-error",
	}

	// First save something to the store
	err = s.Save("test-delete-on-error", fullInventory)
	assert.NoError(t, err)

	inventory := &schema.InventoryRequest{
		Content: &schema.InventoryRequest_Full{Full: fullInventory},
	}

	// This should fail to export and delete the store entry
	err = SinkInventory(cfg, logger, s, inventory, fullInventory)
	assert.Error(t, err)
}

func TestSinkInventory_ClearsRevisionMetadataOnError(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "test-store.db")

	logger := logrus.New()
	cfg := &options.RunOptions{}
	cfg.Facter.Store.Path = storePath
	cfg.Facter.Sink.Output.Type = "file"
	cfg.Facter.Sink.Output.Format = "json"
	cfg.Facter.Sink.Output.OutputDirectory = "/invalid/directory"
	cfg.Facter.Sink.Output.OutputFilename = "output.json"

	s, err := store.NewBoltInventoryStore(cfg.Facter.Store.Path)
	assert.NoError(t, err)

	fullInventory := &schema.HostInventory{
		Hostname: "test-delete-revision-on-error",
	}
	inventory := &schema.InventoryRequest{
		Content: &schema.InventoryRequest_Full{Full: fullInventory},
	}
	revision := &schema.InventoryRevisionEnvelope{
		RevisionId: "test-delete-revision-on-error-1",
		Sequence:   1,
		Hostname:   "test-delete-revision-on-error",
		SourceType: schema.SourceType_SOURCE_TYPE_FULL,
		StateHash:  "hash-1",
		Payload:    &schema.InventoryRevisionEnvelope_Full{Full: fullInventory},
	}

	err = s.Save("test-delete-revision-on-error", fullInventory)
	assert.NoError(t, err)
	err = s.SaveRevision("test-delete-revision-on-error", revision)
	assert.NoError(t, err)

	err = SinkInventory(cfg, logger, s, inventory, fullInventory)
	assert.Error(t, err)

	_, err = s.GetRevision("test-delete-revision-on-error")
	assert.Error(t, err)
}

func TestSinkInventory_DeltaInventory(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "test-store.db")
	outputPath := filepath.Join(tmpDir, "output.json")

	logger := logrus.New()
	cfg := &options.RunOptions{}
	cfg.Facter.Store.Path = storePath
	cfg.Facter.Sink.Output.Type = "file"
	cfg.Facter.Sink.Output.Format = "json"
	cfg.Facter.Sink.Output.OutputDirectory = tmpDir
	cfg.Facter.Sink.Output.OutputFilename = "output.json"

	s, err := store.NewBoltInventoryStore(cfg.Facter.Store.Path)
	assert.NoError(t, err)

	fullInventory := &schema.HostInventory{
		Hostname: "test-delta",
	}

	// Create delta inventory
	deltaInventory := &schema.HostDeltaInventory{
		Hostname: "test-delta",
	}
	inventory := &schema.InventoryRequest{
		Content: &schema.InventoryRequest_Delta{Delta: deltaInventory},
	}

	err = SinkInventory(cfg, logger, s, inventory, fullInventory)
	assert.NoError(t, err)

	// Verify output file was created
	_, err = os.Stat(outputPath)
	assert.NoError(t, err)
}
