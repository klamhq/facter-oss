package initSystem

import (
	"bufio"
	"bytes"
	"os/exec"
	"strconv"
	"strings"

	"github.com/klamhq/facter-oss/pkg/models"
	"github.com/sirupsen/logrus"
)

// GatherSystemdInfo gathers information about systemd services and their dependencies.
func GatherSystemdInfo(logger *logrus.Logger) ([]models.SystemdService, error) {
	services, err := listAllServices(logger)
	if err != nil {
		logger.WithError(err).Error("Failed to list systemd services")
		return nil, err
	}
	systemdServiceInfo := make([]models.SystemdService, 0, len(services))
	for _, service := range services {
		details, err := getServiceDetails(service)
		if err != nil {
			logger.WithError(err).Errorf("Failed to get details for service %s", service)
			continue
		}
		systemdServiceInfo = append(systemdServiceInfo, *details)
	}
	return systemdServiceInfo, nil
}

// listAllServices lists all systemd services on the system.
// It returns a slice of service names or an error if the command fails.
func listAllServices(logger *logrus.Logger) ([]string, error) {
	cmd := exec.Command("systemctl", "list-units", "--type=service", "--no-pager", "--all", "--no-legend")
	output, err := cmd.Output()
	if err != nil {
		logger.WithError(err).Error("Failed to execute systemctl command")
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var services []string

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 0 {
			services = append(services, fields[0])
		}
	}

	return services, nil
}

// getServiceDetails retrieves detailed information about a specific systemd service.
// It returns a SystemdService struct or an error if the command fails.
func getServiceDetails(name string) (*models.SystemdService, error) {
	cmd := exec.Command("systemctl", "show", name)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	service := &models.SystemdService{Name: name}
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		key, val := parts[0], parts[1]

		switch key {
		case "Description":
			service.Description = val
		case "LoadState":
			service.Loaded = val
		case "ActiveState":
			service.Active = val
		case "SubState":
			service.SubState = val
		case "ExecMainPID":
			pid, _ := strconv.Atoi(val)
			service.PID = int64(pid)
		case "TasksCurrent":
			tasks, _ := strconv.Atoi(val)
			service.Tasks = int64(tasks)
		case "MemoryCurrent":
			service.MemoryBytes, _ = strconv.ParseInt(val, 10, 64)
		case "CPUUsageNSec":
			service.CPUUsageNsec, _ = strconv.ParseInt(val, 10, 64)
		case "ControlGroup":
			service.CGroup = val
		case "Wants":
			service.Wants = strings.Fields(val)
		case "Requires":
			service.Requires = strings.Fields(val)
		case "After":
			service.After = strings.Fields(val)
		case "Before":
			service.Before = strings.Fields(val)
		}
	}

	// Enabled ?
	cmd = exec.Command("systemctl", "is-enabled", name)
	if output, err := cmd.Output(); err == nil {
		service.Enabled = strings.TrimSpace(string(output)) == "enabled"
	}

	return service, nil
}
