package overview

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSectionDescriber(t *testing.T) {
	namespace := "default"

	d := NewSectionDescriber(
		newStubDescriber(),
	)

	cache := NewMemoryCache()

	got, err := d.Describe("/prefix", namespace, cache, nil)
	require.NoError(t, err)

	assert.Equal(t, stubbedContent, got)
}
