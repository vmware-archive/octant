package component

import "encoding/json"

// Labels is a component representing key/value based labels
type Labels struct {
	Metadata Metadata     `json:"metadata"`
	Config   LabelsConfig `json:"config"`
}

// LabelsConfig is the contents of Labels
type LabelsConfig struct {
	Labels map[string]string `json:"labels"`
}

// NewLabels creates a labels component
func NewLabels(labels map[string]string) *Labels {
	return &Labels{
		Metadata: Metadata{
			Type: "labels",
		},
		Config: LabelsConfig{
			Labels: labels,
		},
	}
}

// GetMetadata accesses the components metadata. Implements ViewComponent.
func (t *Labels) GetMetadata() Metadata {
	return t.Metadata
}

// IsEmpty specifes whether the component is considered empty. Implements ViewComponent.
func (t *Labels) IsEmpty() bool {
	return len(t.Config.Labels) == 0
}

type labelsMarshal Labels

// MarshalJSON implements json.Marshaler
func (t *Labels) MarshalJSON() ([]byte, error) {
	m := labelsMarshal(*t)
	m.Metadata.Type = "labels"
	return json.Marshal(&m)
}
