/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package v1alpha1

// OperationType is an operation type. It is used in Condition to determine if an element
// is enabled.
type OperationType string

const (
	// OperationTypeString is an equal operation.
	OperationTypeString OperationType = "Equal"
)

// Condition is a condition. Conditions are used to determine if an element is enabled.
type Condition struct {
	LHS string        `json:"lhs"`
	RHS string        `json:"rhs"`
	Op  OperationType `json:"op"`
}

// Element is an element in a section.
type Element struct {
	Name              string      `json:"name"`
	Value             string      `json:"value"`
	Type              string      `json:"type"`
	DisableConditions []Condition `json:"disableConditions,omitempty"`
	Config            interface{} `json:"config"`
}

const (
	elementTypeRadio = "radio"
	elementTypeText  = "text"
)

// RadioValue is an individual radio value.
type RadioValue struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// RadioConfig is configuration for a radio element.
type RadioConfig struct {
	Values []RadioValue `json:"values"`
}

// NewRadioElement creates a radio element.
func NewRadioElement(name, value string, disableConditions []Condition, values map[string]string) Element {
	e := Element{
		Name:              name,
		Value:             value,
		Type:              elementTypeRadio,
		DisableConditions: disableConditions,
	}

	config := RadioConfig{}

	for k, v := range values {
		config.Values = append(config.Values, RadioValue{Label: k, Value: v})
	}
	e.Config = config

	return e
}

// TextConfig is configuration for a text element.
type TextConfig struct {
	Label       string `json:"label"`
	Placeholder string `json:"placeholder"`
}

// NewTextElement creates a text element.
func NewTextElement(name, value string, disableConditions []Condition, placeholder, label string) Element {
	e := Element{
		Name:              name,
		Value:             value,
		Type:              elementTypeText,
		DisableConditions: disableConditions,
		Config: TextConfig{
			Placeholder: placeholder,
			Label:       label,
		},
	}

	return e
}

// PreferenceSection is a section in a preference panel.
type PreferenceSection struct {
	Name     string    `json:"name"`
	Elements []Element `json:"elements,omitempty"`
}

// PreferencePanel is a preference panel.
type PreferencePanel struct {
	Name     string              `json:"name"`
	Sections []PreferenceSection `json:"sections,omitempty"`
}

// Preferences are preferences.
type Preferences struct {
	Version    string            `json:"version"`
	Panels     []PreferencePanel `json:"panels,omitempty"`
	UpdateName string            `json:"updateName"`
}

// Update updates preferences given a set of a values.
func (p *Preferences) Update(values map[string]string) {
	for pID := range p.Panels {
		for sID := range p.Panels[pID].Sections {
			for eID, element := range p.Panels[pID].Sections[sID].Elements {
				for k, v := range values {
					if element.Name == k {
						element.Value = v
						p.Panels[pID].Sections[sID].Elements[eID] = element
					}
				}
			}
		}
	}
}
