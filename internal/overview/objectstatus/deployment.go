package objectstatus

import (
	"fmt"

	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

// deploymentAppsV1 creates status for an v1/apps deployment. This is
// not the final implementation. It is included to generate output.
func deploymentAppsV1(object runtime.Object) (ObjectStatus, error) {
	if object == nil {
		return ObjectStatus{}, errors.Errorf("deployment is nil")
	}

	deployment := &appsv1.Deployment{}

	if err := scheme.Scheme.Convert(object, deployment, 0); err != nil {
		return ObjectStatus{}, errors.Wrap(err, "convert object to apps/v1 deployment")
	}

	status := deployment.Status

	switch {
	case status.Replicas == status.UnavailableReplicas:
		return ObjectStatus{
			NodeStatus: component.NodeStatusError,
		}, nil
	case status.Replicas == status.AvailableReplicas:
		return ObjectStatus{
			NodeStatus: component.NodeStatusOK,
			Details:    component.Title(component.NewText("Deployment is OK"))}, nil
	default:
		return ObjectStatus{
			NodeStatus: component.NodeStatusWarning,
			Details: component.Title(
				component.NewText(
					fmt.Sprintf("Expected %d replicas, but %d are available",
						status.Replicas, status.AvailableReplicas)))}, nil
	}
}
