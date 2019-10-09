/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"encoding/json"
	"sort"
)

// Selector identifies a Component as being a selector flavor.
type Selector interface {
	IsSelector()
	Name() string
}

// SelectorsConfig is the contents of a Selectors
type SelectorsConfig struct {
	Selectors []Selector `json:"selectors"`
}

func (t *SelectorsConfig) UnmarshalJSON(data []byte) error {
	x := struct {
		Selectors []TypedObject
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	for _, to := range x.Selectors {
		i, err := unmarshal(to)
		if err != nil {
			return err
		}

		if s, ok := i.(Selector); ok {
			t.Selectors = append(t.Selectors, s)
		}
	}

	return nil
}

// Selectors contains other Components
type Selectors struct {
	base
	Config SelectorsConfig `json:"config"`
}

// NewSelectors creates a selectors component
func NewSelectors(selectors []Selector) *Selectors {
	return &Selectors{
		base: newBase(typeSelectors, nil),
		Config: SelectorsConfig{
			Selectors: selectors,
		},
	}
}

// GetMetadata accesses the components metadata. Implements Component.
func (t *Selectors) GetMetadata() Metadata {
	return t.Metadata
}

// Add adds additional items to the tail of the selectors.
func (t *Selectors) Add(selectors ...Selector) {
	t.Config.Selectors = append(t.Config.Selectors, selectors...)
}

type selectorsMarshal Selectors

// MarshalJSON implements json.Marshaler
func (t *Selectors) MarshalJSON() ([]byte, error) {
	filtered := &Selectors{}
	for _, s := range t.Config.Selectors {
		if !isInStringSlice(s.Name(), labelsFilteredKeys) {
			filtered.Config.Selectors = append(filtered.Config.Selectors, s)
		}
	}

	filtered.Metadata.Type = typeSelectors
	filtered.Metadata.Title = t.Metadata.Title

	m := selectorsMarshal(*filtered)

	sort.Slice(m.Config.Selectors, func(i, j int) bool {
		a := m.Config.Selectors[i]
		b := m.Config.Selectors[j]
		return a.Name() < b.Name()
	})

	return json.Marshal(&m)
}
