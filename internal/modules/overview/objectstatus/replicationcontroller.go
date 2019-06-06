package objectstatus

import (
	"context"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/heptio/developer-dash/internal/objectstore"
	"github.com/heptio/developer-dash/pkg/view/component"
)

func replicationController(_ context.Context, object runtime.Object, _ objectstore.ObjectStore) (ObjectStatus, error) {
	if object == nil {
		return ObjectStatus{}, errors.Errorf("replication controller is nil")
	}

	rc := &corev1.ReplicationController{}

	if err := scheme.Scheme.Convert(object, rc, 0); err != nil {
		return ObjectStatus{}, errors.Wrap(err, "convert object to replication controller")
	}

	status := rc.Status

	switch {
	case status.ReadyReplicas != status.Replicas:
		return ObjectStatus{
			nodeStatus: component.NodeStatusWarning,
			Details:    []component.Component{component.NewText("Replication Controller pods are not ready")},
		}, nil
	default:
		return ObjectStatus{
			nodeStatus: component.NodeStatusOK,
			Details:    []component.Component{component.NewText("Replication Controller is OK")},
		}, nil
	}
}
