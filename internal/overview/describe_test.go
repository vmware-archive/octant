package overview

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
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

	otf := func(namespace, prefix string, contents *[]Content) func(*metav1beta1.Table) error {
		*contents = append(*contents, theContent)
		return func(*metav1beta1.Table) error {
			return nil
		}
	}

	d := NewListDescriber(thePath, key, listType, objectType, otf)

	contents, err := d.Describe("/path", namespace, cache, fields)
	require.NoError(t, err)

	expected := []Content{theContent}

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

	otf := func(namespace, prefix string, contents *[]Content) func(*metav1beta1.Table) error {
		*contents = append(*contents, theContent)
		return func(*metav1beta1.Table) error {
			return nil
		}
	}

	d := NewObjectDescriber(thePath, key, objectType, otf)

	contents, err := d.Describe("/path", namespace, cache, fields)
	require.NoError(t, err)

	require.Len(t, contents, 2)

	expected := theContent
	assert.Equal(t, expected, contents[0])
	assert.True(t, cache.isSatisfied())
}

func TestSectionDescriber(t *testing.T) {
	namespace := "default"

	d := NewSectionDescriber(
		"/section",
		newStubDescriber("/foo"),
	)

	cache := NewMemoryCache()

	got, err := d.Describe("/prefix", namespace, cache, nil)
	require.NoError(t, err)

	assert.Equal(t, stubbedContent, got)
}
