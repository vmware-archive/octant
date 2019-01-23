package component

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type ContentResponse struct {
	Title          string          `json:"title,omitempty"`
	ViewComponents []ViewComponent `json:"viewComponents"`
}

func (c *ContentResponse) UnmarshalJSON(data []byte) error {
	stage := struct {
		Title          string        `json:"title,omitempty"`
		ViewComponents []typedObject `json:"viewComponents,omitempty"`
	}{}

	if err := json.Unmarshal(data, &stage); err != nil {
		return err
	}

	c.Title = stage.Title

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
	Type  string `json:"type"`
	Title string `json:"title,omitempty"`
}

// ViewComponent is a common interface for the data representation
// of visual components as rendered by the UI.
type ViewComponent interface {
	IsEmpty() bool
	GetMetadata() Metadata
}
