package astilectron

// Position represents a position
type Position struct {
	X, Y int
}

// PositionOptions represents position options
type PositionOptions struct {
	X *int `json:"x,omitempty"`
	Y *int `json:"y,omitempty"`
}

// Size represents a size
type Size struct {
	Height, Width int
}

// SizeOptions represents size options
type SizeOptions struct {
	Height *int `json:"height,omitempty"`
	Width  *int `json:"width,omitempty"`
}

// Rectangle represents a rectangle
type Rectangle struct {
	Position
	Size
}

// RectangleOptions represents rectangle options
type RectangleOptions struct {
	PositionOptions
	SizeOptions
}
