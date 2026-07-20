package options

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunOptions(t *testing.T) {
	opts := RunOptions{}

	// Test default values
	assert.False(t, opts.Facter.Enabled)

	// Test setting values
	opts.Facter.Enabled = true
	opts.Facter.Store.Path = "/tmp/test.db"
	opts.Facter.Logs.DebugMode = true

	assert.True(t, opts.Facter.Enabled)
	assert.Equal(t, "/tmp/test.db", opts.Facter.Store.Path)
	assert.True(t, opts.Facter.Logs.DebugMode)
}

func TestStoreOptions(t *testing.T) {
	store := StoreOptions{
		Path: "/var/lib/facter/store.db",
	}

	assert.Equal(t, "/var/lib/facter/store.db", store.Path)
}

func TestInventoryOptions(t *testing.T) {
	inv := Inventory{
		Packages: PackagesOptions{
			Enabled: true,
		},
		SSH: SSHOptions{
			Enabled: true,
		},
		Networks: NetworksOptions{
			Enabled: true,
			Ports: PortsOptions{
				Enabled: true,
			},
		},
	}

	assert.True(t, inv.Packages.Enabled)
	assert.True(t, inv.SSH.Enabled)
	assert.True(t, inv.Networks.Enabled)
	assert.True(t, inv.Networks.Ports.Enabled)
}

func TestNetworksOptions(t *testing.T) {
	networks := NetworksOptions{
		Enabled: true,
		Ports: PortsOptions{
			Enabled: true,
		},
		Connections: ConnectionsOptions{
			Enabled: true,
		},
		Firewall: FirewallOptions{
			Enabled: true,
		},
		PublicIp: PublicIpOptions{
			Enabled:        true,
			Timeout:        10,
			PublicIpApiUrl: "https://ifconfig.me/",
		},
		GeoIp: GeoIpOptions{
			Enabled:         true,
			Timeout:         10,
			GoogleGeoApikey: "test-key",
			GoogleGeoUrl:    "https://geo.api.com",
		},
	}

	assert.True(t, networks.Enabled)
	assert.True(t, networks.Ports.Enabled)
	assert.True(t, networks.Connections.Enabled)
	assert.True(t, networks.Firewall.Enabled)
	assert.Equal(t, 10, networks.PublicIp.Timeout)
	assert.Equal(t, "https://ifconfig.me/", networks.PublicIp.PublicIpApiUrl)
	assert.Equal(t, "test-key", networks.GeoIp.GoogleGeoApikey)
}

func TestPlatformOptions(t *testing.T) {
	platform := PlatformOptions{
		Enabled: true,
		Hardware: HardwareOptions{
			Enabled: true,
		},
		Kernel: KernelOptions{
			Enabled: true,
		},
		Os: OsOptions{
			Enabled: true,
		},
		Virtualization: VirtualizationOptions{
			Enabled: true,
		},
	}

	assert.True(t, platform.Enabled)
	assert.True(t, platform.Hardware.Enabled)
	assert.True(t, platform.Kernel.Enabled)
	assert.True(t, platform.Os.Enabled)
	assert.True(t, platform.Virtualization.Enabled)
}

func TestSystemdServiceOptions(t *testing.T) {
	systemd := SystemdServiceOptions{
		Enabled: true,
	}

	assert.True(t, systemd.Enabled)
}

func TestProcessOptions(t *testing.T) {
	process := ProcessOptions{
		Enabled: true,
	}

	assert.True(t, process.Enabled)
}

func TestApplicationsOptions(t *testing.T) {
	apps := ApplicationsOptions{
		Enabled: true,
		Docker: docker{
			Enabled: true,
		},
	}

	assert.True(t, apps.Enabled)
	assert.True(t, apps.Docker.Enabled)
}

func TestLogsOptions(t *testing.T) {
	logs := LogsOptions{
		DebugMode: true,
	}

	assert.True(t, logs.DebugMode)
}

