/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import "github.com/vmware-tanzu/octant/internal/util/json"

// Link is a text component that contains a link.
//
// +octant:component
type Link struct {
	Config LinkConfig `json:"config"`
	Base
}

var _ Component = &Link{}

// LinkConfig is the contents of Link
type LinkConfig struct {
	Text string `json:"value"`
	Ref  string `json:"ref"`
	// Status sets the status of the component.
	Status       TextStatus `json:"status,omitempty" tsType:"number"`
	StatusDetail Component  `json:"statusDetail,omitempty"`
}

func (lc *LinkConfig) UnmarshalJSON(data []byte) error {
	x := struct {
		Text         string       `json:"value,omitempty"`
		Ref          string       `json:"ref,omitempty"`
		Status       TextStatus   `json:"status,omitempty"`
		StatusDetail *TypedObject `json:"statusDetail,omitempty"`
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	lc.Text = x.Text
	lc.Ref = x.Ref
	lc.Status = x.Status
	if x.StatusDetail != nil {
		sd, err := x.StatusDetail.ToComponent()
		if err != nil {
			return err
		}
		lc.StatusDetail = sd
	}

	return nil
}

type LinkOption func(l *Link)

// NewLink creates a link component
func NewLink(title, s, ref string, options ...LinkOption) *Link {
	l := &Link{
		Base: newBase(TypeLink, TitleFromString(title)),
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
	m.Metadata.Type = TypeLink
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
