package component

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertEqual asserts two components are equal.
func AssertEqual(t *testing.T, expected, got Component) {
	a, err := json.Marshal(expected)
	require.NoError(t, err)

	b, err := json.Marshal(got)
	require.NoError(t, err)

	assert.JSONEq(t, string(a), string(b))
}
