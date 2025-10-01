package sink

import (
	"os"
	"path"

	"github.com/klamhq/facter-oss/pkg/options"
	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func exportToFile(inventoryMsg *schema.InventoryRequest, logger *logrus.Logger, cfg *options.RunOptions) error {
	bin, err := proto.Marshal(inventoryMsg)
	if err != nil {
		logger.WithError(err).Fatal("Unable to marshal protobuf message")
		return err
	}

	if cfg.Facter.Sink.Enabled {
		if cfg.Facter.Sink.Output.Type == "file" {
			if cfg.Facter.Sink.Output.Format == "json" {
				protobufJSON := protojson.Format(inventoryMsg)
				bin = []byte(protobufJSON)
			}
			dest := path.Join(cfg.Facter.Sink.Output.OutputDirectory, cfg.Facter.Sink.Output.OutputFilename)
			if err := os.WriteFile(dest, bin, 0644); err != nil {
				logger.WithError(err).Fatal("Unable to write message")
				return err
			}
			logger.Infof("File saved to %s", dest)
			return nil
		}

	}
	return nil
}
