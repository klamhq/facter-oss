package inventory

import (
	"context"
	"fmt"
	"os/user"
	"runtime"
	"sync"
	"time"

	"github.com/klamhq/facter-oss/pkg/agent/collect/applications"
	"github.com/klamhq/facter-oss/pkg/agent/collect/compliance"
	"github.com/klamhq/facter-oss/pkg/agent/collect/networks"
	"github.com/klamhq/facter-oss/pkg/agent/collect/packages"
	"github.com/klamhq/facter-oss/pkg/agent/collect/platform"
	"github.com/klamhq/facter-oss/pkg/agent/collect/process"
	"github.com/klamhq/facter-oss/pkg/agent/collect/ssh"
	"github.com/klamhq/facter-oss/pkg/agent/collect/systemservices"
	"github.com/klamhq/facter-oss/pkg/agent/collect/users"
	"github.com/klamhq/facter-oss/pkg/agent/collect/vulnerability"
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

	Platform            platform.PlatformCollector
	Packages            packages.PackagesCollector
	Applications        applications.ApplicationsCollector
	SystemServices      systemservices.SystemServicesCollector
	Networks            networks.NetworksCollector
	Users               users.UsersCollector
	Processes           process.ProcessCollector
	SSHInfos            ssh.SSHInfosCollector
	ComplianceReport    compliance.ComplianceCollector
	VulnerabilityReport vulnerability.VulnerabilityCollector
}

func newInventoryStore(cfg options.RunOptions) (store.InventoryStore, error) {
	s, err := store.NewBoltInventoryStore(cfg.Facter.Store.Path)
	if err != nil {
		return nil, fmt.Errorf("unable to create inventory boltdb store: %w", err)
	}
	return s, nil
}

