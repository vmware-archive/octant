package component

import "encoding/json"

// Text is a component for freetext
type Text struct {
	base
	Config TextConfig `json:"config"`
}

// TextConfig is the contents of Text
type TextConfig struct {
	Text string `json:"value"`
}

// NewText creates a text component
func NewText(s string) *Text {
	return &Text{
		base: newBase(typeText, nil),
		Config: TextConfig{
			Text: s,
		},
	}
}

// SupportsTitle designates this is a TextComponent.
func (t *Text) SupportsTitle() {}

// GetMetadata accesses the components metadata. Implements Component.
func (t *Text) GetMetadata() Metadata {
	return t.Metadata
}

type textMarshal Text

// MarshalJSON implements json.Marshaler
func (t *Text) MarshalJSON() ([]byte, error) {
	m := textMarshal(*t)
	m.Metadata.Type = typeText
	return json.Marshal(&m)
}
