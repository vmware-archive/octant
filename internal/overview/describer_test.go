package overview

import (
	"context"
	"testing"

	"github.com/heptio/developer-dash/internal/cluster/fake"
	"github.com/heptio/developer-dash/internal/content"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/kubernetes/pkg/apis/core"
)

func TestListDescriber(t *testing.T) {
	thePath := "/"
	key := CacheKey{APIVersion: "v1", Kind: "kind"}
	namespace := "default"
	fields := map[string]string{}

	cache := newSpyCache()

	object := map[string]interface{}{
		"kind":       "kind",
		"apiVersion": "v1",
		"metadata": map[string]interface{}{
			"name": "name",
		},
	}

	retrieveKey := CacheKey{Namespace: namespace, APIVersion: "v1", Kind: "kind"}
	cache.spyRetrieve(retrieveKey, []*unstructured.Unstructured{{Object: object}}, nil)

	listType := func() interface{} {
		return &core.EventList{}
	}

	objectType := func() interface{} {
		return &core.Event{}
	}

	theContent := newFakeContent(false)

	otf := func(namespace, prefix string, contents *[]content.Content) func(*metav1beta1.Table) error {
		*contents = append(*contents, theContent)
		return func(*metav1beta1.Table) error {
			return nil
		}
	}

	d := NewListDescriber(thePath, "list", key, listType, objectType, otf)

	scheme := runtime.NewScheme()
	objects := []runtime.Object{}
	clusterClient, err := fake.NewClient(scheme, resources, objects)
	require.NoError(t, err)

	options := DescriberOptions{
		Cache:  cache,
		Fields: fields,
	}

	ctx := context.Background()
	cResponse, err := d.Describe(ctx, "/path", namespace, clusterClient, options)
	require.NoError(t, err)

	expected := ContentResponse{
		Views: []Content{
			{
				Contents: stubbedContent,
			},
		},
		ViewComponents: []content.ViewComponent{
			{
				Metadata: content.Metadata{
					Type:  "list",
					Title: "list",
				},
				Config: ListConfig{
					Items: []content.ViewComponent{
						{},
					},
				},
			},
		},
	}

	assert.Equal(t, expected, cResponse)

	assert.True(t, cache.isSatisfied())
}

func TestObjectDescriber(t *testing.T) {
	thePath := "/"
	key := CacheKey{APIVersion: "v1", Kind: "Pod"}
	namespace := "default"
	fields := map[string]string{}

	cache := newSpyCache()

	object := map[string]interface{}{
		"kind":       "Pod",
		"apiVersion": "v1",
		"metadata": map[string]interface{}{
			"name": "name",
		},
	}

	retrieveKey := CacheKey{Namespace: namespace, APIVersion: "v1", Kind: "Pod"}
	cache.spyRetrieve(retrieveKey, []*unstructured.Unstructured{{Object: object}}, nil)

	objectType := func() interface{} {
		return &core.Pod{}
	}

	viewFac := func(string, string, clock.Clock) View {
		return newFakeView()
	}
	fn := DefaultLoader(key)
	sections := []ContentSection{
		{
			Title: "section 1",
			Views: []ViewFactory{viewFac},
		},
	}
	d := NewObjectDescriber(thePath, "object", fn, objectType, sections)

	scheme := runtime.NewScheme()
	objects := []runtime.Object{}
	clusterClient, err := fake.NewClient(scheme, resources, objects)
	require.NoError(t, err)

	options := DescriberOptions{
		Cache:  cache,
		Fields: fields,
	}

	ctx := context.Background()
	cResponse, err := d.Describe(ctx, "/path", namespace, clusterClient, options)
	require.NoError(t, err)

	expected := ContentResponse{
		Title: "object: name",
		Views: []Content{
			{
				Contents: stubbedContent,
				Title:    "section 1",
			},
		},
	}
	assert.Equal(t, expected, cResponse)
	assert.True(t, cache.isSatisfied())
}

func TestSectionDescriber(t *testing.T) {
	namespace := "default"

	cache := NewMemoryCache()

	scheme := runtime.NewScheme()
	objects := []runtime.Object{}
	clusterClient, err := fake.NewClient(scheme, resources, objects)
	require.NoError(t, err)

	options := DescriberOptions{
		Cache: cache,
	}

	ctx := context.Background()

	tests := []struct {
		name     string
		d        *SectionDescriber
		expected ContentResponse
	}{
		{
			name: "general",
			d: NewSectionDescriber(
				"/section",
				"section",
				newStubDescriber("/foo"),
			),
			expected: ContentResponse{
				Views: []Content{
					{
						Contents: stubbedContent,
						Title:    "section",
					},
				},
				ViewComponents: []content.ViewComponent{
					{
						Metadata: content.Metadata{
							Type:  "list",
							Title: "section",
						},
						Config: ListConfig{
							Items: []content.ViewComponent{
								{},
							},
						},
					},
				},
			},
		},
		{
			name: "empty",
			d: NewSectionDescriber(
				"/section",
				"section",
				newEmptyDescriber("/foo"),
				newEmptyDescriber("/bar"),
			),
			expected: ContentResponse{
				Views: []Content{
					{
						Title: "section",
						Contents: []content.Content{
							&content.Table{
								Type:         "table",
								Title:        "section",
								EmptyContent: "Namespace default does not have any resources of this type",
							},
						},
					},
				},
				ViewComponents: []content.ViewComponent{
					{
						Metadata: content.Metadata{
							Type:  "list",
							Title: "section",
						},
						Config: ListConfig{
							Items: []content.ViewComponent{},
						},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			got, err := tc.d.Describe(ctx, "/prefix", namespace, clusterClient, options)
			require.NoError(t, err)

			assert.Equal(t, tc.expected, got)
		})
	}

}
