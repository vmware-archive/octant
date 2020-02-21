/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	corev1 "k8s.io/api/core/v1"
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

func NewLogs(pod *corev1.Pod) *Logs {

	var containerNames []string

	for _, c := range pod.Spec.InitContainers {
		containerNames = append(containerNames, c.Name)
	}

	for _, c := range pod.Spec.Containers {
		containerNames = append(containerNames, c.Name)
	}

	return &Logs{
		Config: LogsConfig{
			Namespace:  pod.Namespace,
			Name:       pod.Name,
			Containers: containerNames,
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
