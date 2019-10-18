package resourceviewer

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware-tanzu/octant/pkg/view/component"
)

type podGroupNode struct {
	objectStatus ObjectStatus
}

func (pgn *podGroupNode) Create(ctx context.Context, podGroupName string, objects []unstructured.Unstructured) (*component.Node, error) {
	podStatus := component.NewPodStatus()

	for _, object := range objects {
		if !isObjectPod(&object) {
			continue
		}
		pod, err := convertObjectToPod(&object)
		if err != nil {
			return nil, err
		}

		status, err := pgn.objectStatus.Status(ctx, &object)
		if err != nil {
			return nil, err
		}

		podStatus.AddSummary(pod.Name, status.Details, status.Status())
	}

	node := &component.Node{
		Name:       podGroupName,
		APIVersion: "v1",
		Kind:       "Pod",
		Status:     podStatus.Status(),
		Details:    []component.Component{podStatus},
	}
	return node, nil
}
