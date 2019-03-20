package flexlayout_test

import (
	"testing"

	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/heptio/developer-dash/pkg/view/flexlayout"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlexLayout(t *testing.T) {
	fl := flexlayout.New()

	section1 := fl.AddSection()

	t1 := component.NewText("item 1")
	t2 := component.NewText("item 2")
	t3 := component.NewText("item 3")

	require.NoError(t, section1.Add(t1, 12))
	require.NoError(t, section1.Add(t2, 12))
	require.NoError(t, section1.Add(t3, 12))

	section2 := fl.AddSection()

	t4 := component.NewText("item 4")
	t5 := component.NewText("item 4")

	require.NoError(t, section2.Add(t4, 6))
	require.NoError(t, section2.Add(t5, 6))

	got := fl.ToComponent("Title")


	expected := component.NewFlexLayout("Title")
	expected.AddSections([]component.FlexLayoutSection{
		{
			component.FlexLayoutItem{Width: 12, View: t1},
			component.FlexLayoutItem{Width: 12, View: t2},
			component.FlexLayoutItem{Width: 12, View: t3},
		},
		{
			component.FlexLayoutItem{Width: 6, View: t4},
			component.FlexLayoutItem{Width: 6, View: t5},
		},
	}...)

	assert.Equal(t, expected, got)
}

func TestFlexLayout_default_title(t *testing.T) {
	fl := flexlayout.New()

	got := fl.ToComponent("")
	expected := component.NewFlexLayout("Summary")

	assert.Equal(t, expected, got)
}