func TestPerformanceOptions(t *testing.T) {
	perf := PerformanceOptions{
		Enabled: true,
	}

	assert.True(t, perf.Enabled)
}

func TestSinkOptions(t *testing.T) {
	sink := SinkOptions{
		Output: OutputOptions{
			Format:          "json",
			Type:            "file",
			OutputFilename:  "output.json",
			OutputDirectory: "/tmp",
		},
	}

	assert.Equal(t, "json", sink.Output.Format)
	assert.Equal(t, "file", sink.Output.Type)
	assert.Equal(t, "output.json", sink.Output.OutputFilename)
	assert.Equal(t, "/tmp", sink.Output.OutputDirectory)
}

func TestOutputOptions(t *testing.T) {
	output := OutputOptions{
		FacterServer: FacterServerOptions{
			ServerHost:         "localhost",
			ServerPort:         "56230",
			CertificatePath:    "/path/to/cert",
			CertificateKeyPath: "/path/to/key",
			CaPath:             "/path/to/ca",
			SSLHostname:        "test.facter.fr",
		},
		Format:          "proto",
		Type:            "remote",
		OutputFilename:  "export.iya",
		OutputDirectory: "/var/lib/facter",
	}

	assert.Equal(t, "localhost", output.FacterServer.ServerHost)
	assert.Equal(t, "56230", output.FacterServer.ServerPort)
	assert.Equal(t, "proto", output.Format)
	assert.Equal(t, "remote", output.Type)
}

func TestFacterServerOptions(t *testing.T) {
	server := FacterServerOptions{
		ServerHost:         "facter.example.com",
		ServerPort:         "443",
		CertificatePath:    "/etc/ssl/cert.pem",
		CertificateKeyPath: "/etc/ssl/key.pem",
		CaPath:             "/etc/ssl/ca.pem",
		SSLHostname:        "facter.example.com",
	}

	assert.Equal(t, "facter.example.com", server.ServerHost)
	assert.Equal(t, "443", server.ServerPort)
	assert.Equal(t, "/etc/ssl/cert.pem", server.CertificatePath)
	assert.Equal(t, "/etc/ssl/key.pem", server.CertificateKeyPath)
	assert.Equal(t, "/etc/ssl/ca.pem", server.CaPath)
	assert.Equal(t, "facter.example.com", server.SSLHostname)
}

func TestPackagesOptions(t *testing.T) {
	packages := PackagesOptions{
		Enabled: true,
	}

	assert.True(t, packages.Enabled)
}

func TestGeoIpOptions(t *testing.T) {
	geoip := GeoIpOptions{
		Enabled:         true,
		Timeout:         15,
		GoogleGeoApikey: "AIzaSyB9i6aHGUa5Iv0tfDb9sy_HqXUebFmk8wI",
		GoogleGeoUrl:    "https://www.googleapis.com/geolocation/v1/geolocate",
	}

	assert.True(t, geoip.Enabled)
	assert.Equal(t, 15, geoip.Timeout)
	assert.NotEmpty(t, geoip.GoogleGeoApikey)
	assert.Equal(t, "https://www.googleapis.com/geolocation/v1/geolocate", geoip.GoogleGeoUrl)
}

func TestSSHOptions(t *testing.T) {
	ssh := SSHOptions{
		Enabled: true,
	}

	assert.True(t, ssh.Enabled)
}

func TestUserOptions(t *testing.T) {
	user := UserOptions{
		Enabled:    true,
		PasswdFile: "/etc/passwd",
	}

	assert.True(t, user.Enabled)
	assert.Equal(t, "/etc/passwd", user.PasswdFile)
}

func TestFirewallOptions(t *testing.T) {
	firewall := FirewallOptions{
		Enabled: true,
	}

	assert.True(t, firewall.Enabled)
}

