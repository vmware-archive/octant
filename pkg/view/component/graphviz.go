/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

// GraphvizConfig is the contents of Graphviz.
type GraphvizConfig struct {
	DOT string `json:"dot,omitempty"`
}

// Graphviz is a component for displaying graphviz diagrams.
//
// +octant:component
type Graphviz struct {
	Base
	Config GraphvizConfig `json:"config"`
}

// NewGraphviz creates a graphviz component.
func NewGraphviz(dot string) *Graphviz {
	return &Graphviz{
		Base: newBase(TypeGraphviz, nil),
		Config: GraphvizConfig{
			DOT: dot,
		},
	}
}

type graphvizMarshal Graphviz

// MarshalJSON implements json.Marshaler
func (g *Graphviz) MarshalJSON() ([]byte, error) {
	m := graphvizMarshal(*g)
	m.Metadata.Type = TypeGraphviz
	return json.Marshal(&m)
}
