package component

import "encoding/json"

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
	Namespace  string           `json:"namespace,omitempty"`
	APIVersion string           `json:"apiVersion,omitempty"`
	Kind       string           `json:"kind,omitempty"`
	Name       string           `json:"name,omitempty"`
	Port       int              `json:"port,omitempty"`
	Protocol   string           `json:"protocol,omitempty"`
	State      PortForwardState `json:"state,omitempty"`
}

// NewPort creates a port component
func NewPort(namespace, apiVersion, kind, name string, port int, protocol string, pfs PortForwardState) *Port {
	return &Port{
		base: newBase(typePort, nil),
		Config: PortConfig{
			Namespace:  namespace,
			APIVersion: apiVersion,
			Kind:       kind,
			Name:       name,
			Port:       port,
			Protocol:   protocol,
			State:      pfs,
		},
	}
}

// GetMetadata accesses the components metadata. Implements ViewComponent.
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
