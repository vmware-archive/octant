/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"encoding/json"
)

// DropdownType defines what the dropdown source is (UI component that's visible when dropdown is closed)
//
type DropdownType string

const (
	DropdownButton DropdownType = "button"
	DropdownLink   DropdownType = "link"
	DropdownLabel  DropdownType = "label"
	DropdownIcon   DropdownType = "icon"
)

// DropdownPosition denotes a position relative to source button
type DropdownPosition string

const (
	BottomLeft  DropdownPosition = "bottom-left"
	BottomRight DropdownPosition = "bottom-right"
	TopLeft     DropdownPosition = "top-left"
	TopRight    DropdownPosition = "top-right"
	LeftBottom  DropdownPosition = "left-bottom"
	LeftTop     DropdownPosition = "left-top"
	RightTop    DropdownPosition = "right-top"
	RightBottom DropdownPosition = "right-bottom"
)

// Defines the type of dropdown item
type ItemType string

const (
	Header    ItemType = "header"
	PlainText ItemType = "text"
	Url       ItemType = "link"
	Separator ItemType = "separator"
)

type DropdownItemConfig struct {
	Name        string   `json:"name"`
	Type        ItemType `json:"type"`
	Label       string   `json:"label"`
	Url         string   `json:"url,omitempty"`
	Description string   `json:"description,omitempty"`
}

// DropdownConfig defines the contents of a Dropdown
type DropdownConfig struct {
	DropdownPosition DropdownPosition     `json:"position,omitempty"`
	DropdownType     DropdownType         `json:"type"`
	Action           string               `json:"action,omitempty"`
	Selection        string               `json:"selection,omitempty"`
	UseSelection     bool                 `json:"useSelection"`
	Items            []DropdownItemConfig `json:"items"`
}

// Dropdown: displays dropdown component with a list of values
// Used to choose an option or action from a contextual list
// +octant:component
type Dropdown struct {
	Base
	Config DropdownConfig `json:"config"`
}

// NewDropdown creates a new dropdown component
func NewDropdown(title string, dropdownType DropdownType, action string, items ...DropdownItemConfig) *Dropdown {
	dropdownItems := append([]DropdownItemConfig(nil), items...)
	return &Dropdown{
		Base: newBase(TypeDropdown, TitleFromString(title)),
		Config: DropdownConfig{
			DropdownType: dropdownType,
			Action:       action,
			Items:        dropdownItems,
		},
	}
}

// NewDropdownItem  creates a new dropdown item
func NewDropdownItem(name string, itemType ItemType, label string, url string, description string) DropdownItemConfig {
	item := DropdownItemConfig{
		Name:        name,
		Type:        itemType,
		Label:       label,
		Url:         url,
		Description: description,
	}
	return item
}

// AddDropdownItem adds an item to the dropdown
func (t *Dropdown) AddDropdownItem(name string, itemType ItemType, label string, url string, description string) {
	item := NewDropdownItem(name, itemType, label, url, description)
	t.Config.Items = append(t.Config.Items, item)
}

// SetDropdownPosition sets the position of context menu relative to dropdown source.
func (t *Dropdown) SetDropdownPosition(position DropdownPosition) {
	t.Config.DropdownPosition = position
}

// SetSelection specifies the dropdown selected item.
func (t *Dropdown) SetSelection(selection string) {
	t.Config.Selection = selection
}

// SetDropdownUseSelection defines if dropdown title is updated on selection change
func (t *Dropdown) SetDropdownUseSelection(sel bool) {
	t.Config.UseSelection = sel
}

// SupportsTitle designates this is a TextComponent.
func (t *Dropdown) SupportsTitle() {}

// GetMetadata accesses the components metadata
func (t *Dropdown) GetMetadata() Metadata {
	return t.Metadata
}

type dropdownMarshal Dropdown

// MarshalJSON implements json.Marshaler
func (t *Dropdown) MarshalJSON() ([]byte, error) {
	m := dropdownMarshal(*t)
	m.Metadata.Type = TypeDropdown
	return json.Marshal(&m)
}

func (t *DropdownConfig) UnmarshalJSON(data []byte) error {
	x := struct {
		DropdownPosition DropdownPosition     `json:"position,omitempty"`
		DropdownType     DropdownType         `json:"type,omitempty"`
		Action           string               `json:"action,omitempty"`
		Selection        string               `json:"selection,omitempty"`
		UseSelection     bool                 `json:"useSelection,omitempty"`
		Items            []DropdownItemConfig `json:"items"`
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	t.DropdownPosition = x.DropdownPosition
	t.DropdownType = x.DropdownType
	t.Action = x.Action
	t.Selection = x.Selection
	t.UseSelection = x.UseSelection
	t.Items = x.Items
	return nil
}
