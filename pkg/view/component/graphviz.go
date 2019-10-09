/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import "encoding/json"

// GraphvizConfig is the contents of Graphviz.
type GraphvizConfig struct {
	DOT string `json:"dot,omitempty"`
}

// Graphviz is a component for displaying graphviz diagrams.
type Graphviz struct {
	base
	Config GraphvizConfig `json:"config"`
}

func NewGraphviz(dot string) *Graphviz {
	return &Graphviz{
		base: newBase(typeGraphviz, nil),
		Config: GraphvizConfig{
			DOT: dot,
		},
	}
}

type graphvizMarshal Graphviz

// MarshalJSON implements json.Marshaler
func (g *Graphviz) MarshalJSON() ([]byte, error) {
	m := graphvizMarshal(*g)
	m.Metadata.Type = typeGraphviz
	return json.Marshal(&m)
}