func TestPublicIpOptions(t *testing.T) {
	publicip := PublicIpOptions{
		Enabled:        true,
		Timeout:        20,
		PublicIpApiUrl: "https://api.ipify.org",
	}

	assert.True(t, publicip.Enabled)
	assert.Equal(t, 20, publicip.Timeout)
	assert.Equal(t, "https://api.ipify.org", publicip.PublicIpApiUrl)
}

func TestPortsOptions(t *testing.T) {
	ports := PortsOptions{
		Enabled: true,
	}

	assert.True(t, ports.Enabled)
}

func TestComplianceOptions(t *testing.T) {
	compliance := ComplianceOptions{
		Enabled:    true,
		Profile:    "xccdf_org.ssgproject.content_profile_cis_level1_server",
		ResultFile: "/tmp/openscap-results.xml",
	}

	assert.True(t, compliance.Enabled)
	assert.Equal(t, "xccdf_org.ssgproject.content_profile_cis_level1_server", compliance.Profile)
	assert.Equal(t, "/tmp/openscap-results.xml", compliance.ResultFile)
}

func TestVulnerabilitiesOptions(t *testing.T) {
	vuln := VulnerabilitiesOptions{
		Enabled: true,
	}

	assert.True(t, vuln.Enabled)
}

func TestFullRunOptions(t *testing.T) {
	opts := RunOptions{}
	opts.Facter.Enabled = true
	opts.Facter.Store.Path = "/tmp/store.db"
	opts.Facter.Logs.DebugMode = true
	opts.Facter.PerformanceProfiling.Enabled = true

	// Inventory
	opts.Facter.Inventory.Packages.Enabled = true
	opts.Facter.Inventory.SSH.Enabled = true
	opts.Facter.Inventory.User.Enabled = true
	opts.Facter.Inventory.User.PasswdFile = "/etc/passwd"
	opts.Facter.Inventory.Networks.Enabled = true
	opts.Facter.Inventory.Networks.Firewall.Enabled = true
	opts.Facter.Inventory.Platform.Enabled = true
	opts.Facter.Inventory.SystemdService.Enabled = true
	opts.Facter.Inventory.Process.Enabled = true
	opts.Facter.Inventory.Applications.Enabled = true
	opts.Facter.Inventory.Applications.Docker.Enabled = true

	// Sink
	opts.Facter.Sink.Output.Format = "json"
	opts.Facter.Sink.Output.Type = "file"
	opts.Facter.Sink.Output.OutputFilename = "output.json"
	opts.Facter.Sink.Output.OutputDirectory = "/tmp"

	// Compliance
	opts.Facter.Compliance.Enabled = true
	opts.Facter.Compliance.Profile = "test-profile"
	opts.Facter.Compliance.ResultFile = "/tmp/results.xml"

	// Vulnerabilities
	opts.Facter.Vulnerabilities.Enabled = true

	// Verify all settings
	assert.True(t, opts.Facter.Enabled)
	assert.True(t, opts.Facter.Logs.DebugMode)
	assert.True(t, opts.Facter.PerformanceProfiling.Enabled)
	assert.True(t, opts.Facter.Inventory.Packages.Enabled)
	assert.True(t, opts.Facter.Inventory.SSH.Enabled)
	assert.True(t, opts.Facter.Inventory.User.Enabled)
	assert.True(t, opts.Facter.Inventory.Networks.Enabled)
	assert.True(t, opts.Facter.Inventory.Platform.Enabled)
	assert.True(t, opts.Facter.Inventory.SystemdService.Enabled)
	assert.True(t, opts.Facter.Inventory.Process.Enabled)
	assert.True(t, opts.Facter.Inventory.Applications.Enabled)
	assert.True(t, opts.Facter.Inventory.Applications.Docker.Enabled)
	assert.Equal(t, "json", opts.Facter.Sink.Output.Format)
	assert.True(t, opts.Facter.Compliance.Enabled)
	assert.True(t, opts.Facter.Vulnerabilities.Enabled)
}
