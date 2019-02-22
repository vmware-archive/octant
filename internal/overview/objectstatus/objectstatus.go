package objectstatus

import (
	"errors"
	"fmt"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/view/component"
	"k8s.io/apimachinery/pkg/runtime"
)

type statusKey struct {
	apiVersion string
	kind       string
}

type statusFunc func(runtime.Object, cache.Cache) (ObjectStatus, error)

type statusLookup map[statusKey]statusFunc

var (
	defaultStatusLookup = statusLookup{
		{apiVersion: "apps/v1", kind: "DaemonSet"}:             daemonSet,
		{apiVersion: "apps/v1", kind: "Deployment"}:            deploymentAppsV1,
		{apiVersion: "apps/v1", kind: "ReplicaSet"}:            replicaSetAppsV1,
		{apiVersion: "batch/v1", kind: "Job"}:                  runJobStatus,
		{apiVersion: "extensions/v1beta1", kind: "Ingress"}:    runIngressStatus,
		{apiVersion: "extensions/v1beta1", kind: "ReplicaSet"}: replicaSetExtV1Beta1,
	}
)

type ObjectStatus struct {
	nodeStatus component.NodeStatus
	Details    []component.TitleViewComponent
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
func Status(object runtime.Object, c cache.Cache) (ObjectStatus, error) {
	return status(object, c, defaultStatusLookup)
}

func status(object runtime.Object, c cache.Cache, lookup statusLookup) (ObjectStatus, error) {
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

	return fn(object, c)
}
