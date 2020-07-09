/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"encoding/json"
	"errors"
)

const (
	// WidthFull is a full width section.
	WidthFull int = 24
	// WidthQuarter is a quarter width section.
	WidthQuarter int = 6
	// WidthHalf is a half width section.
	WidthHalf int = 12
	// WidthThird is a third width section.
	WidthThird int = 8
)

// FlexLayoutItem is an item in a flex layout.
type FlexLayoutItem struct {
	Width  int       `json:"width,omitempty"`
	Height string    `json:"height,omitempty"`
	Margin string    `json:"margin,omitempty"`
	View   Component `json:"view,omitempty"`
}

func (fli *FlexLayoutItem) UnmarshalJSON(data []byte) error {
	x := struct {
		Width int
		View  TypedObject
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	fli.Width = x.Width
	var err error
	fli.View, err = x.View.ToComponent()
	if err != nil {
		return err
	}

	return nil
}

// FlexLayoutSection is a slice of items group together.
type FlexLayoutSection []FlexLayoutItem

// FlexLayoutConfig is configuration for the flex layout view.
type FlexLayoutConfig struct {
	Sections    []FlexLayoutSection `json:"sections,omitempty"`
	ButtonGroup *ButtonGroup        `json:"buttonGroup,omitempty"`
}

func (f *FlexLayoutConfig) UnmarshalJSON(data []byte) error {
	x := struct {
		Sections    []FlexLayoutSection `json:"sections,omitempty"`
		ButtonGroup *TypedObject        `json:"buttonGroup,omitempty"`
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	if x.ButtonGroup != nil {
		component, err := x.ButtonGroup.ToComponent()
		if err != nil {
			return err
		}

		buttonGroup, ok := component.(*ButtonGroup)
		if !ok {
			return errors.New("item was not a buttonGroup")
		}
		f.ButtonGroup = buttonGroup
	}

	f.Sections = x.Sections

	return nil
}

// FlexLayout is a flex layout view.
type FlexLayout struct {
	base
	Config FlexLayoutConfig `json:"config,omitempty"`
}

// NewFlexLayout creates an instance of FlexLayout.
func NewFlexLayout(title string) *FlexLayout {
	return &FlexLayout{
		base: newBase(typeFlexLayout, TitleFromString(title)),
		Config: FlexLayoutConfig{
			ButtonGroup: NewButtonGroup(),
		},
	}
}

var _ Component = (*FlexLayout)(nil)

// AddSections adds one or more sections to the flex layout.
func (fl *FlexLayout) AddSections(sections ...FlexLayoutSection) {
	fl.Config.Sections = append(fl.Config.Sections, sections...)
}

type flexLayoutMarshal FlexLayout

// MarshalJSON marshals the flex layout to JSON.
func (fl *FlexLayout) MarshalJSON() ([]byte, error) {
	x := flexLayoutMarshal(*fl)
	x.Metadata.Type = typeFlexLayout
	return json.Marshal(&x)
}

func (fl *FlexLayout) SetButtonGroup(group *ButtonGroup) {
	fl.Config.ButtonGroup = group
}

// Tab represents a tab. A tab is a flex layout with a name.
type Tab struct {
	Name     string     `json:"name"`
	Contents FlexLayout `json:"contents"`
}

// NewTabWithContents creates a tab with contents.
func NewTabWithContents(flexLayout FlexLayout) *Tab {
	name, err := TitleFromTitleComponent(flexLayout.Title)
	if err != nil {
		name = ""
	}

	return &Tab{
		Name:     name,
		Contents: flexLayout,
	}
}
