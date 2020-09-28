/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

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
	Name       string      `json:"name,omitempty"`
	APIVersion string      `json:"apiVersion,omitempty"`
	Kind       string      `json:"kind,omitempty"`
	Status     NodeStatus  `json:"status,omitempty"`
	Details    []Component `json:"details,omitempty"`
	Path       *Link       `json:"path,omitempty"`
}

func (n *Node) UnmarshalJSON(data []byte) error {
	x := struct {
		Name       string         `json:"name,omitempty"`
		APIVersion string         `json:"apiVersion,omitempty"`
		Kind       string         `json:"kind,omitempty"`
		Status     NodeStatus     `json:"status,omitempty"`
		Details    []*TypedObject `json:"details,omitempty"`
		Path       *TypedObject   `json:"path,omitempty"`
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	n.Name = x.Name
	n.APIVersion = x.APIVersion
	n.Kind = x.Kind
	n.Status = x.Status

	if x.Details != nil {
		n.Details = make([]Component, len(x.Details))
	}
	for i, detail := range x.Details {
		dc, err := detail.ToComponent()
		if err != nil {
			return errors.Wrap(err, "unmarshal-ing detail")
		}
		n.Details[i] = dc
	}

	if x.Path != nil {
		p, err := x.Path.ToComponent()
		if err != nil {
			return errors.Wrap(err, "unmarshal-ing detail")
		}
		l, ok := p.(*Link)
		if !ok {
			return fmt.Errorf("path must be a link, found %q", x.Path.Metadata.Type)
		}
		n.Path = l
	}

	return nil
}

// ResourceViewerConfig is configuration for a resource viewer.
type ResourceViewerConfig struct {
	Edges    AdjList `json:"edges,omitempty"`
	Nodes    Nodes   `json:"nodes,omitempty"`
	Selected string  `json:"selected,omitempty"`
}

// ResourceView is a resource viewer component.
//
// +octant:component
type ResourceViewer struct {
	Base
	Config ResourceViewerConfig `json:"config,omitempty"`
}

// NewResourceViewer creates a resource viewer component.
func NewResourceViewer(title string) *ResourceViewer {
	return &ResourceViewer{
		Base: newBase(TypeResourceViewer, TitleFromString(title)),
		Config: ResourceViewerConfig{
			Edges: AdjList{},
			Nodes: Nodes{},
		},
	}

}

func (rv *ResourceViewer) AddEdge(nodeID, childID string, edgeType EdgeType) error {
	if _, ok := rv.Config.Nodes[childID]; !ok {
		var nodeIDs []string
		for k := range rv.Config.Nodes {
			nodeIDs = append(nodeIDs, k)
		}
		return errors.Errorf("node %q does not exist in graph. available [%s]",
			childID, strings.Join(nodeIDs, ", "))
	}

	edge := Edge{
		Node: childID,
		Type: edgeType,
	}
	rv.Config.Edges[nodeID] = append(rv.Config.Edges[nodeID], edge)

	return nil
}

func (rv *ResourceViewer) AddNode(id string, node Node) {
	rv.Config.Nodes[id] = node
}

func (rv *ResourceViewer) Select(id string) {
	rv.Config.Selected = id
}

func (rv *ResourceViewer) GetMetadata() Metadata {
	return rv.Metadata
}

func (rv *ResourceViewer) Validate() error {
	for nodeID, edges := range rv.Config.Edges {
		if _, ok := rv.Config.Nodes[nodeID]; !ok {
			var nodes []string
			for node := range rv.Config.Nodes {
				nodes = append(nodes, node)
			}
			return errors.Errorf("node %q in edges does not have a node entry. existing nodes: %s", nodeID, strings.Join(nodes, ", "))
		}

		for _, edge := range edges {
			if _, ok := rv.Config.Nodes[edge.Node]; !ok {
				return errors.Errorf("edge %q from node %q does not have a node entry", edge.Node, nodeID)
			}
		}
	}

	return nil
}

type resourceViewerMarshal ResourceViewer

// MarshalJSON implements json.Marshaler
func (rv *ResourceViewer) MarshalJSON() ([]byte, error) {
	if err := rv.Validate(); err != nil {
		return nil, errors.WithMessage(err, "validate resource viewer component")
	}

	m := resourceViewerMarshal(*rv)
	m.Metadata.Type = TypeResourceViewer
	m.Metadata.Title = rv.Metadata.Title

	return json.Marshal(&m)
}
