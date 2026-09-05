package models

// Connections represents a network connection with its details.
// It includes local and remote ports, process information, protocol type,
// remote and local IP addresses, process ID, package name, state, process name,
// and process path.
type Connections struct {
	LocalPort   uint32 `json:"localPort,omitempty"`
	RemotePort  uint32 `json:"remotePort,omitempty"`
	Process     string `json:"process,omitempty"`
	Protocol    uint32 `json:"protocol,omitempty"`
	RemoteIp    string `json:"remoteIp,omitempty"`
	LocalIp     string `json:"localIp,omitempty"`
	Pid         int32  `json:"pid,omitempty"`
	Package     string `json:"package,omitempty"`
	State       string `json:"state,omitempty"`
	ProcessName string `json:"processName,omitempty"`
	ProcessPath string `json:"processPath,omitempty"`
}
