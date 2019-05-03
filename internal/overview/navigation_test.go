package overview_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	fakestore "github.com/heptio/developer-dash/internal/objectstore/fake"
	"github.com/heptio/developer-dash/internal/overview"
	"github.com/heptio/developer-dash/pkg/objectstoreutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Test_NavigationFactory_Entries(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	o := fakestore.NewMockObjectStore(controller)

	key := objectstoreutil.Key{
		APIVersion: "apiextensions.k8s.io/v1beta1",
		Kind:       "CustomResourceDefinition",
	}

	o.EXPECT().CheckAccess(gomock.Any()).Return(nil)
	o.EXPECT().
		List(gomock.Any(), gomock.Eq(key)).
		Return([]*unstructured.Unstructured{}, nil)

	nf := overview.NewNavigationFactory("", "/content/overview", o)
	ctx := context.Background()
	got, err := nf.Entries(ctx)
	require.NoError(t, err)

	assert.Equal(t, got.Title, "Overview")
	assert.Equal(t, got.Path, "/content/overview/")

	assert.Equal(t, "/content/overview/workloads/cron-jobs", got.Children[0].Children[0].Path)
}

func Test_NavigationFactory_Entries_Namespace(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	o := fakestore.NewMockObjectStore(controller)

	key := objectstoreutil.Key{
		APIVersion: "apiextensions.k8s.io/v1beta1",
		Kind:       "CustomResourceDefinition",
	}

	o.EXPECT().CheckAccess(gomock.Any()).Return(nil)
	o.EXPECT().
		List(gomock.Any(), gomock.Eq(key)).
		Return([]*unstructured.Unstructured{}, nil)

	nf := overview.NewNavigationFactory("default", "/content/overview", o)
	ctx := context.Background()
	got, err := nf.Entries(ctx)
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
			controller := gomock.NewController(t)
			defer controller.Finish()
			o := fakestore.NewMockObjectStore(controller)

			nf := overview.NewNavigationFactory(tc.namespace, tc.path, o)
			assert.Equal(t, tc.expected, nf.Root())
		})
	}
}
