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
	Metadata Metadata   `json:"metadata,omitempty"`
	Config   LogsConfig `json:"config,omitempty"`
}

func NewLogs(namespace, name string, containers []string) *Logs {
	return &Logs{
		Config: LogsConfig{
			Namespace:  namespace,
			Name:       name,
			Containers: containers,
		},
		Metadata: Metadata{
			Type: "logs",
			Title: []TitleViewComponent{
				NewText("Logs"),
			},
		},
	}
}

// GetMetadata accesses the components metadata. Implements ViewComponent.
func (l *Logs) GetMetadata() Metadata {
	return l.Metadata
}

// IsEmpty specifies whether the component is considered empty. Implements ViewComponent.
func (l *Logs) IsEmpty() bool {
	return len(l.Config.Containers) == 0
}

type logsMarshal Logs

func (l *Logs) MarshalJSON() ([]byte, error) {
	m := logsMarshal(*l)
	m.Metadata.Type = "logs"

	return json.Marshal(&m)
}
