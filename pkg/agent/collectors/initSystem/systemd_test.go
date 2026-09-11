package initSystem

import (
	"os"
	"os/exec"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var originalExecCommand = exec.Command

func mockExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := originalExecCommand(os.Args[0], cs...)
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	args := os.Args
	switch args[3] {
	case "systemctl":
		switch args[4] {
		case "list-units":
			os.Stdout.WriteString("nginx.service loaded active running Nginx HTTP Server\nsshd.service loaded active running OpenSSH Daemon\n")
		case "show":
			os.Stdout.WriteString(`Description=Fake Service
LoadState=loaded
ActiveState=active
SubState=running
ExecMainPID=123
TasksCurrent=4
MemoryCurrent=2048
CPUUsageNSec=987654321
ControlGroup=/system.slice/fake.service
Wants=network.target
Requires=basic.target
After=network.target
Before=shutdown.target
`)
		case "is-enabled":
			os.Stdout.WriteString("enabled\n")
		}
	default:
		os.Stderr.WriteString("unknown command")
	}
	os.Exit(0)
}

func TestListAllServices(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	logger := logrus.New()
	services, err := listAllServices(logger)
	assert.NoError(t, err)
	assert.Len(t, services, 2)
	assert.Equal(t, "nginx.service", services[0])
	assert.Equal(t, "sshd.service", services[1])
}

func TestGetServiceDetails(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	svc, err := getServiceDetails("fake.service")
	assert.NoError(t, err)
	assert.Equal(t, "fake.service", svc.Name)
	assert.Equal(t, "Fake Service", svc.Description)
	assert.Equal(t, "loaded", svc.Loaded)
	assert.Equal(t, "active", svc.Active)
	assert.Equal(t, "running", svc.SubState)
	assert.Equal(t, int64(123), svc.PID)
	assert.Equal(t, int64(4), svc.Tasks)
	assert.Equal(t, int64(2048), svc.MemoryBytes)
	assert.Equal(t, int64(987654321), svc.CPUUsageNsec)
	assert.True(t, svc.Enabled)
	assert.Equal(t, []string{"network.target"}, svc.Wants)
	assert.Equal(t, []string{"basic.target"}, svc.Requires)
	assert.Equal(t, []string{"network.target"}, svc.After)
	assert.Equal(t, []string{"shutdown.target"}, svc.Before)
}

func TestGatherSystemdInfo(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	logger := logrus.New()
	services, err := GatherSystemdInfo(logger)
	assert.NoError(t, err)
	assert.Len(t, services, 2)
	assert.Equal(t, "nginx.service", services[0].Name)
}

func TestListAllServices_CommandError(t *testing.T) {
	execCommand = func(command string, args ...string) *exec.Cmd {
		cmd := exec.Command("false")
		return cmd
	}
	defer func() { execCommand = exec.Command }()

	logger := logrus.New()
	services, err := listAllServices(logger)
	assert.Error(t, err)
	assert.Nil(t, services)
}

func TestGetServiceDetails_CommandError(t *testing.T) {
	execCommand = func(command string, args ...string) *exec.Cmd {
		cmd := exec.Command("false")
		return cmd
	}
	defer func() { execCommand = exec.Command }()

	svc, err := getServiceDetails("fake.service")
	assert.Error(t, err)
	assert.Nil(t, svc)
}
