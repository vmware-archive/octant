/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import "encoding/json"

// Stepper component implements json.Marshaler
type Stepper struct {
	base
	Config StepperConfig `json:"config"`
}

type stepperMarshal Stepper

func (t *Stepper) MarshalJSON() ([]byte, error) {
	m := stepperMarshal(*t)
	m.Metadata.Type = typeStepper
	return json.Marshal(&m)
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
