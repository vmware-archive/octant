package component

import "encoding/json"

// Panel contains other ViewComponents
type Panel struct {
	Metadata Metadata    `json:"metadata"`
	Config   PanelConfig `json:"config"`
}

// PanelConfig is the contents of a Panel
type PanelConfig struct {
	Content  ViewComponent `json:"content"`
	Position PanelPosition `json:"position"`
}

// PanelPosition represents the relative location and size of a panel within a grid
type PanelPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

// NewPanel creates a panel component
func NewPanel(title string, content ViewComponent) *Panel {
	return &Panel{
		Metadata: Metadata{
			Type: "panel",
		},
		Config: PanelConfig{
			Content: content,
		},
	}
}

// GetMetadata accesses the components metadata. Implements ViewComponent.
func (t *Panel) GetMetadata() Metadata {
	return t.Metadata
}

// IsEmpty specifes whether the component is considered empty. Implements ViewComponent.
func (t *Panel) IsEmpty() bool {
	return t.Config.Content == nil
}

type panelMarshal Panel

// MarshalJSON implements json.Marshaler
func (t *Panel) MarshalJSON() ([]byte, error) {
	m := panelMarshal(*t)
	m.Metadata.Type = "panel"
	return json.Marshal(&m)
}
