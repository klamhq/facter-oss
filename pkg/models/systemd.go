package models

type SystemdService struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Loaded       string `json:"loaded"`
	Active       string `json:"active"`
	SubState     string `json:"subState"`
	Enabled      bool   `json:"enabled"`
	PID          int64  `json:"pid"`
	Tasks        int64  `json:"tasks"`
	MemoryBytes  int64  `json:"memoryBytes"`
	CPUUsageNsec int64  `json:"cpuUsageNsec"`
	CGroup       string `json:"cgroup"`

	Requires []string `json:"requires"`
	Wants    []string `json:"wants"`
	After    []string `json:"after"`
	Before   []string `json:"before"`
}
