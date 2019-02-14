package component

import (
	"encoding/json"
	"sort"
)

// Selector identifies a ViewComponent as being a selector flavor.
type Selector interface {
	IsSelector()
	Name() string
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

func (t *SelectorsConfig) UnmarshalJSON(data []byte) error {
	x := struct {
		Selectors []typedObject
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	for _, to := range x.Selectors {
		i, err := unmarshal(to)
		if err != nil {
			return err
		}

		if s, ok := i.(Selector); ok {
			t.Selectors = append(t.Selectors, s)
		}
	}

	return nil
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

// Add adds additional items to the tail of the selectors.
func (t *Selectors) Add(selectors ...Selector) {
	t.Config.Selectors = append(t.Config.Selectors, selectors...)
}

type selectorsMarshal Selectors

// MarshalJSON implements json.Marshaler
func (t *Selectors) MarshalJSON() ([]byte, error) {
	filtered := &Selectors{}
	for _, s := range t.Config.Selectors {
		if !isInStringSlice(s.Name(), labelsFilteredKeys) {
			filtered.Config.Selectors = append(filtered.Config.Selectors, s)
		}
	}

	filtered.Metadata.Type = "selectors"
	filtered.Metadata.Title = t.Metadata.Title

	m := selectorsMarshal(*filtered)

	sort.Slice(m.Config.Selectors, func(i, j int) bool {
		a := m.Config.Selectors[i]
		b := m.Config.Selectors[j]
		return a.Name() < b.Name()
	})

	return json.Marshal(&m)
}
