package objectstatus

import (
	"context"
	"errors"
	"fmt"

	"github.com/heptio/developer-dash/internal/objectstore"
	"github.com/heptio/developer-dash/pkg/view/component"
	"k8s.io/apimachinery/pkg/runtime"
)

type statusKey struct {
	apiVersion string
	kind       string
}

type statusFunc func(context.Context, runtime.Object, objectstore.ObjectStore) (ObjectStatus, error)

type statusLookup map[statusKey]statusFunc

var (
	defaultStatusLookup = statusLookup{
		{apiVersion: "apps/v1", kind: "DaemonSet"}:             daemonSet,
		{apiVersion: "apps/v1", kind: "Deployment"}:            deploymentAppsV1,
		{apiVersion: "apps/v1", kind: "ReplicaSet"}:            replicaSetAppsV1,
		{apiVersion: "apps/v1", kind: "StatefulSet"}:           statefulSet,
		{apiVersion: "batch/v1", kind: "Job"}:                  runJobStatus,
		{apiVersion: "v1", kind: "ReplicationController"}:      replicationController,
		{apiVersion: "v1", kind: "Service"}:                    service,
		{apiVersion: "extensions/v1beta1", kind: "Ingress"}:    runIngressStatus,
		{apiVersion: "extensions/v1beta1", kind: "ReplicaSet"}: replicaSetExtV1Beta1,
	}
)

type ObjectStatus struct {
	nodeStatus component.NodeStatus
	Details    []component.TitleComponent
}

func (os *ObjectStatus) AddDetail(detail string) {
	os.Details = append(os.Details, component.NewText(detail))
}

func (os *ObjectStatus) AddDetailf(msg string, args ...interface{}) {
	os.AddDetail(fmt.Sprintf(msg, args...))
}

func (os *ObjectStatus) SetError() {
	os.nodeStatus = component.NodeStatusError
}

func (os *ObjectStatus) SetWarning() {
	if os.nodeStatus != component.NodeStatusError {
		os.nodeStatus = component.NodeStatusWarning
	}
}

func (os *ObjectStatus) Status() component.NodeStatus {
	switch os.nodeStatus {
	case component.NodeStatusWarning,
		component.NodeStatusError,
		component.NodeStatusOK:
		return os.nodeStatus
	default:
		return component.NodeStatusOK
	}
}

// Status creates an ObjectStatus for an object.
func Status(ctx context.Context, object runtime.Object, o objectstore.ObjectStore) (ObjectStatus, error) {
	return status(ctx, object, o, defaultStatusLookup)
}

func status(ctx context.Context, object runtime.Object, o objectstore.ObjectStore, lookup statusLookup) (ObjectStatus, error) {
	if object == nil {
		return ObjectStatus{}, errors.New("object is nil")
	}

	if lookup == nil {
		return ObjectStatus{}, errors.New("status lookup is nil")
	}

	gvk := object.GetObjectKind().GroupVersionKind()
	apiVersion, kind := gvk.ToAPIVersionAndKind()

	fn, ok := lookup[statusKey{apiVersion: apiVersion, kind: kind}]
	if !ok {
		return ObjectStatus{
			nodeStatus: component.NodeStatusOK,
			Details:    component.Title(component.NewText(fmt.Sprintf("%s %s is OK", apiVersion, kind))),
		}, nil

	}

	return fn(ctx, object, o)
}
