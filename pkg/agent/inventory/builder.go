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

func (b *Builder) buildRevisionEnvelope(fullInventory *schema.HostInventory, delta *schema.HostDeltaInventory, previousRevision *schema.InventoryRevisionEnvelope) *schema.InventoryRevisionEnvelope {
	now := b.Now().UTC()
	sequence := uint64(1)
	previousRevisionID := ""
	previousSequence := uint64(0)
	if previousRevision != nil {
		previousRevisionID = previousRevision.RevisionId
		previousSequence = previousRevision.Sequence
		sequence = previousRevision.Sequence + 1
	}

	revisionID := fmt.Sprintf("%s-%d-%d", fullInventory.Hostname, now.UnixNano(), sequence)
	stateHash := fmt.Sprintf("%x", StableHash(fullInventory))
	sourceType := schema.SourceType_SOURCE_TYPE_FULL
	if delta != nil {
		sourceType = schema.SourceType_SOURCE_TYPE_DELTA
	}

	revision := &schema.InventoryRevisionEnvelope{
		RevisionId:         revisionID,
		Sequence:           sequence,
		PreviousRevisionId: previousRevisionID,
		PreviousSequence:   previousSequence,
		CreatedAt:          now.Format(time.RFC3339Nano),
		AgentId:            fullInventory.Hostname,
		Hostname:           fullInventory.Hostname,
		SourceType:         sourceType,
		StateHash:          stateHash,
	}
	if fullInventory.Identifier != nil {
		revision.MachineId = fullInventory.Identifier.MachineId
	}
	if fullInventory.Metadata != nil && fullInventory.Metadata.RunningDate != "" {
		revision.CreatedAt = fullInventory.Metadata.RunningDate
	}

	switch sourceType {
	case schema.SourceType_SOURCE_TYPE_FULL:
		revision.Payload = &schema.InventoryRevisionEnvelope_Full{Full: fullInventory}
	case schema.SourceType_SOURCE_TYPE_DELTA:
		revision.Payload = &schema.InventoryRevisionEnvelope_Delta{Delta: delta}
	}
	return revision
}

func (b *Builder) ManageDelta(fullInventory *schema.HostInventory) (*schema.InventoryRequest, *schema.HostInventory) {
	previous, err := b.Store.Get(fullInventory.Hostname)
	previousRevision, prevRevErr := b.Store.GetRevision(fullInventory.Hostname)
	if prevRevErr != nil {
		b.Log.WithError(prevRevErr).Warn("Unable to load previous revision metadata")
	}

	if err != nil || previous == nil {
		b.Log.Info("No previous inventory, creating full revision")
		revision := b.buildRevisionEnvelope(fullInventory, nil, previousRevision)
		if saveErr := b.Store.SaveRevision(fullInventory.Hostname, revision); saveErr != nil {
			b.Log.WithError(saveErr).Error("Failed to persist revision metadata")
		}
		return &schema.InventoryRequest{Content: &schema.InventoryRequest_Revision{Revision: revision}}, fullInventory
	}

	b.Log.Info("Previous inventory found, computing delta")
	delta := ComputeDelta(previous, fullInventory, b.Log)
	if IsDeltaEmpty(delta) {
		b.Log.Info("No changes detected, nothing to send")
		return nil, nil
	}
	delta.UpdatedAt = time.Now().Format(time.RFC3339)
	revision := b.buildRevisionEnvelope(fullInventory, delta, previousRevision)
	if saveErr := b.Store.SaveRevision(fullInventory.Hostname, revision); saveErr != nil {
		b.Log.WithError(saveErr).Error("Failed to persist revision metadata")
	}
	return &schema.InventoryRequest{Content: &schema.InventoryRequest_Revision{Revision: revision}}, fullInventory
}
