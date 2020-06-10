/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package logviewer

import (
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/internal/util/kubernetes"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// ToComponent converts an object into a log viewer component.
func ToComponent(object runtime.Object) (component.Component, error) {
	if object == nil {
		return nil, errors.Errorf("object is nil")
	}

	pod := &corev1.Pod{}

	switch t := object.(type) {
	case *unstructured.Unstructured:
		if err := kubernetes.FromUnstructured(t, pod); err != nil {
			return nil, err
		}
	case *corev1.Pod:
		pod = t
	default:
		pod = nil
	}

	if pod == nil {
		return nil, errors.Errorf("can't fetch logs from a %T", object)
	}

	containerNames := []string{""}

	for _, c := range pod.Spec.InitContainers {
		containerNames = append(containerNames, c.Name)
	}

	for _, c := range pod.Spec.Containers {
		containerNames = append(containerNames, c.Name)
	}

	for _, c := range pod.Spec.EphemeralContainers {
		containerNames = append(containerNames, c.Name)
	}

	logsComponent := component.NewLogs(pod.Namespace, pod.Name, containerNames...)

	return logsComponent, nil
}
