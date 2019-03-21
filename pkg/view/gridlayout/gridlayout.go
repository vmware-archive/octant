package gridlayout

import "github.com/heptio/developer-dash/pkg/view/component"

const (
	gridWidth = 24
)

// GridLayout is a grid layout manager.
type GridLayout struct {
	sections []*Section
}

// New creates an instance of GridLayout.
func New() *GridLayout {
	return &GridLayout{}
}

// CreateSection creates a new section for the grid layout.
func (gl *GridLayout) CreateSection(height int) *Section {
	section := NewSection(height)
	gl.sections = append(gl.sections, section)

	return section
}

// ToGrid converts the GridLayout to a Grid.
func (gl *GridLayout) ToGrid() *component.Grid {
	row := 0
	col := 0

	var panels []component.Panel

	for _, section := range gl.sections {
		for _, member := range section.Members {
			panel := component.NewPanel("", member.View)
			panel.Position(col, row, member.Width, section.Height)

			col += member.Width
			if col >= gridWidth {
				col = 0
				row += section.Height + 1
			}

			panels = append(panels, *panel)
		}

		col = 0
		row += section.Height + 1
	}

	grid := component.NewGrid("Summary", panels...)

	return grid
}
