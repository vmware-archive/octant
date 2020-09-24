/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"encoding/json"
)

var defaultDurations = []Since{
	{"5 minutes", 300},
	{"10 minutes", 600},
	{"30 minutes", 1800},
	{"1 hour", 3600},
	{"3 hours", 10800},
	{"5 hours", 18000},
	{"Creation", -1},
}

type Since struct {
	Label   string `json:"label,omitempty"`
	Seconds int64  `json:"seconds,omitempty"`
}

type LogsConfig struct {
	Namespace  string   `json:"namespace,omitempty"`
	Name       string   `json:"name,omitempty"`
	Containers []string `json:"containers,omitempty"`
	Durations  []Since  `json:"durations,omitempty"`
}

// Logs is a logs component.
//
// +octant:component
type Logs struct {
	Base
	Config LogsConfig `json:"config,omitempty"`
}

func NewLogs(namespace, name string, containers ...string) *Logs {
	return &Logs{
		Config: LogsConfig{
			Namespace:  namespace,
			Name:       name,
			Containers: containers,
			Durations:  defaultDurations,
		},
		Base: newBase(TypeLogs, TitleFromString("Logs")),
	}
}

// GetMetadata accesses the components metadata. Implements Component.
func (l *Logs) GetMetadata() Metadata {
	return l.Metadata
}

type logsMarshal Logs

func (l *Logs) MarshalJSON() ([]byte, error) {
	m := logsMarshal(*l)
	m.Metadata.Type = TypeLogs

	return json.Marshal(&m)
}
