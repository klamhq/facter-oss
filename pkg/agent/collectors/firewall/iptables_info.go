package firewall

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/coreos/go-iptables/iptables"
	"github.com/klamhq/facter-oss/pkg/utils"
	"github.com/sirupsen/logrus"
)

type RegexpDispatch struct {
	FilterElem []string
	Regex      *regexp.Regexp
	Log        *logrus.Logger
}

type Table struct {
	Name []string
}

type Chain struct {
	NameTable string
	NameChain []string
}

type Rules struct {
	Option        string `json:"option,omitempty"`
	Chain         string `json:"chain,omitempty"`
	Table         string `json:"table,omitempty"`
	Access        string `json:"access,omitempty"`
	Counter       string `json:"counter,omitempty"`
	Option2       string `json:"option2,omitempty"`
	Option3       string `json:"option3,omitempty"`
	Option4       string `json:"option4,omitempty"`
	Option5       string `json:"option5,omitempty"`
	Option6       string `json:"option6,omitempty"`
	Option7       string `json:"option7,omitempty"`
	Option8       string `json:"option8,omitempty"`
	Option9       string `json:"option9,omitempty"`
	Option10      string `json:"option10,omitempty"`
	Option11      string `json:"option11,omitempty"`
	Option2Value  string `json:"option2Value,omitempty"`
	Option3Value  string `json:"option3Value,omitempty"`
	Option4Value  string `json:"option4Value,omitempty"`
	Option5Value  string `json:"option5Value,omitempty"`
	Option6Value  string `json:"option6Value,omitempty"`
	Option7Value  string `json:"option7Value,omitempty"`
	Option8Value  string `json:"option8Value,omitempty"`
	Option9Value  string `json:"option9Value,omitempty"`
	Option10Value string `json:"option10Value,omitempty"`
	Option11Value string `json:"option11Value,omitempty"`
	ValueCount    string `json:"valueCount,omitempty"`
	ValueCount2   string `json:"valueCount2,omitempty"`
}

type IptablesRules struct {
	Init         *iptables.IPTables
	IptablesPath string
	logger       *logrus.Logger
}

// IsApplicable should return true if the current PackageParser is compatible with current system.
func (i *IptablesRules) IsApplicable() bool {
	i.IptablesPath = "iptables"
	_, err := exec.LookPath(i.IptablesPath)
	return err == nil
}

type IptablesRulesStruct struct {
	Chain                string `json:"Chain,omitempty"`                //ex: prerouting input, output, docker, postrouting
	Table                string `json:"Table,omitempty"`                // type of iptables tables: nat, mangle, security, filter, raw
	MethodNegate         string `json:"MethodNegate,omitempty"`         //negate
	MethodDeny           string `json:"MethodDeny,omitempty"`           //Deny
	MethodAccept         string `json:"MethodAccept,omitempty"`         //Accept
	ParamCount           string `json:"ParamCount,omitempty"`           // param -c
	ValueCountInput      string `json:"ValueCountInput,omitempty"`      //value of input paramcount
	ValueCountOutput     string `json:"ValueCountOutput,omitempty"`     //value of output coutner paramacount
	ParamChain           string `json:"ParamChain,omitempty"`           // first param for gat,update, delete, append a iptale chain
	ValueChain           string `json:"ValueChain,omitempty"`           // value of paramchain
	ParamSelectInput     string `json:"ParamSelectInput,omitempty"`     // -i for select interface input network, interface
	ValueSelectInput     string `json:"ValueSelectInput,omitempty"`     // value of param select input
	ParamSelectOutput    string `json:"ParamSelectOutput,omitempty"`    // -o for select interface network output
	ValueSelectOutput    string `json:"ValueSelectOutput,omitempty"`    // value of select output
	ParamJump            string `json:"ParamJump,omitempty"`            // param jump -j
	ValueJump            string `json:"ValueJump,omitempty"`            // value of param jump
	ParamMatch           string `json:"ParamMatch,omitempty"`           // -m
	ValueMatch           string `json:"ValueMatch,omitempty"`           //value of match
	ParamProtocol        string `json:"ParamProtocol,omitempty"`        // param -p
	ValueProtocol        string `json:"ValueProtocol,omitempty"`        // value of -p
	ParamSource          string `json:"ParamSource,omitempty"`          // -s
	ValueSource          string `json:"ValueSource,omitempty"`          // Value of -s
	ParamDestination     string `json:"ParamDestination,omitempty"`     // -d
	ValueDestination     string `json:"ValueDestination,omitempty"`     // value of -d
	ParamDestinationPort string `json:"ParamDestinationPort,omitempty"` // --dport
	ValueDestinationPort string `json:"ValueDestinationPort,omitempty"` // Value of --dport
	ParamSourcePort      string `json:"ParamSourcePort,omitempty"`      // --sport
	ValueSourcePort      string `json:"ValueSourcePort,omitempty"`      // value of --sport
	ParamLimit           string `json:"ParamLimit,omitempty"`           // --limit
	ValueLimit           string `json:"ValueLimit,omitempty"`           // value of --limit
	ParamLimitBurst      string `json:"ParamLimitBurst,omitempty"`      // --limit-burst
	ValueLimitBurst      string `json:"ValueLimitBurst,omitempty"`      // value of --limit-burst
	ParamIcmpType        string `json:"ParamIcmp,omitempty"`            // --icmp-type
	ValueIcmpType        string `json:"ValueIcmp,omitempty"`            // value of --icmp-type
	ParamDestinationType string `json:"ParamDestinationType,omitempty"` // --dst-type
	ValueDestinationType string `json:"ValueDestinationType,omitempty"` // value of dst-type
	ParamCstate          string `json:"ParamCstate,omitempty"`          // --ctstate
	ValueCstate          string `json:"ValueCstate,omitempty"`          // value --cstate

}

