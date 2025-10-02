package sink

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/klamhq/facter-oss/pkg/options"
	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func tempDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "exportToFileTest")
	fmt.Println(dir)
	assert.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestExportToFile_ProtobufFormat(t *testing.T) {
	dir := tempDir(t)
	filename := "test.pb"
	cfg := &options.RunOptions{}
	cfg.Facter.Sink.Enabled = true
	cfg.Facter.Sink.Output.Type = "file"
	cfg.Facter.Sink.Output.Format = "protobuf"
	cfg.Facter.Sink.Output.OutputDirectory = dir
	cfg.Facter.Sink.Output.OutputFilename = filename

	logger := logrus.New()
	inventory := &schema.HostInventory{}
	inventory.Hostname = "test-host"
	inventoryMsg := &schema.InventoryRequest{Content: &schema.InventoryRequest_Full{Full: inventory}}

	err := exportToFile(inventoryMsg, logger, cfg)
	assert.NoError(t, err)

	dest := filepath.Join(dir, filename)
	_, err = os.Stat(dest)
	assert.NoError(t, err)
}

func TestExportToFile_JSONFormat(t *testing.T) {
	dir := tempDir(t)
	filename := "test.json"
	cfg := &options.RunOptions{}
	cfg.Facter.Sink.Enabled = true
	cfg.Facter.Sink.Output.Type = "file"
	cfg.Facter.Sink.Output.Format = "json"
	cfg.Facter.Sink.Output.OutputDirectory = dir
	cfg.Facter.Sink.Output.OutputFilename = filename
	logger := logrus.New()
	inventory := &schema.HostInventory{}
	inventory.Hostname = "json-host"
	inventoryMsg := &schema.InventoryRequest{Content: &schema.InventoryRequest_Full{Full: inventory}}

	err := exportToFile(inventoryMsg, logger, cfg)
	assert.NoError(t, err)

	dest := filepath.Join(dir, filename)
	data, err := os.ReadFile(dest)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "json-host")
}

func TestExportToFile_DisabledSink(t *testing.T) {
	cfg := &options.RunOptions{}
	cfg.Facter.Sink.Enabled = false
	logger := logrus.New()
	inventory := &schema.HostInventory{}
	inventory.Hostname = "test-host"
	inventoryMsg := &schema.InventoryRequest{Content: &schema.InventoryRequest_Full{Full: inventory}}

	err := exportToFile(inventoryMsg, logger, cfg)
	assert.NoError(t, err)
}

func TestExportToFile_InvalidMarshal(t *testing.T) {
	// Pass nil to cause proto.Marshal to fail
	cfg := &options.RunOptions{}
	cfg.Facter.Sink.Enabled = true
	cfg.Facter.Sink.Output.Type = "file"
	cfg.Facter.Sink.Output.Format = "protobuf"
	cfg.Facter.Sink.Output.OutputDirectory = tempDir(t)
	cfg.Facter.Sink.Output.OutputFilename = "invalid.pb"

	logger := logrus.New()
	err := exportToFile(nil, logger, cfg)
	assert.Error(t, err)
}

func TestExportToFile_InvalidPath(t *testing.T) {
	// Pass nil to cause proto.Marshal to fail
	cfg := &options.RunOptions{}
	cfg.Facter.Sink.Enabled = true
	cfg.Facter.Sink.Output.Type = "file"
	cfg.Facter.Sink.Output.Format = "protobuf"
	cfg.Facter.Sink.Output.OutputDirectory = "/invalid/path"
	cfg.Facter.Sink.Output.OutputFilename = "invalid.pb"
	inventory := &schema.HostInventory{}
	inventory.Hostname = "test-host"
	inventoryMsg := &schema.InventoryRequest{Content: &schema.InventoryRequest_Full{Full: inventory}}

	logger := logrus.New()
	err := exportToFile(inventoryMsg, logger, cfg)
	assert.Error(t, err)
}
