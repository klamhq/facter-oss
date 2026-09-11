package compliance

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/gocomply/scap/pkg/scap/models/cdf"
	"github.com/gocomply/scap/pkg/scap/scap_document"
	"github.com/klamhq/facter-oss/pkg/models"
	"github.com/klamhq/facter-oss/pkg/options"
	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
	"github.com/sirupsen/logrus"
)

// collectAllRules recursively collects all rules from groups and stores them in a map.
// This is useful for quickly accessing rules by their ID.
func collectAllRules(groups []cdf.GroupType, rules map[string]*cdf.RuleType) {
	for _, group := range groups {
		for _, rule := range group.Rule {
			rules[rule.Id] = &rule
		}
		if len(group.Group) > 0 {
			collectAllRules(group.Group, rules)
		}
	}
}

func GetDataStreamFile(osName, version string) string {
	osName = strings.ToLower(osName)
	// Rocky, Alma, RHEL : ssg-rl9-ds.xml, ssg-almalinux9-ds.xml
	if osName == "rocky" || osName == "almalinux" || osName == "rhel" {
		parts := strings.Split(version, ".")
		major := parts[0]
		prefix := map[string]string{
			"rocky":     "rl",
			"almalinux": "almalinux",
			"rhel":      "rhel",
		}[osName]
		return fmt.Sprintf("/usr/share/xml/scap/ssg/content/ssg-%s%s-ds.xml", prefix, major)
	}
	// Ubuntu : ssg-ubuntu2404-ds.xml (YYMM sans point)
	if osName == "ubuntu" {
		versionNum := strings.ReplaceAll(version, ".", "")
		return fmt.Sprintf("/usr/share/xml/scap/ssg/content/ssg-ubuntu%s-ds.xml", versionNum)
	}
	// Debian : ssg-debian12-ds.xml
	if osName == "debian" {
		parts := strings.Split(version, ".")
		major := parts[0]
		return fmt.Sprintf("/usr/share/xml/scap/ssg/content/ssg-debian%s-ds.xml", major)
	}
	// Fedora : ssg-fedora-ds.xml (rolling)
	if osName == "fedora" {
		return "/usr/share/xml/scap/ssg/content/ssg-fedora-ds.xml"
	}
	return ""
}

// Oscap runs an OpenSCAP audit and collects the results.
func Oscap(ctx context.Context, cfg *options.ComplianceOptions, operatingSystem *schema.Os, logger *logrus.Logger) (*models.ComplianceReport, error) {
	dataStreamFile := GetDataStreamFile(operatingSystem.Name, operatingSystem.Version)
	if cfg.Profile == "" {
		logger.Warn("No OpenSCAP profile specified, using default 'xccdf_org.ssgproject.content_profile_cis'")
		cfg.Profile = "xccdf_org.ssgproject.content_profile_cis"
	}

	if cfg.ResultFile == "" {
		logger.Warn("No OpenSCAP results file specified, using default 'results.xml'")
		cfg.ResultFile = "results.xml"
	}
	cmd := exec.CommandContext(ctx, "oscap", "xccdf", "eval", "--profile", cfg.Profile, "--results", cfg.ResultFile, dataStreamFile)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				code := status.ExitStatus()
				if code == 1 {
					return nil, fmt.Errorf("OpenSCAP command failed with exit code 1: %s", stderr.String())
				}
				if code == 2 {

				}
			}
		} else {
			logger.Errorf("Error executing oscap command: %v", err)
			return nil, err
		}
	}
	// Read the OpenSCAP document from the results file
	data, err := scap_document.ReadDocumentFromFile(cfg.ResultFile)
	if err != nil {
		logger.Errorf("Error reading OpenSCAP document: %v", err)
		return nil, err
	}

	ruleMap := make(map[string]*cdf.RuleType)
	for _, rule := range data.Rule {
		ruleMap[rule.Id] = &rule
	}
	collectAllRules(data.Group, ruleMap)

	report := models.ComplianceReport{}
	for _, testResult := range data.TestResult {
		var score models.Score
		for _, sc := range testResult.Score {
			score = models.Score{
				Maximum: fmt.Sprintf("%v", sc.Maximum),
				Value:   sc.Text,
			}
		}
		report.Score = score
		report.Profile = cfg.Profile

		for _, ruleResult := range testResult.RuleResult {
			if ruleResult.Result != "notselected" {
				rule := ruleMap[ruleResult.Idref]
				var title, description, fix string
				if rule != nil {
					if len(rule.Title) > 0 {
						title = rule.Title[0].InnerXml // ou .Value, selon ton type
					}
					if len(rule.Description) > 0 {
						description = rule.Description[0].InnerXml // ou .Value
					}
					if len(rule.Fix) > 0 {
						fix = rule.Fix[0].InnerXml // ou .Value
					}
					report.RuleResults = append(report.RuleResults, models.RuleCheckResult{
						ID:          ruleResult.Idref,
						Title:       title,
						Description: description,
						Result:      string(ruleResult.Result),
						Severity:    string(ruleResult.Severity),
						Fix:         fix,
					})
				} else {
					report.RuleResults = append(report.RuleResults, models.RuleCheckResult{
						ID:          ruleResult.Idref,
						Title:       "",
						Description: "",
						Result:      string(ruleResult.Result),
					})
				}
			}
		}
	}

	logger.Infof("Removed OpenSCAP results file: %s", cfg.ResultFile)
	err = os.Remove(cfg.ResultFile)
	if err != nil {
		logger.Errorf("Error removing OpenSCAP results file: %v", err)
	}
	return &report, nil
}
