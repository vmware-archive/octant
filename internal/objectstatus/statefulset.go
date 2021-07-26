/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectstatus

import (
	"context"
	"fmt"

	"github.com/vmware-tanzu/octant/internal/link"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func statefulSet(_ context.Context, object runtime.Object, _ store.Store, _ link.Interface) (ObjectStatus, error) {
	if object == nil {
		return ObjectStatus{}, errors.Errorf("stateful set is nil")
	}

	ss := &appsv1.StatefulSet{}

	if err := scheme.Scheme.Convert(object, ss, 0); err != nil {
		return ObjectStatus{}, errors.Wrap(err, "convert object to stateful set")
	}

	status := ss.Status
	total := fmt.Sprintf("%d", status.Replicas)
	desired := fmt.Sprintf("%d", *ss.Spec.Replicas)

	properties := []component.Property{{Label: "Replicas", Value: component.NewText(fmt.Sprintf("%s Desired / %s Total", desired, total))},
		{Label: "Pod Management Policy", Value: component.NewText(string(ss.Spec.PodManagementPolicy))}}

	switch {
	case status.ReadyReplicas != status.Replicas:
		return ObjectStatus{
			NodeStatus: component.NodeStatusWarning,
			Details:    []component.Component{component.NewText("Stateful Set pods are not ready")},
			Properties: properties,
		}, nil
	default:
		return ObjectStatus{
			NodeStatus: component.NodeStatusOK,
			Details:    []component.Component{component.NewText("Stateful Set is OK")},
			Properties: properties,
		}, nil
	}
}
