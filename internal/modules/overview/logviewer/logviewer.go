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
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(t.Object, pod); err != nil {
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

	logsComponent := component.NewLogs(pod)

	return logsComponent, nil
}
