package packages

import (
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var originalExecCommand = exec.Command

func mockExecCommand(command string, args ...string) *exec.Cmd {
	body := `echo "unknown command" 1>&2; exit 3`

	exe := ""
	if len(args) > 0 {
		exe = args[len(args)-1]
	}

	switch command {
	case "dpkg":
		if len(args) >= 1 && args[0] == "-S" {
			if strings.Contains(exe, "non/existent") {
				body = `echo "dpkg: unable to found file" 1>&2; exit 2`
			} else {
				escaped := strings.ReplaceAll(exe, `'`, `'\''`)
				body = `printf 'apache2: ` + escaped + `\n'`
			}
		}
	case "rpm":
		if len(args) >= 1 && args[0] == "-qf" {
			if strings.Contains(exe, "non/existent") {
				body = `echo "rpm: error" 1>&2; exit 2`
			} else {
				body = `printf 'httpd-2.4.6-80.el7.centos\n'`
			}
		}
	case "pacman":
		if len(args) >= 1 && args[0] == "-Qoq" {
			if strings.Contains(exe, "non/existent") {
				body = `echo "pacman: error" 1>&2; exit 2`
			} else {
				body = `printf 'coolpackage\n'`
			}
		}
	}

	return originalExecCommand("sh", "-c", body)
}

func TestGetPackage_Dpkg_BasedOnString_NoFilesystem(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = originalExecCommand }()

	logger := logrus.New()
	p := &PackageExtractor{
		Bin:  "dpkg",
		Args: "-S",
		PostFn: func(in string) string {
			arr := strings.Split(in, ":")
			if len(arr) == 2 {
				return strings.TrimSpace(arr[0])
			}
			return ""
		},
		logger:       logger,
		PathPkgCache: make(binPathPackageAssociation),
	}

	pkg := p.GetPackage("/usr/bin/apache")
	if pkg != "apache2" {
		t.Fatalf("expected 'apache2', got %q", pkg)
	}

	pkg2 := p.GetPackage("/usr/bin/apache")
	if pkg2 != "apache2" {
		t.Fatalf("expected cached 'apache2', got %q", pkg2)
	}

	pkgErr := p.GetPackage("/non/existent")
	if pkgErr != "unknown" {
		t.Fatalf("expected 'unknown' for /non/existent, got %q", pkgErr)
	}
}

func TestGetPackage_Rpm_BasedOnString_NoFilesystem(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = originalExecCommand }()

	logger := logrus.New()
	p := &PackageExtractor{
		Bin:          "rpm",
		Args:         "-qf",
		PostFn:       func(s string) string { return strings.TrimSpace(s) },
		logger:       logger,
		PathPkgCache: make(binPathPackageAssociation),
	}

	pkg := p.GetPackage("/usr/sbin/httpd")
	if pkg != "httpd-2.4.6-80.el7.centos" {
		t.Fatalf("expected rpm package, got %q", pkg)
	}
}

func TestGetPackage_Pacman_BasedOnString_NoFilesystem(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = originalExecCommand }()

	logger := logrus.New()
	p := &PackageExtractor{
		Bin:          "pacman",
		Args:         "-Qoq",
		PostFn:       func(s string) string { return strings.TrimSpace(s) },
		logger:       logger,
		PathPkgCache: make(binPathPackageAssociation),
	}

	pkg := p.GetPackage("/usr/bin/coolbinary")
	if pkg != "coolpackage" {
		t.Fatalf("expected pacman package 'coolpackage', got %q", pkg)
	}
}

func TestGetPackage_CommandError_ReturnsUnknown_BasedOnString(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = originalExecCommand }()

	logger := logrus.New()
	p := &PackageExtractor{
		Bin:  "dpkg",
		Args: "-S",
		PostFn: func(in string) string {
			arr := strings.Split(in, ":")
			if len(arr) == 2 {
				return strings.TrimSpace(arr[0])
			}
			return ""
		},
		logger:       logger,
		PathPkgCache: make(binPathPackageAssociation),
	}

	pkg := p.GetPackage("/non/existent")
	if pkg != "unknown" {
		t.Fatalf("expected 'unknown' on command error, got %q", pkg)
	}
}

func TestShellAvailable(t *testing.T) {
	if _, err := exec.LookPath("sh"); err != nil {
		t.Skip("no shell available for tests")
	}
}

func TestCacheSpeedQuick(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = originalExecCommand }()

	logger := logrus.New()
	p := &PackageExtractor{
		Bin:  "dpkg",
		Args: "-S",
		PostFn: func(in string) string {
			arr := strings.Split(in, ":")
			if len(arr) == 2 {
				return strings.TrimSpace(arr[0])
			}
			return ""
		},
		logger:       logger,
		PathPkgCache: make(binPathPackageAssociation),
	}

	start := time.Now()
	_ = p.GetPackage("/usr/bin/apache")
	_ = p.GetPackage("/usr/bin/apache")
	if time.Since(start) > time.Second {
		t.Fatalf("cache path unexpectedly slow")
	}
}

func TestNewPackageExtractor(t *testing.T) {
	logger := logrus.New()
	p, err := NewPackageExtractor(logger)
	if runtime.GOOS == "darwin" {
		assert.Nil(t, p)
		assert.Error(t, err)
	} else {
		assert.NotNil(t, p)
		assert.NoError(t, err)
	}
}
