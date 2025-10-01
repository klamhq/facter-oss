package options

// RunOptions run a facter client
type RunOptions struct {
	Facter struct {
		Enabled              bool               `yaml:"enabled"`
		Inventory            Inventory          `yaml:"inventory"`
		Store                StoreOptions       `yaml:"store"`
		Logs                 LogsOptions        `yaml:"logs"`
		Sink                 SinkOptions        `yaml:"sink"`
		PerformanceProfiling PerformanceOptions `yaml:"performanceProfiling"`
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
	Enabled bool          `yaml:"enabled"`
	Output  OutputOptions `yaml:"output"`
}

// OutputOptions contains the options for output export to file
type OutputOptions struct {
	Format          string `yaml:"format"`
	Type            string `yaml:"type"`
	OutputFilename  string `yaml:"outputFilename"`
	OutputDirectory string `yaml:"outputDirectory"`
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

// DefaultNewRunOptions creates a new RunOptions with default parameters
func DefaultNewRunOptions() {
	d := &RunOptions{}
	d.Facter.Enabled = true
	d.Facter.Logs.DebugMode = false
	d.Facter.Inventory.Packages.Enabled = true
	d.Facter.Inventory.Networks.GeoIp.Enabled = true
	d.Facter.Inventory.Networks.GeoIp.Timeout = 10
	d.Facter.Inventory.Networks.GeoIp.GoogleGeoUrl = "https://www.googleapis.com/geolocation/v1/geolocate"
	d.Facter.Inventory.SSH.Enabled = true
	d.Facter.Inventory.User.Enabled = true
	d.Facter.Inventory.User.PasswdFile = "/etc/passwd"
	d.Facter.Inventory.Networks.Ports.Enabled = true
	d.Facter.Inventory.Networks.PublicIp.Enabled = true
	d.Facter.Inventory.Networks.PublicIp.PublicIpApiUrl = "https://ifconfig.me/"
	d.Facter.Inventory.Networks.Firewall.Enabled = true
	d.Facter.Inventory.Platform.Enabled = true
	d.Facter.Inventory.Platform.Virtualization.Enabled = false
	d.Facter.Inventory.Platform.Os.Enabled = true
	d.Facter.Inventory.Platform.System.InitCheckPath = "/proc/1/exe"
	d.Facter.Inventory.Platform.System.MachineID = "/etc/machine-id"
	d.Facter.Inventory.Platform.System.MachineUUID = "/sys/class/dmi/id/product_uuid"
	d.Facter.Inventory.Networks.Connections.Enabled = true
	d.Facter.Inventory.Platform.Kernel.Enabled = true
	d.Facter.Inventory.Platform.Hardware.Enabled = true
	d.Facter.PerformanceProfiling.Enabled = false
}
