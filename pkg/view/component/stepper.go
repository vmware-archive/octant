/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import "encoding/json"

// Stepper component implements json.Marshaler
//
// +octant:component
type Stepper struct {
	Base
	Config StepperConfig `json:"config"`
}

// NewStepper creates a stepper component
func NewStepper(title string, actionName string, steps ...StepConfig) *Stepper {
	s := append([]StepConfig(nil), steps...)
	return &Stepper{
		Base: newBase(TypeStepper, TitleFromString(title)),
		Config: StepperConfig{
			Steps:  s,
			Action: actionName,
		},
	}
}

type stepperMarshal Stepper

func (t *Stepper) MarshalJSON() ([]byte, error) {
	m := stepperMarshal(*t)
	m.Metadata.Type = TypeStepper
	return json.Marshal(&m)
}

// AddStep adds a step to a stepper
func (t *Stepper) AddStep(name string, form Form, title string, description string) {
	step := StepConfig{
		Name:        name,
		Form:        form,
		Title:       title,
		Description: description,
	}
	t.Config.Steps = append(t.Config.Steps, step)
}

type StepperConfig struct {
	Action string       `json:"action"`
	Steps  []StepConfig `json:"steps"`
}

func (t *StepperConfig) UnmarshalJSON(data []byte) error {
	x := struct {
		Action string
		Steps  []StepConfig
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	t.Action = x.Action
	t.Steps = x.Steps

	return nil
}

type StepConfig struct {
	Name        string `json:"name"`
	Form        Form   `json:"form"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
