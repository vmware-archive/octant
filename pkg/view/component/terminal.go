/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"encoding/json"
	"time"
)

type TerminalDetails struct {
	Container string    `json:"container"`
	Command   string    `json:"command"`
	CreatedAt time.Time `json:"createdAt"`
	Active    bool      `json:"active"`
}

// TerminalConfig holds a terminal config.
type TerminalConfig struct {
	Namespace  string          `json:"namespace"`
	Name       string          `json:"name"`
	PodName    string          `json:"podName"`
	Containers []string        `json:"containers"`
	Details    TerminalDetails `json:"terminal"`
}

// Terminal is a terminal component.
//
// +octant:component
type Terminal struct {
	Base
	Config TerminalConfig `json:"config"`
}

// NewTerminal creates a Terminal component.
func NewTerminal(namespace, name string, podName string, containers []string, details TerminalDetails) *Terminal {
	return &Terminal{
		Base: newBase(TypeTerminal, TitleFromString(name)),
		Config: TerminalConfig{
			Namespace:  namespace,
			Name:       name,
			PodName:    podName,
			Containers: containers,
			Details:    details,
		},
	}
}

// GetMetadata accesses the components metadata. Implements Component.
func (t *Terminal) GetMetadata() Metadata {
	return t.Metadata
}

type terminalMarshal Terminal

func (t *Terminal) MarshalJSON() ([]byte, error) {
	m := terminalMarshal(*t)
	m.Metadata.Type = TypeTerminal

	return json.Marshal(&m)
}