func (i *IptablesRules) NewIptablesRules(logger *logrus.Logger) (*IptablesRules, error) {
	init, err := iptables.New()
	if err != nil {
		return nil, fmt.Errorf("initialize failed: %v", err)
	}
	return &IptablesRules{
		Init: init,
	}, nil
}

// Version  Return a string containing various iptables information.
func (i *IptablesRules) Version() string {
	v1, v2, v3 := i.Init.GetIptablesVersion()
	version := fmt.Sprintf("%d.%d.%d", v1, v2, v3)
	return version
}

// GetResults return iptablesRulesStruct
func (i *IptablesRules) GetResults(logger *logrus.Logger, table string) []*IptablesRulesStruct {
	resChain := i.GetListChainForTable(logger, table)
	resTable := []*IptablesRulesStruct{}
	for _, elemChain := range resChain {
		z := i.ListRule(table, elemChain)
		resTable = append(resTable, z...)
	}

	return resTable
}

// GetAvailableTables Return a list containing pre-defined list of table.
func (i *IptablesRules) GetAvailableTables() []string {
	return []string{"nat", "filter", "mangle", "security", "raw"}
}

// GetListChainForTable return a list containing x chain for all iptables's table present in the system
func (i *IptablesRules) GetListChainForTable(logger *logrus.Logger, table string) []string {
	chainListRule, err := i.Init.ListChains(table)
	if err != nil {
		i.logger.Errorf("List Chain failed: %v", err)
	}
	return chainListRule
}

// filterElemCountSize Filter element for get value of kind of iptables chain, ex: prerouting. And Count the number of string split whitespace for build map with good size
func (f *RegexpDispatch) filterElemCountSize(chain string) (string, int) {
	f.FilterElem = strings.Split(chain, " ")
	size := len(f.FilterElem)
	return f.FilterElem[1], size
}

// Regexer regex parser
func (f *RegexpDispatch) Regexer(regexpTable []string, chain string) *Rules {
	chainResult := &Rules{}
	for _, regexpValue := range regexpTable {
		parser := IptablesRules{}
		f.Regex = regexp.MustCompile(regexpValue)
		mapTable, match := createMap(f.Regex, chain)
		chainResult := parser.parseIptablesOutput(match, mapTable, f.Regex)
		// Continue if chain if not parse in one shot
		if chainResult == nil && match == nil {
			continue
		}
		return chainResult
	}
	return chainResult
}

