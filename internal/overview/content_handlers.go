package overview

import (
	"context"
	"sort"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/pkg/cacheutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/kubernetes/staging/src/k8s.io/apimachinery/pkg/util/duration"
)

// loadObject loads a single object from the cache.
func loadObject(ctx context.Context, cache cache.Cache, namespace string, fields map[string]string, cacheKey cacheutil.Key) (*unstructured.Unstructured, error) {
	cacheKey.Namespace = namespace

	if name, ok := fields["name"]; ok && name != "" {
		cacheKey.Name = name
	}

	object, err := cache.Get(ctx, cacheKey)
	if err != nil {
		return nil, err
	}

	return object, nil
}

// loadObjects loads objects from the cache sorted by their name.
func loadObjects(ctx context.Context, cache cache.Cache, namespace string, fields map[string]string, cacheKeys []cacheutil.Key) ([]*unstructured.Unstructured, error) {
	var objects []*unstructured.Unstructured

	for _, cacheKey := range cacheKeys {
		cacheKey.Namespace = namespace

		if name, ok := fields["name"]; ok && name != "" {
			cacheKey.Name = name
		}

		cacheObjects, err := cache.List(ctx, cacheKey)
		if err != nil {
			return nil, err
		}

		objects = append(objects, cacheObjects...)
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
