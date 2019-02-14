package component

import "encoding/json"

// LabelSelectorConfig is the contents of LabelSelector
type LabelSelectorConfig struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// NewLabelSelector creates a labelSelector component
func NewLabelSelector(k, v string) *LabelSelector {
	return &LabelSelector{
		Metadata: Metadata{
			Type: "labelSelector",
		},
		Config: LabelSelectorConfig{
			Key:   k,
			Value: v,
		},
	}
}

// LabelSelector is a component for a single label within a selector
type LabelSelector struct {
	Metadata Metadata            `json:"metadata"`
	Config   LabelSelectorConfig `json:"config"`
}

// Name is the name of the LabelSelector.
func (t *LabelSelector) Name() string {
	return t.Config.Key
}

// GetMetadata accesses the components metadata. Implements ViewComponent.
func (t *LabelSelector) GetMetadata() Metadata {
	return t.Metadata
}

// IsSelector marks the component as selector flavor. Implements Selector.
func (t *LabelSelector) IsSelector() {
}

type labelSelectorMarshal LabelSelector

// MarshalJSON implements json.Marshaler
func (t *LabelSelector) MarshalJSON() ([]byte, error) {
	m := labelSelectorMarshal(*t)
	m.Metadata.Type = "labelSelector"
	return json.Marshal(&m)
}
