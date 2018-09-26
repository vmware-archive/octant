package fake

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimpleClusterOverview_Namespaces(t *testing.T) {
	sco := NewSimpleClusterOverview()
	got, err := sco.Namespaces()
	require.NoError(t, err)

	expected := []string{"default"}
	assert.Equal(t, expected, got)
}

func TestSimpleClusterOverview_Navigation(t *testing.T) {
	sco := NewSimpleClusterOverview()
	err := sco.Navigation()
	require.NoError(t, err)
}

func TestSimpleClusterOverview_Content(t *testing.T) {
	sco := NewSimpleClusterOverview()
	err := sco.Content("/path")
	require.NoError(t, err)
}
