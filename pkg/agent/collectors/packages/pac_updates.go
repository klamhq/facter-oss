package packages

import (
	"context"
	"regexp"
	"strings"

	"github.com/klamhq/facter-oss/pkg/utils"
)

var pacmanRE = regexp.MustCompile(`^(?P<name>[^\s]+)\s+(?P<version>[^\s]+)\s+->\s+(?P<new>[^\s]+)$`)

func parsePacmanOutput(out []byte) map[string]string {
	lines := strings.Split(string(out), "\n")
	result := make(map[string]string)
	for _, line := range lines {
		match := pacmanRE.FindStringSubmatch(line)
		if match != nil {
			name := match[1]
			newVersion := match[3]
			result[name] = newVersion
		}
	}
	return result
}

func GetPacmanUpgradableMap(ctx context.Context) (map[string]string, error) {
	out, err := utils.RunCmd(ctx, "pacman", "-Qu")
	if err != nil {
		return nil, err
	}
	return parsePacmanOutput(out), nil
}
