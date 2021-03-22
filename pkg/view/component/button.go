package component

import (
	"github.com/vmware-tanzu/octant/internal/util/json"

	"github.com/pkg/errors"

	"github.com/vmware-tanzu/octant/pkg/action"
)

const (
	// ButtonStatusSuccess is a green button
	ButtonStatusSuccess ButtonStatus = "success"
	// ButtonStatusInfo is the default blue color
	ButtonStatusInfo ButtonStatus = "info"
	// ButtonStatusDanger is a red button
	ButtonStatusDanger ButtonStatus = "danger"
	// ButtonStatusDisabled is a disabled button
	ButtonStatusDisabled ButtonStatus = "disabled"
)

type ButtonStatus string

const (
	// ButtonSizeBlock is a full width button
	ButtonSizeBlock ButtonSize = "block"
	// ButtonSizeSmall is a small button
	ButtonSizeLarge ButtonSize = "lg"
)

type ButtonSize string

const (
	// ButtonStyleOutline is a transparent button with a colored border
	ButtonStyleOutline ButtonStyle = "outline"
	// ButtonStyleSolid is a button with solid color
	ButtonStyleSolid ButtonStyle = "solid"
	// ButtonStyleFlat is a button with no background color or outline
	ButtonStyleFlat ButtonStyle = "link"
)

type ButtonStyle string

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

// WithLink configures a button with an href
func WithButtonLink(ref string) ButtonOption {
	return func(button *Button) {
		button.Ref = ref
	}
}

// WithButtonStatus configures the button color
func WithButtonStatus(status ButtonStatus) ButtonOption {
	return func(button *Button) {
		button.Status = status
	}
}

// WithButtonSize configures the button size
func WithButtonSize(size ButtonSize) ButtonOption {
	return func(button *Button) {
		button.Size = size
	}
}

// WithButtonStyle configures the button appearance
func WithButtonStyle(style ButtonStyle) ButtonOption {
	return func(button *Button) {
		button.Style = style
	}
}

func (b *Button) UnmarshalJSON(data []byte) error {
	x := struct {
		Name         string         `json:"name"`
		Payload      action.Payload `json:"payload"`
		Confirmation *Confirmation  `json:"confirmation,omitempty"`
		Modal        *TypedObject   `json:"modal,omitempty"`
		Ref          string         `json:"ref,omitempty"`
		Status       ButtonStatus   `json:"status,omitempty"`
		Size         ButtonSize     `json:"size,omitempty"`
		Style        ButtonStyle    `json:"style,omitempty"`
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
	b.Ref = x.Ref
	b.Status = x.Status
	b.Size = x.Size
	b.Style = x.Style

	return nil
}

// Button is a button in a group.
type Button struct {
	Name         string         `json:"name"`
	Payload      action.Payload `json:"payload"`
	Confirmation *Confirmation  `json:"confirmation,omitempty"`
	Modal        Component      `json:"modal,omitempty"`
	Ref          string         `json:"ref,omitempty"`
	Status       ButtonStatus   `json:"status,omitempty"`
	Size         ButtonSize     `json:"size,omitempty"`
	Style        ButtonStyle    `json:"style,omitempty"`
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

// NewButtonLink creates a button component that navigates to a link
func NewButtonLink(name, ref string, options ...ButtonOption) Button {
	button := Button{
		Name: name,
		Ref:  ref,
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
