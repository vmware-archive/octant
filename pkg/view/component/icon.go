package component

import (
	"github.com/vmware-tanzu/octant/internal/util/json"
)

// WARNING: This is using version on clarity 4.0.12 but
// this is built using the API of version 5 is not expected to change the the version
// but is good to know :)

type Icon struct {
	Base
	Config IconConfig `json:"config"`
}

type IconConfig struct {
	Shape     string    `json:"shape"`
	Size      string    `json:"size"`
	Direction Direction `json:"direction"`
	Flip      Flip      `json:"flip"`
	Solid     bool      `json:"solid"`
	Status    Status    `json:"status"`
	Inverse   bool      `json:"inverse"`
	Badge     Badge     `json:"badge"`
	Color     string    `json:"color"`
}

type Direction string

const (
	DirectionUp    Direction = "up"
	DirectionDown  Direction = "down"
	DirectionLeft  Direction = "left"
	DirectionRight Direction = "right"
)

type Flip string

const (
	FlipHorizontal Flip = "horizontal"
	FlipVertical   Flip = "vertical"
)

type Status string

const (
	StatusInfo    Status = "info"
	StatusSuccess Status = "success"
	StatusWarning Status = "warning"
	StatusDanger  Status = "danger"
)

type Badge string

const (
	BadgeInfo            Badge = "info"
	BadgeSuccess         Badge = "success"
	BadgeWarningTriangle Badge = "warning-triangle"
	BadgeDanger          Badge = "danger"
	// BadgeWarning         Badge = "warning" just for clarity 5
	// BadgeInherit         Badge = "inherit" just for clarity 5
	// BadgeInheritTriangle Badge = "inherit-triangle" just for clarity 5
)

func NewIcon(shape string) *Icon {
	return &Icon{
		Base: newBase(TypeIcon, nil),
		Config: IconConfig{
			Shape: shape,
		},
	}
}

type iconMarshal Icon

func (i *Icon) MarshalJSON() ([]byte, error) {
	m := iconMarshal{
		Base:   i.Base,
		Config: i.Config,
	}
	m.Metadata.Type = TypeIcon

	return json.Marshal(&m)
}
