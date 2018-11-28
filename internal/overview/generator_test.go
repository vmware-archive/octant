package overview

import (
	"context"
	"testing"

	"github.com/heptio/developer-dash/internal/queryer"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/view/component"

	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/cluster/fake"
	"github.com/heptio/developer-dash/internal/content"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func Test_realGenerator_Generate(t *testing.T) {
	key := cache.Key{Namespace: "default"}

	textOther := component.NewText("other")
	textFoo := component.NewText("foo")
	textSub := component.NewText("sub")

	describers := []Describer{
		newStubDescriber("/other", textOther),
		newStubDescriber("/foo", textFoo),
		newStubDescriber("/sub/(?P<name>.*?)", textSub),
	}

	var pathFilters []pathFilter
	for _, d := range describers {
		pathFilters = append(pathFilters, d.PathFilters("default")...)
	}

	cases := []struct {
		name      string
		path      string
		initCache func(*spyCache)
		expected  component.ContentResponse
		isErr     bool
	}{
		{
			name: "dynamic content",
			path: "/foo",
			initCache: func(c *spyCache) {
				c.spyRetrieve(key, []*unstructured.Unstructured{}, nil)
			},
			expected: component.ContentResponse{ViewComponents: []component.ViewComponent{textFoo}},
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
				subKey := cache.Key{Namespace: key.Namespace, Name: "foo"}
				c.spyRetrieve(subKey, []*unstructured.Unstructured{}, nil)
			},
			expected: component.ContentResponse{
				ViewComponents: []component.ViewComponent{textSub},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := newSpyCache()
			if tc.initCache != nil {
				tc.initCache(c)
			}

			q := queryer.New(c, nil)

			scheme := runtime.NewScheme()
			objects := []runtime.Object{}
			clusterClient, err := fake.NewClient(scheme, resources, objects)
			require.NoError(t, err)

			g, err := newGenerator(c, q, pathFilters, clusterClient, nil)
			require.NoError(t, err)

			ctx := context.Background()
			cResponse, err := g.Generate(ctx, tc.path, "/prefix", "default", GeneratorOptions{})
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
	store map[cache.Key][]*unstructured.Unstructured
	errs  map[cache.Key]error
	used  map[cache.Key]bool
}

func newSpyCache() *spyCache {
	return &spyCache{
		store: make(map[cache.Key][]*unstructured.Unstructured),
		errs:  make(map[cache.Key]error),
		used:  make(map[cache.Key]bool),
	}
}

func (c *spyCache) Store(obj *unstructured.Unstructured) error {
	return nil
}

func (c *spyCache) spyRetrieve(key cache.Key, objects []*unstructured.Unstructured, err error) {
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

func (c *spyCache) List(key cache.Key) ([]*unstructured.Unstructured, error) {
	c.used[key] = true

	objs := c.store[key]
	err := c.errs[key]

	return objs, err
}

func (c *spyCache) Get(key cache.Key) (*unstructured.Unstructured, error) {
	c.used[key] = true

	var obj *unstructured.Unstructured
	objs := c.store[key]
	if len(objs) > 0 {
		obj = objs[0]
	}
	err := c.errs[key]

	return obj, err
}

func (c *spyCache) Delete(obj *unstructured.Unstructured) error {
	return nil
}

type stubDescriber struct {
	path       string
	components []component.ViewComponent
}

func newStubDescriber(p string, components ...component.ViewComponent) *stubDescriber {
	return &stubDescriber{
		path:       p,
		components: components,
	}
}

func newEmptyDescriber(p string) *stubDescriber {
	return &stubDescriber{
		path: p,
	}
}

func (d *stubDescriber) Describe(context.Context, string, string, cluster.ClientInterface, DescriberOptions) (component.ContentResponse, error) {
	return component.ContentResponse{
		ViewComponents: d.components,
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

func (c *fakeContent) ViewComponent() content.ViewComponent {
	return content.ViewComponent{}
}

func (c fakeContent) MarshalJSON() ([]byte, error) {
	return []byte(`{"type":"stubbed"}`), nil
}
