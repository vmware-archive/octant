package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestMemoryCache_Store(t *testing.T) {
	c := NewMemoryCache()

	o := &unstructured.Unstructured{}
	o.SetNamespace("default")
	o.SetAPIVersion("foo/v1")
	o.SetKind("Kind")
	o.SetName("name")

	assert.Len(t, c.store, 0)

	err := c.Store(o)
	require.NoError(t, err)

	assert.Len(t, c.store, 1)

	c.Reset()
	assert.Len(t, c.store, 0)
}

func TestMemoryCache_List(t *testing.T) {

	cases := []struct {
		name        string
		key         Key
		expectedLen int
		isErr       bool
	}{
		{
			name: "ns, apiVersion, kind, name",
			key: Key{
				Namespace:  "default",
				APIVersion: "foo/v1",
				Kind:       "Kind",
				Name:       "name",
			},
			expectedLen: 1,
			isErr:       true,
		},
		{
			name: "ns, apiVersion, kind",
			key: Key{
				Namespace:  "default",
				APIVersion: "foo/v1",
				Kind:       "Kind",
			},
			expectedLen: 2,
		},
		{
			name: "ns, apiVersion",
			key: Key{
				Namespace:  "default",
				APIVersion: "foo/v1",
			},
			isErr: true,
		},
		{
			name: "ns",
			key: Key{
				Namespace: "default",
			}, expectedLen: 4,
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewMemoryCache()

			for _, obj := range genObjectsSeed() {
				err := c.Store(obj)
				require.NoError(t, err)
			}

			objs, err := c.List(tc.key)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Len(t, objs, tc.expectedLen)
		})
	}
}

func TestMemoryCache_Delete(t *testing.T) {
	c := NewMemoryCache()

	for _, obj := range genObjectsSeed() {
		err := c.Store(obj)
		require.NoError(t, err)
	}

	l := len(c.store)

	o := &unstructured.Unstructured{}
	o.SetNamespace("default")
	o.SetAPIVersion("foo/v1")
	o.SetKind("Kind")
	o.SetName("foo1")

	err := c.Delete(o)
	require.NoError(t, err)

	assert.Equal(t, l-1, len(c.store))
}

func genObjectsSeed() []*unstructured.Unstructured {
	var objects []*unstructured.Unstructured

	type source struct {
		ns, apiVersion, kind, name string
	}

	sources := []source{
		{"app-1", "foo/v1", "Kind", "foo1"},
		{"default", "foo/v1", "Kind", "foo1"},
		{"default", "foo/v1", "Kind", "foo2"},
		{"default", "foo/v1", "Other", "other1"},
		{"default", "bar/v1", "Bar", "bar1"},
	}

	for _, src := range sources {
		o := &unstructured.Unstructured{}
		o.SetNamespace(src.ns)
		o.SetAPIVersion(src.apiVersion)
		o.SetKind(src.kind)
		o.SetName(src.name)

		objects = append(objects, o)
	}

	return objects
}
