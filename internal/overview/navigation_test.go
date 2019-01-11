package overview_test

import (
	"testing"

	"github.com/heptio/developer-dash/internal/overview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NavigationFactory_Entries(t *testing.T) {
	nf := overview.NewNavigationFactory("/content/overview")
	got, err := nf.Entries()
	require.NoError(t, err)

	assert.Equal(t, got.Title, "Overview")
	assert.Equal(t, got.Path, "/content/overview/")
}

func Test_NavigationFactory_Root(t *testing.T) {
	cases := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "without trailing slash",
			path:     "/content/overview",
			expected: "/content/overview/",
		},
		{
			name:     "with trailing slash",
			path:     "/content/overview/",
			expected: "/content/overview/",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			nf := overview.NewNavigationFactory(tc.path)
			assert.Equal(t, tc.expected, nf.Root())
		})
	}
}
