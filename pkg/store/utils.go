package store

import (
	"context"
	"sort"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/internal/util/kubernetes"
)

func EventsForObject(ctx context.Context, object runtime.Object, o Store) (*corev1.EventList, error) {
	accessor := meta.NewAccessor()
	namespace, err := accessor.Namespace(object)
	if err != nil {
		return nil, errors.Wrap(err, "get namespace for object")
	}

	apiVersion, err := accessor.APIVersion(object)
	if err != nil {
		return nil, errors.Wrap(err, "Get apiVersion for object")
	}

	kind, err := accessor.Kind(object)
	if err != nil {
		return nil, errors.Wrap(err, "get kind for object")
	}

	name, err := accessor.Name(object)
	if err != nil {
		return nil, errors.Wrap(err, "get name for object")
	}

	key := Key{
		Namespace:  namespace,
		APIVersion: "v1",
		Kind:       "Event",
	}

	list, _, err := o.List(ctx, key)
	if err != nil {
		return nil, errors.Wrap(err, "list events for object")
	}

	eventList := &corev1.EventList{}

	for _, unstructuredEvent := range list.Items {
		event := &corev1.Event{}
		err := kubernetes.FromUnstructured(&unstructuredEvent, event)
		if err != nil {
			return nil, err
		}

		involvedObject := event.InvolvedObject
		if involvedObject.APIVersion == "autoscaling/v2beta2" || involvedObject.APIVersion == "autoscaling/v2beta1" {
			involvedObject.APIVersion = "autoscaling/v1"
		}

		if involvedObject.Namespace == namespace &&
			involvedObject.APIVersion == apiVersion &&
			involvedObject.Kind == kind &&
			involvedObject.Name == name {
			eventList.Items = append(eventList.Items, *event)
		}
	}

	sort.SliceStable(eventList.Items, func(i, j int) bool {
		a := eventList.Items[i]
		b := eventList.Items[j]

		if b.LastTimestamp.After(a.LastTimestamp.Time) {
			return true
		}

		return a.LastTimestamp.After(b.LastTimestamp.Time)
	})

	return eventList, nil
}
