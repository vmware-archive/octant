package overview

import (
	"context"
	"testing"
	"time"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/overview/printer"
	"github.com/heptio/developer-dash/internal/view"
	"github.com/heptio/developer-dash/internal/view/component"

	"github.com/heptio/developer-dash/internal/cluster/fake"
	"github.com/heptio/developer-dash/internal/content"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
)

func TestListDescriber(t *testing.T) {
	thePath := "/"
	key := cache.Key{APIVersion: "v1", Kind: "Pod"}
	namespace := "default"
	fields := map[string]string{}

	c := newSpyCache()

	object := map[string]interface{}{
		"kind":       "Pod",
		"apiVersion": "v1",
		"metadata": map[string]interface{}{
			"name":              "name",
			"creationTimestamp": "2019-01-14T13:34:56+00:00",
		},
	}

	retrieveKey := cache.Key{Namespace: namespace, APIVersion: "v1", Kind: "Pod"}
	c.spyRetrieve(retrieveKey, []*unstructured.Unstructured{{Object: object}}, nil)

	listType := func() interface{} {
		return &corev1.PodList{}
	}

	objectType := func() interface{} {
		return &corev1.Pod{}
	}

	otf := func(namespace, prefix string, contents *[]content.Content) func(*metav1beta1.Table) error {
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
		Cache:   c,
		Fields:  fields,
		Printer: printer.NewResource(c),
	}

	ctx := context.Background()
	cResponse, err := d.Describe(ctx, "/path", namespace, clusterClient, options)
	require.NoError(t, err)

	list := component.NewList("", nil)

	tableCols := component.NewTableCols("Name", "Labels", "Age")
	table := component.NewTable("*v1.PodList", tableCols)
	table.Add(component.TableRow{
		"Age":    component.NewTimestamp(time.Unix(1547472896, 0)),
		"Labels": component.NewLabels(nil),
		"Name":   component.NewText("", "name"),
	})
	list.Add(table)

	expected := component.ContentResponse{
		ViewComponents: []component.ViewComponent{list},
	}

	assert.Equal(t, expected, cResponse)

	assert.True(t, c.isSatisfied(), "cache was not satisfied")
}

func TestObjectDescriber(t *testing.T) {
	thePath := "/"
	key := cache.Key{APIVersion: "v1", Kind: "Pod"}
	namespace := "default"
	fields := map[string]string{}

	c := newSpyCache()

	object := map[string]interface{}{
		"kind":       "Pod",
		"apiVersion": "v1",
		"metadata": map[string]interface{}{
			"name": "name",
		},
	}

	retrieveKey := cache.Key{Namespace: namespace, APIVersion: "v1", Kind: "Pod"}
	c.spyRetrieve(retrieveKey, []*unstructured.Unstructured{{Object: object}}, nil)

	objectType := func() interface{} {
		return &corev1.Pod{}
	}

	viewFac := func(string, string, clock.Clock) view.View {
		return newFakeView()
	}
	fn := DefaultLoader(key)
	sections := []ContentSection{
		{
			Title: "section 1",
			Views: []view.ViewFactory{viewFac},
		},
	}
	d := NewObjectDescriber(thePath, "object", fn, objectType, sections)

	scheme := runtime.NewScheme()
	objects := []runtime.Object{}
	clusterClient, err := fake.NewClient(scheme, resources, objects)
	require.NoError(t, err)

	options := DescriberOptions{
		Cache:   c,
		Fields:  fields,
		Printer: printer.NewResource(c),
	}

	ctx := context.Background()
	cResponse, err := d.Describe(ctx, "/path", namespace, clusterClient, options)
	require.NoError(t, err)

	expected := component.ContentResponse{
		Title: "object: name",
		ViewComponents: []component.ViewComponent{
			component.NewText("", "*v1.Pod"),
		},
	}
	assert.Equal(t, expected, cResponse)
	assert.True(t, c.isSatisfied())
}

func TestSectionDescriber(t *testing.T) {
	namespace := "default"

	c := cache.NewMemoryCache()

	scheme := runtime.NewScheme()
	objects := []runtime.Object{}
	clusterClient, err := fake.NewClient(scheme, resources, objects)
	require.NoError(t, err)

	options := DescriberOptions{
		Cache: c,
	}

	ctx := context.Background()

	tests := []struct {
		name     string
		d        *SectionDescriber
		expected component.ContentResponse
	}{
		{
			name: "general",
			d: NewSectionDescriber(
				"/section",
				"section",
				newStubDescriber("/foo"),
			),
			expected: component.ContentResponse{
				Title: "section",
				ViewComponents: []component.ViewComponent{
					component.NewList("", nil),
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
			expected: component.ContentResponse{
				Title: "section",
				ViewComponents: []component.ViewComponent{
					component.NewList("", nil),
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
