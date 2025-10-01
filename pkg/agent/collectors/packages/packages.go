package packages

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/klamhq/facter-oss/pkg/models"
	"github.com/klamhq/facter-oss/pkg/utils"
	"github.com/sirupsen/logrus"
)

// NewPackagesGrabber return correct PackageGrabber instance for current system
func NewPackagesGrabber(ctx context.Context, logger *logrus.Logger) ([]*models.Package, error) {
	// When adding new package extract, mind to update this array.
	registeredExtract := []func(ctx context.Context, logger *logrus.Logger) ([]*models.Package, error){
		NewPackageRpmConfig,
		NewPackageDebConfig,
		NewPackagePacConfig,
		NewPackageHomebrewConfig,
		NewPackageAptConfig,
	}
	for _, fn := range registeredExtract {

		foundPkgs, err := fn(ctx, logger)
		if err != nil {
			logger.Error("Error executing command: ", err)
			logger.WithError(err).Debug("unable to read packages list")
			continue
		}
		return foundPkgs, nil
	}

	return nil, fmt.Errorf("unable to find a package extract usable on this system")
}

type PackageExtractorInterface interface {
	Extract() ([]*models.Package, error)
	IsApplicable() bool
}

type PackageExtractRegexpBased struct {
	config *PackageParserConfig
}

// Extract is responsible of running an external program and extract useful content defined by ParseRegExp.
func (e *PackageExtractRegexpBased) Extract(ctx context.Context, logger *logrus.Logger) ([]*models.Package, error) {
	if !e.config.IsApplicable() {
		return nil, fmt.Errorf("%s is not present", e.config.PackageManagerBin)
	}

	var command *exec.Cmd
	args := []string{}
	if e.config.PackageManagerArgs != "" {
		args = append(args, e.config.PackageManagerArgs)
	}
	if len(e.config.PackageManagerAdditionalArgs) > 0 {
		args = append(args, e.config.PackageManagerAdditionalArgs...)
	}
	command = exec.CommandContext(ctx, e.config.PackageManagerBin, args...)
	command.Env = append(os.Environ(), "LANG=C")

	// Exécution
	output, err := command.CombinedOutput()
	logger.Debugf("Package extractor command output: %s", output)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")

	// Canal pour collecter les résultats
	ch := make(chan *models.Package, len(lines))
	wg := sync.WaitGroup{}

	for _, line := range lines {
		wg.Add(1)
		go func(line string) {
			defer wg.Done()
			if extracted := e.parseSingleLine(line); extracted != nil {
				ch <- extracted
			}
		}(line)
	}

	// Fermeture du canal quand toutes les goroutines sont terminées
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Collecte des résultats dans un slice local
	result := []*models.Package{}
	for pkg := range ch {
		result = append(result, pkg)
	}

	return result, nil
}

// parseSingleLine is responsible of parsing a single line extracted from the package manager.
func (e *PackageExtractRegexpBased) parseSingleLine(line string) *models.Package {

	item := models.Package{}
	if e.config.PackageIgnoreCondition != nil {
		if !e.config.PackageIgnoreCondition(line) {
			return nil
		}
	}

	if matches := e.config.ParseRegExp.FindStringSubmatch(line); matches != nil {
		for i := range e.config.ParseRegExp.SubexpNames() {
			if i == 0 {
				// First elemnt is the full catch, we have to ignore this part.
				continue
			}
			/*
				This part is used to set dynamically structure's fields by ascending order of extracted value by the regexp.
				Be carefull when defining a new type of structure.
			*/

			reflect.ValueOf(&item).Elem().Field(i - 1).SetString(matches[i])
		}
	}
	if len(item.Name) == 0 {
		return nil
	}
	return &item
}

// PackageParserConfig is used to transport config for a package manager.
type PackageParserConfig struct {
	ParseRegExp                  *regexp.Regexp
	PackageManagerBin            string
	PackageManagerArgs           string
	PackageManagerAdditionalArgs []string
	// This is used to define a anonymous function used to ignore some line.
	// Example: in case of debian package, we have to ignore non beginning line by `ii`
	PackageIgnoreCondition func(string) bool
}

// IsApplicable should return true if the current PackageParser is compatible with current system.
func (e *PackageParserConfig) IsApplicable() bool {
	_, err := exec.LookPath(e.PackageManagerBin)
	return err == nil
}

type PackageParser interface {
}

func addUpgradablePackage(pkgs []*models.Package, upgradableMap map[string]string) []*models.Package {
	for i, pkg := range pkgs {
		if upgradeVer, ok := upgradableMap[pkg.Name]; ok {
			pkgs[i].UpgradableVersion = upgradeVer
			pkgs[i].IsUpToDate = false

		} else {
			pkgs[i].IsUpToDate = true
		}
	}
	return pkgs
}

// NewPackageHomebrewConfig provide a configuration for homebrew based package system
func NewPackageHomebrewConfig(ctx context.Context, logger *logrus.Logger) ([]*models.Package, error) {
	config := newHomebrewConfig()
	parser := &PackageExtractRegexpBased{
		config: config,
	}
	if !utils.IsRoot() {
		pkgs, err := parser.Extract(ctx, logger)
		if err != nil {
			return nil, err
		}
		upgradableMap, err := GetBrewUpgradableMap(ctx)
		if err != nil {
			return nil, err
		}
		pkgs = addUpgradablePackage(pkgs, upgradableMap)

		return pkgs, nil
	}
	return nil, fmt.Errorf("homebrew packages can be extracted only when not root")
}

