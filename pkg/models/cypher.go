package models

// CypherNode représente un nœud générique retourné par Cypher
type CypherNode struct {
	ID         string                 `json:"id"`
	Labels     []string               `json:"labels"`
	Properties map[string]interface{} `json:"properties"`
}

// CypherRelation représente une relation générique retournée par Cypher
type CypherRelation struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	StartNode  string                 `json:"startNode"`
	EndNode    string                 `json:"endNode"`
	Properties map[string]interface{} `json:"properties"`
}

// CypherResult regroupe les nœuds et relations retournés
type CypherResult struct {
	Nodes         []CypherNode     `json:"nodes"`
	Relationships []CypherRelation `json:"relationships"`
}
