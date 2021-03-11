/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import "github.com/vmware-tanzu/octant/internal/util/json"

// LabelSelectorConfig is the contents of LabelSelector
type LabelSelectorConfig struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// LabelSelector is a component for a single label within a selector
//
// +octant:component
type LabelSelector struct {
	Base
	Config LabelSelectorConfig `json:"config"`
}

// NewLabelSelector creates a labelSelector component
func NewLabelSelector(k, v string) *LabelSelector {
	return &LabelSelector{
		Base: newBase(TypeLabelSelector, nil),
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
	m.Metadata.Type = TypeLabelSelector
	return json.Marshal(&m)
}
