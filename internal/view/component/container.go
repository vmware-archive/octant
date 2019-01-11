package component

import "encoding/json"

// Containers is a component wrapping multiple docker container definitions
type Containers struct {
	Metadata Metadata         `json:"metadata"`
	Config   ContainersConfig `json:"config"`
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
		Metadata: Metadata{
			Type: "containers",
		},
		Config: ContainersConfig{},
	}
}

// GetMetadata accesses the components metadata. Implements ViewComponent.
func (t *Containers) GetMetadata() Metadata {
	return t.Metadata
}

// IsEmpty specifies whether the component is considered empty. Implements ViewComponent.
func (t *Containers) IsEmpty() bool {
	return len(t.Config.Containers) == 0
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
