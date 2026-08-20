package models

// Revision represents an append-only host revision envelope persisted in AGE.
type Revision struct {
	RevisionID         string `json:"revision_id"`
	Sequence           uint64 `json:"sequence"`
	PreviousRevisionID string `json:"previous_revision_id,omitempty"`
	PreviousSequence   uint64 `json:"previous_sequence,omitempty"`
	CreatedAt          string `json:"created_at"`
	AgentID            string `json:"agent_id"`
	Hostname           string `json:"hostname"`
	MachineID          string `json:"machine_id,omitempty"`
	SourceType         string `json:"source_type"`
	StateHash          string `json:"state_hash"`
	FullPayload        string `json:"full_payload,omitempty"`
	DeltaPayload       string `json:"delta_payload,omitempty"`
}
