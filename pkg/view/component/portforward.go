package component

import (
	"encoding/json"
)

// PortForward is a component for freetext
type PortForward struct {
	base
	Config PortForwardConfig `json:"config"`
}

// PortForwardConfig is the contents of PortForward
type PortForwardConfig struct {
	Text   string                `json:"text"`
	ID     string                `json:"id"`     // ID of an optionally running portforward
	Action PortForwardAction     `json:"action"` // The type of action to take when interacting with the component
	Status PortForwardStatus     `json:"status"` // The status of the component (and its associated port forward)
	Ports  []PortForwardPortSpec `json:"ports"`  // If active, ports forwarded by associated port forwarder

	Target PortForwardTarget `json:"target"`
}

// PortForwardPortSpec represents a local->remote port mapping
type PortForwardPortSpec struct {
	Local  uint16 `json:"local"`
	Remote uint16 `json:"remote"`
}

type PortForwardTarget struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Namespace  string `json:"namespace"`
	Name       string `json:"name"`
}

// NewPortForwardPorts is a helper for creating a single-port portforward port spec list
func NewPortForwardPorts(local, remote uint16) []PortForwardPortSpec {
	return []PortForwardPortSpec{
		PortForwardPortSpec{
			Local:  local,
			Remote: remote,
		},
	}
}

// PortForwardAction is an enumeration of possible actions to take when interacting with the component
type PortForwardAction string

// PortForwardStatus is an enumeration of possible states the component may be in
type PortForwardStatus string

const (
	// PortForwardActionCreate indicates the component currently allows a create action
	PortForwardActionCreate PortForwardAction = "create"
	// PortForwardActionDelete indicates the component currently allows a delete action
	PortForwardActionDelete PortForwardAction = "delete"

	// PortForwardStatusInitial indicates the component is in the initial state - i.e.
	// there are no associated port forwards running
	PortForwardStatusInitial PortForwardStatus = "initial"
	// PortForwardStatusRunning indicates the component is in the running state - i.e.
	// its associated port forward is running
	PortForwardStatusRunning PortForwardStatus = "running"
)

// NewPortForward creates a portforward component
func NewPortForward(s string) *PortForward {
	return &PortForward{
		base: newBase(typePortForward, nil),
		Config: PortForwardConfig{
			Text:  s,
			Ports: []PortForwardPortSpec{},
		},
	}
}

// NewPortForwardCreator creates a portforward component representing
// a non-running port forward that can be created.
func NewPortForwardCreator(text string, ports []PortForwardPortSpec, target PortForwardTarget) *PortForward {
	pf := NewPortForward(text)
	pf.Config.Ports = make([]PortForwardPortSpec, len(ports))
	copy(pf.Config.Ports, ports)
	pf.Config.Status = PortForwardStatusInitial
	pf.Config.Action = PortForwardActionCreate
	pf.Config.Target = target
	return pf
}

// NewPortForwardDeleter creates a portforward component representing
// a running port forward that can be deleted.
func NewPortForwardDeleter(text, id string, ports []PortForwardPortSpec) *PortForward {
	pf := NewPortForward(text)
	pf.Config.ID = id
	pf.Config.Ports = make([]PortForwardPortSpec, len(ports))
	copy(pf.Config.Ports, ports)
	pf.Config.Status = PortForwardStatusRunning
	pf.Config.Action = PortForwardActionDelete
	return pf
}

// SupportsTitle designates this is a TextComponent.
func (t *PortForward) SupportsTitle() {}

// GetMetadata accesses the components metadata. Implements Component.
func (t *PortForward) GetMetadata() Metadata {
	return t.Metadata
}

type portForwardMarshal PortForward

// MarshalJSON implements json.Marshaler
func (t *PortForward) MarshalJSON() ([]byte, error) {
	m := portForwardMarshal(*t)
	m.Metadata.Type = typePortForward
	return json.Marshal(&m)
}
