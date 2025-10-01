package packages

import (
	"context"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// Map: packageName -> upgradable version
func GetAptUpgradableMap(ctx context.Context) (map[string]string, error) {
	cmd := exec.CommandContext(ctx, "apt", "list", "--upgradable")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "LANG=C")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`^(?P<name>[^\s]+)/[^\s]+\s+(?:(?:\d+:)?)(?P<version>[^\s]+)\s+[^\s]+\s+\[upgradable from: (?:(?:\d+:)?)(?P<from_version>[^\]]+)\]$`)
	lines := strings.Split(string(output), "\n")

	result := map[string]string{}
	for _, line := range lines {
		match := re.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		name := match[1]
		version := match[2]
		result[name] = version
	}

	return result, nil
}
