package component

import "encoding/json"

// Link is a text component that contains a link.
type Link struct {
	Metadata Metadata   `json:"metadata"`
	Config   LinkConfig `json:"config"`
}

// LinkConfig is the contents of Link
type LinkConfig struct {
	Text string `json:"value"`
	Ref  string `json:"ref"`
}

// NewLink creates a link component
func NewLink(title, s, ref string) *Link {
	return &Link{
		Metadata: Metadata{
			Type:  "link",
			Title: []TitleViewComponent{NewText(title)},
		},
		Config: LinkConfig{
			Text: s,
			Ref:  ref,
		},
	}
}

// GetMetadata accesses the components metadata. Implements ViewComponent.
func (t *Link) GetMetadata() Metadata {
	return t.Metadata
}

type linkMarshal Link

// MarshalJSON implements json.Marshaler
func (t *Link) MarshalJSON() ([]byte, error) {
	m := linkMarshal(*t)
	m.Metadata.Type = "link"
	return json.Marshal(&m)
}
