package models

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

type System struct {
	Memory mem.VirtualMemoryStat `json:"memory,omitempty"`
	Host   host.InfoStat         `json:"host,omitempty"`
	CPU    []cpu.InfoStat        `json:"cpu,omitempty"`
	Load   load.AvgStat          `json:"load,omitempty"`
	Disk   []Disk                `json:"disk,omitempty"`
	Uptime uint64                `json:"uptime,omitempty"`
}

type SystemPaths struct {
	InitCheckPath string
	MachineID     string
	MachineUUID   string
}

type SystemFamily struct {
	Name string
}

type MachineID struct {
	UUID      string `json:"uuid"`
	MachineId string `json:"machine_id"`
}

type CPU struct {
	Model string
	Core  int64
	Mhz   float64
}

type OperatingSystem struct {
	Name       string `json:"name,omitempty"`
	Version    string `json:"version,omitempty"`
	InitSystem string `json:"initSystem,omitempty"`
	Family     string `json:"family,omitempty"`
}

type DiskPartition struct {
	Mountpoint  string
	Fstype      string
	Total       uint64
	Free        uint64
	Used        uint64
	UsedPercent float64
}

type Disk struct {
	UUID       string
	Device     string
	Partitions []DiskPartition `json:"partitions"`
}
