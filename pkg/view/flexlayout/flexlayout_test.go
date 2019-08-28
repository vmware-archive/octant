/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package flexlayout_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vmware/octant/pkg/action"
	"github.com/vmware/octant/pkg/view/component"
	"github.com/vmware/octant/pkg/view/flexlayout"
)

func TestFlexLayout(t *testing.T) {
	fl := flexlayout.New()

	section1 := fl.AddSection()

	t1 := component.NewText("item 1")
	t2 := component.NewText("item 2")
	t3 := component.NewText("item 3")

	require.NoError(t, section1.Add(t1, component.WidthFull))
	require.NoError(t, section1.Add(t2, component.WidthFull))
	require.NoError(t, section1.Add(t3, component.WidthFull))

	section2 := fl.AddSection()

	t4 := component.NewText("item 4")
	t5 := component.NewText("item 4")

	require.NoError(t, section2.Add(t4, component.WidthHalf))
	require.NoError(t, section2.Add(t5, component.WidthHalf))

	got := fl.ToComponent("Title")

	expected := component.NewFlexLayout("Title")
	expected.AddSections([]component.FlexLayoutSection{
		{
			component.FlexLayoutItem{Width: component.WidthFull, View: t1},
			component.FlexLayoutItem{Width: component.WidthFull, View: t2},
			component.FlexLayoutItem{Width: component.WidthFull, View: t3},
		},
		{
			component.FlexLayoutItem{Width: component.WidthHalf, View: t4},
			component.FlexLayoutItem{Width: component.WidthHalf, View: t5},
		},
	}...)

	component.AssertEqual(t, expected, got)
}

func TestFlexLayout_default_title(t *testing.T) {
	fl := flexlayout.New()

	got := fl.ToComponent("")
	expected := component.NewFlexLayout("Summary")

	component.AssertEqual(t, expected, got)
}

func TestFlexLayout_add_button(t *testing.T) {
	fl := flexlayout.New()
	fl.AddButton("button", action.Payload{})

	got := fl.ToComponent("Title")
	expected := component.NewFlexLayout("Title")
	buttonGroup := component.NewButtonGroup()
	buttonGroup.AddButton(component.NewButton("button", action.Payload{}))
	expected.SetButtonGroup(buttonGroup)

	component.AssertEqual(t, expected, got)
}
