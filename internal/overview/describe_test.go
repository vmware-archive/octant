package overview

import (
	"testing"

	"github.com/heptio/developer-dash/internal/cluster/fake"
	"github.com/heptio/developer-dash/internal/content"
	"github.com/heptio/developer-dash/internal/view"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
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
	clusterClient, err := fake.NewClient(scheme, objects)
	require.NoError(t, err)

	options := DescriberOptions{
		Cache:  cache,
		Fields: fields,
	}

	contents, _, err := d.Describe("/path", namespace, clusterClient, options)
	require.NoError(t, err)

	expected := []content.Content{theContent}

	assert.Equal(t, expected, contents)

	assert.True(t, cache.isSatisfied())
}

func TestObjectDescriber(t *testing.T) {
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

	d := NewObjectDescriber(thePath, "object", key, objectType, otf, []view.View{})

	scheme := runtime.NewScheme()
	objects := []runtime.Object{}
	clusterClient, err := fake.NewClient(scheme, objects)
	require.NoError(t, err)

	options := DescriberOptions{
		Cache:  cache,
		Fields: fields,
	}

	contents, title, err := d.Describe("/path", namespace, clusterClient, options)
	require.NoError(t, err)
	require.Len(t, contents, 2)
	assert.Equal(t, title, "object: name")

	expected := theContent
	assert.Equal(t, expected, contents[0])
	assert.True(t, cache.isSatisfied())
}

func TestSectionDescriber(t *testing.T) {
	namespace := "default"

	d := NewSectionDescriber(
		"/section",
		"section",
		newStubDescriber("/foo"),
	)

	cache := NewMemoryCache()

	scheme := runtime.NewScheme()
	objects := []runtime.Object{}
	clusterClient, err := fake.NewClient(scheme, objects)
	require.NoError(t, err)

	options := DescriberOptions{
		Cache: cache,
	}

	got, title, err := d.Describe("/prefix", namespace, clusterClient, options)
	require.NoError(t, err)

	assert.Equal(t, stubbedContent, got)
	assert.Equal(t, title, "section")
}
