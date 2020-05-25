package astilectron

// Display event names
const (
	EventNameDisplayEventAdded          = "display.event.added"
	EventNameDisplayEventMetricsChanged = "display.event.metrics.changed"
	EventNameDisplayEventRemoved        = "display.event.removed"
)

// Display represents a display
// https://github.com/electron/electron/blob/v1.8.1/docs/api/structures/display.md
type Display struct {
	o       *DisplayOptions
	primary bool
}

// DisplayOptions represents display options
// https://github.com/electron/electron/blob/v1.8.1/docs/api/structures/display.md
type DisplayOptions struct {
	Bounds       *RectangleOptions `json:"bounds,omitempty"`
	ID           *int64            `json:"id,omitempty"`
	Rotation     *int              `json:"rotation,omitempty"` // 0, 90, 180 or 270
	ScaleFactor  *float64          `json:"scaleFactor,omitempty"`
	Size         *SizeOptions      `json:"size,omitempty"`
	TouchSupport *string           `json:"touchSupport,omitempty"` // available, unavailable or unknown
	WorkArea     *RectangleOptions `json:"workArea,omitempty"`
	WorkAreaSize *SizeOptions      `json:"workAreaSize,omitempty"`
}

// newDisplay creates a displays
func newDisplay(o *DisplayOptions, primary bool) *Display { return &Display{o: o, primary: primary} }

// Bounds returns the display bounds
func (d Display) Bounds() Rectangle {
	return Rectangle{
		Position: Position{X: *d.o.Bounds.X, Y: *d.o.Bounds.Y},
		Size:     Size{Height: *d.o.Bounds.Height, Width: *d.o.Bounds.Width},
	}
}

// ID returns the display's ID
func (d Display) ID() int64 {
	return *d.o.ID
}

// IsPrimary checks whether the display is the primary display
func (d Display) IsPrimary() bool {
	return d.primary
}

// IsTouchAvailable checks whether touch is available on this display
func (d Display) IsTouchAvailable() bool {
	return *d.o.TouchSupport == "available"
}

// Rotation returns the display rotation
func (d Display) Rotation() int {
	return *d.o.Rotation
}

// ScaleFactor returns the display scale factor
func (d Display) ScaleFactor() float64 {
	return *d.o.ScaleFactor
}

// Size returns the display size
func (d Display) Size() Size {
	return Size{Height: *d.o.Size.Height, Width: *d.o.Size.Width}
}

// WorkArea returns the display work area
func (d Display) WorkArea() Rectangle {
	return Rectangle{
		Position: Position{X: *d.o.WorkArea.X, Y: *d.o.WorkArea.Y},
		Size:     Size{Height: *d.o.WorkArea.Height, Width: *d.o.WorkArea.Width},
	}
}

// WorkAreaSize returns the display work area size
func (d Display) WorkAreaSize() Size {
	return Size{Height: *d.o.WorkAreaSize.Height, Width: *d.o.WorkAreaSize.Width}
}
