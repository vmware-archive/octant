/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package flexlayout

import (
	"github.com/vmware/octant/pkg/action"
	"github.com/vmware/octant/pkg/view/component"
)

// FlexLayout is a flex layout manager.
type FlexLayout struct {
	sections    []*Section
	buttonGroup *component.ButtonGroup
}

// New creates an instance of FlexLayout.
func New() *FlexLayout {
	return &FlexLayout{
		buttonGroup: component.NewButtonGroup(),
	}
}

// AddSection adds a new section to the flex layout.
func (fl *FlexLayout) AddSection() *Section {
	section := NewSection()
	fl.sections = append(fl.sections, section)
	return section
}

// AddButton adds a button the button group for a flex layout.
func (fl *FlexLayout) AddButton(name string, payload action.Payload, buttonOptions ...component.ButtonOption) {
	button := component.NewButton(name, payload, buttonOptions...)
	fl.buttonGroup.AddButton(button)
}

// ToComponent converts the FlexLayout to a FlexLayout.
func (fl *FlexLayout) ToComponent(title string) *component.FlexLayout {
	var sections []component.FlexLayoutSection

	for _, section := range fl.sections {
		layoutSection := component.FlexLayoutSection{}

		for _, member := range section.Members {
			item := component.FlexLayoutItem{
				Width: member.Width,
				View:  member.View,
			}

			layoutSection = append(layoutSection, item)
		}

		sections = append(sections, layoutSection)
	}

	if title == "" {
		title = "Summary"
	}

	view := component.NewFlexLayout(title)
	view.AddSections(sections...)
	view.SetButtonGroup(fl.buttonGroup)

	return view
}
