package clustereye

// Navigation is a set of navigation entries.
type Navigation struct {
	Title    string       `json:"title,omitempty"`
	Path     string       `json:"path,omitempty"`
	Children []Navigation `json:"children,omitempty"`
}

// NewNavigation creates a Navigation.
func NewNavigation(title, path string) *Navigation {
	return &Navigation{Title: title, Path: path}
}
