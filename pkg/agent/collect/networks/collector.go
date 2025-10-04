package networks

import (
	"context"
	"fmt"
	"strings"

	"github.com/klamhq/facter-oss/pkg/agent/collectors/external"
	"github.com/klamhq/facter-oss/pkg/agent/collectors/firewall"
	"github.com/klamhq/facter-oss/pkg/agent/collectors/network"
	"github.com/klamhq/facter-oss/pkg/options"
	"github.com/klamhq/facter-oss/pkg/utils"
	"github.com/sirupsen/logrus"

	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
)

type NetworksCollectorImpl struct {
	log *logrus.Logger
	cfg *options.NetworksOptions
}

func New(log *logrus.Logger, cfg *options.NetworksOptions) *NetworksCollectorImpl {

	return &NetworksCollectorImpl{
		log: log,
		cfg: cfg,
	}
}
func (c *NetworksCollectorImpl) CollectNetworks(ctx context.Context) (*schema.Network, error) {
	c.log.Info("Crafting networks")
	networks := &schema.Network{}

	getDnsInfo := network.GetDnsConf()
	networks.DnsInfo = &schema.DnsInfo{
		Nameservers:   strings.Join(getDnsInfo.Servers, ","),
		SearchDomains: strings.Join(getDnsInfo.Search, ","),
		Port:          getDnsInfo.Port,
	}
	if c.cfg.PublicIp.Enabled {
		getExternalIp, err := external.GetIpInfo(c.cfg.PublicIp.PublicIpApiUrl, c.cfg.PublicIp.Timeout)
		if err != nil {
			c.log.Error("failed to get external ip addr: ", err)
		}

		if getExternalIp != nil {
			externalIp := &schema.ExternalIp{
				Ip:        getExternalIp.Ip,
				Forwarded: getExternalIp.Forwarded,
			}
			networks.ExternalIp = externalIp
		}
	}

	networks.GeoipInfo = &schema.GeoIpInfo{
		Longitude: 0,
		Latitude:  0,
		Accuracy:  0,
	}
	if c.cfg.GeoIp.Enabled {
		if c.cfg.GeoIp.GoogleGeoApikey == "" && c.cfg.GeoIp.GoogleGeoUrl == "" {
			c.log.Error("Google Api Key and Url not define")
		} else {
			// Get geoIp information
			getGeoIpInfo, err := external.GetGeoIpLocalisation(c.cfg.GeoIp.GoogleGeoApikey, c.cfg.GeoIp.GoogleGeoUrl, c.cfg.GeoIp.Timeout)
			if err != nil {
				c.log.Error("failed to get geoIp localisation information: ", err)
			}
			networks.GeoipInfo = &schema.GeoIpInfo{
				Longitude: getGeoIpInfo.GeoIpInfoLocationLatitude,
				Latitude:  getGeoIpInfo.GeoIpInfoLocationLongitude,
				Accuracy:  getGeoIpInfo.GeoIpInfoAccuracy,
			}
		}
	}

	netIfs, err := network.GetNetworkInterfaces()
	if err != nil {
		return nil, err
	}
	for _, netIf := range netIfs {

		ifProto := &schema.Interface{
			Name:         netIf.Name,
			HardwareAddr: netIf.HardwareAddress,
		}
		for _, ip := range netIf.IP {
			protoIP := &schema.Ip{}
			protoIP.Addr = ip.Addr
			protoIP.Version = ip.Version
			protoIP.Cidr = ip.CIDR
			ifProto.Ips = append(ifProto.Ips, protoIP)
		}

		networks.Interfaces = append(networks.Interfaces, ifProto)
	}
	if c.cfg.Connections.Enabled {
		err := c.craftConnections(networks)
		if err != nil {
			c.log.Error("Error when we craft Connection:", err)
		}
	}
	if c.cfg.Firewall.Enabled {
		if utils.IsRoot() {
			err = c.craftFirewall(networks)
			if err == nil {
				c.log.Debugf("%d firewall rules parsed", len(networks.Firewall.Rules))
			} else {
				c.log.Error("Unable to fetch iptables rules: ", err)
			}
		} else {
			c.log.Warn("Unable to fetch firewall config without root privileges.")
		}
	}
	return networks, nil
}

