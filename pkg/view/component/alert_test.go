package component

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlert(t *testing.T) {
	got := NewAlert(AlertTypeSuccess, "message")
	expected := Alert{
		Type:    AlertTypeSuccess,
		Message: "message",
	}

	assert.Equal(t, got, expected)
}
