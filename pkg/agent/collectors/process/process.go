package process

import (
	"strings"

	"github.com/klamhq/facter-oss/pkg/agent/collectors/packages"
	"github.com/klamhq/facter-oss/pkg/models"
	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
	"github.com/shirou/gopsutil/process"
	"github.com/sirupsen/logrus"
)

// Processes return all connections in protobuf schema
func Processes(logger *logrus.Logger) ([]*schema.Process, error) {
	pkgExtractor, err := packages.NewPackageExtractor(logger)
	if err != nil {
		logger.Errorf("Error getting package extractor: %v", err)
	}
	procs, err := getProcess(logger, pkgExtractor)
	if err != nil {
		return nil, err
	}
	processes := make([]*schema.Process, 0, len(procs))
	for _, p := range procs {
		proc := schema.Process{
			Pid:           p.PID,
			Name:          p.Name,
			Package:       &schema.Package{Name: p.Package.Name},
			Username:      p.Username,
			Cmdline:       p.Cmdline,
			Terminal:      p.Terminal,
			Exe:           p.Exe,
			CreateTime:    p.CreateTime,
			Parent:        p.Parent,
			Status:        p.Status,
			MemoryPercent: p.MemPercent,
			CpuPercent:    p.CpuPercent,
		}
		processes = append(processes, &proc)
	}
	return processes, nil
}

// getProcess returns all processes in protobuf schema
// It retrieves process information such as PID, name, and package association.
func getProcess(logger *logrus.Logger, pkgExtractor *packages.PackageExtractor) ([]*models.Process, error) {
	proc, err := process.Processes()
	if err != nil {
		logger.Errorf("Error retrieving processes: %v", err)
		return nil, err
	}

	processes := make([]*models.Process, 0, len(proc))
	for _, p := range proc {
		pid := p.Pid

		name, err := p.Name()
		if err != nil {
			name = "unknown"
		}
		name = strings.Fields(name)[0]

		exe, err := p.Exe()
		if err != nil {
			exe = "unknown"
		}
		createTime, err := p.CreateTime()
		if err != nil {
			createTime = 0
		}
		cmdline, err := p.Cmdline()
		if err != nil {
			cmdline = "unknown"
		}

		username, err := p.Username()
		if err != nil {
			username = "unknown"
		}
		status, err := p.Status()
		if err != nil {
			status = "unknown"
		}
		parent, err := p.Parent()
		if err != nil {
			parent = nil
		}
		ppid := int32(0)
		if parent != nil {
			ppid = parent.Pid
		}
		terminal, err := p.Terminal()
		if err != nil {
			terminal = ""
		}
		memPercent, err := p.MemoryPercent()
		if err != nil {
			memPercent = 0.0
		}

		cpuPercent, err := p.CPUPercent()
		if err != nil {
			cpuPercent = 0.0
		}

		packageName := "unknown"
		if exe != "unknown" && pkgExtractor != nil {
			packageName = pkgExtractor.GetPackage(exe)
		}

		processes = append(processes, &models.Process{
			Name:       name,
			PID:        int64(pid),
			Username:   username,
			Status:     status,
			CreateTime: createTime,
			Parent:     int64(ppid),
			Cmdline:    cmdline,
			Terminal:   terminal,
			Exe:        exe,
			Package: models.Package{
				Name: packageName,
			},
			CpuPercent: cpuPercent,
			MemPercent: float64(memPercent),
		})

	}
	return processes, nil
}
