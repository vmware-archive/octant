package component

import (
	"encoding/json"
)

type ModalSize string

const (
	// ModalSizeSmall is the smallest modal
	ModalSizeSmall ModalSize = "sm"
	// ModalSizeLarge is a large modal
	ModalSizeLarge ModalSize = "lg"
	// ModalSizeExtraLarge is the largest modal
	ModalSizeExtraLarge ModalSize = "xl"
)

// ModalConfig is a configuration for the modal component.
type ModalConfig struct {
	Body      Component `json:"body,omitempty"`
	Form      *Form     `json:"form,omitempty"`
	Opened    bool      `json:"opened"`
	ModalSize ModalSize `json:"size,omitempty"`
	Buttons   []Button  `json:"buttons,omitempty"`
}

// UnmarshalJSON unmarshals a modal config from JSON.
func (m *ModalConfig) UnmarshalJSON(data []byte) error {
	x := struct {
		Body      *TypedObject `json:"body,omitempty"`
		Form      *Form        `json:"form,omitempty"`
		Opened    bool         `json:"opened"`
		ModalSize ModalSize    `json:"size,omitempty"`
		Buttons   []Button     `json:"buttons,omitempty"`
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	if x.Body != nil {
		var err error
		m.Body, err = x.Body.ToComponent()
		if err != nil {
			return err
		}
	}

	m.Form = x.Form
	m.Opened = x.Opened
	m.ModalSize = x.ModalSize
	m.Buttons = x.Buttons
	return nil
}

// Modal is a modal component.
//
// +octant:component
type Modal struct {
	Base
	Config ModalConfig `json:"config"`
}

// NewModal creates a new modal.
func NewModal(title []TitleComponent) *Modal {
	return &Modal{
		Base: newBase(TypeModal, title),
	}
}

var _ Component = (*Modal)(nil)

// SetBody sets the body of a modal.
func (m *Modal) SetBody(body Component) {
	m.Config.Body = body
}

// AddForm adds a form to a modal. It is added after the body.
func (m *Modal) AddForm(form Form) {
	m.Config.Form = &form
}

// SetSize sets the size of a modal. Size is medium by default.
func (m *Modal) SetSize(size ModalSize) {
	m.Config.ModalSize = size
}

// AddButton is a helper to add a custom button
func (m *Modal) AddButton(button Button) {
	m.Config.Buttons = append(m.Config.Buttons, button)
}

// Open opens a modal. A modal is closed by default.
func (m *Modal) Open() {
	m.Config.Opened = true
}

// Close closes a modal.
func (m *Modal) Close() {
	m.Config.Opened = false
}

type modalMarshal Modal

// MarshalJSON marshal a modal to JSON.
func (m *Modal) MarshalJSON() ([]byte, error) {
	k := modalMarshal(*m)
	k.Metadata.Type = TypeModal
	return json.Marshal(&k)
}