// Dispatcher manage all regex
func (f *RegexpDispatch) Dispatcher(size int, elementToDispatch, chain string) *Rules {
	for range elementToDispatch {
		switch size {
		case 2:
			regexpTable := []string{`^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))$`, `^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z-0-9]+))$`}
			res := f.Regexer(regexpTable, chain)
			return res
		case 6:
			regexpTable := []string{`(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))\s(?P<access>(?:[A-Z]+))\s(?P<counter>(?:-)\S)\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))$`}
			res := f.Regexer(regexpTable, chain)
			return res
		case 7:
			regexpTable := []string{`^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z+-A-Z]+))\s(?P<counter>(?:-)\S)\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option2>(?:-)\S)\s(?P<option2Value>(?:[\\w\s\S]+))$`, `^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))\s(?P<counter>(?:-)\S)\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option2>(?:-)\S)\s(?P<option2Value>(?:[\\w\s\S]+))$`}
			res := f.Regexer(regexpTable, chain)
			return res
		case 9:
			regexpTable := []string{`^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))\s(?P<option2>(?:-.))\s(?P<option2Value>(?:\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b/\d{1,2}))\s(?P<counter>(?:-.))\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option4>(?:-[a-z]+))\s(?P<option4Value>(?:[A-Z]+))$`, `^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))\s(?P<option2>(?:-.))\s(?P<option2Value>(?:[a-z])+)\s(?P<counter>(?:-)\S)\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option3>(?:-.))\s(?P<option3Value>(?:[A-Z]+))$`, `^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))\s(?P<option2>(?:-.))\s(?P<option2Value>(?:\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b/\d{1,2}))\s(?P<counter>(?:-)\S)\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option3>(?:-.))\s(?P<option3Value>(?:[A-Z]+))$`, `^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))\s(?P<option2>(?:-)\S)\s(?P<option2Value>(?:[a-z-0-9]+))\s(?P<counter>(?:-.))\s(?P<valueCount>(?:[\d]))\s(?P<valueCount2>(?:[\d]))\s(?P<option3>(?:(-.)+))\s(?P<option3Value>(?:[A-Z]+))$`, `^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))\s(?P<option2>(?:-.))\s(?P<option2Value>(?:[a-z-0-9]+))\s(?P<counter>(?:-.))\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option3>(?:-[a-z]+))\s(?P<option3Value>(?:[A-Z]+))$`, `^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z+-A-Z]+))\s(?P<option2>(?:-[a-z]))\s(?P<option2Value>(?:[a-z-0-9]+))\s(?P<counter>(?:-.))\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option3>(?:-)\S)\s(?P<option3Value>(?:[\\w\s\S]+))$`}
			res := f.Regexer(regexpTable, chain)
			return res
		case 11:
			regexpTable := []string{`^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))\s(?P<option2>(?:-.))\s(?P<option2Value>(?:[a-z-0-9]+))\s(?P<option3>(?:-.))\s(?P<option3Value>(?:[a-z-0-9]+))\s(?P<counter>(?:-.))\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option4>(?:-[a-z]+))\s(?P<option4Value>(?:[A-Z]+))$`, `^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))\s(?P<option2>(?:-.))\s(?P<option2Value>(?:[a-z0-9]+))\s(?P<option3>(?:-[a-z]+))\s(?P<option3Value>(?:[a-z0-9]+))\s(?P<counter>(?:-.))\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option4>(?:-[a-z]+))\s(?P<option4Value>(?:[A-Z]+))$`, `^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))\s(?P<option2>(?:-.))\s(?P<option2Value>(?:\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b/\d{1,2}))\s(?P<option3>(?:-[a-z]+))\s(?P<option3Value>(?:[a-z0-9]+))\s(?P<counter>(?:-)\S)\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option4>(?:-.))\s(?P<option4Value>(?:[A-Z]+))$`, `^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))\s(?P<option2>(?:-.))\s(?P<option2Value>(?:[a-z]+))\s(?P<option3>(?:--[a-z]+-[a-z]+))\s(?P<option3Value>(?:[A-Z]+))\s(?P<counter>(?:-.))\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option4>(?:-[a-z]+))\s(?P<option4Value>(?:[A-Z]+))$`, `^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))\s(?P<option2>(?:-.))\s(?P<option2Value>(?:[a-z]+))\s(?P<option3>(?:--[a-z]+-[a-z]+))\s(?P<option3Value>(?:[A-Z]+))\s(?P<counter>(?:-.))\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option4>(?:-[a-z]+))\s(?P<option4Value>(?:[A-Z]+))$`}
			res := f.Regexer(regexpTable, chain)
			return res
		case 12:
			regexpTable := []string{`^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))\s(?P<option2>(?:-.))\s(?P<option2Value>(?:\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b/\d{1,2}))\s(?P<access>(?:!))\s(?P<option3>(?:-.))\s(?P<option3Value>(?:[a-z-0-9]+))\s(?P<counter>(?:-.))\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option4>(?:-[a-z]+))\s(?P<option4Value>(?:[A-Z]+))$`, `^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))\s(?P<option2>(?:-.))\s(?P<option2Value>(?:[a-z-0-9]+))\s(?P<access>(?:!))\s(?P<option3>(?:-.))\s(?P<option3Value>(?:[a-z-0-9]+))\s(?P<counter>(?:-.))\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option4>(?:-.))\s(?P<option4Value>(?:[A-Z]+))$`, `^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+)-[A-Z]+-[A-Z]+-[\d])\s(?P<option2>(?:-.)+)\s(?P<option2Value>(?:[a-z]+-[a-z0-9]+))\s(?P<access>(?:!))\s(?P<option3>(?:-.))\s(?P<option3Value>(?:[a-z-0-9]+))\s(?P<counter>(?:-.))\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option4>(?:-.))\s(?P<option4Value>(?:[A-Z]+)-[A-Z]+-[A-Z]+-[\d])$`, `^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+)-[A-Z]+-[A-Z]+-[\d])\s(?P<option2>(?:-.)+)\s(?P<option2Value>(?:[a-z])+[\d])\s(?P<access>(?:(!)))\s(?P<option3>(?:-.)+)\s(?P<option3Value>(?:[a-z])+[\d])\s(?P<counter>(?:-.)+)\s(?P<valueCount>(?:[\d])+)\s(?P<valueCount2>(?:[\d])+)\s(?P<option4>(?:-.)+)\s(?P<option4Value>(?:[A-Z]+)-[A-Z]+-[A-Z]+-[\d])$`}
			res := f.Regexer(regexpTable, chain)
			return res
			//f.Regex = regexp.MustCompile(`^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))\s(?P<option2>(?:-.))\s(?P<option2Value>(?:\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b/\d{1,2}))\s(?P<option3>(?:!))\s(?P<option4>(?:-.))\s(?P<option4Value>(?:[a-z-0-9]+))\s(?P<counter>(?:-.))\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option5>(?:-[a-z]+))\s(?P<option5Value>(?:[A-Z]+))$`)
		case 13:
			regexpTable := []string{`^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))\s(?P<option2>(?:-.))\s(?P<option2Value>(?:[a-z])+)\s(?P<option3>(?:-[a-z]))\s(?P<option3Value>(?:[a-z])+)\s(?P<option4>(?:[--]+[a-z-a-z]+))\s(?P<option4Value>(?:[\d])+)\s(?P<counter>(?:-)\S)\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option5>(?:-.))\s(?P<option5Value>(?:[A-Z]+))$`, `^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))\s(?P<option2>(?:-.))\s(?P<option2Value>(?:[a-z0-9])+)\s(?P<option3>(?:-.))\s(?P<option3Value>(?:[a-z])+)\s(?P<option4>(?:[--]+[a-z]+))\s(?P<option4Value>(?:[\d])+)\s(?P<counter>(?:-)\S)\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option5>(?:-.))\s(?P<option5Value>(?:[A-Z]+))$`, `^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))\s(?P<option2>(?:-.))\s(?P<option2Value>(?:\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b/\d{1,2}))\s(?P<option3>(?:-[a-z]+))\s(?P<option3Value>(?:[a-z0-9]+))\s(?P<option4>(?:-[a-z]+))\s(?P<option4Value>(?:[a-z]+))\s(?P<counter>(?:-)\S)\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option5>(?:-.))\s(?P<option5Value>(?:[A-Z]+))$`, `^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))\s(?P<option2>(?:-.))\s(?P<option2Value>(?:[a-z-0-9]+))\s(?P<option3>(?:-.))\s(?P<option3Value>(?:[a-z]+))\s(?P<option4>(?:--[a-z]+))\s(?P<option4Value>(?:[A-Z,]+))\s(?P<counter>(?:-.))\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option5>(?:-[a-z]+))\s(?P<option5Value>(?:[A-Z]+))$`}
			res := f.Regexer(regexpTable, chain)
			return res
		case 14:
			regexpTable := []string{`^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))\s(?P<access>(?:!))\s(?P<option2>(?:-[a-z]+))\s(?P<option2Value>(?:\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b/\d{1,2}))\s(?P<option3>(?:-.))\s(?P<option3Value>(?:[a-z]+))\s(?P<option4>(?:--[a-z]+-[a-z]+))\s(?P<option4Value>(?:[A-Z]+))\s(?P<counter>(?:-.))\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option5>(?:-[a-z]+))\s(?P<option5Value>(?:[A-Z]+))$`}
			res := f.Regexer(regexpTable, chain)
			return res
		case 15:
			regexpTable := []string{`^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))\s(?P<option2>(?:-.))\s(?P<option2Value>(?:[a-z0-9])+)\s(?P<option3>(?:-.))\s(?P<option3Value>(?:[a-z])+)\s(?P<option4>(?:-.))\s(?P<option4Value>(?:[a-z])+)\s(?P<option5>(?:[--]+[a-z]+))\s(?P<option5Value>(?:[\d])+)\s(?P<counter>(?:-)\S)\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option6>(?:-.))\s(?P<option6Value>(?:[A-Z]+))$`}
			res := f.Regexer(regexpTable, chain)
			return res
		case 17:
			regexpTable := []string{`^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))\s(?P<option2>(?:-.))\s(?P<option2Value>(?:\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b/\d{1,2}))\s(?P<option3>(?:-.))\s(?P<option3Value>(?:\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b/\d{1,2}))\s(?P<option4>(?:-.))\s(?P<option4Value>(?:[a-z]+))\s(?P<option5>(?:-.))\s(?P<option5Value>(?:[a-z]+))\s(?P<option6>(?:--[a-z]+))\s(?P<option6Value>(?:\d+))\s(?P<counter>(?:-.))\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option7>(?:-[a-z]+))\s(?P<option7Value>(?:[A-Z]+))$`, `^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))\s(?P<access>(?:!))\s(?P<option2>(?:-[a-z]+))\s(?P<option2Value>(?:[a-z]+-[a-z0-9]+))\s(?P<option3>(?:-[a-z]))\s(?P<option3Value>(?:[a-z]+))\s(?P<option4>(?:-[a-z]))\s(?P<option4Value>(?:[a-z]+))\s(?P<option5>(?:--[a-z]+))\s(?P<option5Value>(?:[\d]+))\s(?P<counter>(?:-)\S)\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option6>(?:-[a-z]))\s(?P<option6Value>(?:[A-Z]+))\s(?P<option7>(?:--[a-z]+-[a-z]+))\s(?P<option7Value>(?:\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b:\d+))\s(?P<option8>(?:-[a-z]+))\s(?P<option8Value>(?:[A-Z]+))$`}
			res := f.Regexer(regexpTable, chain)
			return res
		case 18:
			regexTable := []string{`^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))\s(?P<access>(?:!))\s(?P<option2>(?:-.))\s(?P<option2Value>(?:[a-z-0-9]+))\s(?P<option3>(?:-.))\s(?P<option3Value>(?:[a-z]+))\s(?P<option4>(?:-.))\s(?P<option4Value>(?:[a-z]+))\s(?P<option5>(?:--[a-z]+))\s(?P<option5Value>(?:[\d]+))\s(?P<counter>(?:-)\S)\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option6>(?:-.))\s(?P<option6Value>(?:[A-Z]+))\s(?P<option7>(?:--[a-z-a-z]+))\s(?P<option7Value>(?:[\d]+.[\d]+.[\d]+.[\d]+:[\d]+))$`}
			res := f.Regexer(regexTable, chain)
			return res
		case 19:
			regexpTable := []string{`^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))\s(?P<option2>(?:-.))\s(?P<option2Value>(?:[a-z0-9]+))\s(?P<option3>(?:-.))\s(?P<option3Value>(?:[a-z])+)\s(?P<option4>(?:-.))\s(?P<option4Value>(?:[a-z])+)\s(?P<option5>(?:--[a-z]+))\s(?P<option5Value>(?:[\d,])+)\s(?P<option6>(?:-[a-z]+))\s(?P<option6Value>(?:[a-z])+)\s(?P<option7>(?:--[a-z]+))\s(?P<option7Value>(?:[A-Z,])+)\s(?P<counter>(?:-)\S)\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option9>(?:-.))\s(?P<option9Value>(?:[A-Z]+))$`, `^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))\s(?P<option2>(?:-.))\s(?P<option2Value>(?:[a-z0-9])+)\s(?P<option3>(?:-.))\s(?P<option3Value>(?:[a-z])+)\s(?P<option4>(?:[--]+[a-z]+))\s(?P<option4Value>(?:[\d])+)\s(?P<option5>(?:-[a-z])+)\s(?P<option5Value>(?:[a-z])+)\s(?P<option6>(?:--[a-z]+))\s(?P<option6Value>(?:[\d/a-z])+)\s(?P<option7>(?:[--]+[a-z-a-z]+)\S)\s(?P<option7Value>(?:[\d])+)\s(?P<counter>(?:-)\S)\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option8>(?:-.))\s(?P<option8Value>(?:[A-Z]+))$`, `^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))\s(?P<option2>(?:-.))\s(?P<option2Value>(?:[a-z0-9])+)\s(?P<option3>(?:-.))\s(?P<option3Value>(?:[a-z])+)\s(?P<option4>(?:-.))\s(?P<option4Value>(?:[a-z])+)\s(?P<option5>(?:[--]+[a-z]+))\s(?P<option5Value>(?:[\d])+)\s(?P<option6>(?:-[a-z]+))\s(?P<option6Value>(?:[a-z])+)\s(?P<option7>(?:[--]+[a-z]+)\S)\s(?P<option7Value>(?:[A-Z,])+)\s(?P<counter>(?:-)\S)\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option8>(?:-.))\s(?P<option8Value>(?:[A-Z]+))$`, `^^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))\s(?P<option2>(?:-.))\s(?P<option2Value>(?:[a-z0-9])+)\s(?P<option3>(?:-.))\s(?P<option3Value>(?:[a-z])+)\s(?P<option4>(?:--[a-z]+))\s(?P<option4Value>(?:[0-9])+)\s(?P<option5>(?:-.))\s(?P<option5Value>(?:[a-z-a-z])+)\s(?P<option6>(?:--[a-z]+))\s(?P<option6Value>(?:[0-9/a-z])+)\s(?P<option7>(?:--[a-z-a-z]+))\s(?P<option7Value>(?:[0-9])+)\s(?P<counter>(?:-)\S)\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option8>(?:-.))\s(?P<option8Value>(?:[A-Z]+))`}
			res := f.Regexer(regexpTable, chain)
			return res
		case 20:
			regexpTable := []string{`^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))\s(?P<option2>(?:-.))\s(?P<option2Value>(?:\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b/\d{1,2}))\s(?P<access>(?:!))\s(?P<option3>(?:[a-z-0-9]+))\s(?P<option3Value>(?:[a-z-0-9]+))\s(?P<option4>(?:-.))\s(?P<option4Value>(?:[a-z-0-9]+))\s(?P<option5>(?:-.))\s(?P<option5Value>(?:[a-z]+))\s(?P<option6>(?:-.))\s(?P<option6Value>(?:[a-z]+))\s(?P<option7>(?:--[a-z]+))\s(?P<option7Value>(?:[\d]+))\s(?P<counter>(?:-)\S)\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option8>(?:-.))\s(?P<option8Value>(?:[A-Z]+))$`}
			res := f.Regexer(regexpTable, chain)
			return res
		case 21:
			regexpTable := []string{`^(?P<option>(?:-.))\s(?P<chain>(?:[A-Z]+))\s(?P<option2>(?:-.))\s(?P<option2Value>(?:\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b/\d{1,2}))\s(?P<option3>(?:-.))\s(?P<option3Value>(?:[a-z0-9])+)\s(?P<option4>(?:-.))\s(?P<option4Value>(?:[a-z])+)\s(?P<option5>(?:-.))\s(?P<option5Value>(?:[a-z])+)\s(?P<option6>(?:[--]+[a-z]+))\s(?P<option6Value>(?:[\d])+)\s(?P<option7>(?:-)\S)\s(?P<option7Value>(?:[a-z])+)\s(?P<option8>(?:--[a-z]+))\s(?P<option8Value>(?:[A-Z,])+)\s(?P<counter>(?:-)\S)\s(?P<valueCount>(?:[\d]+))\s(?P<valueCount2>(?:[\d]+))\s(?P<option9>(?:-.))\s(?P<option9Value>(?:[A-Z]+))$`}
			res := f.Regexer(regexpTable, chain)
			return res
		default:
			logrus.Debugf("Unable to parse firewall rules %s", chain)
		}
	}
	return nil
}

