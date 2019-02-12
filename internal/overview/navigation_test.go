package overview_test

import (
	"testing"

	"github.com/heptio/developer-dash/internal/overview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NavigationFactory_Entries(t *testing.T) {
	nf := overview.NewNavigationFactory("", "/content/overview")
	got, err := nf.Entries()
	require.NoError(t, err)

	assert.Equal(t, got.Title, "Overview")
	assert.Equal(t, got.Path, "/content/overview/")

	assert.Equal(t, "/content/overview/workloads/cron-jobs", got.Children[0].Children[0].Path)
}

func Test_NavigationFactory_Entries_Namespace(t *testing.T) {
	nf := overview.NewNavigationFactory("default", "/content/overview")
	got, err := nf.Entries()
	require.NoError(t, err)

	assert.Equal(t, got.Title, "Overview")
	assert.Equal(t, got.Path, "/content/overview/namespace/default/")

	assert.Equal(t, "/content/overview/namespace/default/workloads/cron-jobs", got.Children[0].Children[0].Path)
}

func Test_NavigationFactory_Root(t *testing.T) {
	cases := []struct {
		name      string
		path      string
		namespace string
		expected  string
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
		{
			name:      "without trailing slash (namespaced)",
			path:      "/content/overview",
			namespace: "default",
			expected:  "/content/overview/namespace/default/",
		},
		{
			name:      "with trailing slash (namespaced)",
			path:      "/content/overview/",
			namespace: "default",
			expected:  "/content/overview/namespace/default/",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			nf := overview.NewNavigationFactory(tc.namespace, tc.path)
			assert.Equal(t, tc.expected, nf.Root())
		})
	}
}
