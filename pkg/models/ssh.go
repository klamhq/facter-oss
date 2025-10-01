package models

type SshInfo struct {
	SshInfo []SshKeyInfo
}

// SshKeyInfo represents information about an SSH key.
type SshKeyInfo struct {
	Fingerprint        string   `json:"fingerprint,omitempty"`
	Type               string   `json:"type,omitempty"`
	Length             int64    `json:"length,omitempty"`
	Comment            string   `json:"comment,omitempty"`
	Path               string   `json:"path,omitempty"`
	Name               string   `json:"name,omitempty"`
	Owner              string   `json:"owner,omitempty"`
	FromAuthorizedKeys bool     `json:"fromAuthorizedKeys,omitempty"`
	Options            []string `json:"options,omitempty"`
}

// KnownHost represents a known host entry in SSH configuration.
type KnownHost struct {
	Hostname    string `json:"hostname,omitempty"`
	Type        string `json:"type,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
	Owner       string `json:"owner,omitempty"`
}

// SshKeyAccess represents access information for an SSH key.
type SshKeyAccess struct {
	Fingerprint string `json:"fingerprint,omitempty"`
	AsUser      string `json:"asUser,omitempty"`
}

type SshKeyCanConnectAsRoot struct {
	Count int64 `json:"count,omitempty"`
}
