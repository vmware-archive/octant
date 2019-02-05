package content

var _ Content = (*DAG)(nil)

// DAG represents a directed acyclic graph, although we do not currently check for absense of cycles :)
type DAG struct {
	Type     string  `json:"type"`
	Selected string  `json:"selected"`
	Edges    AdjList `json:"adjacencyList"`
	Nodes    Nodes   `json:"objects"`
}

// Nodes is a set of graph nodes
type Nodes map[string]*Node

// Node is a node in a graph, representing a kubernetes object
// IsNetwork is a hint to the layout engine.
type Node struct {
	Name       string     `json:"name,omitempty"`
	APIVersion string     `json:"apiVersion,omitempty"`
	Kind       string     `json:"kind,omitempty"`
	Status     NodeStatus `json:"status,omitempty"`
	IsNetwork  bool       `json:"isNetwork,omitempty"`
	Views      []Content  `json:"views,omitempty"`
}

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

// NewDAG constructs an empty directed acyclic graph
func NewDAG() *DAG {
	return &DAG{
		Type:  "resourceviewer",
		Edges: AdjList{},
		Nodes: Nodes{},
	}
}

// IsEmpty implements Content.IsEmpty.
// Returns true if the graph is empty.
func (d *DAG) IsEmpty() bool {
	return len(d.Nodes) == 0
}

func (d *DAG) ViewComponent() ViewComponent {
	return ViewComponent{}
}
