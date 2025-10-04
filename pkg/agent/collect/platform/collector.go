// pkg/collect/platform/collector.go
package platform

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/klamhq/facter-oss/pkg/agent/collectors/initSystem"
	"github.com/klamhq/facter-oss/pkg/agent/collectors/machineIdentifier"
	"github.com/klamhq/facter-oss/pkg/models"
	"github.com/klamhq/facter-oss/pkg/options"
	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
)

type PlatformCollectorImpl struct {
	log          *logrus.Logger
	cfg          *options.PlatformOptions
	paths        models.SystemPaths
	systemGather *models.System
}

func New(log *logrus.Logger, cfg *options.PlatformOptions, paths models.SystemPaths, systemGather *models.System) *PlatformCollectorImpl {

	return &PlatformCollectorImpl{
		log:          log,
		cfg:          cfg,
		paths:        paths,
		systemGather: systemGather,
	}
}
func (c *PlatformCollectorImpl) CollectPlatform(ctx context.Context) (*schema.Platform, error) {
	c.log.Info("Crafting platform")

	p := &schema.Platform{
		Uptime:     c.systemGather.Uptime,
		InitSystem: initSystem.GetSystemInit(c.log, c.paths.InitCheckPath),
	}

	// Hardware
	if c.cfg.Hardware.Enabled {
		c.fillHardware(p)
	}

	// Kernel
	if c.cfg.Kernel.Enabled {
		p.Kernel = &schema.Kernel{
			Kernel: c.systemGather.Host.KernelVersion,
		}
	}

	// OS
	if c.cfg.Os.Enabled {
		p.Os = &schema.Os{
			Name:    c.systemGather.Host.Platform,
			Version: c.systemGather.Host.PlatformVersion,
			Family:  c.systemGather.Host.PlatformFamily,
		}
	}

	// Virtualization
	if c.cfg.Virtualization.Enabled {
		p.Virtualization = &schema.Virtualization{
			System: c.systemGather.Host.VirtualizationSystem,
			Role:   c.systemGather.Host.VirtualizationRole,
		}
	}

	// Identifier
	if id, err := machineIdentifier.GetMachineID(
		c.log, c.paths.MachineID, c.paths.MachineUUID,
	); err != nil {
		c.log.Error("fetching MachineID failed: ", err, ". Setting identifier to unknown.")
		p.Identifier = &schema.Identifier{
			MachineId: "unknown",
			Uuid:      "unknown",
		}
	} else {
		p.Identifier = &schema.Identifier{
			MachineId: id.MachineId,
			Uuid:      id.UUID,
		}
	}
	return p, nil
}

// ---- helpers privés ----

func (c *PlatformCollectorImpl) fillHardware(p *schema.Platform) {
	mem := c.systemGather.Memory
	cpus := c.systemGather.CPU
	disks := c.systemGather.Disk

	if p.Hardware == nil {
		p.Hardware = &schema.Hardware{}
	}

	// CPU
	if len(cpus) > 0 {
		main := cpus[0]
		p.Hardware.Cpu = &schema.Cpu{
			Model: main.ModelName,
			Core:  uint32(len(cpus)),
			Mhz:   float32(main.Mhz),
		}
	}

	// Memory
	p.Hardware.Memory = &schema.Memory{
		Total: mem.Total,
		Used:  mem.Used,
		Swap:  mem.SwapTotal,
	}

	// Disks
	out := make([]*schema.Disk, 0, len(disks))
	for _, d := range disks {
		var parts []*schema.DiskPartition
		for _, p := range d.Partitions {
			parts = append(parts, &schema.DiskPartition{
				Mountpoint:  p.Mountpoint,
				FsType:      p.Fstype,
				Total:       p.Total,
				Used:        p.Used,
				Free:        p.Free,
				UsedPercent: p.UsedPercent,
			})
		}
		out = append(out, &schema.Disk{
			Device:     d.Device,
			Uuid:       d.UUID,
			Partitions: parts,
		})
	}
	p.Hardware.Disk = out
}
