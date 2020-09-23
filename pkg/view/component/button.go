package component

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/vmware-tanzu/octant/pkg/action"
)

// Confirmation is configuration for a confirmation dialog.
type Confirmation struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// ButtonOption is a function for configuring a Button.
type ButtonOption func(button *Button)

// WithButtonConfirmation configured a button with a confirmation.
func WithButtonConfirmation(title, body string) ButtonOption {
	return func(button *Button) {
		button.Confirmation = &Confirmation{
			Title: title,
			Body:  body,
		}
	}
}

// WithModal configures a button to open a modal
func WithModal(modal *Modal) ButtonOption {
	return func(button *Button) {
		button.Modal = modal
	}
}

func (b *Button) UnmarshalJSON(data []byte) error {
	x := struct {
		Name         string         `json:"name"`
		Payload      action.Payload `json:"payload"`
		Confirmation *Confirmation  `json:"confirmation,omitempty"`
		Modal        *TypedObject   `json:"modal,omitempty"`
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	if x.Modal != nil {
		component, err := x.Modal.ToComponent()
		if err != nil {
			return err
		}

		modal, ok := component.(*Modal)
		if !ok {
			return errors.New("item was not a modal")
		}
		b.Modal = modal
	}

	b.Name = x.Name
	b.Payload = x.Payload
	b.Confirmation = x.Confirmation

	return nil
}

// Button is a button in a group.
type Button struct {
	Name         string         `json:"name"`
	Payload      action.Payload `json:"payload"`
	Confirmation *Confirmation  `json:"confirmation,omitempty"`
	Modal        Component      `json:"modal,omitempty"`
}

// NewButton creates an instance of Button.
func NewButton(name string, payload action.Payload, options ...ButtonOption) Button {
	button := Button{
		Name:    name,
		Payload: payload,
	}

	for _, option := range options {
		option(&button)
	}

	return button
}

// ButtonGroupConfig is configuration for a button group.
type ButtonGroupConfig struct {
	// Buttons are buttons in the group.
	Buttons []Button `json:"buttons"`
}

// ButtonGroup is a group of buttons.
//
// +octant:component
type ButtonGroup struct {
	Base
	Config ButtonGroupConfig `json:"config"`
}

// NewButtonGroup creates an instance of ButtonGroup.
func NewButtonGroup() *ButtonGroup {
	return &ButtonGroup{
		Base: newBase(TypeButtonGroup, nil),
	}
}

// AddButton adds a button to the ButtonGroup.
func (bg *ButtonGroup) AddButton(button Button) {
	bg.Config.Buttons = append(bg.Config.Buttons, button)
}

type buttonGroupMarshal ButtonGroup

// MarshalJSON marshals a button group.
func (bg *ButtonGroup) MarshalJSON() ([]byte, error) {
	m := buttonGroupMarshal(*bg)
	m.Metadata.Type = TypeButtonGroup
	return json.Marshal(&m)
}
