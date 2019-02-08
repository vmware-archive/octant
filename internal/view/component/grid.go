package component

import "encoding/json"

// Grid contains other ViewComponents
type Grid struct {
	Metadata Metadata   `json:"metadata"`
	Config   GridConfig `json:"config"`
}

// GridConfig is the contents of a Grid
type GridConfig struct {
	Panels []Panel `json:"panels"`
}

// NewGrid creates a grid component
func NewGrid(title string, panels ...Panel) *Grid {
	p := append([]Panel(nil), panels...) // Make a copy
	return &Grid{
		Metadata: Metadata{
			Type:  "grid",
			Title: []TitleViewComponent{NewText(title)},
		},
		Config: GridConfig{
			Panels: p,
		},
	}
}

// GetMetadata accesses the components metadata. Implements ViewComponent.
func (t *Grid) GetMetadata() Metadata {
	return t.Metadata
}

// IsEmpty specifies whether the component is considered empty. Implements ViewComponent.
func (t *Grid) IsEmpty() bool {
	return len(t.Config.Panels) == 0
}

// Add adds additional panels to the grid
func (t *Grid) Add(panels ...Panel) {
	t.Config.Panels = append(t.Config.Panels, panels...)
}

type gridMarshal Grid

// MarshalJSON implements json.Marshaler
func (t *Grid) MarshalJSON() ([]byte, error) {
	m := gridMarshal(*t)
	m.Metadata.Type = "grid"
	return json.Marshal(&m)
}
