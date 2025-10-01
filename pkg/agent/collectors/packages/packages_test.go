package packages

import (
	"context"
	"runtime"
	"testing"

	"github.com/klamhq/facter-oss/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewPackagesGrabber(t *testing.T) {
	var factory utils.LoggerFactory = &utils.DefaultLoggerFactory{}
	logger := factory.New(logrus.ErrorLevel)
	ctx := context.Background()
	if utils.IsRoot() && runtime.GOOS == "darwin" {
		grabber, err := NewPackagesGrabber(ctx, logger)
		assert.Error(t, err)
		assert.Nil(t, grabber)
		return
	}
	grabber, err := NewPackagesGrabber(ctx, logger)
	assert.NoError(t, err)
	assert.NotEmpty(t, grabber)
}

func TestNewPackageDebConfig(t *testing.T) {
	fixture := `ii  tzdata:amd64                  2018g-0+deb9u1        all          time zone and daylight-saving time data -libs`
	config := newDebConfig()
	parser := &PackageExtractRegexpBased{
		config: config,
	}
	pkg := *parser.parseSingleLine(fixture)
	assert.Equal(t, "tzdata", pkg.Name)
	assert.Equal(t, "2018g-0+deb9u1", pkg.Version)
	assert.Equal(t, "all", pkg.Architecture)
	assert.Equal(t, "time zone and daylight-saving time data -libs", pkg.Description)
}

func TestNewAltPackageDebConfig(t *testing.T) {
	fixture := `rpm-plugin-selinux-4.14.2.1-2.fc29.x86_64`
	config := newRpmConfig()
	parser := &PackageExtractRegexpBased{
		config: config,
	}
	pkg := *parser.parseSingleLine(fixture)
	assert.Equal(t, "rpm-plugin-selinux", pkg.Name)
	assert.Equal(t, "4.14.2.1-2.fc29", pkg.Version)
	assert.Equal(t, "x86_64", pkg.Architecture)
	assert.Equal(t, "", pkg.Description)
}

func TestNewHomebrewPackageConfig(t *testing.T) {
	fixture := `ansible 2.9.2`
	config := newHomebrewConfig()
	parser := &PackageExtractRegexpBased{
		config: config,
	}
	pkg := *parser.parseSingleLine(fixture)
	assert.Equal(t, "ansible", pkg.Name)
	assert.Equal(t, "2.9.2", pkg.Version)
	assert.Equal(t, "", pkg.Architecture)
	assert.Equal(t, "", pkg.Description)
}

func TestNewArchPackageConfig(t *testing.T) {
	fixture := `ansible 2.9.2`
	config := newPacConfig()
	parser := &PackageExtractRegexpBased{
		config: config,
	}
	pkg := *parser.parseSingleLine(fixture)
	assert.Equal(t, "ansible", pkg.Name)
	assert.Equal(t, "2.9.2", pkg.Version)
	assert.Equal(t, "", pkg.Architecture)
	assert.Equal(t, "", pkg.Description)
}
