package models

type Hosts struct {
	Hostname        string `json:"hostname"`
	SwapTotal       int64
	MemoryTotal     int64
	MemoryUsed      int64
	OperatingSystem string
	Kernel          string
	Disk            []Disk
	Vulnerability   []MatchedVuln `json:"vulnerabilities"`
	Uptime          uint64        `json:"uptime"`
}
