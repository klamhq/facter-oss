package network

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/klamhq/facter-oss/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func writeTempResolvConf(t *testing.T, content string) string {
	t.Helper()
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "resolv.conf")
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write temp resolv.conf: %v", err)
	}
	return tmpFile
}

func TestGetDnsConf_WithEnvPath(t *testing.T) {
	content := "nameserver 1.1.1.1\n"
	tmpFile := writeTempResolvConf(t, content)
	os.Setenv("DNSCONF_PATH", tmpFile)
	defer os.Unsetenv("DNSCONF_PATH")

	conf := GetDnsConf()
	assert.NotNil(t, conf)
	assert.Contains(t, conf.Servers, "1.1.1.1")
}

func TestGetDnsConf_DefaultPath(t *testing.T) {
	// Backup and replace /etc/resolv.conf if possible
	origPath := defaultPath
	origContent, _ := os.ReadFile(origPath)
	tmpFile := writeTempResolvConf(t, "nameserver 8.8.8.8\n")
	// Move the temp file to /etc/resolv.conf if running as root, else skip
	if utils.IsRoot() {
		defer os.WriteFile(origPath, origContent, 0644)
		os.Rename(tmpFile, origPath)
		conf := GetDnsConf()
		assert.NotNil(t, conf)
		assert.True(t,
			strings.Contains(strings.Join(conf.Servers, ","), "8.8.8.8") || strings.Contains(strings.Join(conf.Servers, ","), "169.254.169.254"),
			"expected string to contain either %q or %q, got %q", "abc", "edf", conf.Servers,
		)
	} else {
		t.Skip("Skipping test that modifies /etc/resolv.conf (requires root)")
	}
}

func TestGetDnsConf_InvalidPath(t *testing.T) {
	os.Setenv("DNSCONF_PATH", "/nonexistent/path/resolv.conf")
	defer os.Unsetenv("DNSCONF_PATH")

	conf := GetDnsConf()
	assert.Nil(t, conf)
}
