package component

import "github.com/vmware-tanzu/octant/internal/util/json"

type MfComponent struct {
	Base
	Config MfComponentConfig `json:"config"`
}

type MfComponentMarshal MfComponent

func (m *MfComponent) MarshalJSON() ([]byte, error) {
	c := MfComponentMarshal(*m)
	c.Metadata.Type = TypeMFComponent
	return json.Marshal(&c)
}

type MfComponentConfig struct {
	Name          string `json:"name"`
	RemoteEntry   string `json:"remoteEntry"`
	RemoteName    string `json:"remoteName"`
	ExposedModule string `json:"exposedModule"`
	ElementName   string `json:"elementName"`
}

// NewMFComponent creates a new module federation component
func NewMFComponent(mfc MfComponentConfig, options ...func(*MfComponent)) *MfComponent {
	t := &MfComponent{
		Base:   newBase(TypeMFComponent, nil),
		Config: mfc,
	}

	for _, option := range options {
		option(t)
	}

	return t
}
