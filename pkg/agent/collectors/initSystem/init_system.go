package initSystem

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// https://www.debian-fr.org/t/systemd-sysvinit-savoir-lequel-est/64018

const systemdReturn = "systemd"
const otherReturn = "sysvinit"

// GetSystemInit Read the symlink of /sbin/init and return value of constant `systemdReturn` or `otherReturn`.
func GetSystemInit(logger *logrus.Logger, path string) string {
	path, err := os.Readlink(path)
	if err != nil {
		logger.WithError(err).Warnf("unable to find init system, %v", err)
		return ""
	}
	if strings.Contains(path, "systemd") {
		return systemdReturn
	}
	return otherReturn
}
