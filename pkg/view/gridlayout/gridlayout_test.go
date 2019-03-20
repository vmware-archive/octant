package gridlayout_test

import (
	"testing"

	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/heptio/developer-dash/pkg/view/gridlayout"
	"github.com/stretchr/testify/assert"
)

func TestGridLayout(t *testing.T) {
	gl := gridlayout.New()

	section1 := gl.CreateSection(8)

	t1 := component.NewText("panel1")
	t2 := component.NewText("panel2")
	t3 := component.NewText("panel3")

	section1.Add(t1, 12)
	section1.Add(t2, 12)
	section1.Add(t3, 12)

	section2 := gl.CreateSection(2)

	t4 := component.NewText("panel4")
	t5 := component.NewText("panel5")
	t6 := component.NewText("panel6")

	section2.Add(t4, 8)
	section2.Add(t5, 8)
	section2.Add(t6, 8)

	got := gl.ToGrid()

	p1 := component.NewPanel("", t1)
	p1.Position(0, 0, 12, 8)

	p2 := component.NewPanel("", t2)
	p2.Position(12, 0, 12, 8)

	p3 := component.NewPanel("", t3)
	p3.Position(0, 9, 12, 8)

	p4 := component.NewPanel("", t4)
	p4.Position(0, 18, 8, 2)

	p5 := component.NewPanel("", t5)
	p5.Position(8, 18, 8, 2)

	p6 := component.NewPanel("", t6)
	p6.Position(16, 18, 8, 2)

	panels := []component.Panel{
		*p1, *p2, *p3, *p4, *p5, *p6,
	}
	expected := component.NewGrid("Summary", panels...)

	assert.Equal(t, expected, got)
}
