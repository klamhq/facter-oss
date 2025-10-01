package initSystem

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/klamhq/facter-oss/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetSystemDInit(t *testing.T) {

	path := "/tmp/systemInit/systemd"
	target := "systemd"
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Printf("could not create dir %s: %v", path, err)
	}
	if err := os.WriteFile(filepath.Join(path, target), []byte("systemd\n"), 0644); err != nil {
		log.Printf("could not write file %s: %v", filepath.Join(path, target), err)
	}
	symlink := filepath.Join(path, "symlink")
	if err := os.Symlink(path, symlink); err != nil {
		log.Printf("could not create symlink from %s to %s: %v", symlink, path, err)
	}
	var factory utils.LoggerFactory = &utils.DefaultLoggerFactory{}
	logger := factory.New(logrus.ErrorLevel)
	res := GetSystemInit(logger, "/tmp/systemInit/systemd/symlink")
	assert.Equal(t, res, "systemd")
}

func TestGetSysvInit(t *testing.T) {
	path := "/tmp/systemInit/sysvinit"
	target := "sysvinit"
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Printf("could not create dir %s: %v", path, err)
	}
	if err := os.WriteFile(filepath.Join(path, target), []byte("sysvinit\n"), 0644); err != nil {
		log.Printf("could not write file %s: %v", filepath.Join(path, target), err)
	}
	symlink := filepath.Join(path, "symlink")
	if err := os.Symlink(path, symlink); err != nil {
		log.Printf("could not create symlink from %s to %s: %v", symlink, path, err)
	}
	var factory utils.LoggerFactory = &utils.DefaultLoggerFactory{}
	logger := factory.New(logrus.ErrorLevel)
	res := GetSystemInit(logger, "/tmp/systemInit/sysvinit/symlink")
	assert.Equal(t, res, "sysvinit")
}

func TestGetSystemInit(t *testing.T) {
	const initCheckPathError = ""
	var factory utils.LoggerFactory = &utils.DefaultLoggerFactory{}
	logger := factory.New(logrus.ErrorLevel)
	res := GetSystemInit(logger, initCheckPathError)
	assert.Empty(t, res)
}
