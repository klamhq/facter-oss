package packages

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/shirou/gopsutil/host"
	"github.com/sirupsen/logrus"
)

// binPathPackageAssociation act as a cache.
// Each time we want too fetch corresponding package from a binary, it's took 200ms.
// So this cache should speedup the process.
// Exemple
// path(/usr/bin/apache) -> string(Apache)
type binPathPackageAssociation map[string]string

type PackageExtractor struct {
	Bin    string
	Args   string
	PostFn func(string) string
	logger *logrus.Logger
	// Cache layer, association between binary path and installed package.
	PathPkgCache binPathPackageAssociation
}

func NewPackageExtractor(logger *logrus.Logger) (*PackageExtractor, error) {
	hostInfo, _ := host.Info()
	pkgExtract := &PackageExtractor{logger: logger, PathPkgCache: make(binPathPackageAssociation)}
	switch os := hostInfo.Platform; os {
	case "ubuntu":
		{
			pkgExtract.Bin = "dpkg"
			pkgExtract.Args = "-S"
			pkgExtract.PostFn = func(in string) string {
				arr := strings.Split(string(in), ":")
				if len(arr) == 2 {
					return arr[0]
				}
				return ""
			}

		}
	case "debian":
	case "linuxmint":
		{
			pkgExtract.Bin = "dpkg"
			pkgExtract.Args = "-S"
			pkgExtract.PostFn = func(in string) string {
				arr := strings.Split(string(in), ":")
				if len(arr) == 2 {
					return arr[0]
				}
				return ""
			}

		}
	case "fedora":
	case "rocky":
		{
			pkgExtract.Bin = "rpm"
			pkgExtract.Args = "-qf"
			pkgExtract.PostFn = func(str string) string {
				return str
			}
		}
	case "centos":
		{
			pkgExtract.Bin = "rpm"
			pkgExtract.Args = "-qf"
			pkgExtract.PostFn = func(str string) string {
				return str
			}
		}
	case "arch":
		{
			pkgExtract.Bin = "pacman"
			pkgExtract.Args = "-Qoq"
			pkgExtract.PostFn = func(str string) string {
				return str
			}
		}
	default:
		return nil, fmt.Errorf("no package extractor for %s system", os)
	}
	return pkgExtract, nil
}

func (p *PackageExtractor) GetPackage(exe string) string {
	start := time.Now()
	var pkgName string
	// Return cached value
	if pkgName, inCache := p.PathPkgCache[exe]; inCache {
		p.logger.WithFields(logrus.Fields{"package": pkgName}).Debugf("cached value returned")
		return pkgName
	}

	logPath := strings.Join([]string{p.Bin, p.Args, exe}, " ")
	cmd := exec.Command(p.Bin, p.Args, exe)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "LANG=C")
	var rawOutput []byte
	rawOutput, err := cmd.CombinedOutput()
	if err != nil {
		p.logger.Debugf("Encountered error during command %s execution, output %s with error %s:", cmd, rawOutput, err)
		pkgName = "unknown"
	} else {
		// Apply system specific post-processing function
		if p.PostFn != nil {
			pkgName = p.PostFn(string(rawOutput))
		}
	}

	p.PathPkgCache[exe] = pkgName

	elapsed := time.Since(start).Round(time.Microsecond)
	p.logger.WithFields(logrus.Fields{"elapsed": elapsed, "cmd": logPath}).Debugf("fetching package for exe %s", exe)

	if pkgName == "" {
		p.logger.WithFields(logrus.Fields{"elapsed": elapsed, "cmd": logPath}).Warnf("Unable to get package name from exe %s", exe)
		pkgName = "unknown"
	}
	return pkgName
}
