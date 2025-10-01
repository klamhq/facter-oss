package packages

import (
	"context"

	"github.com/klamhq/facter-oss/pkg/agent/collectors/packages"
	"github.com/klamhq/facter-oss/pkg/options"
	"github.com/sirupsen/logrus"

	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
)

type PackagesCollectorImpl struct {
	log *logrus.Logger
	cfg *options.PackagesOptions
}

func New(log *logrus.Logger, cfg *options.PackagesOptions) *PackagesCollectorImpl {

	return &PackagesCollectorImpl{
		log: log,
		cfg: cfg,
	}
}

func (c *PackagesCollectorImpl) Name() string { return "packages" }

func (c *PackagesCollectorImpl) CollectPackages(ctx context.Context) ([]*schema.Package, error) {
	c.log.Info("Crafting packages")

	items, err := packages.NewPackagesGrabber(ctx, c.log)
	if err != nil {
		return nil, err
	}
	pkgs := make([]*schema.Package, 0, len(items))
	for _, pkg := range items {

		pkgProto := &schema.Package{
			Name:              pkg.Name,
			Version:           pkg.Version,
			Architecture:      pkg.Architecture,
			Description:       pkg.Description,
			UpgradableVersion: pkg.UpgradableVersion,
			IsUpToDate:        pkg.IsUpToDate,
		}
		pkgs = append(pkgs, pkgProto)

	}
	return pkgs, nil

}
