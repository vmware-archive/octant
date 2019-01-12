package component

import "encoding/json"

var labelsFilteredKeys = []string{
	"controller-revision-hash",
	"controller-uid",
	"pod-template-generation",
	"pod-template-hash",
	"statefulset.kubernetes.io/pod-name",
	"job-name",
}

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

// IsEmpty specifies whether the component is considered empty. Implements ViewComponent.
func (t *Labels) IsEmpty() bool {
	return len(t.Config.Labels) == 0
}

type labelsMarshal Labels

// MarshalJSON implements json.Marshaler. It will filter
// label keys specified in `labelsFilteredKeys`.
func (t *Labels) MarshalJSON() ([]byte, error) {
	filtered := &Labels{Config: LabelsConfig{Labels: make(map[string]string)}}
	for k, v := range t.Config.Labels {
		if !isInStringSlice(k, labelsFilteredKeys) {
			filtered.Config.Labels[k] = v
		}
	}

	m := labelsMarshal(*filtered)
	m.Metadata.Type = "labels"
	return json.Marshal(&m)
}

func isInStringSlice(s string, sl []string) bool {
	for i := range sl {
		if sl[i] == s {
			return true
		}
	}

	return false
}
