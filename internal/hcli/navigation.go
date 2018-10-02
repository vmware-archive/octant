package hcli

// Navigation is a set of navigation entries.
type Navigation struct {
	Title    string        `json:"title,omitempty"`
	Path     string        `json:"path,omitempty"`
	Children []*Navigation `json:"children,omitempty"`
}
