package packages

import (
	"context"
	"regexp"
	"strings"

	"github.com/klamhq/facter-oss/pkg/utils"
)

func GetBrewUpgradableMap(ctx context.Context) (map[string]string, error) {
	out, err := utils.RunCmd(ctx, "brew", "outdated", "--verbose")
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`^(?P<name>[^\s]+)\s+\([^)]+\)\s+(?:<|!=|>|<=|>=|==)\s+(?P<version>[^\s]+)$`)
	lines := strings.Split(string(out), "\n")
	result := make(map[string]string)

	for _, line := range lines {
		match := re.FindStringSubmatch(line)
		if match != nil {
			name := match[1]
			version := match[2]
			result[name] = version
		}
	}
	return result, nil
}
