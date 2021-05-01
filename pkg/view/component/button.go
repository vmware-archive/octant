/*
Copyright (c) 2021 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"github.com/pkg/errors"

	"github.com/vmware-tanzu/octant/internal/util/json"
	"github.com/vmware-tanzu/octant/pkg/action"
)

const (
	// ButtonStatusSuccess is a green button
	ButtonStatusSuccess ButtonStatus = "success"
	// ButtonStatusInfo is the default blue color
	ButtonStatusInfo ButtonStatus = "primary"
	// ButtonStatusDanger is a red button
	ButtonStatusDanger ButtonStatus = "danger"
	// ButtonStatusDisabled is a disabled button
	ButtonStatusDisabled ButtonStatus = "disabled"
)

// ButtonStatus is the status of a Button
type ButtonStatus string

const (
	// ButtonSizeBlock is a full width button
	ButtonSizeBlock ButtonSize = "block"
	// ButtonSizeMedium is a clarity small button
	ButtonSizeMedium ButtonSize = "md"
)

// ButtonSize defines the size of a Button
type ButtonSize string

const (
	// ButtonStyleOutline is a transparent button with a colored border
	ButtonStyleOutline ButtonStyle = "outline"
	// ButtonStyleSolid is a button with solid color
	ButtonStyleSolid ButtonStyle = "solid"
	// ButtonStyleFlat is a button with no background color or outline
	ButtonStyleFlat ButtonStyle = "flat"
)

// ButtonStyle is the style of a Button
type ButtonStyle string

// WithButtonStatus configures the button color
func WithButtonStatus(status ButtonStatus) ButtonOption {
	return func(button *Button) {
		button.Config.Status = status
	}
}

// WithButtonSize configures the button size
func WithButtonSize(size ButtonSize) ButtonOption {
	return func(button *Button) {
		button.Config.Size = size
	}
}

// WithButtonStyle configures the button appearance
func WithButtonStyle(style ButtonStyle) ButtonOption {
	return func(button *Button) {
		button.Config.Style = style
	}
}

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
		button.Config.Confirmation = &Confirmation{
			Title: title,
			Body:  body,
		}
	}
}

// WithModal configures a button to open a modal
func WithModal(modal *Modal) ButtonOption {
	return func(button *Button) {
		button.Config.Modal = modal
	}
}

var _ Component = (*Button)(nil)

type buttonMarshal Button

func (b *ButtonConfig) UnmarshalJSON(data []byte) error {
	x := struct {
		Name         string         `json:"name"`
		Payload      action.Payload `json:"payload"`
		Confirmation *Confirmation  `json:"confirmation,omitempty"`
		Modal        *TypedObject   `json:"modal,omitempty"`
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
	b.Status = x.Status
	b.Size = x.Size
	b.Style = x.Style

	return nil
}

// Button is a component for a Button
// +octant:component
type Button struct {
	Base
	Config ButtonConfig `json:"config"`
}

// NewButton creates an instance of Button.
func NewButton(name string, payload action.Payload, options ...ButtonOption) *Button {
	button := Button{
		Base: newBase(TypeButton, nil),
		Config: ButtonConfig{
			Name:    name,
			Payload: payload,
		},
	}

	for _, option := range options {
		option(&button)
	}

	return &button
}

// ButtonConfig is the contents of a Button
type ButtonConfig struct {
	Name         string         `json:"name"`
	Payload      action.Payload `json:"payload"`
	Confirmation *Confirmation  `json:"confirmation,omitempty"`
	Modal        Component      `json:"modal,omitempty"`
	Status       ButtonStatus   `json:"status,omitempty"`
	Size         ButtonSize     `json:"size,omitempty"`
	Style        ButtonStyle    `json:"style,omitempty"`
}

// MarshalJSON implements json.Marshaler
func (b *Button) MarshalJSON() ([]byte, error) {
	m := buttonMarshal(*b)
	m.Metadata.Type = TypeButton
	return json.Marshal(&m)
}
