package systemservices

import (
	"context"

	"github.com/klamhq/facter-oss/pkg/agent/collectors/initSystem"
	"github.com/klamhq/facter-oss/pkg/options"
	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
	"github.com/sirupsen/logrus"
)

type SystemServicesCollectorImpl struct {
	log *logrus.Logger
	cfg *options.SystemdServiceOptions
}

func New(log *logrus.Logger, cfg *options.SystemdServiceOptions) *SystemServicesCollectorImpl {

	return &SystemServicesCollectorImpl{
		log: log,
		cfg: cfg,
	}
}

func (c *SystemServicesCollectorImpl) CollectSystemServices(ctx context.Context, initsystem string) ([]*schema.SystemdService, error) {
	c.log.Info("Crafting initsystem services")
	systemServices := []*schema.SystemdService{}

	if initsystem == "systemd" {
		items, err := initSystem.GatherSystemdInfo(c.log)
		if err != nil {
			c.log.WithError(err).Error("Failed to gather systemd services")
			return nil, err
		}
		for _, service := range items {
			systemServices = append(systemServices, &schema.SystemdService{
				Name:         service.Name,
				Description:  service.Description,
				Loaded:       service.Loaded,
				Active:       service.Active,
				SubState:     service.SubState,
				Enabled:      service.Enabled,
				Pid:          service.PID,
				Tasks:        service.Tasks,
				MemoryBytes:  service.MemoryBytes,
				CpuUsageNsec: service.CPUUsageNsec,
				Cgroup:       service.CGroup,
				Requires:     service.Requires,
				After:        service.After,
				Before:       service.Before,
				Wants:        service.Wants,
			})
		}
	}
	return systemServices, nil
}
