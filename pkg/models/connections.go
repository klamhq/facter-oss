package models

// Connections represents a network connection with its details.
// It includes local and remote ports, process information, protocol type,
// remote and local IP addresses, process ID, package name, state, process name,
// and process path.
type Connections struct {
	LocalPort   uint32
	RemotePort  uint32
	Process     string
	Protocol    uint32
	RemoteIp    string
	LocalIp     string
	Pid         int32
	Package     string
	State       string
	ProcessName string
	ProcessPath string
}
