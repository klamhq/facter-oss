package process

import (
	"context"

	"github.com/klamhq/facter-oss/pkg/agent/collectors/process"
	"github.com/klamhq/facter-oss/pkg/options"
	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
	"github.com/sirupsen/logrus"
)

type ProcessCollectorImpl struct {
	log *logrus.Logger
	cfg *options.ProcessOptions
}

func New(log *logrus.Logger, cfg *options.ProcessOptions) *ProcessCollectorImpl {

	return &ProcessCollectorImpl{
		log: log,
		cfg: cfg,
	}
}

func (c *ProcessCollectorImpl) CollectProcess(ctx context.Context) ([]*schema.Process, error) {
	c.log.Info("Crafting process")
	collectedProcess, err := process.Processes(c.log)
	if err != nil {
		c.log.Errorf("Error during crafting processes %v", err)
	}
	return collectedProcess, nil
}
