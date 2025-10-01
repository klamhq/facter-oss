package packages

import (
	"context"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func GetBrewUpgradableMap(ctx context.Context) (map[string]string, error) {
	cmd := exec.CommandContext(ctx, "brew", "outdated", "--verbose")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "LANG=C")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	// Ex: zeromq (4.3.4) < 4.3.5_2
	re := regexp.MustCompile(`^(?P<name>[^\s]+)\s+\([^)]+\)\s+(<|!=|>|<=|>=|==)\s+(?P<version>[^\s]+)$`)
	lines := strings.Split(string(out), "\n")
	result := make(map[string]string)

	for _, line := range lines {
		match := re.FindStringSubmatch(line)
		if match != nil {
			name := match[1]
			version := match[3]
			result[name] = version
		}
	}

	return result, nil
}
