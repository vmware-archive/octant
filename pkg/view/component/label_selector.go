/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import "encoding/json"

// LabelSelectorConfig is the contents of LabelSelector
type LabelSelectorConfig struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// LabelSelector is a component for a single label within a selector
type LabelSelector struct {
	base
	Config LabelSelectorConfig `json:"config"`
}

// NewLabelSelector creates a labelSelector component
func NewLabelSelector(k, v string) *LabelSelector {
	return &LabelSelector{
		base: newBase(typeLabelSelector, nil),
		Config: LabelSelectorConfig{
			Key:   k,
			Value: v,
		},
	}
}

// Name is the name of the LabelSelector.
func (t *LabelSelector) Name() string {
	return t.Config.Key
}

// GetMetadata accesses the components metadata. Implements Component.
func (t *LabelSelector) GetMetadata() Metadata {
	return t.Metadata
}

// IsSelector marks the component as selector flavor. Implements Selector.
func (t *LabelSelector) IsSelector() {
}

type labelSelectorMarshal LabelSelector

// MarshalJSON implements json.Marshaler
func (t *LabelSelector) MarshalJSON() ([]byte, error) {
	m := labelSelectorMarshal(*t)
	m.Metadata.Type = typeLabelSelector
	return json.Marshal(&m)
}
