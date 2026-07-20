package platform

import (
	"context"
	"testing"

	"github.com/klamhq/facter-oss/pkg/models"
	"github.com/klamhq/facter-oss/pkg/options"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCollectPlatform_WithHardware(t *testing.T) {
	cfg := options.RunOptions{}
	cfg.Facter.Inventory.Platform.Hardware.Enabled = true
	ctx := context.Background()
	
	path := models.SystemPaths{}
	sys := &models.System{
		CPU: []cpu.InfoStat{
			{
				ModelName: "Test CPU",
				Mhz:       2400.0,
			},
		},
		Memory: mem.VirtualMemoryStat{
			Total:     16000000000,
			Used:      8000000000,
			SwapTotal: 2000000000,
		},
		Disk: []models.Disk{
			{
				Device: "/dev/sda",
				UUID:   "test-uuid",
				Partitions: []models.DiskPartition{
					{
						Mountpoint:  "/",
						Fstype:      "ext4",
						Total:       100000000000,
						Used:        50000000000,
						Free:        50000000000,
						UsedPercent: 50.0,
					},
				},
			},
		},
		Host: host.InfoStat{},
	}
	
	p := New(logrus.New(), &cfg.Facter.Inventory.Platform, path, sys)
	res, err := p.CollectPlatform(ctx)
	
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.Hardware)
	assert.NotNil(t, res.Hardware.Cpu)
	assert.Equal(t, "Test CPU", res.Hardware.Cpu.Model)
	assert.NotNil(t, res.Hardware.Memory)
	assert.NotEmpty(t, res.Hardware.Disk)
}

func TestCollectPlatform_WithKernel(t *testing.T) {
	cfg := options.RunOptions{}
	cfg.Facter.Inventory.Platform.Kernel.Enabled = true
	ctx := context.Background()
	
	path := models.SystemPaths{}
	sys := &models.System{
		Host: host.InfoStat{
			KernelVersion: "5.4.0-42-generic",
		},
	}
	
	p := New(logrus.New(), &cfg.Facter.Inventory.Platform, path, sys)
	res, err := p.CollectPlatform(ctx)
	
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.Kernel)
	assert.Equal(t, "5.4.0-42-generic", res.Kernel.Kernel)
}

func TestCollectPlatform_WithOS(t *testing.T) {
	cfg := options.RunOptions{}
	cfg.Facter.Inventory.Platform.Os.Enabled = true
	ctx := context.Background()
	
	path := models.SystemPaths{}
	sys := &models.System{
		Host: host.InfoStat{
			Platform:        "ubuntu",
			PlatformVersion: "20.04",
			PlatformFamily:  "debian",
		},
	}
	
	p := New(logrus.New(), &cfg.Facter.Inventory.Platform, path, sys)
	res, err := p.CollectPlatform(ctx)
	
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.Os)
	assert.Equal(t, "ubuntu", res.Os.Name)
	assert.Equal(t, "20.04", res.Os.Version)
	assert.Equal(t, "debian", res.Os.Family)
}

func TestCollectPlatform_WithVirtualization(t *testing.T) {
	cfg := options.RunOptions{}
	cfg.Facter.Inventory.Platform.Virtualization.Enabled = true
	ctx := context.Background()
	
	path := models.SystemPaths{}
	sys := &models.System{
		Host: host.InfoStat{
			VirtualizationSystem: "kvm",
			VirtualizationRole:   "guest",
		},
	}
	
	p := New(logrus.New(), &cfg.Facter.Inventory.Platform, path, sys)
	res, err := p.CollectPlatform(ctx)
	
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.Virtualization)
	assert.Equal(t, "kvm", res.Virtualization.System)
	assert.Equal(t, "guest", res.Virtualization.Role)
}

func TestCollectPlatform_AllFeaturesEnabled(t *testing.T) {
	cfg := options.RunOptions{}
	cfg.Facter.Inventory.Platform.Hardware.Enabled = true
	cfg.Facter.Inventory.Platform.Kernel.Enabled = true
	cfg.Facter.Inventory.Platform.Os.Enabled = true
	cfg.Facter.Inventory.Platform.Virtualization.Enabled = true
	ctx := context.Background()
	
	path := models.SystemPaths{
		InitCheckPath: "/sbin/init",
		MachineID:     "/etc/machine-id",
		MachineUUID:   "/sys/class/dmi/id/product_uuid",
	}
	sys := &models.System{
		Uptime: 123456,
		CPU: []cpu.InfoStat{
			{
				ModelName: "Intel Core i7",
				Mhz:       3200.0,
			},
		},
		Memory: mem.VirtualMemoryStat{
			Total: 16000000000,
			Used:  8000000000,
		},
		Host: host.InfoStat{
			KernelVersion:        "5.4.0",
			Platform:             "ubuntu",
			PlatformVersion:      "20.04",
			PlatformFamily:       "debian",
			VirtualizationSystem: "kvm",
			VirtualizationRole:   "guest",
		},
	}
	
	p := New(logrus.New(), &cfg.Facter.Inventory.Platform, path, sys)
	res, err := p.CollectPlatform(ctx)
	
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.Hardware)
	assert.NotNil(t, res.Kernel)
	assert.NotNil(t, res.Os)
	assert.NotNil(t, res.Virtualization)
	assert.NotNil(t, res.Identifier)
	assert.Equal(t, uint64(123456), res.Uptime)
}

func TestCollectPlatform_EmptySystem(t *testing.T) {
	cfg := options.RunOptions{}
	cfg.Facter.Inventory.Platform.Hardware.Enabled = true
	ctx := context.Background()
	
	path := models.SystemPaths{}
	sys := &models.System{}
	
	p := New(logrus.New(), &cfg.Facter.Inventory.Platform, path, sys)
	res, err := p.CollectPlatform(ctx)
	
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestFillHardware_NoCPU(t *testing.T) {
	cfg := options.RunOptions{}
	cfg.Facter.Inventory.Platform.Hardware.Enabled = true
	ctx := context.Background()
	
	path := models.SystemPaths{}
	sys := &models.System{
		CPU: []cpu.InfoStat{}, // Empty CPU list
		Memory: mem.VirtualMemoryStat{
			Total: 8000000000,
		},
	}
	
	p := New(logrus.New(), &cfg.Facter.Inventory.Platform, path, sys)
	res, err := p.CollectPlatform(ctx)
	
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.Hardware)
	// CPU should be nil when no CPU info available
	assert.Nil(t, res.Hardware.Cpu)
}

func TestFillHardware_WithMultipleCPUs(t *testing.T) {
	cfg := options.RunOptions{}
	cfg.Facter.Inventory.Platform.Hardware.Enabled = true
	ctx := context.Background()
	
	path := models.SystemPaths{}
	sys := &models.System{
		CPU: []cpu.InfoStat{
			{ModelName: "CPU 1", Mhz: 2400.0},
			{ModelName: "CPU 2", Mhz: 2400.0},
			{ModelName: "CPU 3", Mhz: 2400.0},
		},
		Memory: mem.VirtualMemoryStat{},
	}
	
	p := New(logrus.New(), &cfg.Facter.Inventory.Platform, path, sys)
	res, err := p.CollectPlatform(ctx)
	
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.Hardware.Cpu)
	assert.Equal(t, uint32(3), res.Hardware.Cpu.Core)
	assert.Equal(t, "CPU 1", res.Hardware.Cpu.Model)
}
