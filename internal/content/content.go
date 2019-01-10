package content

type Metadata struct {
	Type  string `json:"type,omitempty"`
	Title string `json:"title,omitempty"`
}

type ViewComponent struct {
	Metadata Metadata    `json:"metadata,omitempty"`
	Config   interface{} `json:"config,omitempty"`
}

// Content is content served by the overview API.
type Content interface {
	IsEmpty() bool

	ViewComponent() ViewComponent
}
