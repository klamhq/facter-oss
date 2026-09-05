package inventory

import (
	"testing"
	"time"

	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestComputeDelta(t *testing.T) {
	oldInv := &schema.HostInventory{
		Hostname:  "test-host",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		Packages: []*schema.Package{
			{Name: "pkg1", Version: "1.0.0"},
			{Name: "pkg2", Version: "1.0.0"},
		},
		Users: []*schema.User{
			{Username: "user1"},
		},
	}

	newInv := &schema.HostInventory{
		Hostname:  "test-host",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		Packages: []*schema.Package{
			{Name: "pkg1", Version: "1.0.0"},
			{Name: "pkg3", Version: "1.0.0"},
		},
		Users: []*schema.User{
			{Username: "user1"},
			{Username: "user2"},
		},
	}

	logger := logrus.New()
	delta := ComputeDelta(oldInv, newInv, logger)

	assert.NotNil(t, delta)
	assert.Equal(t, "test-host", delta.Hostname)
	assert.NotEmpty(t, delta.UpdatedAt)
	assert.Len(t, delta.PackagesAdded, 1)
	assert.Len(t, delta.PackagesRemoved, 1)
	assert.Len(t, delta.UsersAdded, 1)
	assert.Len(t, delta.UsersRemoved, 0)
}

func TestIsDeltaEmpty(t *testing.T) {
	delta := &schema.HostDeltaInventory{
		PackagesAdded:   []*schema.Package{},
		PackagesRemoved: []*schema.Package{},
		UsersAdded:      []*schema.User{},
		UsersRemoved:    []*schema.User{},
	}

	assert.True(t, IsDeltaEmpty(delta))

	delta.PackagesAdded = append(delta.PackagesAdded, &schema.Package{Name: "pkg1"})
	assert.False(t, IsDeltaEmpty(delta))
}

func TestStableHash(t *testing.T) {
	p1 := &schema.Package{Name: "pkg1", Version: "1.0.0"}
	p2 := &schema.Package{Name: "pkg1", Version: "1.0.0"}
	p3 := &schema.Package{Name: "pkg2", Version: "1.0.0"}

	hash1 := StableHash(p1)
	hash2 := StableHash(p2)
	hash3 := StableHash(p3)

	assert.Equal(t, hash1, hash2, "Hashes for identical packages should match")
	assert.NotEqual(t, hash1, hash3, "Hashes for different packages should not match")
}
func TestDiffByHash(t *testing.T) {
	oldList := []*schema.Package{
		{Name: "pkg1", Version: "1.0.0"},
		{Name: "pkg2", Version: "1.0.0"},
		{Name: "pkg3", Version: "1.0.0"},
	}

	newList := []*schema.Package{
		{Name: "pkg1", Version: "1.0.0"},
		{Name: "pkg3", Version: "2.0.0"},
		{Name: "pkg4", Version: "1.0.0"},
	}

	added, removed, changed := DiffGenericByHash(oldList, newList, func(p *schema.Package) string { return p.Name })

	assert.Len(t, changed, 1)
	assert.Equal(t, "pkg3", changed[0].Name)
	assert.Len(t, added, 1)
	assert.Equal(t, "pkg4", added[0].Name)
	assert.Len(t, removed, 1)
	assert.Equal(t, "pkg2", removed[0].Name)
}

func TestDiffByHashIgnoresUpdatedAt(t *testing.T) {
	oldList := []*schema.User{{Username: "user1", Shell: "/bin/sh", UpdatedAt: "2024-01-01T00:00:00Z"}}
	newList := []*schema.User{{Username: "user1", Shell: "/bin/sh", UpdatedAt: "2024-01-02T00:00:00Z"}}

	added, removed, changed := DiffGenericByHash(oldList, newList, func(user *schema.User) string { return user.Username })

	assert.Empty(t, added)
	assert.Empty(t, removed)
	assert.Empty(t, changed)
}

func TestComputeDelta_PackagesChanged(t *testing.T) {
	old := &schema.HostInventory{
		Packages: []*schema.Package{{Name: "pkg1", Version: "1.0"}},
	}
	new := &schema.HostInventory{
		Packages: []*schema.Package{{Name: "pkg1", Version: "2.0"}},
	}
	logger := logrus.New()
	delta := ComputeDelta(old, new, logger)
	assert.Len(t, delta.PackagesChanged, 1)
}

func TestComputeDelta_UsersChanged(t *testing.T) {
	old := &schema.HostInventory{
		Users: []*schema.User{{Username: "user1", Shell: "/bin/sh"}},
	}
	new := &schema.HostInventory{
		Users: []*schema.User{{Username: "user1", Shell: "/bin/bash"}},
	}
	logger := logrus.New()
	delta := ComputeDelta(old, new, logger)
	assert.Len(t, delta.UsersChanged, 1)
}

func TestComputeDelta_DockerChanged(t *testing.T) {

	old := &schema.HostInventory{
		Application: []*schema.Application{
			{Docker: &schema.Docker{
				Containers: []*schema.Containers{{Id: "container1", Name: "Container One"}, {Id: "container2", Name: "Container Two-a"}},
				Images:     []*schema.ContainersImages{{Id: "image1", RepoTags: []string{"repo1:tag1"}}, {Id: "image2", RepoTags: []string{"repo2:tag2"}}},
				Networks:   []*schema.DockerNetworks{{Id: "network1", Name: "network1"}, {Id: "network2", Name: "network2"}},
			}},
		},
	}
	new := &schema.HostInventory{
		Application: []*schema.Application{
			{Docker: &schema.Docker{
				Containers: []*schema.Containers{{Id: "container2", Name: "Container Two-b"}, {Id: "container3", Name: "Container Three"}},
				Images:     []*schema.ContainersImages{{Id: "image1", RepoTags: []string{"repo1:tag2"}}, {Id: "image3", RepoTags: []string{"repo3:tag3"}}},
				Networks:   []*schema.DockerNetworks{{Id: "network1", Name: "network1-b"}, {Id: "network3", Name: "network3"}},
			}},
		},
	}
	logger := logrus.New()
	delta := ComputeDelta(old, new, logger)
	assert.Len(t, delta.DockerDelta.ContainersAdded, 1)
	assert.Len(t, delta.DockerDelta.ContainersChanged, 1)
	assert.Len(t, delta.DockerDelta.ContainersRemoved, 1)
	assert.Len(t, delta.DockerDelta.ImagesAdded, 1)
	assert.Len(t, delta.DockerDelta.ImagesRemoved, 1)
	assert.Len(t, delta.DockerDelta.NetworksAdded, 1)
	assert.Len(t, delta.DockerDelta.NetworksRemoved, 1)
	assert.Len(t, delta.DockerDelta.NetworksChanged, 1)
}

func TestComputeDelta_ProcessesChanged(t *testing.T) {
	old := &schema.HostInventory{
		Processes: []*schema.Process{{Pid: 1234, Name: "process1", Cmdline: "old command"}, {Pid: 5678, Name: "process2", Cmdline: "command2"}},
	}
	new := &schema.HostInventory{
		Processes: []*schema.Process{{Pid: 1234, Name: "process1", Cmdline: "new command"}, {Pid: 1234, Name: "process3", Cmdline: "command3"}},
	}
	logger := logrus.New()
	delta := ComputeDelta(old, new, logger)
	assert.Len(t, delta.ProcessesChanged, 1)
	assert.Len(t, delta.ProcessesAdded, 1)
	assert.Len(t, delta.ProcessesRemoved, 1)
}

func TestComputeDelta_SystemdServicesChanged(t *testing.T) {
	old := &schema.HostInventory{
		SystemdService: []*schema.SystemdService{{Name: "service1", SubState: "active"}, {Name: "service2", SubState: "inactive"}},
	}
	new := &schema.HostInventory{
		SystemdService: []*schema.SystemdService{{Name: "service1", SubState: "inactive"}, {Name: "service3", SubState: "inactive"}},
	}
	logger := logrus.New()
	delta := ComputeDelta(old, new, logger)
	assert.Len(t, delta.SystemdservicesChanged, 1)
	assert.Len(t, delta.SystemdservicesAdded, 1)
	assert.Len(t, delta.SystemdservicesRemoved, 1)
}

func TestComputeDelta_KnownHostsChanged(t *testing.T) {
	old := &schema.HostInventory{
		KnownHost: []*schema.KnownHost{{Hostname: "host1", Fingerprint: "fp1", Type: "rsa"}, {Hostname: "host2", Fingerprint: "fp2", Type: "ed25519"}},
	}
	new := &schema.HostInventory{
		KnownHost: []*schema.KnownHost{{Hostname: "host1", Fingerprint: "fp1", Type: "ed25519"}, {Hostname: "host3", Fingerprint: "fp3", Type: "rsa"}},
	}
	logger := logrus.New()
	delta := ComputeDelta(old, new, logger)
	assert.Len(t, delta.KnownhostsChanged, 1)
	assert.Len(t, delta.KnownhostsAdded, 1)
	assert.Len(t, delta.KnownhostsRemoved, 1)
}

func TestComputeDelta_SshKeyInfoChanged(t *testing.T) {
	old := &schema.HostInventory{
		SshKeyInfo: []*schema.SshKeyInfo{{Fingerprint: "fp1", Type: "rsa", Comment: "old"}, {Fingerprint: "fp2", Type: "ed25519", Comment: "abcd"}},
	}
	new := &schema.HostInventory{
		SshKeyInfo: []*schema.SshKeyInfo{{Fingerprint: "fp1", Type: "rsa", Comment: "new"}, {Fingerprint: "fp3", Type: "ed25519", Comment: "new"}},
	}
	logger := logrus.New()
	delta := ComputeDelta(old, new, logger)
	assert.Len(t, delta.SshkeyinfoChanged, 1)
	assert.Len(t, delta.SshkeyinfoAdded, 1)
	assert.Len(t, delta.SshkeyinfoRemoved, 1)
}

func TestComputeDelta_SshKeysIgnoreUpdatedAt(t *testing.T) {
	old := &schema.HostInventory{
		SshKeyInfo: []*schema.SshKeyInfo{{
			Fingerprint: "fp1",
			Type:        "rsa",
			UpdatedAt:   "2024-01-01T00:00:00Z",
		}},
		SshKeyAccess: []*schema.SshKeyAccess{{
			Fingerprint: "fp1",
			AsUser:      "alice",
			UpdatedAt:   "2024-01-01T00:00:00Z",
		}},
	}
	newInventory := &schema.HostInventory{
		SshKeyInfo: []*schema.SshKeyInfo{{
			Fingerprint: "fp1",
			Type:        "rsa",
			UpdatedAt:   "2024-01-02T00:00:00Z",
		}},
		SshKeyAccess: []*schema.SshKeyAccess{{
			Fingerprint: "fp1",
			AsUser:      "alice",
			UpdatedAt:   "2024-01-02T00:00:00Z",
		}},
	}

	delta := ComputeDelta(old, newInventory, logrus.New())

	assert.Empty(t, delta.SshkeyinfoAdded)
	assert.Empty(t, delta.SshkeyinfoChanged)
	assert.Empty(t, delta.SshkeyinfoRemoved)
	assert.Empty(t, delta.SshkeyaccessAdded)
	assert.Empty(t, delta.SshkeyaccessRemoved)
}

func TestComputeDelta_SshKeysIgnoreDuplicateEntries(t *testing.T) {
	old := &schema.HostInventory{
		SshKeyInfo: []*schema.SshKeyInfo{{
			Fingerprint: "fp1",
			Type:        "ssh-rsa",
			Comment:     "istia@comment",
		}},
		SshKeyAccess: []*schema.SshKeyAccess{{
			Fingerprint: "fp1",
			AsUser:      "root",
		}},
	}
	newInventory := &schema.HostInventory{
		SshKeyInfo: []*schema.SshKeyInfo{
			{Fingerprint: "fp1", Type: "ssh-rsa", Comment: "istia@comment"},
			{Fingerprint: "fp1", Type: "ssh-rsa", Comment: "istia@comment", Name: "authorized_keys"},
		},
		SshKeyAccess: []*schema.SshKeyAccess{
			{Fingerprint: "fp1", AsUser: "root"},
			{Fingerprint: "fp1", AsUser: "root"},
		},
	}

	delta := ComputeDelta(old, newInventory, logrus.New())

	assert.Empty(t, delta.SshkeyinfoAdded)
	assert.Empty(t, delta.SshkeyinfoChanged)
	assert.Empty(t, delta.SshkeyaccessAdded)
}

func TestComputeDelta_SshKeyAccessChanged(t *testing.T) {
	old := &schema.HostInventory{
		SshKeyAccess: []*schema.SshKeyAccess{{Fingerprint: "fp1", AsUser: "user1"}, {Fingerprint: "fp2", AsUser: "user2"}},
	}
	new := &schema.HostInventory{
		SshKeyAccess: []*schema.SshKeyAccess{{Fingerprint: "fp3", AsUser: "user3"}},
	}
	logger := logrus.New()
	delta := ComputeDelta(old, new, logger)
	assert.Len(t, delta.SshkeyaccessAdded, 1)
	assert.Len(t, delta.SshkeyaccessRemoved, 2)
}

func TestComputeDelta_VulnerabilityChanged(t *testing.T) {
	old := &schema.HostInventory{
		VulnerabilityReport: &schema.VulnerabilityReport{
			Matches: []*schema.PackageVulnMatch{
				{PackageName: "pkg1", InstalledVersion: "1.0.0"}, {PackageName: "pkg2", InstalledVersion: "1.0.0"},
			},
		},
	}
	new := &schema.HostInventory{
		VulnerabilityReport: &schema.VulnerabilityReport{
			Matches: []*schema.PackageVulnMatch{
				{PackageName: "pkg1", InstalledVersion: "2.0.0"}, {PackageName: "pkg3", InstalledVersion: "1.0.0"},
			},
		},
	}
	logger := logrus.New()
	delta := ComputeDelta(old, new, logger)
	assert.NotNil(t, delta.VulnerabilityDelta)
	assert.Len(t, delta.VulnerabilityDelta.MatchesChanged, 1)
	assert.Len(t, delta.VulnerabilityDelta.MatchesAdded, 1)
	assert.Len(t, delta.VulnerabilityDelta.MatchesRemoved, 1)
}

func TestComputeDelta_ComplianceChanged(t *testing.T) {
	old := &schema.HostInventory{
		ComplianceReport: &schema.ComplianceReport{
			Score:   &schema.Score{Maximum: "100", Value: "80"},
			Profile: "default",
			RuleResults: []*schema.RuleCheckResult{
				{Id: "rule1", Result: "PASS"},
			},
		},
	}
	new := &schema.HostInventory{
		ComplianceReport: &schema.ComplianceReport{
			Score:   &schema.Score{Maximum: "100", Value: "90"},
			Profile: "strict",
			RuleResults: []*schema.RuleCheckResult{
				{Id: "rule1", Result: "FAIL"},
			},
		},
	}
	logger := logrus.New()
	delta := ComputeDelta(old, new, logger)
	assert.NotNil(t, delta.ComplianceReportDelta)
	assert.Len(t, delta.ComplianceReportDelta.RulesChanged, 1)
	assert.Len(t, delta.ComplianceReportDelta.RulesAdded, 0)
	assert.Len(t, delta.ComplianceReportDelta.RulesRemoved, 0)
	assert.Equal(t, "100", delta.ComplianceReportDelta.GetScoreDelta().GetMaximum())
	assert.Equal(t, "90", delta.ComplianceReportDelta.GetScoreDelta().GetValue())
	assert.Equal(t, "strict", delta.ComplianceReportDelta.GetProfile())
}

func TestComputeDelta_ComplianceUnchanged(t *testing.T) {
	report := &schema.ComplianceReport{
		Score:   &schema.Score{Maximum: "100", Value: "90"},
		Profile: "strict",
		RuleResults: []*schema.RuleCheckResult{
			{Id: "rule1", Result: "PASS"},
		},
	}

	delta := ComputeDelta(
		&schema.HostInventory{ComplianceReport: report},
		&schema.HostInventory{ComplianceReport: report},
		logrus.New(),
	)

	assert.Nil(t, delta.ComplianceReportDelta)
}

func TestComputeDelta_PlatformChanged(t *testing.T) {
	old := &schema.HostInventory{
		Platform: &schema.Platform{Os: &schema.Os{Name: "Linux", Version: "1.0"}},
	}
	new := &schema.HostInventory{
		Platform: &schema.Platform{Os: &schema.Os{Name: "Linux", Version: "2.0"}},
	}
	logger := logrus.New()
	delta := ComputeDelta(old, new, logger)
	assert.NotNil(t, delta.Platform)
	assert.Equal(t, "2.0", delta.Platform.Os.Version)
}

// func TestComputeDelta_NetworkChanged(t *testing.T) {
// 	old := &schema.HostInventory{
// 		Network: &schema.Network{Connections: []*schema.ConnectionState{{Process: &schema.Process{Name: "proc1"}}}},
// 	}
// 	new := &schema.HostInventory{
// 		Network: &schema.Network{Connections: []*schema.ConnectionState{{Process: &schema.Process{Name: "proc1"}},
// 	}
// 	logger := logrus.New()
// 	delta := ComputeDelta(old, new, logger)
// 	assert.NotNil(t, delta.Network)
// 	assert.Equal(t, "new-host", delta.Network.Hostname)
// }
