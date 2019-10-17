/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectstatus

import (
	"context"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func statefulSet(_ context.Context, object runtime.Object, _ store.Store) (ObjectStatus, error) {
	if object == nil {
		return ObjectStatus{}, errors.Errorf("stateful set is nil")
	}

	ss := &appsv1.StatefulSet{}

	if err := scheme.Scheme.Convert(object, ss, 0); err != nil {
		return ObjectStatus{}, errors.Wrap(err, "convert object to stateful set")
	}

	status := ss.Status

	switch {
	case status.ReadyReplicas != status.Replicas:
		return ObjectStatus{
			nodeStatus: component.NodeStatusWarning,
			Details:    []component.Component{component.NewText("Stateful Set pods are not ready")},
		}, nil
	default:
		return ObjectStatus{
			nodeStatus: component.NodeStatusOK,
			Details:    []component.Component{component.NewText("Stateful Set is OK")},
		}, nil
	}
}