// runConverterStruct
func (i *IptablesRulesStruct) runConverterStruct(chain *Rules) *IptablesRulesStruct {
	i.defineElem(chain)
	return i
}

// RunDispatcher
func RunDispatcher(listChain []string) []*IptablesRulesStruct {
	var factory utils.LoggerFactory = &utils.DefaultLoggerFactory{}
	logger := factory.New(logrus.ErrorLevel)
	item := make([]*IptablesRulesStruct, 0)
	for _, element := range listChain {

		t := &IptablesRulesStruct{}
		// filter the second value of the chain for filter the regexp
		filter := RegexpDispatch{}
		// Count the number of value to parse in regexp and get the size
		toDispatch, size := filter.filterElemCountSize(element)
		results := filter.Dispatcher(size, toDispatch, element)
		if results.Chain == "" {
			logger.Debugf("Unable to parse firewall rule, ignoring: %s", element)
			continue
		}
		res := t.runConverterStruct(results)

		item = append(item, res)

	}
	return item
}

func (i *IptablesRules) ListRule(table, chain string) []*IptablesRulesStruct {

	list, err := i.Init.ListWithCounters(table, chain)

	if err != nil {
		logrus.Errorf("list Chain failed, missing root privileges ? \n: %v", err)
	}
	result := RunDispatcher(list)

	return result
}

