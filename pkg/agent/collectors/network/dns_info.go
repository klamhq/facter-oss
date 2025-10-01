package network

import (
	"os"

	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
)

const (
	// defaultPath is the default path to the resolv.conf that contains information to resolve DNS. See Path().
	defaultPath = "/etc/resolv.conf"
)

// GetDnsConf reads the DNS configuration from the specified file path.
func GetDnsConf() *dns.ClientConfig {
	envPath := os.Getenv("DNSCONF_PATH")
	if envPath == "" {
		envPath = defaultPath
	}
	dnsConfig, err := dns.ClientConfigFromFile(envPath)
	if err != nil {
		logrus.Errorf("Failed to read resolv.conf file %s", err)
	}
	return dnsConfig
}
