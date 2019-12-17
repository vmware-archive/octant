/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// EmptyContentResponse is an empty content response.
var EmptyContentResponse = ContentResponse{}

// ContentResponse is a a content response. It contains a
// title and one or more components.
type ContentResponse struct {
	Title      []TitleComponent `json:"title,omitempty"`
	Components []Component      `json:"viewComponents"`
	IconName   string           `json:"iconName,omitempty"`
	IconSource string           `json:"iconSource,omitempty"`
}

// NewContentResponse creates an instance of ContentResponse.
func NewContentResponse(title []TitleComponent) *ContentResponse {
	return &ContentResponse{
		Title: title,
	}
}

// Add adds zero or more components to a content response.
func (c *ContentResponse) Add(components ...Component) {
	c.Components = append(c.Components, components...)
}

// UnmarshalJSON unmarshals a content response from JSON.
func (c *ContentResponse) UnmarshalJSON(data []byte) error {
	stage := struct {
		Title      []TypedObject `json:"title,omitempty"`
		Components []TypedObject `json:"viewComponents,omitempty"`
	}{}

	if err := json.Unmarshal(data, &stage); err != nil {
		return err
	}

	for _, t := range stage.Title {
		title, err := getTitleByUnmarshalInterface(t.Config)
		if err != nil {
			return err
		}

		c.Title = Title(NewText(title))
	}

	for _, to := range stage.Components {
		vc, err := to.ToComponent()
		if err != nil {
			return err
		}

		c.Components = append(c.Components, vc)
	}

	return nil
}

func getTitleByUnmarshalInterface(config json.RawMessage) (string, error) {
	var objmap map[string]interface{}
	if err := json.Unmarshal(config, &objmap); err != nil {
		return "", err
	}

	if value, ok := objmap["value"].(string); ok {
		return value, nil
	}

	return "", fmt.Errorf("title does not have a value")
}

type TypedObject struct {
	Config   json.RawMessage `json:"config,omitempty"`
	Metadata Metadata        `json:"metadata,omitempty"`
}

func (to *TypedObject) ToComponent() (Component, error) {
	o, err := unmarshal(*to)
	if err != nil {
		return nil, err
	}

	vc, ok := o.(Component)
	if !ok {
		return nil, errors.Errorf("unable to convert %T to Component",
			o)
	}

	return vc, nil
}

// Metadata collects common fields describing Components
type Metadata struct {
	Type     string           `json:"type"`
	Title    []TitleComponent `json:"title,omitempty"`
	Accessor string           `json:"accessor,omitempty"`
}

// SetTitleText sets the title using text components.
func (m *Metadata) SetTitleText(parts ...string) {
	var titleComponents []TitleComponent

	for _, part := range parts {
		titleComponents = append(titleComponents, NewText(part))
	}

	m.Title = titleComponents
}

func (m *Metadata) UnmarshalJSON(data []byte) error {
	x := struct {
		Type     string        `json:"type,omitempty"`
		Title    []TypedObject `json:"title,omitempty"`
		Accessor string        `json:"accessor,omitempty"`
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	m.Type = x.Type
	m.Accessor = x.Accessor

	for _, title := range x.Title {
		vc, err := title.ToComponent()
		if err != nil {
			return errors.Wrap(err, "unmarshal-ing title")
		}

		tvc, ok := vc.(TitleComponent)
		if !ok {
			return errors.New("component in title isn't a title view component")
		}

		m.Title = append(m.Title, tvc)
	}

	return nil
}

// Component is a common interface for the data representation
// of visual components as rendered by the UI.
type Component interface {
	json.Marshaler

	// GetMetadata returns metadata for the component.
	GetMetadata() Metadata
	// SetAccessor sets the accessor for the component.
	SetAccessor(string)
	// IsEmpty returns true if the component is "empty".
	IsEmpty() bool
	// String returns a string representation of the component.
	String() string
	// LessThan returns true if the components value is less than the other value.
	LessThan(other interface{}) bool
}

// TitleComponent is a view component that can be used for a title.
type TitleComponent interface {
	Component

	SupportsTitle()
}

// Title is a convenience method for creating a title.
func Title(components ...TitleComponent) []TitleComponent {
	return components
}

// TitleFromString is a convenience method for create a title from a string.
func TitleFromString(s string) []TitleComponent {
	return Title(NewText(s))
}

// TitleFromTitleComponent gets a title from a TitleComponent
func TitleFromTitleComponent(tc []TitleComponent) (string, error) {
	if len(tc) != 1 {
		return "", errors.New("exactly one title component can be converted")
	}
	return tc[0].String(), nil
}
