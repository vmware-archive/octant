package objectstore

import (
	"context"

	"github.com/heptio/developer-dash/pkg/objectstoreutil"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	kcache "k8s.io/client-go/tools/cache"
)

//go:generate mockgen -source=objectstore.go -destination=./fake/mock_objectstore.go -package=fake github.com/heptio/developer-dash/internal/objectstore ObjectStore

// ObjectStore stores Kubernetes objects.
type ObjectStore interface {
	List(ctx context.Context, key objectstoreutil.Key) ([]*unstructured.Unstructured, error)
	Get(ctx context.Context, key objectstoreutil.Key) (*unstructured.Unstructured, error)
	Watch(key objectstoreutil.Key, handler kcache.ResourceEventHandler) error
}

// GetAs gets an object from the objectstore by key.
func GetAs(ctx context.Context, o ObjectStore, key objectstoreutil.Key, as interface{}) error {
	u, err := o.Get(ctx, key)
	if err != nil {
		return errors.Wrap(err, "get object from objectstore")
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
