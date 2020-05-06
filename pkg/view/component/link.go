/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import "encoding/json"

// Link is a text component that contains a link.
type Link struct {
	base
	Config LinkConfig `json:"config"`
}

// LinkConfig is the contents of Link
type LinkConfig struct {
	Text string `json:"value"`
	Ref  string `json:"ref"`
	// Status sets the status of the component.
	Status       TextStatus `json:"status,omitempty"`
	StatusDetail Component  `json:"statusDetail,omitempty"`
}

type LinkOption func(l *Link)

// NewLink creates a link component
func NewLink(title, s, ref string, options ...LinkOption) *Link {
	l := &Link{
		base: newBase(typeLink, TitleFromString(title)),
		Config: LinkConfig{
			Text: s,
			Ref:  ref,
		},
	}

	for _, option := range options {
		option(l)
	}

	return l
}

// SupportsTitle designates this is a TextComponent.
func (t *Link) SupportsTitle() {}

// GetMetadata accesses the components metadata. Implements Component.
func (t *Link) GetMetadata() Metadata {
	return t.Metadata
}

// Text returns the link's text.
func (t *Link) Text() string {
	return t.Config.Text
}

// Ref returns the link's ref.
func (t *Link) Ref() string {
	return t.Config.Ref
}

// SetStatus sets the status of the text component.
func (t *Link) SetStatus(status TextStatus, detail Component) {
	t.Config.Status = status
	t.Config.StatusDetail = detail
}

type linkMarshal Link

// MarshalJSON implements json.Marshaler
func (t *Link) MarshalJSON() ([]byte, error) {
	m := linkMarshal(*t)
	m.Metadata.Type = typeLink
	return json.Marshal(&m)
}

// String returns the link's text.
func (t *Link) String() string {
	return t.Config.Text
}

// LessThan returns true if this component's value is less than the argument supplied.
func (t *Link) LessThan(i interface{}) bool {
	v, ok := i.(*Link)
	if !ok {
		return false
	}

	return t.Config.Text < v.Config.Text

}
