package objectstatus

import (
	"github.com/heptio/developer-dash/internal/view/component"
	"k8s.io/apimachinery/pkg/runtime"
)

type ObjectStatus struct {
	NodeStatus component.NodeStatus
}

func Status(object runtime.Object) ObjectStatus {
	gvk := object.GetObjectKind().GroupVersionKind()
	apiVersion, kind := gvk.ToAPIVersionAndKind()

	switch {
	case apiVersion == "apps/v1" && kind == "Deployment":
		return ObjectStatus{
			NodeStatus: component.NodeStatusError,
		}
	case apiVersion == "extensions/v1beta1" && kind == "ReplicaSet":
		return ObjectStatus{
			NodeStatus: component.NodeStatusWarning,
		}
	}

	return ObjectStatus{
		NodeStatus: component.NodeStatusOK,
	}
}
