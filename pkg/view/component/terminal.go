/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"encoding/json"
	"fmt"
)

// TerminalConfig holds a terminal config.
type TerminalConfig struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Container string `json:"container"`
	UUID      string `json:"uuid"`
}

type Terminal struct {
	base
	Config TerminalConfig `json:"config"`
}

// NewTerminal creates a terminal component.
func NewTerminal(namespace, name, container, command, uuid string) *Terminal {
	return &Terminal{
		base: newBase(typeTerminal, TitleFromString(fmt.Sprintf("%s / %s", container, command))),
		Config: TerminalConfig{
			Namespace: namespace,
			Name:      name,
			Container: container,
			UUID:      uuid,
		},
	}
}

// GetMetadata accesses the components metadata. Implements Component.
func (l *Terminal) GetMetadata() Metadata {
	return l.Metadata
}

type terminalMarshall Terminal

func (t *Terminal) MarshalJSON() ([]byte, error) {
	m := terminalMarshall(*t)
	m.Metadata.Type = typeTerminal

	return json.Marshal(&m)
}
