/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectstatus

import (
	"context"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func pod(ctx context.Context, object runtime.Object, o store.Store) (ObjectStatus, error) {
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
	status.	nodeStatus = component.NodeStatusWarning
	}

	status.Details = []component.Component{
		component.NewText(pod.Status.Message),
	}

	return status, nil
}
