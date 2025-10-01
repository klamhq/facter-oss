package network

import (
	"net"
	"strings"

	"github.com/klamhq/facter-oss/pkg/models"
)

// GetNetworkInterfaces retrieves the network interfaces on the system.
// It returns a slice of NetworkInterface structs or an error if the retrieval fails.
func GetNetworkInterfaces() ([]models.NetworkInterface, error) {
	iFaces, err := net.Interfaces()
	var result []models.NetworkInterface
	if err != nil {
		return result, err
	}

	for _, i := range iFaces {
		if i.Name == "lo" || i.Name == "lo0" {
			continue
		}
		if i.HardwareAddr == nil {
			unknownHardwareAddr := net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
			i.HardwareAddr = unknownHardwareAddr
		}
		networkInterface := models.NetworkInterface{
			Name:            i.Name,
			HardwareAddress: i.HardwareAddr.String(),
		}

		networkInterface.Flags = append(networkInterface.Flags, strings.Split(i.Flags.String(), "|")...)

		// Grabbing ip addresses
		ips, err := i.Addrs()
		if err != nil {
			return result, err
		}
		for _, addr := range ips {
			if addr == nil {
				continue
			}

			var ipNet *net.IPNet
			var ipAddr net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ipAddr = v.IP
				ipNet = v
			case *net.IPAddr:
				ipAddr = v.IP
				// créer un réseau fictif si IPAddr uniquement
				mask := net.CIDRMask(32, 32)
				if ipAddr.To4() == nil {
					mask = net.CIDRMask(128, 128)
				}
				ipNet = &net.IPNet{IP: ipAddr, Mask: mask}
			default:
				continue
			}

			// Calculer l'adresse réseau correcte (network address)
			networkIP := ipAddr.Mask(ipNet.Mask)
			networkCIDR := &net.IPNet{
				IP:   networkIP,
				Mask: ipNet.Mask,
			}

			// Déterminer la version
			version := "unknown"
			if ipAddr.To4() != nil {
				version = "4"
			} else if ipAddr.To16() != nil {
				version = "6"
			}

			ipInfo := models.IP{
				Addr:    ipAddr.String(),
				CIDR:    networkCIDR.String(), // Correct network address, ex: 192.168.1.0/24
				Version: version,
			}

			networkInterface.IP = append(networkInterface.IP, ipInfo)
		}

		result = append(result, networkInterface)
	}
	return result, nil

}

/*func (n *NetworkInterface) ToProtobuf() *schema.NetworkInterface {
	netIf := schema.NetworkInterface{}
	netIf.Name = n.Name
	netIf.HardwareAddr = n.HardwareAddress
	// TODO Guess IP version
	for _, ip := range n.IP {
		protoIp := schema.IP{Value: ip}
		netIf.Addresses = append(netIf.Addresses, &protoIp)
	}
	return &netIf
}*/
