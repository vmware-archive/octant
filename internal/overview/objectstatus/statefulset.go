package objectstatus

import (
	"context"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

func statefulSet(_ context.Context, object runtime.Object, _ cache.Cache) (ObjectStatus, error) {
	if object == nil {
		return ObjectStatus{}, errors.Errorf("stateful set is nil")
	}

	ss := &appsv1.StatefulSet{}

	if err := scheme.Scheme.Convert(object, ss, 0); err != nil {
		return ObjectStatus{}, errors.Wrap(err, "convert object to stateful set")
	}

	status := ss.Status

	switch {
	case status.ReadyReplicas != status.Replicas:
		return ObjectStatus{
			nodeStatus: component.NodeStatusWarning,
			Details:    component.TitleFromString("Stateful Set pods are not ready"),
		}, nil
	default:
		return ObjectStatus{
			nodeStatus: component.NodeStatusOK,
			Details:    component.TitleFromString("Stateful Set is OK"),
		}, nil
	}
}
