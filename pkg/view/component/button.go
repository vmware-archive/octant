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
}

// MarshalJSON implements json.Marshaler
func (b *Button) MarshalJSON() ([]byte, error) {
	m := buttonMarshal(*b)
	m.Metadata.Type = TypeButton
	return json.Marshal(&m)
}