func createMap(regex *regexp.Regexp, chain string) (map[string]string, []string) {
	match := regex.FindStringSubmatch(chain)
	// Get the size of the regexp
	sizeOfTable := len(regex.SubexpNames())
	// Create the table with te good size
	mapTable := make(map[string]string, sizeOfTable)
	return mapTable, match
}

func (i *IptablesRules) parseIptablesOutput(match []string, mapTable map[string]string, regex *regexp.Regexp) *Rules {
	for i, name := range regex.SubexpNames() {

		if i >= len(match) {
			return nil
		}
		mapTable[name] = match[i]

	}

	if len(mapTable) == 0 {
		return nil
	}

	return &Rules{
		Chain:         mapTable["chain"],
		Table:         mapTable["table"],
		Access:        mapTable["access"],
		Counter:       mapTable["counter"],
		Option:        mapTable["option"],
		Option2:       mapTable["option2"],
		Option3:       mapTable["option3"],
		Option4:       mapTable["option4"],
		Option5:       mapTable["option5"],
		Option6:       mapTable["option6"],
		Option7:       mapTable["option7"],
		Option8:       mapTable["option8"],
		Option9:       mapTable["option9"],
		Option10:      mapTable["option10"],
		Option11:      mapTable["option11"],
		Option2Value:  mapTable["option2Value"],
		Option3Value:  mapTable["option3Value"],
		Option4Value:  mapTable["option4Value"],
		Option5Value:  mapTable["option5Value"],
		Option6Value:  mapTable["option6Value"],
		Option7Value:  mapTable["option7Value"],
		Option8Value:  mapTable["option8Value"],
		Option9Value:  mapTable["option9Value"],
		Option10Value: mapTable["option10Value"],
		Option11Value: mapTable["option11Value"],
		ValueCount:    mapTable["valueCount"],
		ValueCount2:   mapTable["valueCount2"],
	}

}

