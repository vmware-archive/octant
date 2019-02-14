package component

import (
	"encoding/json"
	"fmt"
)

// QuadrantPosition denotes a position within a quadrant
type QuadrantPosition int

const (
	// QuadNW denotes the north-west position within a quadrant
	QuadNW QuadrantPosition = iota
	// QuadNE denotes the north-east position within a quadrant
	QuadNE
	// QuadSE denotes the south-east position within a quadrant
	QuadSE
	// QuadSW denotes the south-west position within a quadrant
	QuadSW

	// quadrantType is the component type for a quadrant
	quadrantType = "quadrant"
)

// Quadrant contains other ViewComponents
type Quadrant struct {
	Metadata Metadata       `json:"metadata"`
	Config   QuadrantConfig `json:"config"`
}

type QuadrantValue struct {
	Value string `json:"value,omitempty"`
	Label string `json:"label,omitempty"`
}

// QuadrantConfig is the contents of a Quadrant
type QuadrantConfig struct {
	NW QuadrantValue `json:"nw,omitempty"`
	NE QuadrantValue `json:"ne,omitempty"`
	SE QuadrantValue `json:"se,omitempty"`
	SW QuadrantValue `json:"sw,omitempty"`
}

// NewQuadrant creates a quadrant component
func NewQuadrant() *Quadrant {
	return &Quadrant{
		Metadata: Metadata{
			Type: quadrantType,
		},
		Config: QuadrantConfig{},
	}
}

// GetMetadata accesses the components metadata. Implements ViewComponent.
func (t *Quadrant) GetMetadata() Metadata {
	return t.Metadata
}

// Set adds additional panels to the quadrant
func (t *Quadrant) Set(pos QuadrantPosition, label, value string) error {
	qv := QuadrantValue{Label: label, Value: value}
	switch pos {
	case QuadNW:
		t.Config.NW = qv
	case QuadNE:
		t.Config.NE = qv
	case QuadSE:
		t.Config.SE = qv
	case QuadSW:
		t.Config.SW = qv
	default:
		return fmt.Errorf("invalid quadrant position: %v", pos)
	}
	return nil
}

type quadrantMarshal Quadrant

// MarshalJSON implements json.Marshaler
func (t *Quadrant) MarshalJSON() ([]byte, error) {
	m := quadrantMarshal(*t)
	m.Metadata.Type = quadrantType
	return json.Marshal(&m)
}
