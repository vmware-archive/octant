package component

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// ContentResponse is a a content response. It contains a
// title and one or more components.
type ContentResponse struct {
	Title          []TitleViewComponent `json:"title,omitempty"`
	ViewComponents []ViewComponent      `json:"viewComponents"`
}

// NewContentResponse creates an instance of ContentResponse.
func NewContentResponse(title []TitleViewComponent) *ContentResponse {
	return &ContentResponse{
		Title: title,
	}
}

// Add adds zero or more components to a content response.
func (c *ContentResponse) Add(components ...ViewComponent) {
	c.ViewComponents = append(c.ViewComponents, components...)
}

// UnmarshalJSON unarmshals a content response from JSON.
func (c *ContentResponse) UnmarshalJSON(data []byte) error {
	stage := struct {
		Title          string        `json:"title,omitempty"`
		ViewComponents []typedObject `json:"viewComponents,omitempty"`
	}{}

	if err := json.Unmarshal(data, &stage); err != nil {
		return err
	}

	c.Title = []TitleViewComponent{
		NewText(stage.Title),
	}

	for _, to := range stage.ViewComponents {
		vc, err := to.ToViewComponent()
		if err != nil {
			return err
		}

		c.ViewComponents = append(c.ViewComponents, vc)
	}

	return nil
}

type typedObject struct {
	Config   json.RawMessage `json:"config,omitempty"`
	Metadata Metadata        `json:"metadata,omitempty"`
}

func (to *typedObject) ToViewComponent() (ViewComponent, error) {
	o, err := unmarshal(*to)
	if err != nil {
		return nil, err
	}

	vc, ok := o.(ViewComponent)
	if !ok {
		return nil, errors.Errorf("unable to convert %T to ViewComponent",
			o)
	}

	return vc, nil
}

// Metadata collects common fields describing ViewComponents
type Metadata struct {
	Type  string               `json:"type"`
	Title []TitleViewComponent `json:"title,omitempty"`
}

// SetTitleText sets the title using text components.
func (m *Metadata) SetTitleText(parts ...string) {
	var titleViewComponents []TitleViewComponent

	for _, part := range parts {
		titleViewComponents = append(titleViewComponents, NewText(part))
	}

	m.Title = titleViewComponents
}

func (m *Metadata) UnmarshalJSON(data []byte) error {
	x := struct {
		Type  string        `json:"type,omitempty"`
		Title []typedObject `json:"title,omitempty"`
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	m.Type = x.Type

	for _, title := range x.Title {
		vc, err := title.ToViewComponent()
		if err != nil {
			return errors.Wrap(err, "unmarshaling title")
		}

		tvc, ok := vc.(TitleViewComponent)
		if !ok {
			return errors.New("component in title isn't a title view component")
		}

		m.Title = append(m.Title, tvc)
	}

	return nil
}

// ViewComponent is a common interface for the data representation
// of visual components as rendered by the UI.
type ViewComponent interface {
	IsEmpty() bool
	GetMetadata() Metadata
}

// TitleViewComponent is a view component that can be used for a title.
type TitleViewComponent interface {
	ViewComponent

	SupportsTitle()
}
