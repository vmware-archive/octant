/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectstatus

import (
	"context"
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

type statusKey struct {
	apiVersion string
	kind       string
}

type statusFunc func(context.Context, runtime.Object, store.Store) (ObjectStatus, error)

type statusLookup map[statusKey]statusFunc

var (
	defaultStatusLookup = statusLookup{
		{apiVersion: "batch/v1beta1", kind: "CronJob"}:      cronJob,
		{apiVersion: "apps/v1", kind: "DaemonSet"}:          daemonSet,
		{apiVersion: "apps/v1", kind: "Deployment"}:         deploymentAppsV1,
		{apiVersion: "apps/v1", kind: "ReplicaSet"}:         replicaSetAppsV1,
		{apiVersion: "apps/v1", kind: "StatefulSet"}:        statefulSet,
		{apiVersion: "batch/v1", kind: "Job"}:               runJobStatus,
		{apiVersion: "v1", kind: "Pod"}:                     pod,
		{apiVersion: "v1", kind: "ReplicationController"}:   replicationController,
		{apiVersion: "v1", kind: "Service"}:                 service,
		{apiVersion: "extensions/v1beta1", kind: "Ingress"}: runIngressStatus,
	}
)

type ObjectStatus struct {
	nodeStatus component.NodeStatus
	Details    []component.Component
}

func (os *ObjectStatus) AddDetail(detail string) {
	os.Details = append(os.Details, component.NewText(detail))
}

func (os *ObjectStatus) AddDetailf(msg string, args ...interface{}) {
	os.AddDetail(fmt.Sprintf(msg, args...))
}

func (os *ObjectStatus) SetError() {
	os.nodeStatus = component.NodeStatusError
}

func (os *ObjectStatus) SetWarning() {
	if os.nodeStatus != component.NodeStatusError {
		os.nodeStatus = component.NodeStatusWarning
	}
}

func (os *ObjectStatus) Status() component.NodeStatus {
	switch os.nodeStatus {
	case component.NodeStatusWarning,
		component.NodeStatusError,
		component.NodeStatusOK:
		return os.nodeStatus
	default:
		return component.NodeStatusOK
	}
}

// Status creates an ObjectStatus for an object.
func Status(ctx context.Context, object runtime.Object, o store.Store) (ObjectStatus, error) {
	return status(ctx, object, o, defaultStatusLookup)
}

func status(ctx context.Context, object runtime.Object, o store.Store, lookup statusLookup) (ObjectStatus, error) {
	if object == nil {
		return ObjectStatus{}, errors.New("object is nil")
	}

	gvk := object.GetObjectKind().GroupVersionKind()
	apiVersion, kind := gvk.ToAPIVersionAndKind()

	accessor, err := meta.Accessor(object)
	if err != nil {
		return ObjectStatus{}, err
	}

	if accessor.GetDeletionTimestamp() != nil {
		return ObjectStatus{
			nodeStatus: component.NodeStatusWarning,
			Details: []component.Component{
				component.NewTextf("%s is being deleted", kind),
			},
		}, nil
	}

	if lookup == nil {
		return ObjectStatus{}, errors.New("status lookup is nil")
	}

	fn, ok := lookup[statusKey{apiVersion: apiVersion, kind: kind}]
	if !ok {
		return ObjectStatus{
			nodeStatus: component.NodeStatusOK,
			Details:    []component.Component{component.NewText(fmt.Sprintf("%s %s is OK", apiVersion, kind))},
		}, nil

	}

	return fn(ctx, object, o)
}
