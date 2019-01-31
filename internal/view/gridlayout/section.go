package gridlayout

import "github.com/heptio/developer-dash/internal/view/component"

// SectionMember is a member of a section.
type SectionMember struct {
	View  component.ViewComponent
	Width int
}

// Section defines a section in a grid layout. Sections can contain
// multiple views.
type Section struct {
	Height  int
	Members []SectionMember
}

// NewSection create an instance of Section.
func NewSection(height int) *Section {
	return &Section{
		Height: height,
	}
}

// Add adds a view to the section with a width.
func (s *Section) Add(view component.ViewComponent, width int) {
	member := SectionMember{View: view, Width: width}
	s.Members = append(s.Members, member)
}
