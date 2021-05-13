/*
Copyright (c) 2021 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"fmt"

	"github.com/vmware-tanzu/octant/internal/util/json"
)

type TabsView struct {
	Base
	Config TabsViewConfig `json:"config"`
}

// TabsOrientation is the direction of the Tabs
type TabsOrientation string

const (
	// VerticalTabs are tabs organized vertically
	VerticalTabs TabsOrientation = "vertical"
	// HorizontalTabs are tabs organized horizontally
	HorizontalTabs TabsOrientation = "horizontal"
)

type TabsViewConfig struct {
	// Tabs are an array of Tab structs
	Tabs []SingleTab `json:"tabs"`
	// Orientation is the direction of the tabs
	Orientation TabsOrientation `json:"orientation,omitempty"`
}

func NewTabs(orientation TabsOrientation, tabs []SingleTab) *TabsView {
	return &TabsView{
		Base: newBase(TypeTabsView, nil),
		Config: TabsViewConfig{
			Tabs:        tabs,
			Orientation: orientation,
		},
	}
}

type tabsMarshal TabsView

// MarshalJSON marshals a button group.
func (t *TabsView) MarshalJSON() ([]byte, error) {
	m := tabsMarshal(*t)
	m.Metadata.Type = TypeTabsView
	return json.Marshal(&m)
}

var _ Component = (*TabsView)(nil)

type SingleTab struct {
	Name     string     `json:"name"`
	Contents FlexLayout `json:"contents"`
}

func (t *TabsViewConfig) UnmarshalJSON(data []byte) error {
	x := struct {
		Orientation TabsOrientation `json:"orientation,omitempty"`
		Tabs        []struct {
			Name     string      `json:"name"`
			Contents TypedObject `json:"contents"`
		} `json:"tabs"`
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	for _, tab := range x.Tabs {
		c, err := tab.Contents.ToComponent()
		if err != nil {
			return err
		}
		fl, ok := c.(*FlexLayout)
		if !ok {
			return fmt.Errorf("item was not a FlexLayout")
		}
		st := SingleTab{
			Name:     tab.Name,
			Contents: *fl,
		}
		t.Tabs = append(t.Tabs, st)
	}

	t.Orientation = x.Orientation

	return nil
}
