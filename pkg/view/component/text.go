/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import "encoding/json"

// TextStatus is the status of a text component
type TextStatus int

const (
	// TextStatusOK
	TextStatusOK TextStatus = 1
	// TextStatusWarning
	TextStatusWarning TextStatus = 2
	// TextStatusError
	TextStatusError TextStatus = 3
)

// Text is a component for text
type Text struct {
	base
	Config TextConfig `json:"config"`
}

// TextConfig is the contents of Text
type TextConfig struct {
	// Text is the text that will be displayed.
	Text string `json:"value"`
	// IsMarkdown sets if the component has markdown text.
	IsMarkdown bool `json:"isMarkdown,omitempty"`
	// Status sets the status of the component.
	Status TextStatus `json:"status,omitempty"`
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

// NewMarkdownText creates a text component styled with markdown.
func NewMarkdownText(s string) *Text {
	t := NewText(s)
	t.Config.IsMarkdown = true

	return t
}

// IsMarkdown returns if this component is markdown.
func (t *Text) IsMarkdown() bool {
	return t.Config.IsMarkdown
}

// EnableMarkdown enables markdown for this text component.
func (t *Text) EnableMarkdown() {
	t.Config.IsMarkdown = true
}

// DisableMarkdown disables markdown for this text component.
func (t *Text) DisableMarkdown() {
	t.Config.IsMarkdown = false
}

// SetStatus sets the status of the text component.
func (t *Text) SetStatus(status TextStatus) {
	t.Config.Status = status
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
