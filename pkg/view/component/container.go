/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import "encoding/json"

// Containers is a component wrapping multiple docker container definitions
type Containers struct {
	base
	Config ContainersConfig `json:"config"`
}

// ContainersConfig is the contents of a Containers wrapper
type ContainersConfig struct {
	Containers []ContainerDef `json:"containers"`
}

// ContainerDef defines an individual docker container
type ContainerDef struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

// NewContainers creates a containers component
func NewContainers() *Containers {
	return &Containers{
		base:   newBase("containers", nil),
		Config: ContainersConfig{},
	}
}

// GetMetadata accesses the components metadata. Implements Component.
func (t *Containers) GetMetadata() Metadata {
	return t.Metadata
}

// Add adds additional items to the tail of the containers.
func (t *Containers) Add(name string, image string) {
	t.Config.Containers = append(t.Config.Containers, ContainerDef{Name: name, Image: image})
}

type containersMarshal Containers

// MarshalJSON implements json.Marshaler
func (t *Containers) MarshalJSON() ([]byte, error) {
	m := containersMarshal(*t)
	m.Metadata.Type = "containers"
	return json.Marshal(&m)
}
