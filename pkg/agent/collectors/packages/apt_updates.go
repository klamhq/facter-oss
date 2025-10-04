package packages

import (
	"context"
	"regexp"
	"strings"

	"github.com/klamhq/facter-oss/pkg/utils"
)

var aptLineRE = regexp.MustCompile(`^(?P<name>[^\s]+)/[^\s]+\s+(?:(?:\d+:)?)(?P<version>[^\s]+)\s+[^\s]+\s+\[upgradable from: (?:(?:\d+:)?)(?P<from_version>[^\]]+)\]$`)

func parseAptUpgradableOutput(output []byte) map[string]string {
	lines := strings.Split(string(output), "\n")
	result := map[string]string{}
	for _, line := range lines {
		match := aptLineRE.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		name := match[1]
		version := match[2]
		result[name] = version
	}
	return result
}

func GetAptUpgradableMap(ctx context.Context) (map[string]string, error) {
	output, err := utils.RunCmd(ctx, "apt", "list", "--upgradable")
	if err != nil {
		return nil, err
	}
	return parseAptUpgradableOutput(output), nil
}
