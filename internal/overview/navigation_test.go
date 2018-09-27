package overview

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_navigationEntries(t *testing.T) {
	got, err := navigationEntries()
	require.NoError(t, err)

	assert.Equal(t, got.Title, "Overview")
}
