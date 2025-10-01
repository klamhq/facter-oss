package firewall

import (
	"bufio"
	"os"
	"testing"

	"github.com/klamhq/facter-oss/pkg/utils"
	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetIptablesResult(t *testing.T) {
	if _, runInCI := os.LookupEnv("CI"); runInCI {
		t.SkipNow()
	}
	if utils.CheckBinInstalled(logrus.New(), "iptables") == false {
		t.Skip("iptables is not installed")
	}
	var factory utils.LoggerFactory = &utils.DefaultLoggerFactory{}
	logger := factory.New(logrus.ErrorLevel)
	if utils.IsRoot() {
		init := IptablesRules{}
		iptables, err := init.NewIptablesRules(logger)
		if err != nil {
			logrus.Errorf("Iptables struct doesn't initialize")
		}
		assert.NoError(t, err)
		assert.NotNil(t, iptables)
		assert.NotEmpty(t, iptables)
		assert.IsType(t, &IptablesRules{}, iptables)

		// test  get version
		version := iptables.Version()
		assert.NotNil(t, version)
		assert.IsType(t, string("1.8.2"), version)
		assert.NotEmpty(t, version)

		//// check list chain
		for _, tables := range iptables.GetAvailableTables() {
			assert.NotEmpty(t, tables)
			assert.IsType(t, string("nat"), tables)
			value := iptables.GetResults(logger, tables)
			assert.IsType(t, []*IptablesRulesStruct{}, value)

			for _, val := range value {
				//	fmt.Println(val.Chain)
				rule := &schema.FirewallRule{}

				rule.Chain = val.Chain
				assert.IsType(t, string("INPUT"), rule.Chain)

				rule.MethodNegate = val.MethodNegate
				assert.IsType(t, string("!"), rule.MethodNegate)

				rule.MethodDeny = val.MethodDeny
				assert.IsType(t, string("DENY"), rule.MethodDeny)

				rule.MethodAccept = val.MethodAccept
				assert.IsType(t, string("ACCEPT"), rule.MethodAccept)

				rule.ParamCount = val.ParamCount
				assert.IsType(t, string("-c"), rule.ParamCount)

				rule.ValueCountInput = val.ValueCountInput
				assert.IsType(t, string("0"), rule.ValueCountInput)

				rule.ValueCountOutput = val.ValueCountOutput
				assert.IsType(t, string("0"), rule.ValueCountOutput)

				rule.ParamChain = val.ParamChain
				assert.IsType(t, string("-A"), rule.ParamChain)

				rule.ValueChain = val.ValueChain
				assert.IsType(t, string("PREROUTING"), rule.ValueChain)

				rule.ParamSelectInput = val.ParamSelectInput
				assert.IsType(t, string("-i"), rule.ParamSelectInput)

				rule.ValueSelectInput = val.ValueSelectInput
				assert.IsType(t, string("eth1"), rule.ValueSelectInput)

				rule.ParamSelectOutput = val.ParamSelectOutput
				assert.IsType(t, string("-o"), rule.ParamSelectOutput)

				rule.ValueSelectOutput = val.ValueSelectOutput
				assert.IsType(t, string("eth0"), rule.ValueSelectOutput)

				rule.ParamJump = val.ParamJump
				assert.IsType(t, string("-j"), rule.ParamJump)

				rule.ValueJump = val.ValueJump
				assert.IsType(t, string("DROP"), rule.ValueJump)

				rule.ParamMatch = val.ParamMatch
				assert.IsType(t, string("-m"), rule.ParamMatch)

				rule.ValueMatch = val.ValueMatch
				assert.IsType(t, string("tcp"), rule.ValueMatch)

				rule.ParamProtocol = val.ParamProtocol
				assert.IsType(t, string("-p"), rule.ParamProtocol)

				rule.ValueProtocol = val.ValueProtocol
				assert.IsType(t, string("tcp"), rule.ValueProtocol)

				rule.ParamSource = val.ParamSource
				assert.IsType(t, string("-s"), rule.ParamSource)

				rule.ValueSource = val.ValueSource
				assert.IsType(t, string("192.168.1.1"), rule.ValueSource)

				rule.ParamDestination = val.ParamDestination
				assert.IsType(t, string("--dst"), rule.ParamDestination)

				rule.ValueDestination = val.ValueDestination
				assert.IsType(t, string("172.16.1.1"), rule.ValueDestination)

				rule.ParamDestinationPort = val.ParamDestinationPort
				assert.IsType(t, string("--port"), rule.ParamDestinationPort)

				rule.ValueDestinationPort = val.ValueDestinationPort
				assert.IsType(t, string("80"), rule.ValueDestinationPort)

				rule.ParamDestinationType = val.ParamDestinationType
				assert.IsType(t, string("--dst-type"), rule.ParamDestinationType)

				rule.ValueDestinationType = val.ValueDestinationType
				assert.IsType(t, string("UNICAST"), rule.ValueDestinationType)

				rule.ParamCstate = val.ParamCstate
				assert.IsType(t, string("--cstate"), rule.ParamCstate)

				rule.ValueCstate = val.ValueCstate
				assert.IsType(t, string("ESTABLISHED,NEW"), rule.ValueCstate)

				rule.ParamSourcePort = val.ParamSourcePort
				assert.IsType(t, string("-s"), rule.ParamSourcePort)

				rule.ValueSourcePort = val.ValueSourcePort
				assert.IsType(t, string("22"), rule.ValueSourcePort)

				rule.ParamLimit = val.ParamLimit
				assert.IsType(t, string("--limit"), rule.ParamLimit)

				rule.ValueLimit = val.ValueLimit
				assert.IsType(t, string("10"), rule.ValueLimit)

				rule.ParamLimitBurst = val.ParamLimitBurst
				assert.IsType(t, string("--limit-burst"), rule.ParamLimitBurst)

				rule.ValueLimitBurst = val.ValueLimitBurst
				assert.IsType(t, string("100"), rule.ValueLimitBurst)

				rule.ParamIcmpType = val.ParamIcmpType
				assert.IsType(t, string("--icmp-type"), rule.ParamIcmpType)

				rule.ValueIcmpType = val.ValueIcmpType
				assert.IsType(t, string("echo"), rule.ValueIcmpType)

			}
		}

	}

}

func TestParserIptables(t *testing.T) {
	if _, runInCI := os.LookupEnv("CI"); runInCI {
		t.SkipNow()
	}
	var factory utils.LoggerFactory = &utils.DefaultLoggerFactory{}
	logger := factory.New(logrus.ErrorLevel)
	basePath := "./export_rules"
	file, err := os.Open(basePath)
	if err != nil {
		logger.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Fatal(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	var rulesTable []string
	for scanner.Scan() {
		z := scanner.Text()
		rulesTable = append(rulesTable, z)
		assert.IsType(t, []string{}, rulesTable)
		assert.NotEmpty(t, rulesTable)
	}
	// parse file with regexp use by iptable_info
	results := RunDispatcher(rulesTable)
	assert.IsType(t, []*IptablesRulesStruct{}, results)
	assert.NotEmpty(t, results)

	if err := scanner.Err(); err != nil {
		logger.Fatal(err)
	}
}
