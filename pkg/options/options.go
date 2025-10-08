package options

// RunOptions run a facter client
type RunOptions struct {
	Facter struct {
		Enabled              bool                   `yaml:"enabled"`
		Inventory            Inventory              `yaml:"inventory"`
		Store                StoreOptions           `yaml:"store"`
		Logs                 LogsOptions            `yaml:"logs"`
		Sink                 SinkOptions            `yaml:"sink"`
		PerformanceProfiling PerformanceOptions     `yaml:"performanceProfiling"`
		Compliance           ComplianceOptions      `yaml:"compliance"`
		Vulnerabilities      VulnerabilitiesOptions `yaml:"vulnerabilities"`
	} `yaml:"facter"`
}

type StoreOptions struct {
	Path string `yaml:"path"`
}

type Inventory struct {
	Packages       PackagesOptions       `yaml:"packages"`
	SSH            SSHOptions            `yaml:"ssh"`
	User           UserOptions           `yaml:"user"`
	Networks       NetworksOptions       `yaml:"networks"`
	Platform       PlatformOptions       `yaml:"platform"`
	SystemdService SystemdServiceOptions `yaml:"systemdService"`
	Process        ProcessOptions        `yaml:"process"`
	Applications   ApplicationsOptions   `yaml:"applications"`
}

type NetworksOptions struct {
	Enabled     bool               `yaml:"enabled"`
	Ports       PortsOptions       `yaml:"ports"`
	Connections ConnectionsOptions `yaml:"connections"`
	Firewall    FirewallOptions    `yaml:"firewall"`
	PublicIp    PublicIpOptions    `yaml:"publicIp"`
	GeoIp       GeoIpOptions       `yaml:"geoIp"`
}

// PlatformOptions contains the options for fetch platform, os, configuration
type PlatformOptions struct {
	Enabled        bool                  `yaml:"enabled"`
	System         systemOptions         `yaml:"system"`
	Hardware       HardwareOptions       `yaml:"hardware"`
	Kernel         KernelOptions         `yaml:"kernel"`
	Os             OsOptions             `yaml:"os"`
	Virtualization VirtualizationOptions `yaml:"virtualization"`
}

// SystemdServiceOptions contains the options for fetch systemd service information
type SystemdServiceOptions struct {
	Enabled bool `yaml:"enabled"`
}

// ProcessOptions contains the options for fetch process information
type ProcessOptions struct {
	Enabled bool `yaml:"enabled"`
}

// systemOptions contains the options for detect init system and unique machine
type systemOptions struct {
	InitCheckPath string `yaml:"initCheckPath"`
	MachineUUID   string `yaml:"machineUUID"`
	MachineID     string `yaml:"machineID"`
}

// HardwareOptions contains the options for fetch hardware configuration
type HardwareOptions struct {
	Enabled bool `yaml:"enabled"`
}

// KernelOptions contains the options for fetch kernel configuration
type KernelOptions struct {
	Enabled bool `yaml:"enabled"`
}

// OsOptions contains the options for fetch os configuration
type OsOptions struct {
	Enabled bool `yaml:"enabled"`
}

// ConnectionsOptions contains the options for fetch logs configurations
type ConnectionsOptions struct {
	Enabled bool `yaml:"enabled"`
}

// ApplicationsOptions contains the options for fetch applications configurations
type ApplicationsOptions struct {
	Enabled bool   `yaml:"enabled"`
	Docker  docker `yaml:"docker"`
}

// docker define options for enable gather facts of docker (networks, images, containers, configuration)
// Use docker environment variable for configure docker host connection (env DOCKER_HOST, etc ...)
type docker struct {
	Enabled bool `yaml:"enabled"`
}

// VirtualizationOptions contains the options for fetch logs configurations
type VirtualizationOptions struct {
	Enabled bool `yaml:"enabled"`
}

// LogsOptions contains the options for logs configurations
type LogsOptions struct {
	DebugMode bool `yaml:"debugMode"`
}

// PerformanceOptions enable performanceProfiling mode
type PerformanceOptions struct {
	Enabled bool `yaml:"enabled"`
}

// SinkOptions contains the options for output sink
type SinkOptions struct {
	Output OutputOptions `yaml:"output"`
}

// OutputOptions contains the options for output export to file
type OutputOptions struct {
	FacterServer    FacterServerOptions `yaml:"facterServer"`
	Format          string              `yaml:"format"`
	Type            string              `yaml:"type"`
	OutputFilename  string              `yaml:"outputFilename"`
	OutputDirectory string              `yaml:"outputDirectory"`
}

// FacterServerOptions contains the options for facterServer data upload from client
type FacterServerOptions struct {
	ServerHost         string `yaml:"serverHost"`
	ServerPort         string `yaml:"serverPort"`
	CertificatePath    string `yaml:"certificatePath"`
	CertificateKeyPath string `yaml:"certificateKeyPath"`
	CaPath             string `yaml:"caPath"`
	SSLHostname        string `yaml:"sslHostname"`
}

// PackagesOptions contains the options for fetch installed package
type PackagesOptions struct {
	Enabled bool `yaml:"enabled"`
}

// GeoIpOptions contains the options for fetch  Geo IP location
type GeoIpOptions struct {
	Enabled         bool   `yaml:"enabled"`
	Timeout         int    `yaml:"timeout"`
	GoogleGeoApikey string `yaml:"googleGeoApikey"`
	GoogleGeoUrl    string `yaml:"googleGeoUrl"`
}

// SSHOptions contains the options for fetch ssh configurations
type SSHOptions struct {
	Enabled bool `yaml:"enabled"`
}

// UserOptions contains the options for fetch ssh configurations
type UserOptions struct {
	Enabled    bool   `yaml:"enabled"`
	PasswdFile string `yaml:"passwdFile"`
}

// FirewallOptions contains the options for fetch firewall configuration
type FirewallOptions struct {
	Enabled bool `yaml:"enabled"`
}

// PublicIpOptions contains the options for fetch public ip
type PublicIpOptions struct {
	Enabled        bool   `yaml:"enabled"`
	Timeout        int    `yaml:"timeout"`
	PublicIpApiUrl string `yaml:"publicIpApiUrl"`
}

// PortsOptions contains the options for fetch open ports
type PortsOptions struct {
	Enabled bool `yaml:"enabled"`
}

// ComplianceOptions contains the options for fetch compliance information
type ComplianceOptions struct {
	Enabled    bool   `yaml:"enabled"`
	Profile    string `yaml:"profile"`
	ResultFile string `yaml:"resultFile"`
}

// VulnerabilitiesOptions contains the options for fetch installed vulnerabilities
type VulnerabilitiesOptions struct {
	Enabled bool `yaml:"enabled"`
}
