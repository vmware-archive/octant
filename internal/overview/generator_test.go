package overview

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Test_realGenerator_Generate(t *testing.T) {
	key := CacheKey{Namespace: "default"}

	describers := []Describer{
		newStubDescriber("/other"),
		newStubDescriber("/foo"),
		newStubDescriber("/sub/(?P<name>.*?)"),
	}

	var pathFilters []pathFilter
	for _, d := range describers {
		pathFilters = append(pathFilters, d.PathFilters()...)
	}

	cases := []struct {
		name      string
		path      string
		initCache func(*spyCache)
		expected  []Content
		isErr     bool
	}{
		{
			name: "dynamic content",
			path: "/foo",
			initCache: func(c *spyCache) {
				c.spyRetrieve(key, []*unstructured.Unstructured{}, nil)
			},
			expected: stubbedContent,
		},
		{
			name:     "stubbed content",
			path:     "/rbac/roles",
			expected: stubContent("/rbac/roles"),
		},
		{
			name:  "invalid path",
			path:  "/missing",
			isErr: true,
		},
		{
			name: "sub path",
			path: "/sub/foo",
			initCache: func(c *spyCache) {
				subKey := CacheKey{Namespace: key.Namespace, Name: "foo"}
				c.spyRetrieve(subKey, []*unstructured.Unstructured{}, nil)
			},
			expected: stubbedContent,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cache := newSpyCache()
			if tc.initCache != nil {
				tc.initCache(cache)
			}

			g := newGenerator(cache, pathFilters)

			contents, err := g.Generate(tc.path, "/prefix", "default")
			if tc.isErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, tc.expected, contents)
		})
	}
}

type spyCache struct {
	store map[CacheKey][]*unstructured.Unstructured
	errs  map[CacheKey]error
	used  map[CacheKey]bool
}

func newSpyCache() *spyCache {
	return &spyCache{
		store: make(map[CacheKey][]*unstructured.Unstructured),
		errs:  make(map[CacheKey]error),
		used:  make(map[CacheKey]bool),
	}
}

func (c *spyCache) Store(obj *unstructured.Unstructured) error {
	return nil
}

func (c *spyCache) spyRetrieve(key CacheKey, objects []*unstructured.Unstructured, err error) {
	c.store[key] = objects
	c.errs[key] = err
}

func (c *spyCache) isSatisfied() bool {
	for k := range c.store {
		isUsed, ok := c.used[k]
		if !ok {
			return false
		}

		if !isUsed {
			return false
		}
	}

	return true
}

func (c *spyCache) Retrieve(key CacheKey) ([]*unstructured.Unstructured, error) {
	c.used[key] = true

	objs := c.store[key]
	err := c.errs[key]

	return objs, err
}

func (c *spyCache) Delete(obj *unstructured.Unstructured) error {
	return nil
}

func (c *spyCache) Events(obj *unstructured.Unstructured) ([]*unstructured.Unstructured, error) {
	return nil, errors.New("Events not implemented")
}

type stubDescriber struct {
	path string
}

func newStubDescriber(p string) *stubDescriber {
	return &stubDescriber{
		path: p,
	}
}

func (d *stubDescriber) Describe(string, string, Cache, map[string]string) ([]Content, error) {
	return stubbedContent, nil
}

func (d *stubDescriber) PathFilters() []pathFilter {
	return []pathFilter{
		*newPathFilter(d.path, d),
	}
}

var stubbedContent = []Content{"content"}
