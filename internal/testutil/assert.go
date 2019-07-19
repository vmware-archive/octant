package testutil

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertJSONEqual asserts two object's generated JSON is equal.
func AssertJSONEqual(t *testing.T, expected, actual interface{}) {
	a, err := json.Marshal(expected)
	require.NoError(t, err)

	b, err := json.Marshal(actual)
	require.NoError(t, err)

	assert.JSONEq(t, string(a), string(b))
}
