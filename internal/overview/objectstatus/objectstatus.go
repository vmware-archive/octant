package objectstatus

import (
	"fmt"

	"github.com/heptio/developer-dash/internal/view/component"
	"k8s.io/apimachinery/pkg/runtime"
)

type ObjectStatus struct {
	NodeStatus component.NodeStatus
	Details    []component.TitleViewComponent
}

func Status(object runtime.Object) (ObjectStatus, error) {
	gvk := object.GetObjectKind().GroupVersionKind()
	apiVersion, kind := gvk.ToAPIVersionAndKind()

	switch {
	case apiVersion == "apps/v1" && kind == "Deployment":
		return deploymentAppsV1(object)
	case apiVersion == "extensions/v1beta1" && kind == "ReplicaSet":
		return replicaSetExtV1Beta1(object)
	case apiVersion == "apps/v1" && kind == "ReplicaSet":
		return replicaSetAppsV1Beta1(object)
	default:
		return ObjectStatus{
			NodeStatus: component.NodeStatusOK,
			Details:    component.Title(component.NewText(fmt.Sprintf("%s %s is OK", apiVersion, kind))),
		}, nil
	}
}

func PodGroupStatus() {

}
