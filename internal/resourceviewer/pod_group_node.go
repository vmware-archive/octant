package resourceviewer

import (
	"context"
	"sort"

	"github.com/vmware-tanzu/octant/internal/link"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware-tanzu/octant/pkg/view/component"
)

type podGroupNode struct {
	objectStatus ObjectStatus
}

func (pgn *podGroupNode) Create(ctx context.Context, podGroupName string, objects []unstructured.Unstructured, link link.Interface) (*component.Node, error) {
	podStatus := component.NewPodStatus()
	var podProperties []component.Property

	sort.Slice(objects, func(i, j int) bool {
		return objects[i].GetName() < objects[j].GetName()
	})

	for i, object := range objects {
		if !isObjectPod(&object) {
			continue
		}
		pod, err := convertObjectToPod(&object)
		if err != nil {
			return nil, err
		}

		status, err := pgn.objectStatus.Status(ctx, &object, link)
		if err != nil {
			return nil, err
		}

		podStatus.AddSummary(pod.Name, status.Details, status.Properties, status.Status())
		if i == 0 { // we add Properties that are common for all pods and only once
			podProperties = append(podProperties, status.Properties...)
		}
	}

	node := &component.Node{
		Name:       podGroupName,
		APIVersion: "v1",
		Kind:       "Pod",
		Status:     podStatus.Status(),
		Properties: podProperties,
		Details:    []component.Component{podStatus},
	}
	return node, nil
}