func (i *IptablesRulesStruct) multiRegexer(ruleOption, ruleOptionValue string) {
	switch ruleOption {
	case "-m":
		i.ParamMatch = ruleOption
		i.ValueMatch = ruleOptionValue
	case "-d":
		i.ParamDestination = ruleOption
		i.ValueDestination = ruleOptionValue
	case "-s":
		i.ParamSource = ruleOption
		i.ValueSource = ruleOptionValue
	case "-i":
		i.ParamSelectInput = ruleOption
		i.ValueSelectInput = ruleOptionValue
	case "-j":
		i.ParamJump = ruleOption
		i.ValueJump = ruleOptionValue
	case "-o":
		i.ParamSelectOutput = ruleOption
		i.ValueSelectOutput = ruleOptionValue
	case "--dst-type":
		i.ParamDestinationType = ruleOption
		i.ValueDestinationType = ruleOptionValue
	case "!":
		i.MethodNegate = ruleOptionValue
	case "-p":
		i.ParamProtocol = ruleOption
		i.ValueProtocol = ruleOptionValue
	case "--ctstate":
		i.ParamCstate = ruleOption
		i.ValueCstate = ruleOptionValue
	case "--dport", "--dports":
		i.ParamDestinationPort = ruleOption
		i.ValueDestinationPort = ruleOptionValue
	case "--sport", "--sports":
		i.ParamSourcePort = ruleOption
		i.ValueSourcePort = ruleOptionValue
	case "--limit":
		i.ParamLimit = ruleOption
		i.ValueLimit = ruleOptionValue
	case "--limit-burst":
		i.ParamLimitBurst = ruleOption
		i.ValueLimitBurst = ruleOptionValue
	case "--icmp-type":
		i.ParamIcmpType = ruleOption
		i.ValueIcmpType = ruleOptionValue
	}
}

