package cache

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kLabels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	kcache "k8s.io/client-go/tools/cache"
)

//go:generate mockgen -destination=./fake/mock_cache.go -package=fake github.com/heptio/developer-dash/internal/cache Cache

// Cache stores Kubernetes objects.
type Cache interface {
	List(key Key) ([]*unstructured.Unstructured, error)
	Get(key Key) (*unstructured.Unstructured, error)
	Watch(key Key, handler kcache.ResourceEventHandler) error
}

// Key is a key for the cache.
type Key struct {
	Namespace  string
	APIVersion string
	Kind       string
	Name       string
	Selector   kLabels.Selector
}

func (k Key) String() string {
	var sb strings.Builder

	sb.WriteString("CacheKey[")
	if k.Namespace != "" {
		fmt.Fprintf(&sb, "Namespace=%q, ", k.Namespace)
	}
	fmt.Fprintf(&sb, "APIVersion=%q, ", k.APIVersion)
	fmt.Fprintf(&sb, "Kind=%q", k.Kind)

	if k.Name != "" {
		fmt.Fprintf(&sb, ", Name=%q", k.Name)
	}

	if k.Selector != nil {
		fmt.Fprintf(&sb, ", Selector=%q", k.Selector.String())
	}

	sb.WriteString("]")

	return sb.String()
}

// GetAs gets an object from the cache by key.
func GetAs(c Cache, key Key, as interface{}) error {
	u, err := c.Get(key)
	if err != nil {
		return errors.Wrap(err, "get object from cache")
	}

	if u == nil {
		return nil
	}

	err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, as)
	if err != nil {
		return err
	}

	if err := copyObjectMeta(as, u); err != nil {
		return errors.Wrap(err, "copy object metadata")
	}

	return nil
}

// TODO: see if all the other versions of this function could be replaced
func copyObjectMeta(to interface{}, from *unstructured.Unstructured) error {
	object, ok := to.(metav1.Object)
	if !ok {
		return errors.Errorf("%T is not an object", to)
	}

	t, err := meta.TypeAccessor(object)
	if err != nil {
		return errors.Wrapf(err, "accessing type meta")
	}
	t.SetAPIVersion(from.GetAPIVersion())
	t.SetKind(from.GetObjectKind().GroupVersionKind().Kind)

	object.SetNamespace(from.GetNamespace())
	object.SetName(from.GetName())
	object.SetGenerateName(from.GetGenerateName())
	object.SetUID(from.GetUID())
	object.SetResourceVersion(from.GetResourceVersion())
	object.SetGeneration(from.GetGeneration())
	object.SetSelfLink(from.GetSelfLink())
	object.SetCreationTimestamp(from.GetCreationTimestamp())
	object.SetDeletionTimestamp(from.GetDeletionTimestamp())
	object.SetDeletionGracePeriodSeconds(from.GetDeletionGracePeriodSeconds())
	object.SetLabels(from.GetLabels())
	object.SetAnnotations(from.GetAnnotations())
	object.SetInitializers(from.GetInitializers())
	object.SetOwnerReferences(from.GetOwnerReferences())
	object.SetClusterName(from.GetClusterName())
	object.SetFinalizers(from.GetFinalizers())

	return nil
}
