package component

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlert(t *testing.T) {
	buttonGroup := NewButtonGroup()
	got := NewAlert(AlertStatusSuccess, AlertTypeDefault, "message", true, buttonGroup)
	expected := Alert{
		Status:      AlertStatusSuccess,
		Type:        AlertTypeDefault,
		Message:     "message",
		Closable:    true,
		ButtonGroup: buttonGroup,
	}

	assert.Equal(t, got, expected)
}
