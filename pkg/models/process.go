package models

type Process struct {
	Name       string  `json:"name,omitempty"`
	PID        int64   `json:"pid,omitempty"`
	Username   string  `json:"user,omitempty"`
	Status     string  `json:"status,omitempty"` // R: Running S: Sleep T: Stop I: Idle Z: Zombie W: Wait L: Lock
	CreateTime int64   `json:"createTime,omitempty"`
	Parent     int64   `json:"parent,omitempty"`
	Cmdline    string  `json:"cmdLine,omitempty"`
	Terminal   string  `json:"Terminal,omitempty"`
	Exe        string  `json:"exe,omitempty"`
	Package    Package `json:"package,omitempty"`    // Package information associated with the process
	CpuPercent float64 `json:"cpuPercent,omitempty"` // CPU usage percentage
	MemPercent float64 `json:"memPercent,omitempty"` // Memory usage percentage
}
