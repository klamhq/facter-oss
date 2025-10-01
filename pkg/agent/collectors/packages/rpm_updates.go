package packages

import (
	"context"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// Map: packageName -> upgradable version
func GetRpmUpgradableMap(ctx context.Context) (map[string]string, error) {
	// Try dnf first, fallback to yum
	var output []byte
	var err error

	cmd := exec.CommandContext(ctx, "dnf", "check-update")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "LANG=C")

	output, err = cmd.CombinedOutput()
	if err != nil && !isDNFNoUpdates(err) {
		// fallback
		cmd = exec.CommandContext(ctx, "yum", "check-update")
		output, err = cmd.CombinedOutput()
		if err != nil && !isDNFNoUpdates(err) {
			return nil, err
		}
	}

	re := regexp.MustCompile(`^(?P<name>[^\s]+)\s+(?P<version>[^\s]+)\s+`)

	lines := strings.Split(string(output), "\n")
	result := map[string]string{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		match := re.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		// Extract name and version
		rawName := match[1]
		name := strings.Split(rawName, ".")[0] // remove .x86_64
		version := match[2]
		result[name] = version
	}

	return result, nil
}

// isDNFNoUpdates dnf/yum return exit code 100 if updates are available (not an error in our case)
func isDNFNoUpdates(err error) bool {
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode() == 100
	}
	return false
}
