/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"encoding/json"
	"fmt"
	"time"
)

type TerminalDetails struct {
	Container string    `json:"container"`
	Command   string    `json:"command"`
	UUID      string    `json:"uuid"`
	CreatedAt time.Time `json:"createdAt"`
	Active    bool      `json:"active"`
}

// TerminalConfig holds a terminal config.
type TerminalConfig struct {
	Namespace string          `json:"namespace"`
	Name      string          `json:"name"`
	Details   TerminalDetails `json:"terminal"`
}

type Terminal struct {
	base
	Config TerminalConfig `json:"config"`
}

// NewTerminal creates a terminal component.
func NewTerminal(namespace, name string, details TerminalDetails) *Terminal {
	return &Terminal{
		base: newBase(typeTerminal, TitleFromString(fmt.Sprintf("%s / %s", namespace, name))),
		Config: TerminalConfig{
			Namespace: namespace,
			Name:      name,
			Details:   details,
		},
	}
}

// GetMetadata accesses the components metadata. Implements Component.
func (t *Terminal) GetMetadata() Metadata {
	return t.Metadata
}

type terminalMarshall Terminal

func (t *Terminal) MarshalJSON() ([]byte, error) {
	m := terminalMarshall(*t)
	m.Metadata.Type = typeTerminal

	return json.Marshal(&m)
}
