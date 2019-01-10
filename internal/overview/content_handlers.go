package overview

import (
	"context"
	"sort"

	"github.com/heptio/developer-dash/internal/cache"

	"github.com/heptio/developer-dash/internal/content"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/kubernetes/staging/src/k8s.io/apimachinery/pkg/util/duration"
)

type lookupFunc func(namespace, prefix string, cell interface{}) content.Text

// loadObjects loads objects from the cache sorted by their name.
func loadObjects(ctx context.Context, cache cache.Cache, namespace string, fields map[string]string, cacheKeys []cache.Key) ([]*unstructured.Unstructured, error) {
	var objects []*unstructured.Unstructured

	for _, cacheKey := range cacheKeys {
		cacheKey.Namespace = namespace

		if name, ok := fields["name"]; ok && name != "" {
			cacheKey.Name = name
		}

		cacheObjects, err := cache.Retrieve(cacheKey)
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

func eventsForObject(object *unstructured.Unstructured, cache cache.Cache, prefix, namespace string, cl clock.Clock) (*content.Table, error) {
	eventObjects, err := cache.Events(object)
	if err != nil {
		return nil, err
	}

	eventsTable := newEventTable(namespace, object)
	for _, obj := range eventObjects {
		event := &corev1.Event{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Object, event)
		if err != nil {
			return nil, err
		}

		eventsTable.AddRow(printEvent(event, prefix, namespace, cl))
	}

	return &eventsTable, nil
}
