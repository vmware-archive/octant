/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"encoding/json"
)

// List contains other Components
type List struct {
	base
	Config ListConfig `json:"config"`
}

// ListConfig is the contents of a List
type ListConfig struct {
	Items []Component `json:"items"`
}

func (t *ListConfig) UnmarshalJSON(data []byte) error {
	x := struct {
		Items []TypedObject
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	for _, item := range x.Items {
		listItem, err := item.ToComponent()
		if err != nil {
			return err
		}
		t.Items = append(t.Items, listItem)
	}

	return nil
}

// NewList creates a list component
func NewList(title string, items []Component) *List {
	return &List{
		base: newBase(typeList, TitleFromString(title)),
		Config: ListConfig{
			Items: items,
		},
	}
}

// GetMetadata accesses the components metadata. Implements Component.
func (t *List) GetMetadata() Metadata {
	return t.Metadata
}

// Add adds additional items to the tail of the list.
func (t *List) Add(items ...Component) {
	t.Config.Items = append(t.Config.Items, items...)
}

type listMarshal List

// MarshalJSON implements json.Marshaler
func (t *List) MarshalJSON() ([]byte, error) {
	m := listMarshal(*t)
	m.Metadata.Type = typeList
	return json.Marshal(&m)
}
