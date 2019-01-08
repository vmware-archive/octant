package component

// Metadata collects common fields describing ViewComponents
type Metadata struct {
	Type  string `json:"type"`
	Title string `json:"title,omitempty"`
}

// ViewComponent is a common interface for the data representation
// of visual components as rendered by the UI.
type ViewComponent interface {
	IsEmpty() bool
	GetMetadata() Metadata
}

//
