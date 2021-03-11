package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/util/json"
)

// AssertJSONEqual asserts two object's generated JSON is equal.
func AssertJSONEqual(t *testing.T, expected, actual interface{}) {
	a, err := json.MarshalIndent(expected, "", "  ")
	require.NoError(t, err)

	b, err := json.MarshalIndent(actual, "", "  ")
	require.NoError(t, err)

	assert.JSONEq(t, string(a), string(b))
}
