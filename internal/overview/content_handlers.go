package overview

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/kubernetes/staging/src/k8s.io/apimachinery/pkg/util/duration"
)

type lookupFunc func(namespace, prefix string, cell interface{}) text

func loadObjects(cache Cache, namespace string, fields map[string]string, cacheKeys []CacheKey) ([]*unstructured.Unstructured, error) {
	var objects []*unstructured.Unstructured

	for _, cacheKey := range cacheKeys {
		cacheKey.Namespace = namespace

		if name, ok := fields["name"]; ok && name != "" {
			cacheKey.Name = name
		}

		objs, err := cache.Retrieve(cacheKey)
		if err != nil {
			return nil, err
		}

		objects = append(objects, objs...)
	}

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

func eventsForObject(object *unstructured.Unstructured, cache Cache, prefix, namespace string, cl clock.Clock) (table, error) {
	eventObjects, err := cache.Events(object)
	if err != nil {
		return table{}, err
	}

	eventsTable := newEventTable("Events")
	for _, obj := range eventObjects {
		event := &corev1.Event{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Object, event)
		if err != nil {
			return table{}, err
		}

		eventsTable.AddRow(printEvent(event, prefix, namespace, cl))
	}

	return eventsTable, nil
}
