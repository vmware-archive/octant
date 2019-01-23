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
)

// Quadrant contains other ViewComponents
type Quadrant struct {
	Metadata Metadata       `json:"metadata"`
	Config   QuadrantConfig `json:"config"`
}

// QuadrantConfig is the contents of a Quadrant
type QuadrantConfig struct {
	NW ViewComponent
	NE ViewComponent
	SE ViewComponent
	SW ViewComponent
}

func (t *QuadrantConfig) UnmarshalJSON(data []byte) error {
	x := struct {
		NW typedObject
		NE typedObject
		SW typedObject
		SE typedObject
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	var err error

	t.NW, err = x.NW.ToViewComponent()
	if err != nil {
		return err
	}
	t.NE, err = x.NE.ToViewComponent()
	if err != nil {
		return err
	}
	t.SW, err = x.SW.ToViewComponent()
	if err != nil {
		return err
	}
	t.SE, err = x.SE.ToViewComponent()
	if err != nil {
		return err
	}

	return nil
}

// NewQuadrant creates a quadrant component
func NewQuadrant() *Quadrant {
	return &Quadrant{
		Metadata: Metadata{
			Type: "quadrant",
		},
		Config: QuadrantConfig{},
	}
}

// GetMetadata accesses the components metadata. Implements ViewComponent.
func (t *Quadrant) GetMetadata() Metadata {
	return t.Metadata
}

// IsEmpty specifes whether the component is considered empty. Implements ViewComponent.
func (t *Quadrant) IsEmpty() bool {
	return t.Config.NW == nil &&
		t.Config.NE == nil &&
		t.Config.SE == nil &&
		t.Config.SW == nil
}

// Set adds additional panels to the quadrant
func (t *Quadrant) Set(pos QuadrantPosition, content ViewComponent) error {
	switch pos {
	case QuadNW:
		t.Config.NW = content
	case QuadNE:
		t.Config.NE = content
	case QuadSE:
		t.Config.SE = content
	case QuadSW:
		t.Config.SW = content
	default:
		return fmt.Errorf("invalid quadrant position: %v", pos)
	}
	return nil
}

type quadrantMarshal Quadrant

// MarshalJSON implements json.Marshaler
func (t *Quadrant) MarshalJSON() ([]byte, error) {
	m := quadrantMarshal(*t)
	m.Metadata.Type = "quadrant"
	return json.Marshal(&m)
}