func (c *NetworksCollectorImpl) craftConnections(networks *schema.Network) error {
	c.log.Info("Crafting connections")
	connections, err := network.Connections(c.log)
	if err != nil {
		c.log.Errorf("Error during crafting connections %v", err)
	}
	networks.Connections = connections
	return nil
}

func (c *NetworksCollectorImpl) craftFirewall(networks *schema.Network) error {
	c.log.Info("Crafting firewall")
	// if facter is run with root user, he gets iptables rules
	// Initialize IptablesRules struct for begin iptables parser job
	iptablesInit := firewall.IptablesRules{}
	//Check if iptables is present, if is true, that means that iptables is present on the host
	if !iptablesInit.IsApplicable() {
		return fmt.Errorf("incompatible with current system iptables not found in path")
	}
	parser, err := iptablesInit.NewIptablesRules(c.log)
	if err != nil {
		c.log.Error("Iptables struct doesn't initialize")
	}

	networks.Firewall = &schema.Firewall{}

	// register iptables version
	networks.Firewall.Version = parser.Version()

	//register all field in rules get with iptables_info.go
	for _, table := range parser.GetAvailableTables() {
		value := parser.GetResults(c.log, table)

		for _, val := range value {
			//fmt.Println(val.Chain)
			// if val.Chain is empty, chain is not parsed
			//if val.Chain != "" {
			rule := &schema.FirewallRule{}
			rule.Chain = val.Chain
			rule.MethodNegate = val.MethodNegate
			rule.MethodDeny = val.MethodDeny
			rule.MethodAccept = val.MethodAccept
			rule.ParamCount = val.ParamCount
			rule.ValueCountInput = val.ValueCountInput
			rule.ValueCountOutput = val.ValueCountOutput
			rule.ParamChain = val.ParamChain
			rule.ValueChain = val.ValueChain
			rule.ParamSelectInput = val.ParamSelectInput
			rule.ValueSelectInput = val.ValueSelectInput
			rule.ParamSelectOutput = val.ParamSelectOutput
			rule.ValueSelectOutput = val.ValueSelectOutput
			rule.ParamJump = val.ParamJump
			rule.ValueJump = val.ValueJump
			rule.ParamMatch = val.ParamMatch
			rule.ValueMatch = val.ValueMatch
			rule.ParamProtocol = val.ParamProtocol
			rule.ValueProtocol = val.ValueProtocol
			rule.ParamSource = val.ParamSource
			rule.ValueSource = val.ValueSource
			rule.ParamDestination = val.ParamDestination
			rule.ValueDestination = val.ValueDestination
			rule.ParamDestinationPort = val.ParamDestinationPort
			rule.ValueDestinationPort = val.ValueDestinationPort
			rule.ParamDestinationType = val.ParamDestinationType
			rule.ValueDestinationType = val.ValueDestinationType
			rule.ParamCstate = val.ParamCstate
			rule.ValueCstate = val.ValueCstate
			rule.ParamSourcePort = val.ParamSourcePort
			rule.ValueSourcePort = val.ValueSourcePort
			rule.ParamLimit = val.ParamLimit
			rule.ValueLimit = val.ValueLimit
			rule.ParamLimitBurst = val.ParamLimitBurst
			rule.ValueLimitBurst = val.ValueLimitBurst
			rule.ParamIcmpType = val.ParamIcmpType
			rule.ValueIcmpType = val.ValueIcmpType
			networks.Firewall.Rules = append(networks.Firewall.Rules, rule)
			//}
		}
	}
	// No error
	return nil

}
