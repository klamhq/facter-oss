package system

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/klamhq/facter-oss/pkg/models"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

// GetSystem fill System structure
func GetSystem() *models.System {
	return &models.System{
		Memory: getMemory(),
		Host:   getInfoStat(),
		CPU:    getCPU(),
		Load:   getLoad(),
		Disk:   getDisk(),
		Uptime: getUptime(),
	}
}

// getUptime retrieves the system's uptime in seconds using gopsutil.
func getUptime() uint64 {
	uptime, _ := host.Uptime()
	return uptime
}

// getMemory retrieves the system's memory statistics using gopsutil.
func getMemory() mem.VirtualMemoryStat {
	memoryStat, _ := mem.VirtualMemory()
	return *memoryStat
}

// getInfoStat retrieves the system's host information using gopsutil.
func getInfoStat() host.InfoStat {
	infoStat, _ := host.Info()
	return *infoStat
}

// getCPU retrieves the CPU information using gopsutil.
func getCPU() []cpu.InfoStat {
	cpuStats, _ := cpu.Info()
	return cpuStats
}

// getLoad retrieves the system's load average statistics using gopsutil.
func getLoad() load.AvgStat {
	avgStat, _ := load.Avg()
	return *avgStat
}

// getDiskUUIDMap retrieves the UUIDs of disks on the system.
// It handles both Linux and macOS systems, using appropriate commands to gather the UUIDs.
func getDiskUUIDMap() map[string]string {
	diskUUIDs := make(map[string]string)

	if runtime.GOOS == "linux" {
		files, err := os.ReadDir("/dev/disk/by-uuid/")
		if err != nil {
			return diskUUIDs
		}

		for _, file := range files {
			linkPath := "/dev/disk/by-uuid/" + file.Name()
			resolved, err := os.Readlink(linkPath)
			if err != nil {
				continue
			}

			device := filepath.Join("/dev", filepath.Base(resolved))
			diskUUIDs[device] = file.Name()
		}
	} else if runtime.GOOS == "darwin" {
		// macOS: diskutil list -plist | grep VolumeUUID
		output, err := exec.Command("diskutil", "info", "-all").Output()
		if err != nil {
			return diskUUIDs
		}
		lines := strings.Split(string(output), "\n")
		var currentDisk string
		for _, line := range lines {
			if strings.HasPrefix(line, "   Device Node:") {
				currentDisk = strings.TrimSpace(strings.TrimPrefix(line, "   Device Node:"))
			}
			if strings.HasPrefix(line, "   Volume UUID:") {
				uuid := strings.TrimSpace(strings.TrimPrefix(line, "   Volume UUID:"))
				diskUUIDs[currentDisk] = uuid
			}
		}
	}
	return diskUUIDs
}

// getDisk retrieves the disk information, including partitions and their usage statistics.
// It uses the gopsutil disk package to gather partition information and disk usage statistics.
func getDisk() []models.Disk {
	partitions, _ := disk.Partitions(false)
	usageMap := make(map[string]*disk.UsageStat)
	disksMap := make(map[string]*models.Disk)
	diskUUIDMap := getDiskUUIDMap()

	for _, p := range partitions {
		if !strings.Contains(p.Device, "loop") {
			usage, err := disk.Usage(p.Mountpoint)
			if err == nil {
				usageMap[p.Device] = usage
			}

			if _, exists := disksMap[p.Device]; !exists {
				disksMap[p.Device] = &models.Disk{
					Device:     p.Device,
					UUID:       diskUUIDMap[p.Device],
					Partitions: []models.DiskPartition{},
				}
			}
		}

	}

	for _, p := range partitions {
		usage := usageMap[p.Device]
		if usage == nil {
			continue
		}
		disk := disksMap[p.Device]
		disk.Partitions = append(disk.Partitions, models.DiskPartition{
			Mountpoint:  p.Mountpoint,
			Fstype:      p.Fstype,
			Total:       usage.Total,
			Free:        usage.Free,
			Used:        usage.Used,
			UsedPercent: usage.UsedPercent,
		})
	}

	var result []models.Disk
	for _, d := range disksMap {
		result = append(result, *d)
	}
	return result
}
