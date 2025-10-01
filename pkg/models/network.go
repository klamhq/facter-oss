package models

type NetworkInterface struct {
	Name            string   `json:"name,omitempty"`
	IP              []IP     `json:"ip,omitempty"`
	HardwareAddress string   `json:"hardware_address,omitempty"`
	Flags           []string `json:"flags,omitempty"`
}

type DnsInfo struct {
	NameServers   string `json:"nameservers,omitempty"`
	SearchDomains string `json:"searchdomains,omitempty"`
	Port          string `json:"port,omitempty"`
}

type ConnectionState struct {
	Protocol string      `json:"protocol"`
	State    string      `json:"state"`
	Local    PortState   `json:"local"`
	Remote   PortState   `json:"remote"`
	Process  PortProcess `json:"process"`
}

type PortProcess struct {
	Name    string  `json:"name"`
	Pid     uint32  `json:"pid"`
	Package Package `json:"package"`
}

type IP struct {
	Addr         string `json:"addr"`
	Version      string `json:"version"`
	CIDR         string `json:"cidr"`
	ExternalIp   string `json:"externalIp"`
	ForwardedFor string `json:"forwardedFor"`
}

type PortState struct {
	IP IP
	//port uint16
}

// Port is used to store relevant informations about tcp and upd open ports.
// This structure is returned by method `Grab`
type Port struct {
	Protocol      string `json:"protocol,omitempty"`
	State         string `json:"state,omitempty"`
	LocalAddress  string `json:"local_address,omitempty"`
	LocalPort     string `json:"local_port,omitempty"`
	RemoteAddress string `json:"remote_address,omitempty"`
	RemotePort    string `json:"remote_port,omitempty"`
}

type ConnectionCount struct {
	Source          string `json:"source,omitempty"`
	Destination     string `json:"destination,omitempty"`
	ConnectionCount int    `json:"count,omitempty"`
}

type ProcessConnectionCount struct {
	Process         string `json:"process,omitempty"`
	ConnectionCount int    `json:"count,omitempty"`
}

type NetworkGroup struct {
	Network     string              `json:"network"`
	Connections []NetworkConnection `json:"connections"`
}

type NetworkConnection struct {
	Host      string `json:"host"`
	Interface string `json:"interface"`
	MAC       string `json:"mac"`
	IP        string `json:"ip"`
	Version   string `json:"version"`
}
