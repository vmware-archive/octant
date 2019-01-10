package overview

import (
	"context"
	"github.com/heptio/developer-dash/internal/cache"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/client-go/kubernetes/scheme"
)

func createObject(name string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": map[string]interface{}{
				"name": name,
			},
		},
	}
}

func Test_loadObjects(t *testing.T) {
	sampleObjects := []*unstructured.Unstructured{
		createObject("omega"),
		createObject("alpha"),
	}
	sortedSampleObjects := []*unstructured.Unstructured{
		createObject("alpha"),
		createObject("omega"),
	}

	cases := []struct {
		name      string
		initCache func(*spyCache)
		fields    map[string]string
		keys      []cache.CacheKey
		expected  []*unstructured.Unstructured
		isErr     bool
	}{
		{
			name: "without name",
			initCache: func(c *spyCache) {
				c.spyRetrieve(cache.CacheKey{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "kind"},
					sampleObjects, nil)
			},
			fields:   map[string]string{},
			keys:     []cache.CacheKey{{APIVersion: "v1", Kind: "kind"}},
			expected: sortedSampleObjects,
		},
		{
			name: "name",
			initCache: func(c *spyCache) {
				c.spyRetrieve(cache.CacheKey{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "kind",
					Name:       "name"},
					[]*unstructured.Unstructured{}, nil)
			},
			fields: map[string]string{"name": "name"},
			keys:   []cache.CacheKey{{APIVersion: "v1", Kind: "kind"}},
		},
		{
			name: "cache retrieve error",
			initCache: func(c *spyCache) {
				c.spyRetrieve(cache.CacheKey{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "kind"},
					nil, errors.New("error"))
			},
			fields: map[string]string{},
			keys:   []cache.CacheKey{{APIVersion: "v1", Kind: "kind"}},
			isErr:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cache := newSpyCache()
			if tc.initCache != nil {
				tc.initCache(cache)
			}

			namespace := "default"

			ctx := context.Background()
			got, err := loadObjects(ctx, cache, namespace, tc.fields, tc.keys)
			if tc.isErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.True(t, cache.isSatisfied())
			assert.Equal(t, tc.expected, got)
		})
	}

}

func Test_translateTimestamp(t *testing.T) {
	ti := time.Unix(1538828130, 0)
	c := clock.NewFakeClock(ti)

	cases := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "zero",
			expected: "<unknown>",
		},
		{
			name:     "not zero",
			time:     time.Unix(1538828100, 0),
			expected: "30s",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ts := metav1.NewTime(tc.time)

			got := translateTimestamp(ts, c)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func loadType(t *testing.T, path string) runtime.Object {
	data, err := ioutil.ReadFile(filepath.Join("testdata", path))
	require.NoError(t, err)

	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode(data, nil, nil)
	require.NoError(t, err)

	return obj
}

func loadUnstructured(t *testing.T, cache cache.Cache, namespace, path string) {
	obj := loadType(t, path)
	m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	require.NoError(t, err)

	u := &unstructured.Unstructured{
		Object: m,
	}
	u.Object = m
	u.SetNamespace(namespace)

	err = cache.Store(u)
	require.NoError(t, err)
}
