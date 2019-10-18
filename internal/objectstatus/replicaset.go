/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectstatus

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// replicaSetExtV1Beta1 creates status for an ext/v1beta1 replica set. This is
// not the final implementation. It is included to generate output.
func replicaSetExtV1Beta1(_ context.Context, object runtime.Object, _ store.Store) (ObjectStatus, error) {
	if object == nil {
		return ObjectStatus{}, errors.Errorf("replica set is nil")
	}

	replicaSet := &extv1beta1.ReplicaSet{}

	if err := scheme.Scheme.Convert(object, replicaSet, 0); err != nil {
		return ObjectStatus{}, errors.Wrap(err, "convert object to ext/v1beta1 replica set")
	}

	status := replicaSet.Status

	switch {
	case status.Replicas == 0:
		return ObjectStatus{
			nodeStatus: component.NodeStatusError,
			Details:    []component.Component{component.NewText(fmt.Sprintf("Replica Set has no replicas available"))},
		}, nil
	case status.Replicas == status.AvailableReplicas:
		return ObjectStatus{
			nodeStatus: component.NodeStatusOK,
			Details:    []component.Component{component.NewText(fmt.Sprintf("Replica Set is OK"))},
		}, nil
	default:
		return ObjectStatus{
			nodeStatus: component.NodeStatusWarning,
			Details:    []component.Component{component.NewText(fmt.Sprintf("Expected %d replicas, but %d are available", status.Replicas, status.AvailableReplicas))},
		}, nil
	}

}

// replicaSetExtV1Beta1 creates status for an v1/apps replica set. This is
// not the final implementation. It is included to generate output.
func replicaSetAppsV1(_ context.Context, object runtime.Object, _ store.Store) (ObjectStatus, error) {
	if object == nil {
		return ObjectStatus{}, errors.Errorf("replica set is nil")
	}

	replicaSet := &appsv1.ReplicaSet{}

	if err := scheme.Scheme.Convert(object, replicaSet, 0); err != nil {
		return ObjectStatus{}, errors.Wrap(err, "convert object to apps/v1 replica set")
	}

	status := replicaSet.Status

	switch {
	case status.Replicas == 0:
		return ObjectStatus{
			nodeStatus: component.NodeStatusError,
			Details:    []component.Component{component.NewText("Replica Set has no replicas available")},
		}, nil
	case status.Replicas == status.AvailableReplicas:
		return ObjectStatus{nodeStatus: component.NodeStatusOK,
			Details: []component.Component{component.NewText("Replica Set is OK")},
		}, nil
	default:
		return ObjectStatus{
			nodeStatus: component.NodeStatusWarning,
			Details:    []component.Component{component.NewText(fmt.Sprintf("Expected %d replicas, but %d are available", status.Replicas, status.AvailableReplicas))},
		}, nil
	}

}
