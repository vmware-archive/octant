package overview

import (
	"fmt"
	"math/rand"
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

func TestMemoryCache_Retrieve(t *testing.T) {

	cases := []struct {
		name        string
		key         CacheKey
		expectedLen int
	}{
		{
			name: "ns, apiVersion, kind, name",
			key: CacheKey{
				Namespace:  "default",
				APIVersion: "foo/v1",
				Kind:       "Kind",
				Name:       "foo1",
			},
			expectedLen: 1,
		},
		{
			name: "ns, apiVersion, kind",
			key: CacheKey{
				Namespace:  "default",
				APIVersion: "foo/v1",
				Kind:       "Kind",
			},
			expectedLen: 2,
		},
		{
			name: "ns, apiVersion",
			key: CacheKey{
				Namespace:  "default",
				APIVersion: "foo/v1",
			},
			expectedLen: 3,
		},
		{
			name: "ns",
			key: CacheKey{
				Namespace: "default",
			}, expectedLen: 4,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewMemoryCache()

			for _, obj := range genObjectsSeed() {
				err := c.Store(obj)
				require.NoError(t, err)
			}

			objs, err := c.Retrieve(tc.key)
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

func TestMemoryCache_Events(t *testing.T) {
	cases := []struct {
		name         string
		eventFactory func(*unstructured.Unstructured) []*unstructured.Unstructured
		expected     int
	}{
		{
			name: "with matches",
			eventFactory: func(obj *unstructured.Unstructured) []*unstructured.Unstructured {
				return []*unstructured.Unstructured{
					genEvent(obj),
				}
			},
			expected: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewMemoryCache()

			o := genObject("test")
			err := c.Store(o)
			require.NoError(t, err)

			if tc.eventFactory != nil {
				for _, event := range tc.eventFactory(o) {
					err = c.Store(event)
					require.NoError(t, err)
				}
			}

			events, err := c.Events(o)
			require.NoError(t, err)

			assert.Len(t, events, tc.expected)
		})
	}

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

func genEvent(u *unstructured.Unstructured) *unstructured.Unstructured {
	o := &unstructured.Unstructured{}
	o.SetNamespace("default")
	o.SetAPIVersion("v1")
	o.SetKind("Event")
	o.SetName(fmt.Sprintf("event.%d", rand.Intn(100)))

	o.Object["involvedObject"] = map[string]interface{}{
		"apiVersion": u.GetAPIVersion(),
		"kind":       u.GetKind(),
		"name":       u.GetName(),
		"namespace":  "default",
	}

	return o
}

func genObject(name string) *unstructured.Unstructured {
	o := &unstructured.Unstructured{}
	o.SetNamespace("default")
	o.SetAPIVersion("foo/v1")
	o.SetKind("Kind")
	o.SetName(name)

	return o
}
