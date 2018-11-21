package overview

import (
	"context"
	"testing"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecretData_InvalidObject(t *testing.T) {
	assertViewInvalidObject(t, NewSecretData())
}

func TestSecretData(t *testing.T) {
	v := NewSecretData()

	ctx := context.Background()
	cache := NewMemoryCache()

	secret := loadFromFile(t, "secret-1.yaml")
	secret = convertToInternal(t, secret)

	got, err := v.Content(ctx, secret, cache)
	require.NoError(t, err)

	dataSection := content.NewSection()
	dataSection.AddText("ca.crt", "1025 bytes")
	dataSection.AddText("namespace", "8 bytes")
	dataSection.AddText("token", "token")

	dataSummary := content.NewSummary("Data", []content.Section{dataSection})

	expected := []content.Content{
		&dataSummary,
	}

	assert.Equal(t, got, expected)
}
