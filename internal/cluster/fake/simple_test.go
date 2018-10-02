package fake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimpleClusterOverview_Navigation(t *testing.T) {
	sco := NewSimpleClusterOverview()
	_, err := sco.Navigation("/root")
	require.NoError(t, err)
}
