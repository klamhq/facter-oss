package utils

import (
	"os/exec"

	"github.com/sirupsen/logrus"
)

func CheckBinInstalled(logger *logrus.Logger, bin string) bool {
	_, err := exec.LookPath(bin)
	if err != nil {
		logger.Errorf("%s is not installed", bin)
		return false
	}
	return true
}
