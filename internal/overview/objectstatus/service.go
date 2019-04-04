package objectstatus

import (
	"context"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/pkg/cacheutil"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

func service(ctx context.Context, object runtime.Object, c cache.Cache) (ObjectStatus, error) {
	if object == nil {
		return ObjectStatus{}, errors.Errorf("service is nil")
	}

	service := &corev1.Service{}

	if err := scheme.Scheme.Convert(object, service, 0); err != nil {
		return ObjectStatus{}, errors.Wrap(err, "convert object to service")
	}

	key := cacheutil.Key{
		Namespace:  service.Namespace,
		APIVersion: "v1",
		Kind:       "Endpoints",
		Name:       service.Name,
	}

	endpoints := &corev1.Endpoints{}

	if err := cache.GetAs(ctx, c, key, endpoints); err != nil {
		return ObjectStatus{}, errors.Wrapf(err, "get endpoints for service %s", service.Name)
	}

	addressCount := 0

	for _, subset := range endpoints.Subsets {
		addressCount += len(subset.Addresses)
	}

	if addressCount == 0 {
		return ObjectStatus{
			nodeStatus: component.NodeStatusWarning,
			Details:    component.TitleFromString("Service has no endpoints"),
		}, nil
	}

	return ObjectStatus{
		nodeStatus: component.NodeStatusOK,
		Details:    component.TitleFromString("Service is OK"),
	}, nil
}
