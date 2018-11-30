package overview

import (
	"context"
	"testing"

	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/cluster/fake"
	"github.com/heptio/developer-dash/internal/content"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
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
		pathFilters = append(pathFilters, d.PathFilters("default")...)
	}

	cases := []struct {
		name      string
		path      string
		initCache func(*spyCache)
		expected  ContentResponse
		isErr     bool
	}{
		{
			name: "dynamic content",
			path: "/foo",
			initCache: func(c *spyCache) {
				c.spyRetrieve(key, []*unstructured.Unstructured{}, nil)
			},
			expected: ContentResponse{
				Views: []Content{
					{
						Contents: stubbedContent,
						Title:    "section content",
					},
				},
			},
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
			expected: ContentResponse{
				Views: []Content{
					{
						Contents: stubbedContent,
						Title:    "section content",
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cache := newSpyCache()
			if tc.initCache != nil {
				tc.initCache(cache)
			}

			scheme := runtime.NewScheme()
			objects := []runtime.Object{}
			clusterClient, err := fake.NewClient(scheme, resources, objects)
			require.NoError(t, err)

			g := newGenerator(cache, pathFilters, clusterClient)

			ctx := context.Background()
			cResponse, err := g.Generate(ctx, tc.path, "/prefix", "default")
			if tc.isErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, tc.expected, cResponse)
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
	return []*unstructured.Unstructured{}, nil
}

type stubDescriber struct {
	path     string
	contents []content.Content
}

func newStubDescriber(p string) *stubDescriber {
	return &stubDescriber{
		path:     p,
		contents: []content.Content{newFakeContent(false)},
	}
}

func newEmptyDescriber(p string) *stubDescriber {
	return &stubDescriber{
		path:     p,
		contents: []content.Content{newFakeContent(true)},
	}
}

func (d *stubDescriber) Describe(context.Context, string, string, cluster.ClientInterface, DescriberOptions) (ContentResponse, error) {
	return ContentResponse{
		Views: []Content{
			{
				Contents: d.contents,
				Title:    "section content",
			},
		},
	}, nil
}

func (d *stubDescriber) PathFilters(namespace string) []pathFilter {
	return []pathFilter{
		*newPathFilter(d.path, d),
	}
}

var stubbedContent = []content.Content{newFakeContent(false)}

type fakeContent struct {
	isEmpty bool
}

func newFakeContent(isEmpty bool) *fakeContent {
	return &fakeContent{
		isEmpty: isEmpty,
	}
}

func (c *fakeContent) IsEmpty() bool {
	return c.isEmpty
}

func (c fakeContent) MarshalJSON() ([]byte, error) {
	return []byte(`{"type":"stubbed"}`), nil
}

type fakeView struct{}

var _ View = (*fakeView)(nil)

func newFakeView() *fakeView {
	return &fakeView{}
}

func (v *fakeView) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	return stubbedContent, nil
}
