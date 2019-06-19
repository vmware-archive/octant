/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package flexlayout

import (
	"github.com/vmware/octant/pkg/view/component"
)

// FlexLayout is a flex layout manager.
type FlexLayout struct {
	sections []*Section
}

// New creates an instance of FlexLayout.
func New() *FlexLayout {
	return &FlexLayout{}
}

// AddSection adds a new section to the grid layout.
func (fl *FlexLayout) AddSection() *Section {
	section := NewSection()
	fl.sections = append(fl.sections, section)
	return section
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

	return view
}
