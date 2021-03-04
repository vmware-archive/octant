/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

// Loading is a component for a spinner
//
// +octant:component
type Loading struct {
	Base
	Config LoadingConfig `json:"config"`
}

// LoadingConfig is the contents of Loading
type LoadingConfig struct {
	Text string `json:"value"`
}

// NewLoading creates a loading component
func NewLoading(title []TitleComponent, message string) *Loading {
	return &Loading{
		Base: newBase(TypeLoading, title),
		Config: LoadingConfig{
			Text: message,
		},
	}
}

// SupportsTitle denotes this is a LoadingComponent.
func (t *Loading) SupportsTitle() {}

type loadingMarshal Loading

// MarshalJSON implements json.Marshaler
func (t *Loading) MarshalJSON() ([]byte, error) {
	m := loadingMarshal(*t)
	m.Metadata.Type = TypeLoading
	return json.Marshal(&m)
}

// String returns the text content of the component.
func (t *Loading) String() string {
	return t.Config.Text
}
