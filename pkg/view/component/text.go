/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import "encoding/json"

// Text is a component for freetext
type Text struct {
	base
	Config TextConfig `json:"config"`
}

// TextConfig is the contents of Text
type TextConfig struct {
	Text string `json:"value"`
}

// NewText creates a text component
func NewText(s string) *Text {
	return &Text{
		base: newBase(typeText, nil),
		Config: TextConfig{
			Text: s,
		},
	}
}

// SupportsTitle denotes this is a TextComponent.
func (t *Text) SupportsTitle() {}

type textMarshal Text

// MarshalJSON implements json.Marshaler
func (t *Text) MarshalJSON() ([]byte, error) {
	m := textMarshal(*t)
	m.Metadata.Type = typeText
	return json.Marshal(&m)
}

// String returns the text content of the component.
func (t *Text) String() string {
	return t.Config.Text
}

// LessThan returns true if this component's value is less than the argument supplied.
func (t *Text) LessThan(i interface{}) bool {
	v, ok := i.(*Text)
	if !ok {
		return false
	}

	return t.Config.Text < v.Config.Text

}
