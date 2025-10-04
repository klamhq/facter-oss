package network

import (
	"strings"

	"github.com/klamhq/facter-oss/pkg/agent/collectors/packages"
	"github.com/klamhq/facter-oss/pkg/models"
	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
	"github.com/shirou/gopsutil/net"
	pproc "github.com/shirou/gopsutil/process"
	"github.com/sirupsen/logrus"
)

// Connections return all connections in protobuf schema
func Connections(logger *logrus.Logger) ([]*schema.ConnectionState, error) {
	cnx, err := getConnections(logger)
	if err != nil {
		return nil, err
	}
	connections := make([]*schema.ConnectionState, 0, len(cnx))
	var state schema.State
	for _, c := range cnx {
		if c.State == "ESTABLISHED" {
			state = schema.State_STATE_ESTABLISHED
		} else {
			state = schema.State_STATE_LISTENING
		}
		connection := schema.ConnectionState{
			Protocol: schema.Protocol_PROTOCOL_TCP,
			State:    state,
			Local: &schema.IpPort{
				Ip: &schema.Ip{
					Addr:    c.LocalIp,
					Version: "4",
				},
				Port: c.LocalPort,
			},
			Remote: &schema.IpPort{
				Ip: &schema.Ip{
					Addr:    c.RemoteIp,
					Version: "4",
				},
				Port: c.RemotePort,
			},
			Process: &schema.Process{
				Pid:  int64(c.Pid),
				Name: c.ProcessName,
				Package: &schema.Package{
					Name: c.Package,
				},
			},
		}
		connections = append(connections, &connection)
	}
	return connections, nil
}

// GetConnections get listen and established tcp connections and get if possible the associated package, no available for darwin
func getConnections(logger *logrus.Logger) ([]models.Connections, error) {
	conns, err := net.Connections("all")
	if err != nil {
		logger.Errorf("Error during fetching connections: %v", err)
		return nil, err
	}

	results := make([]models.Connections, 0, len(conns))
	pkgExtract, err := packages.NewPackageExtractor(logger)
	if err != nil {
		logger.Errorf("Unable to extract package from exe path: %v", err)
	}
	for _, c := range conns {
		if c.Status != "LISTEN" && c.Status != "ESTABLISHED" {
			continue
		}

		info := models.Connections{
			Protocol:   c.Type,
			LocalIp:    c.Laddr.IP,
			LocalPort:  c.Laddr.Port,
			RemoteIp:   c.Raddr.IP,
			RemotePort: c.Raddr.Port,
			State:      c.Status,
			Pid:        c.Pid,
		}

		if c.Pid != 0 {
			if proc, err := pproc.NewProcess(c.Pid); err == nil {
				if name, err := proc.Name(); err == nil {
					namePart := strings.Fields(name)
					info.ProcessName = namePart[0]
				}
				if pkgExtract != nil {
					if exe, err := proc.Exe(); err == nil {
						info.ProcessPath = exe
						info.Package = pkgExtract.GetPackage(exe)
					}
				}
			}
		}

		results = append(results, info)
	}
	return results, nil
}
