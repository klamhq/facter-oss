package inventory

import (
	"context"
	"os/user"
	"runtime"
	"time"

	"github.com/klamhq/facter-oss/pkg/agent/collect/applications"
	"github.com/klamhq/facter-oss/pkg/agent/collect/networks"
	"github.com/klamhq/facter-oss/pkg/agent/collect/packages"
	"github.com/klamhq/facter-oss/pkg/agent/collect/platform"
	"github.com/klamhq/facter-oss/pkg/agent/collect/process"
	"github.com/klamhq/facter-oss/pkg/agent/collect/ssh"
	"github.com/klamhq/facter-oss/pkg/agent/collect/systemservices"
	"github.com/klamhq/facter-oss/pkg/agent/collect/users"
	"github.com/klamhq/facter-oss/pkg/agent/store"
	"github.com/klamhq/facter-oss/pkg/models"
	"github.com/klamhq/facter-oss/pkg/options"
	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type Builder struct {
	Log          *logrus.Logger
	Cfg          options.RunOptions
	maxParallel  int
	SystemGather *models.System

	Now    func() time.Time
	WhoAmI func() (string, error)
	Store  store.InventoryStore

	Platform       platform.PlatformCollector
	Packages       packages.PackagesCollector
	Applications   applications.ApplicationsCollector
	SystemServices systemservices.SystemServicesCollector
	Networks       networks.NetworksCollector
	Users          users.UsersCollector
	Processes      process.ProcessCollector
	SSHInfos       ssh.SSHInfosCollector
}

func NewBuilder(cfg options.RunOptions, systemGather *models.System, logger *logrus.Logger) *Builder {
	// If users collection is disabled, disable SSH collection too
	if !cfg.Facter.Inventory.User.Enabled {
		cfg.Facter.Inventory.SSH.Enabled = false
	}
	b := &Builder{
		Log:          logger,
		Cfg:          cfg,
		SystemGather: systemGather,
		maxParallel:  runtime.NumCPU(),
		Now:          time.Now,
		Store: func() store.InventoryStore {
			s, err := store.NewBoltInventoryStore(cfg.Facter.Store.Path)
			if err != nil {
				logrus.WithError(err).Fatal("Unable to create inventory boltdb store")
			}
			return s
		}(),
		// Default collectors set to nil, will be initialized later if enabled in the config
		Platform: nil,
		WhoAmI: func() (string, error) {
			u, err := user.Current()
			if err != nil {
				return "", err
			}
			return u.Name, nil
		},
	}

	return b
}

func (b *Builder) Build(ctx context.Context) (*schema.HostInventory, error) {
	inv := &schema.HostInventory{
		CreatedAt: time.Now().Format(time.RFC3339),
		Network:   &schema.Network{},
		Metadata:  &schema.Metadata{FacterVersion: "0.1.0", RunningDate: time.Now().Format(time.RFC3339)},
	}
	inv.Hostname = b.SystemGather.Host.Hostname
	if u, err := user.Current(); err == nil {
		inv.Metadata.RunningUser = u.Name
	}
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(b.maxParallel)

	// Collectors are called in a specific order to handle dependencies
	// between them (e.g. users -> ssh, platform -> systemd services)
	// The order is as follows:
	// 1) Platform
	var (
		platform     *schema.Platform
		users        []*schema.User
		pkgs         []*schema.Package
		services     []*schema.SystemdService
		processes    []*schema.Process
		sshKeyAccess []*schema.SshKeyAccess
		sshKeyInfos  []*schema.SshKeyInfo
		knownHosts   []*schema.KnownHost
		networks     *schema.Network
		apps         []*schema.Application
		err          error
	)

	g.Go(func() error {
		if b.Cfg.Facter.Inventory.Platform.Enabled && b.Platform != nil {
			start := time.Now()
			platform, err = b.Platform.CollectPlatform(ctx)
			if err != nil {
				b.Log.WithError(err).Error("platform")
			}
			defer func(n string) { b.Log.WithField("collector", n).WithField("dur", time.Since(start)).Info("done") }("platform")
		}
		g.Go(func() error {
			if b.Cfg.Facter.Inventory.SystemdService.Enabled {
				start := time.Now()
				services, err = b.SystemServices.CollectSystemServices(ctx, platform.InitSystem)
				if err != nil {
					b.Log.WithError(err).Error("initsystem services")
				}
				defer func(n string) { b.Log.WithField("collector", n).WithField("dur", time.Since(start)).Info("done") }("initsystem services")

			}
			return nil
		})
		return nil
	})

	// 2) Packages
	g.Go(func() error {
		if b.Cfg.Facter.Inventory.Packages.Enabled {
			start := time.Now()
			pkgs, err = b.Packages.CollectPackages(ctx)
			if err != nil {
				b.Log.WithError(err).Error("packages")
			}
			defer func(n string) { b.Log.WithField("collector", n).WithField("dur", time.Since(start)).Info("done") }("packages")

		}
		return nil
	})

	// 3) Applications
	g.Go(func() error {
		if b.Cfg.Facter.Inventory.Applications.Enabled {
			start := time.Now()
			apps, err = b.Applications.CollectApplications(ctx)
			if err != nil {
				b.Log.WithError(err).Error("applications")
			}
			defer func(n string) { b.Log.WithField("collector", n).WithField("dur", time.Since(start)).Info("done") }("applications")
		}
		return nil
	})

	// 4) initsystem services
	// g.Go(func() error {
	// 	if b.Cfg.Facter.Inventory.SystemdService.Enabled {
	// 		start := time.Now()
	// 		services, err = b.SystemServices.CollectSystemServices(ctx, platform.InitSystem)
	// 		if err != nil {
	// 			b.Log.WithError(err).Error("initsystem services")
	// 		}
	// 		defer func(n string) { b.Log.WithField("collector", n).WithField("dur", time.Since(start)).Info("done") }("initsystem services")

	// 	}
	// 	return nil
	// })

	// 5) Networks
	g.Go(func() error {
		if b.Cfg.Facter.Inventory.Networks.Enabled {
			start := time.Now()
			networks, err = b.Networks.CollectNetworks(ctx)
			if err != nil {
				b.Log.WithError(err).Error("networks")
			}
			defer func(n string) { b.Log.WithField("collector", n).WithField("dur", time.Since(start)).Info("done") }("networks")

		}
		return nil
	})

	// 6) Users
	g.Go(func() error {
		if b.Cfg.Facter.Inventory.User.Enabled {
			start := time.Now()
			users, err := b.Users.CollectUsers(ctx)
			if err != nil {
				b.Log.WithError(err).Error("users")
			}
			// SSH
			if b.Cfg.Facter.Inventory.SSH.Enabled {
				b.SSHInfos = ssh.New(b.Log, &b.Cfg.Facter.Inventory.SSH)
				sshKeyAccess, knownHosts, sshKeyInfos, err = b.SSHInfos.CollectSSHInfos(ctx, users)
				if err != nil {
					b.Log.WithError(err).Error("ssh")
				}
			}
			defer func(n string) { b.Log.WithField("collector", n).WithField("dur", time.Since(start)).Info("done") }("users and ssh")

		}
		return nil
	})

	// 7) Processes
	g.Go(func() error {
		if b.Cfg.Facter.Inventory.Process.Enabled {
			start := time.Now()
			processes, err = b.Processes.CollectProcess(ctx)
			if err != nil {
				b.Log.WithError(err).Error("processes")
			}
			defer func(n string) { b.Log.WithField("collector", n).WithField("dur", time.Since(start)).Info("done") }("processes")

		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return nil, err
	}

	inv.Platform = platform
	inv.Application = apps
	inv.Packages = pkgs
	inv.Users = users
	inv.Network = networks
	inv.Processes = processes
	inv.SshKeyAccess = sshKeyAccess
	inv.KnownHost = knownHosts
	inv.SshKeyInfo = sshKeyInfos
	inv.SystemdService = services

	return inv, nil
}

func (b *Builder) ManageDelta(fullInventory *schema.HostInventory) (*schema.InventoryRequest, *schema.HostInventory) {
	// Retrieve the old inventory from BoltDB
	previous, err := b.Store.Get(fullInventory.Hostname)
	var result *schema.InventoryRequest

	// Check if previous inventory exists, compute delta and send it else send full inventory
	if err != nil || previous == nil {
		b.Log.Info("No previous inventory, computing full inventory")
		result = &schema.InventoryRequest{
			Content: &schema.InventoryRequest_Full{Full: fullInventory},
		}
		return result, fullInventory
	} else {
		b.Log.Info("Previous inventory found, computing delta")
		delta := ComputeDelta(previous, fullInventory, b.Log)
		if IsDeltaEmpty(delta) {
			b.Log.Info("No changes detected, nothing to send")
			return nil, nil
		}
		delta.UpdatedAt = time.Now().Format(time.RFC3339)
		result = &schema.InventoryRequest{
			Content: &schema.InventoryRequest_Delta{Delta: delta},
		}
		b.Log.Debugf("Send this delta %s", result)
		return result, fullInventory
	}
}
