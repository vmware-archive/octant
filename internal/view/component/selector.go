package component

import "encoding/json"

// Selector identifies a ViewComponent as being a selector flavor.
type Selector interface {
	IsSelector()
}

// Selectors contains other ViewComponents
type Selectors struct {
	Metadata Metadata        `json:"metadata"`
	Config   SelectorsConfig `json:"config"`
}

// SelectorsConfig is the contents of a Selectors
type SelectorsConfig struct {
	Selectors []Selector `json:"selectors"`
}

// NewSelectors creates a selectors component
func NewSelectors(selectors []Selector) *Selectors {
	return &Selectors{
		Metadata: Metadata{
			Type: "selectors",
		},
		Config: SelectorsConfig{
			Selectors: selectors,
		},
	}
}

// GetMetadata accesses the components metadata. Implements ViewComponent.
func (t *Selectors) GetMetadata() Metadata {
	return t.Metadata
}

// IsEmpty specifes whether the component is considered empty. Implements ViewComponent.
func (t *Selectors) IsEmpty() bool {
	return len(t.Config.Selectors) == 0
}

// Add adds additional items to the tail of the selectors.
func (t *Selectors) Add(selectors ...Selector) {
	t.Config.Selectors = append(t.Config.Selectors, selectors...)
}

type selectorsMarshal Selectors

// MarshalJSON implements json.Marshaler
func (t *Selectors) MarshalJSON() ([]byte, error) {
	m := selectorsMarshal(*t)
	m.Metadata.Type = "selectors"
	return json.Marshal(&m)
}
