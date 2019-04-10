package overview

import (
	"context"
	"sort"

	"github.com/heptio/developer-dash/internal/objectstore"
	"github.com/heptio/developer-dash/pkg/objectstoreutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/kubernetes/staging/src/k8s.io/apimachinery/pkg/util/duration"
)

// loadObject loads a single object from the objectstore.
func loadObject(ctx context.Context, objectstore objectstore.ObjectStore, namespace string, fields map[string]string, objectStoreKey objectstoreutil.Key) (*unstructured.Unstructured, error) {
	objectStoreKey.Namespace = namespace

	if name, ok := fields["name"]; ok && name != "" {
		objectStoreKey.Name = name
	}

	object, err := objectstore.Get(ctx, objectStoreKey)
	if err != nil {
		return nil, err
	}

	return object, nil
}

// loadObjects loads objects from the objectstore sorted by their name.
func loadObjects(ctx context.Context, objectstore objectstore.ObjectStore, namespace string, fields map[string]string, objectStoreKeys []objectstoreutil.Key) ([]*unstructured.Unstructured, error) {
	var objects []*unstructured.Unstructured

	for _, objectStoreKey := range objectStoreKeys {
		objectStoreKey.Namespace = namespace

		if name, ok := fields["name"]; ok && name != "" {
			objectStoreKey.Name = name
		}

		storedObjects, err := objectstore.List(ctx, objectStoreKey)
		if err != nil {
			return nil, err
		}

		objects = append(objects, storedObjects...)
	}

	sort.SliceStable(objects, func(i, j int) bool {
		a, b := objects[i], objects[j]
		return a.GetName() < b.GetName()
	})

	return objects, nil
}

// translateTimestamp returns the elapsed time since timestamp in
// human-readable approximation.
func translateTimestamp(timestamp metav1.Time, c clock.Clock) string {
	if timestamp.IsZero() {
		return "<unknown>"
	}

	return duration.ShortHumanDuration(c.Since(timestamp.Time))
}
