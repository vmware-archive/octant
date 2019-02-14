package component

import "encoding/json"

// AdjList is an adjacency list - it maps nodes to edges
type AdjList map[string][]Edge

// Edge represents a directed edge in a graph
type Edge struct {
	Node string   `json:"node"`
	Type EdgeType `json:"edge"`
}

// Add adds a directed edge to the adjacency list
func (al AdjList) Add(src string, edge Edge) {
	edges, ok := al[src]
	if !ok || edges == nil {
		edges = make([]Edge, 0)
	}

	edges = append(edges, edge)
	al[src] = edges
}

type NodeStatus string

const (
	// NodeStatusOK means a node is in a health state
	NodeStatusOK NodeStatus = "ok"
	// NodeStatusWarning means ...
	NodeStatusWarning NodeStatus = "warning"
	// NodeStatusError means ...
	NodeStatusError NodeStatus = "error"
)

// EdgeType represents whether a relationship between resources is implicit or explicit
type EdgeType string

const (
	// EdgeTypeImplicit is an implicit edge
	EdgeTypeImplicit = "implicit"
	// EdgeTypeExplicit is an explicit edge
	EdgeTypeExplicit = "explicit"
)

// Nodes is a set of graph nodes
type Nodes map[string]Node

// Node is a node in a graph, representing a kubernetes object
// IsNetwork is a hint to the layout engine.
type Node struct {
	Name       string               `json:"name,omitempty"`
	APIVersion string               `json:"apiVersion,omitempty"`
	Kind       string               `json:"kind,omitempty"`
	Status     NodeStatus           `json:"status,omitempty"`
	Details    []TitleViewComponent `json:"details,omitempty"`
	// TODO: delete this
	StatusMessages []string `json:"statusMessages,omitempty"`
}

// ResourceViewerConfig is configuration for a resource viewer.
type ResourceViewerConfig struct {
	Edges AdjList `json:"edges,omitempty"`
	Nodes Nodes   `json:"nodes,omitempty"`
}

// ResourceView is a resource viewer component.
type ResourceViewer struct {
	Metadata Metadata             `json:"metadata,omitempty"`
	Config   ResourceViewerConfig `json:"config,omitempty"`
}

// NewResourceViewer creates a resource viewer component.
func NewResourceViewer(title string) *ResourceViewer {
	return &ResourceViewer{
		Metadata: Metadata{
			Type:  "resourceViewer",
			Title: []TitleViewComponent{NewText(title)},
		},
		Config: ResourceViewerConfig{
			Edges: AdjList{},
			Nodes: Nodes{},
		},
	}

}

func (rv *ResourceViewer) AddEdge(nodeID, childID string, edgeType EdgeType) {
	edge := Edge{
		Node: childID,
		Type: edgeType,
	}
	rv.Config.Edges[nodeID] = append(rv.Config.Edges[nodeID], edge)
}

func (rv *ResourceViewer) AddNode(id string, node Node) {
	rv.Config.Nodes[id] = node
}

func (rv *ResourceViewer) GetMetadata() Metadata {
	return rv.Metadata
}

// IsEmpty specifies whether the component is considered empty. Implements ViewComponent.
func (rv *ResourceViewer) IsEmpty() bool {
	return len(rv.Config.Nodes) == 0
}

type resourceViewerMarshal ResourceViewer

// MarshalJSON implements json.Marshaler
func (rv *ResourceViewer) MarshalJSON() ([]byte, error) {
	m := resourceViewerMarshal(*rv)
	m.Metadata.Type = "resourceViewer"
	m.Metadata.Title = rv.Metadata.Title

	return json.Marshal(&m)
}
