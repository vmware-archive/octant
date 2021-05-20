/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectstatus

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/vmware-tanzu/octant/internal/link"

	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func pod(ctx context.Context, object runtime.Object, o store.Store, link link.Interface) (ObjectStatus, error) {
	if object == nil {
		return ObjectStatus{}, errors.Errorf("pod is nil")
	}

	pod := &corev1.Pod{}

	if err := scheme.Scheme.Convert(object, pod, 0); err != nil {
		return ObjectStatus{}, errors.Wrap(err, "convert object to pod")
	}

	status := ObjectStatus{}

	switch pod.Status.Phase {
	case corev1.PodRunning:
		status.nodeStatus = component.NodeStatusOK
	case corev1.PodUnknown:
		status.nodeStatus = component.NodeStatusError
	default:
		status.nodeStatus = component.NodeStatusWarning
	}

	status.Details = []component.Component{
		component.NewText(pod.Status.Message),
	}

	if len(pod.Spec.EphemeralContainers) > 0 {
		status.nodeStatus = component.NodeStatusWarning
		status.Details = append(status.Details, component.NewText("Ephemeral container is running"))
	}

	if link != nil {
		serviceAccountLink, _ := link.ForGVK(pod.Namespace, "v1", "ServiceAccount", pod.Spec.ServiceAccountName, pod.Spec.ServiceAccountName)
		status.Properties = []component.Property{{Label: "ServiceAccount", Value: serviceAccountLink}}

		if nodeName := pod.Spec.NodeName; nodeName != "" {
			nodeLink, _ := link.ForGVK("", "v1", "Node", pod.Spec.NodeName, pod.Spec.NodeName)
			status.Properties = append(status.Properties, component.Property{Label: "Node", Value: nodeLink})
		}

		ownerReference := metav1.GetControllerOf(pod)
		if ownerReference != nil {
			controlledBy, err := link.ForGVK(
				pod.Namespace,
				ownerReference.APIVersion,
				ownerReference.Kind,
				ownerReference.Name,
				ownerReference.Name,
			)
			if err == nil {
				status.Properties = append(status.Properties, component.Property{Label: "Controlled By", Value: controlledBy})
			}
		}
	}

	return status, nil
}
