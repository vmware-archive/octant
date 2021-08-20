package component

import (
	"github.com/vmware-tanzu/octant/internal/util/json"
)

type Icon struct {
	Base
	Config IconConfig `json:"config"`
}

type IconConfig struct {
	Shape      string         `json:"shape"`
	Size       string         `json:"size"`
	Direction  Direction      `json:"direction"`
	Flip       Flip           `json:"flip"`
	Solid      bool           `json:"solid"`
	Status     Status         `json:"status"`
	Inverse    bool           `json:"inverse"`
	Badge      Badge          `json:"badge"`
	Color      string         `json:"color"`
	BadgeColor string         `json:"badgeColor"`
	Label      string         `json:"label"`
	Tooltip    *TooltipConfig `json:"tooltip,omitempty"`
	CustomSvg  string         `json:"customSvg"`
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
	BadgeWarning         Badge = "warning"
	BadgeInherit         Badge = "inherit"
	BadgeInheritTriangle Badge = "inherit-triangle"
)

type IconOption func(icon *Icon)

type TooltipConfig struct {
	Message  string          `json:"message"`
	Size     TooltipSize     `json:"size"`
	Position TooltipPosition `json:"position"`
}

type TooltipSize string

const (
	TooltipExtraSmall TooltipSize = "xs"
	TooltipSmall      TooltipSize = "sm"
	TooltipMedium     TooltipSize = "md"
	TooltipLarge      TooltipSize = "lg"
)

type TooltipPosition string

const (
	TooltipLeft        TooltipPosition = "left"
	TooltipRight       TooltipPosition = "right"
	TooltipTopLeft     TooltipPosition = "top-left"
	TooltipTopRight    TooltipPosition = "top-right"
	TooltipBottomLeft  TooltipPosition = "bottom-left"
	TooltipBottomRight TooltipPosition = "bottom-right"
)

func NewIcon(shape string, options ...func(*Icon)) *Icon {
	i := &Icon{
		Base: newBase(TypeIcon, nil),
		Config: IconConfig{
			Shape: shape,
		},
	}

	for _, option := range options {
		option(i)
	}

	return i
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

// WithTooltip configures an icon with a tooltip
func WithTooltip(message string, size TooltipSize, position TooltipPosition) IconOption {
	return func(icon *Icon) {
		icon.Config.Tooltip = &TooltipConfig{
			Message:  message,
			Size:     size,
			Position: position,
		}
	}
}

// AddLabel adds an aria-label for screen readers
func (i *Icon) AddLabel(label string) {
	i.Config.Label = label
}

// SetColor sets the color of an icon. A color from status has priority over a set color.
func (i *Icon) SetColor(color string) {
	i.Config.Color = color
}

// SetBadgeColor sets the color of a badge. A set badge color has priority over a badge status.
func (i *Icon) SetBadgeColor(color string) {
	i.Config.BadgeColor = color
}

// SetSize sets the size of a badge. The size can me sm, md, lg, xl, xxl, or an integer for N x N pixels.
func (i *Icon) SetSize(size string) {
	i.Config.Size = size
}

// SetCustomSvg sets an SVG for an Icon.
func (i *Icon) SetCustomSvg(customSvg string) {
	i.Config.CustomSvg = customSvg
}
