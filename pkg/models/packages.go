package models

// Package is used to store relevant information about installed by package manager.
// This structure is returned by method `Grab`
type Package struct {
	Name              string `json:"name,omitempty"`
	Version           string `json:"version,omitempty"`
	Architecture      string `json:"architecture,omitempty"`
	Description       string `json:"description,omitempty"`
	UpgradableVersion string `json:"upgradable_version,omitempty"`
	IsUpToDate        bool   `json:"is_up_to_date,omitempty"`
}

type PackageVulnMatch struct {
	PackageName      string        `json:"name"`
	InstalledVersion string        `json:"version"`
	Vulnerabilities  []MatchedVuln `json:"vulnerabilities"`
	Matched          bool
}

type MatchedVuln struct {
	VulnerabilityId string `json:"id"`
	Severity        string `json:"severity"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	FixedVersion    string `json:"fixedVersion"`
}

type PackageVulnSummary struct {
	PackageName      string `json:"name"`
	InstalledVersion string `json:"version"`
	VulnCount        int    `json:"vulnCount"`
}

type DivergingPackage struct {
	Name     string
	Versions []string
	Hosts    []string
}

type SimilarPackage struct {
	Name    string
	Version string
	Hosts   []string
}

type UpgradablePackage struct {
	Name    string
	Version string
}

type PackageCount struct {
	Total int64  `json:"total"`
	Name  string `json:"name"`
}
