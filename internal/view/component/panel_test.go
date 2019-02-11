package component

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Panel_Position(t *testing.T) {
	p := NewPanel("", NewText("text"))

	p.Position(1, 2, 3, 4)

	expected := PanelPosition{X: 1, Y: 2, W: 3, H: 4}
	assert.Equal(t, expected, p.Config.Position)

}
