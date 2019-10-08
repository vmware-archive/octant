/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"encoding/json"
)

type LogsConfig struct {
	Namespace  string   `json:"namespace,omitempty"`
	Name       string   `json:"name,omitempty"`
	Containers []string `json:"containers,omitempty"`
}

type Logs struct {
	base
	Config LogsConfig `json:"config,omitempty"`
}

func NewLogs(namespace, name string, containers []string) *Logs {
	return &Logs{
		Config: LogsConfig{
			Namespace:  namespace,
			Name:       name,
			Containers: containers,
		},
		base: newBase(typeLogs, TitleFromString("Logs")),
	}
}

// GetMetadata accesses the components metadata. Implements Component.
func (l *Logs) GetMetadata() Metadata {
	return l.Metadata
}

type logsMarshal Logs

func (l *Logs) MarshalJSON() ([]byte, error) {
	m := logsMarshal(*l)
	m.Metadata.Type = typeLogs

	return json.Marshal(&m)
}
