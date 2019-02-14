package objectstatus

import (
	"errors"
	"fmt"

	"github.com/heptio/developer-dash/internal/view/component"
	"k8s.io/apimachinery/pkg/runtime"
)

type statusKey struct {
	apiVersion string
	kind       string
}

type statusFunc func(runtime.Object) (ObjectStatus, error)

type statusLookup map[statusKey]statusFunc

var (
	defaultStatusLookup = statusLookup{
		{apiVersion: "apps/v1", kind: "Deployment"}:            deploymentAppsV1,
		{apiVersion: "extensions/v1beta1", kind: "ReplicaSet"}: replicaSetExtV1Beta1,
		{apiVersion: "apps/v1", kind: "ReplicaSet"}:            replicaSetAppsV1,
	}
)

type ObjectStatus struct {
	NodeStatus component.NodeStatus
	Details    []component.TitleViewComponent
}

// Status creates an ObjectStatus for an object.
func Status(object runtime.Object) (ObjectStatus, error) {
	return status(object, defaultStatusLookup)
}

func status(object runtime.Object, lookup statusLookup) (ObjectStatus, error) {
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
			NodeStatus: component.NodeStatusOK,
			Details:    component.Title(component.NewText(fmt.Sprintf("%s %s is OK", apiVersion, kind))),
		}, nil

	}

	return fn(object)
}
