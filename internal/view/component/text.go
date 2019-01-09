package component

import "encoding/json"

// Text is a component for freetext
type Text struct {
	Metadata Metadata   `json:"metadata"`
	Config   TextConfig `json:"config"`
}

// TextConfig is the contents of Text
type TextConfig struct {
	Text string `json:"value"`
}

// NewText creates a text component
func NewText(s string) *Text {
	return &Text{
		Metadata: Metadata{
			Type: "text",
		},
		Config: TextConfig{
			Text: s,
		},
	}
}

// GetMetadata accesses the components metadata. Implements ViewComponent.
func (t *Text) GetMetadata() Metadata {
	return t.Metadata
}

// IsEmpty specifes whether the component is considered empty. Implements ViewComponent.
func (t *Text) IsEmpty() bool {
	return t.Config.Text == ""
}

type textMarshal Text

// MarshalJSON implements json.Marshaler
func (t *Text) MarshalJSON() ([]byte, error) {
	m := textMarshal(*t)
	m.Metadata.Type = "text"
	return json.Marshal(&m)
}
