package models

type ComplianceReport struct {
	Score       Score             `json:"score"`
	Profile     string            `json:"profile"`
	RuleResults []RuleCheckResult `json:"rule_results"`
}

type Score struct {
	Maximum string `json:"maximum"`
	Value   string `json:"value"`
}

type RuleCheckResult struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Result      string `json:"result"`
	Severity    string `json:"severity"`
	Fix         string `json:"fix"`
}