func NewBuilder(cfg options.RunOptions, systemGather *models.System, logger *logrus.Logger) (*Builder, error) {
	// If users collection is disabled, disable SSH collection too
	if !cfg.Facter.Inventory.User.Enabled {
		cfg.Facter.Inventory.SSH.Enabled = false
	}
	s, err := newInventoryStore(cfg)
	if err != nil {
		return nil, fmt.Errorf("initializing inventory store: %w", err)
	}
	b := &Builder{
		Log:          logger,
		Cfg:          cfg,
		SystemGather: systemGather,
		maxParallel:  runtime.NumCPU(),
		Now:          time.Now,
		Store:        s,
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

	return b, nil
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
		platform            *schema.Platform
		users               []*schema.User
		pkgs                []*schema.Package
		services            []*schema.SystemdService
		processes           []*schema.Process
		sshKeyAccess        []*schema.SshKeyAccess
		sshKeyInfos         []*schema.SshKeyInfo
		knownHosts          []*schema.KnownHost
		networks            *schema.Network
		apps                []*schema.Application
		complianceReport    *schema.ComplianceReport
		vulnerabilityReport *schema.VulnerabilityReport
		mu                  sync.Mutex
	)

	g.Go(func() error {
		if b.Cfg.Facter.Inventory.Platform.Enabled && b.Platform != nil {
			start := time.Now()
			p, e := b.Platform.CollectPlatform(ctx)
			if e != nil {
				b.Log.WithError(e).Error("platform")
			}
			mu.Lock()
			platform = p
			mu.Unlock()
			defer func(n string) { b.Log.WithField("collector", n).WithField("duration", time.Since(start)).Info("done") }("platform")
		}
		g.Go(func() error {
			if b.Cfg.Facter.Inventory.SystemdService.Enabled {
				pLocal := platform
				start := time.Now()
				s, se := b.SystemServices.CollectSystemServices(ctx, pLocal.InitSystem)
				if se != nil {
					b.Log.WithError(se).Error("initsystem services")
				}
				mu.Lock()
				services = s
				mu.Unlock()
				defer func(n string) { b.Log.WithField("collector", n).WithField("duration", time.Since(start)).Info("done") }("initsystem services")

			}
			return nil
		})
		return nil
	})

	// 2) Packages
	g.Go(func() error {
		if b.Cfg.Facter.Inventory.Packages.Enabled {
			start := time.Now()
			pk, pkerr := b.Packages.CollectPackages(ctx)
			if pkerr != nil {
				b.Log.WithError(pkerr).Error("packages")
			}
			mu.Lock()
			pkgs = pk
			mu.Unlock()
			defer func(n string) { b.Log.WithField("collector", n).WithField("duration", time.Since(start)).Info("done") }("packages")

		}
		return nil
	})

	// 3) Applications
	g.Go(func() error {
		if b.Cfg.Facter.Inventory.Applications.Enabled {
			start := time.Now()
			a, apperr := b.Applications.CollectApplications(ctx)
			if apperr != nil {
				b.Log.WithError(apperr).Error("applications")
			}
			mu.Lock()
			apps = a
			mu.Unlock()
			defer func(n string) { b.Log.WithField("collector", n).WithField("duration", time.Since(start)).Info("done") }("applications")
		}
		return nil
	})

	// 4) Networks
	g.Go(func() error {
		if b.Cfg.Facter.Inventory.Networks.Enabled {
			start := time.Now()
			net, neterr := b.Networks.CollectNetworks(ctx)
			if neterr != nil {
				b.Log.WithError(neterr).Error("networks")
			}
			mu.Lock()
			networks = net
			mu.Unlock()
			defer func(n string) { b.Log.WithField("collector", n).WithField("duration", time.Since(start)).Info("done") }("networks")

		}
		return nil
	})

	// 5) Users
	g.Go(func() error {
		if b.Cfg.Facter.Inventory.User.Enabled {
			start := time.Now()
			u, uerr := b.Users.CollectUsers(ctx)
			if uerr != nil {
				b.Log.WithError(uerr).Error("users")
			}
			mu.Lock()
			users = u
			mu.Unlock()
			// SSH
			if b.Cfg.Facter.Inventory.SSH.Enabled {
				b.SSHInfos = ssh.New(b.Log, &b.Cfg.Facter.Inventory.SSH)
				ska, kh, ski, ssherr := b.SSHInfos.CollectSSHInfos(ctx, users)
				if ssherr != nil {
					b.Log.WithError(ssherr).Error("ssh")
				}
				mu.Lock()
				sshKeyAccess = ska
				knownHosts = kh
				sshKeyInfos = ski
				mu.Unlock()
			}
			defer func(n string) { b.Log.WithField("collector", n).WithField("duration", time.Since(start)).Info("done") }("users and ssh")

		}
		return nil
	})

	// 6) Processes
	g.Go(func() error {
		if b.Cfg.Facter.Inventory.Process.Enabled {
			start := time.Now()
			proc, procerr := b.Processes.CollectProcess(ctx)
			if procerr != nil {
				b.Log.WithError(procerr).Error("processes")
			}
			mu.Lock()
			processes = proc
			mu.Unlock()
			defer func(n string) { b.Log.WithField("collector", n).WithField("duration", time.Since(start)).Info("done") }("processes")

		}
		return nil
	})

	// 7) Compliance
	g.Go(func() error {
		if b.Cfg.Facter.Compliance.Enabled {
			start := time.Now()
			cReport, cReportErr := b.ComplianceReport.CollectCompliance(ctx)
			if cReportErr != nil {
				b.Log.WithError(cReportErr).Error("compliance report")
			}
			mu.Lock()
			complianceReport = cReport
			mu.Unlock()
			defer func(n string) { b.Log.WithField("collector", n).WithField("duration", time.Since(start)).Info("done") }("compliance report")

		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// 8) Vulnerabilities
	if b.Cfg.Facter.Vulnerabilities.Enabled {
		start := time.Now()
		var vReportErr error
		if runtime.GOOS == "darwin" {
			b.Log.
				WithField("collector", "vulnerability report").
				Warn("Skipping vulnerability scan: unsupported on macOS (darwin)")
		} else {
			vulnerabilityReport, vReportErr = b.VulnerabilityReport.CollectVulnerability(ctx, pkgs)
			if vReportErr != nil {
				b.Log.WithError(vReportErr).Error("vulnerability report")
			}
		}

		defer func(n string) { b.Log.WithField("collector", n).WithField("duration", time.Since(start)).Info("done") }("vulnerability report")
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
	inv.ComplianceReport = complianceReport
	inv.VulnerabilityReport = vulnerabilityReport

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
