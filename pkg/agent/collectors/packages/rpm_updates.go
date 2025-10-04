package packages

import (
	"context"
	"regexp"
	"strings"

	"github.com/klamhq/facter-oss/pkg/utils"
)

func GetRpmUpgradableMap(ctx context.Context) (map[string]string, error) {
	var (
		out []byte
		err error
	)

	out, err = utils.RunCmd(ctx, "dnf", "check-update")
	if err != nil && !isDNFNoUpdates(err) {
		out, err = utils.RunCmd(ctx, "yum", "check-update")
		if err != nil && !isDNFNoUpdates(err) {
			return nil, err
		}
	}

	re := regexp.MustCompile(`^(?P<name>[^\s]+)\s+(?P<version>[^\s]+)\s+`)
	lines := strings.Split(string(out), "\n")
	result := make(map[string]string)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		match := re.FindStringSubmatch(line)
		if match == nil {
			continue
		}

		rawName := match[1]
		name := strings.Split(rawName, ".")[0] // remove .x86_64, .noarch, etc.
		version := match[2]
		result[name] = version
	}

	return result, nil
}

// isDNFNoUpdates : dnf/yum return exit code 100 when update are available
func isDNFNoUpdates(err error) bool {
	type exitCoder interface{ ExitCode() int }
	if ec, ok := err.(exitCoder); ok {
		return ec.ExitCode() == 100
	}
	return false
}