func (i *IptablesRulesStruct) defineElem(rules *Rules) *IptablesRulesStruct {

	i.Chain = rules.Chain
	if rules.Option != "" {
		i.ParamChain = rules.Option

	}
	if rules.Access != "" {
		switch rules.Access {
		case "!":
			i.MethodNegate = rules.Access
		case "ACCEPT":
			i.MethodAccept = rules.Access
		case "DENY":
			i.MethodDeny = rules.Access

		}
	}
	if rules.Counter != "" {
		switch rules.Counter {
		case "-c":
			i.ParamCount = rules.Counter
			i.ValueCountInput = rules.ValueCount
			i.ValueCountOutput = rules.ValueCount2

		}
	}
	// run multiRegexer for filed rules.OptionX
	if rules.Option2 != "" {
		i.multiRegexer(rules.Option2, rules.Option2Value)
	}
	if rules.Option3 != "" {
		i.multiRegexer(rules.Option3, rules.Option3Value)
	}
	if rules.Option4 != "" {
		i.multiRegexer(rules.Option4, rules.Option4Value)
	}
	if rules.Option5 != "" {
		i.multiRegexer(rules.Option5, rules.Option5Value)
	}
	if rules.Option6 != "" {
		i.multiRegexer(rules.Option6, rules.Option6Value)
	}
	if rules.Option7 != "" {
		i.multiRegexer(rules.Option7, rules.Option7Value)
	}
	if rules.Option8 != "" {
		i.multiRegexer(rules.Option8, rules.Option8Value)
	}
	if rules.Option9 != "" {
		i.multiRegexer(rules.Option9, rules.Option9Value)
	}

	return i
}