// NewPackagePacConfig provide a configuration for pacman based package system.
func NewPackagePacConfig(ctx context.Context, logger *logrus.Logger) ([]*models.Package, error) {
	config := newPacConfig()
	parser := &PackageExtractRegexpBased{
		config: config,
	}
	pkgs, err := parser.Extract(ctx, logger)
	if err != nil {
		return nil, err
	}

	// Ajoute les versions à jour
	upgradableMap, err := GetPacmanUpgradableMap(ctx)
	if err != nil {
		return nil, err
	}
	pkgs = addUpgradablePackage(pkgs, upgradableMap)
	return pkgs, nil
}

// newHomebrewConfig return homebrew package parser config
func newHomebrewConfig() *PackageParserConfig {
	config := &PackageParserConfig{
		ParseRegExp:                  regexp.MustCompile(`^(?P<name>[\S]*)\s(?P<version>[\S]*)$`),
		PackageIgnoreCondition:       nil,
		PackageManagerBin:            "brew",
		PackageManagerArgs:           "ls",
		PackageManagerAdditionalArgs: []string{"--versions", "--quiet"},
	}
	return config
}

// newPacConfig return arch package parser config
func newPacConfig() *PackageParserConfig {
	config := &PackageParserConfig{
		ParseRegExp:            regexp.MustCompile(`^(?P<name>[\S]*)\s(?P<version>[\S]*)$`),
		PackageIgnoreCondition: nil,
		PackageManagerBin:      "pacman",
		PackageManagerArgs:     "-Q",
	}
	return config
}

// NewPackageDebConfig provide a configuration for deb based package system
func NewPackageDebConfig(ctx context.Context, logger *logrus.Logger) ([]*models.Package, error) {
	config := newDebConfig()
	parser := &PackageExtractRegexpBased{
		config: config,
	}

	pkgs, err := parser.Extract(ctx, logger)
	if err != nil {
		return nil, err
	}

	// Ajout des versions upgradables
	upgradableMap, err := GetAptUpgradableMap(ctx)
	if err != nil {
		return nil, err
	}
	pkgs = addUpgradablePackage(pkgs, upgradableMap)

	return pkgs, nil
}

// NewDebConfig return dpkg package parser config
func newDebConfig() *PackageParserConfig {
	config := &PackageParserConfig{
		//		ParseRegExp:            regexp.MustCompile(`ii\s*(?P<name>[\w-.]*)(?::amd64)?\s*(?P<version>[\w.\d-:+~]*)\s*(?P<architecture>[\w]*)\s*(?P<description>.*)$`),
		ParseRegExp:                  regexp.MustCompile(`ii\s*(?P<name>[\w-.]*)(?::amd64)?\s*(?P<version>[\w.\d-:+~]*)\s*(?P<architecture>[\w]*)\s*(?P<description>.*)$`),
		PackageIgnoreCondition:       nil,
		PackageManagerBin:            "dpkg",
		PackageManagerArgs:           "-l",
		PackageManagerAdditionalArgs: []string{"--no-pager"},
	}
	return config
}

// NewPackageRpmConfig provide a configuration for rpm based package system
func NewPackageRpmConfig(ctx context.Context, logger *logrus.Logger) ([]*models.Package, error) {
	config := newRpmConfig()
	parser := &PackageExtractRegexpBased{
		config: config,
	}

	pkgs, err := parser.Extract(ctx, logger)
	if err != nil {
		return nil, err
	}

	upgradableMap, err := GetRpmUpgradableMap(ctx)
	if err != nil {
		return nil, err
	}
	pkgs = addUpgradablePackage(pkgs, upgradableMap)

	return pkgs, nil
}

// newRpmConfig return rpm package parser config
func newRpmConfig() *PackageParserConfig {
	config := &PackageParserConfig{
		ParseRegExp:            regexp.MustCompile(`^(?P<name>.*)-(?P<version>\d{1,3}(?:\.|-)\d{1,3}(?:-\S{1,3})?\.\S{2,})\.(?P<architecture>\S+)$`),
		PackageIgnoreCondition: nil,
		PackageManagerBin:      "rpm",
		PackageManagerArgs:     "-qa",
	}
	return config
}

// NewPackageAptConfig provide a configuration for apt based package system
func NewPackageAptConfig(ctx context.Context, logger *logrus.Logger) ([]*models.Package, error) {
	config := newAptConfig()
	parser := &PackageExtractRegexpBased{
		config: config,
	}

	pkgs, err := parser.Extract(ctx, logger)
	if err != nil {
		return nil, err
	}

	// Ajout des versions upgradables
	upgradableMap, err := GetAptUpgradableMap(ctx)
	if err != nil {
		return nil, err
	}
	pkgs = addUpgradablePackage(pkgs, upgradableMap)

	return pkgs, nil
}

// newAptConfig return apt package parser config
func newAptConfig() *PackageParserConfig {
	config := &PackageParserConfig{
		ParseRegExp:            regexp.MustCompile(`^(?P<name>[^\s]+)/[^\s]+\s+(?:(?:\d+:)?)(?P<version>[^\s]+)\s+(?P<architecture>[^\s]+)$`),
		PackageIgnoreCondition: nil,
		PackageManagerBin:      "apt",
		PackageManagerArgs:     "list",
	}
	return config
}
