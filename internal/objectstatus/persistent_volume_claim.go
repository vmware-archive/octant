/*
Copyright (c) 2021 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectstatus

import (
	"context"

	"github.com/vmware-tanzu/octant/pkg/view/component"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/pkg/errors"

	"github.com/vmware-tanzu/octant/internal/link"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/pkg/store"
)

func persistentVolumeClaim(_ context.Context, object runtime.Object, _ store.Store, _ link.Interface) (ObjectStatus, error) {
	if object == nil {
		return ObjectStatus{}, errors.Errorf("cronjob is nil")
	}

	pvc := &corev1.PersistentVolumeClaim{}

	if err := scheme.Scheme.Convert(object, pvc, 0); err != nil {
		return ObjectStatus{}, errors.Wrap(err, "convert object to v1 PersistentVolumeClaim")
	}

	if pvc.Status.Phase == "Pending" {
		return ObjectStatus{
			NodeStatus: component.NodeStatusWarning,
			Details:    []component.Component{component.NewText("PVC cannot be found")}}, nil
	}

	return ObjectStatus{
		NodeStatus: component.NodeStatusOK,
		Details:    []component.Component{component.NewText("v1 PersistentVolumeClaim is OK")}}, nil
}
