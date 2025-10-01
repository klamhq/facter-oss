package models

type PropertyInfo struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Mandatory bool   `json:"mandatory"`
}

type SchemaNode struct {
	ID     string         `json:"id"`
	Labels []string       `json:"labels"`
	Props  []PropertyInfo `json:"props"`
}

type SchemaRelationship struct {
	Type      string         `json:"type"`
	StartNode string         `json:"startNode"`
	EndNode   string         `json:"endNode"`
	Props     []PropertyInfo `json:"props"`
}

type SchemaVisualization struct {
	Nodes         []SchemaNode         `json:"nodes"`
	Relationships []SchemaRelationship `json:"relationships"`
}
