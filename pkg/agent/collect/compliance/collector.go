package compliance

import (
	"context"
	"fmt"

	"github.com/klamhq/facter-oss/pkg/agent/collectors/compliance"
	"github.com/klamhq/facter-oss/pkg/agent/collectors/system"
	"github.com/klamhq/facter-oss/pkg/options"
	"github.com/klamhq/facter-oss/pkg/utils"
	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
	"github.com/sirupsen/logrus"
)

type ComplianceCollectorImpl struct {
	log *logrus.Logger
	cfg *options.ComplianceOptions
}

func New(log *logrus.Logger, cfg *options.ComplianceOptions) *ComplianceCollectorImpl {

	return &ComplianceCollectorImpl{
		log: log,
		cfg: cfg,
	}
}

func (c *ComplianceCollectorImpl) CollectCompliance(ctx context.Context) (*schema.ComplianceReport, error) {
	if !utils.CheckBinInstalled(c.log, "oscap") {
		c.log.Error("OpenSCAP is not installed, compliance report will not be generated, see here: https://github.com/ComplianceAsCode/content?tab=readme-ov-file#installation for installation instructions")
		return nil, fmt.Errorf("openscap is not installed, compliance report will not be generated")
	}

	operatingSystem := system.GetSystem()
	os := &schema.Os{
		Name:    operatingSystem.Host.Platform,
		Version: operatingSystem.Host.PlatformVersion,
	}

	openscapReport, err := compliance.Oscap(ctx, c.cfg, os, c.log)
	if err != nil {
		return nil, err
	}
	complianceReport := &schema.ComplianceReport{}
	complianceReport.Score = &schema.Score{}
	complianceReport.RuleResults = make([]*schema.RuleCheckResult, 0, len(openscapReport.RuleResults))
	complianceReport.Score.Maximum = openscapReport.Score.Maximum
	complianceReport.Score.Value = openscapReport.Score.Value
	complianceReport.Profile = openscapReport.Profile

	for _, ruleResult := range openscapReport.RuleResults {
		ruleCheckResult := &schema.RuleCheckResult{
			Id:          ruleResult.ID,
			Title:       ruleResult.Title,
			Description: ruleResult.Description,
			Result:      ruleResult.Result,
			Severity:    ruleResult.Severity,
			Fix:         ruleResult.Fix,
		}
		complianceReport.RuleResults = append(complianceReport.RuleResults, ruleCheckResult)
	}

	return complianceReport, nil
}
