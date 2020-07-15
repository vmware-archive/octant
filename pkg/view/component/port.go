/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"encoding/json"

	"github.com/vmware-tanzu/octant/pkg/action"
)

type PortForwardState struct {
	IsForwardable bool   `json:"isForwardable,omitempty"`
	IsForwarded   bool   `json:"isForwarded,omitempty"`
	Port          int    `json:"port,omitempty"`
	ID            string `json:"id,omitempty"`
}

// Port is a component for a port
type Port struct {
	base
	Config PortConfig `json:"config"`
}

// PortConfig is the contents of Port
type PortConfig struct {
	Port           int              `json:"port,omitempty"`
	Protocol       string           `json:"protocol,omitempty"`
	TargetPort     int              `json:"targetPort,omitempty"`
	TargetPortName string           `json:"targetPortName,omitempty"`
	State          PortForwardState `json:"state,omitempty"`
	Button         *ButtonGroup     `json:"buttonGroup,omitempty"`
}

// NewPort creates a port component
func NewPort(namespace, apiVersion, kind, name string, port int, protocol string, pfs PortForwardState) *Port {
	return &Port{
		base: newBase(typePort, nil),
		Config: PortConfig{
			Port:     port,
			Protocol: protocol,
			State:    pfs,
			Button:   describeButton(namespace, apiVersion, kind, name, port, pfs),
		},
	}
}

// NewPort creates a port component
func NewServicePort(namespace, apiVersion, kind, name string, port int, protocol string, targetPort int, targetPortName string, pfs PortForwardState) *Port {
	return &Port{
		base: newBase(typePort, nil),
		Config: PortConfig{
			Port:           port,
			Protocol:       protocol,
			TargetPort:     targetPort,
			TargetPortName: targetPortName,
			State:          pfs,
			Button:         describeButton(namespace, apiVersion, kind, name, targetPort, pfs),
		},
	}
}

// GetMetadata accesses the components metadata. Implements Component.
func (t *Port) GetMetadata() Metadata {
	return t.Metadata
}

type portMarshal Port

// MarshalJSON implements json.Marshaler
func (t *Port) MarshalJSON() ([]byte, error) {
	m := portMarshal(*t)
	m.Metadata.Type = typePort
	return json.Marshal(&m)
}

type PortsConfig struct {
	Ports []Port `json:"ports,omitempty"`
}

type Ports struct {
	base
	Config PortsConfig `json:"config,omitempty"`
}

func NewPorts(ports []Port) *Ports {
	return &Ports{
		base: newBase(typePorts, nil),
		Config: PortsConfig{
			Ports: ports,
		},
	}
}

func (t *Ports) GetMetadata() Metadata {
	return t.Metadata
}

type portsMarshal Ports

func (t *Ports) MarshalJSON() ([]byte, error) {
	m := portsMarshal(*t)
	m.Metadata.Type = typePorts
	return json.Marshal(&m)
}

func describeButton(namespace, apiVersion, kind, name string, port int, pfs PortForwardState) *ButtonGroup {
	buttonGroup := NewButtonGroup()
	var buttonText, actionName string
	var payload action.Payload

	if pfs.IsForwarded {
		buttonText = "Stop port forward"
		actionName = "overview/stopPortForward"
		payload = action.Payload{
			"id": pfs.ID,
		}
	} else {
		buttonText = "Start port forward"
		actionName = "overview/startPortForward"
		payload = action.Payload{
			"apiVersion": apiVersion,
			"kind":       kind,
			"name":       name,
			"namespace":  namespace,
			"port":       port,
		}
	}

	buttonGroup.AddButton(
		NewButton(buttonText, action.CreatePayload(actionName, payload)),
	)

	return buttonGroup
}
