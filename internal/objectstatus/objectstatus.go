/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectstatus

import (
	"context"
	"errors"
	"fmt"

	"github.com/vmware-tanzu/octant/internal/link"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

type statusKey struct {
	apiVersion string
	kind       string
}

type statusFunc func(context.Context, runtime.Object, store.Store, link.Interface) (ObjectStatus, error)

type statusLookup map[statusKey]statusFunc

var (
	defaultStatusLookup = statusLookup{
		{apiVersion: "batch/v1beta1", kind: "CronJob"}:                cronJob,
		{apiVersion: "apps/v1", kind: "DaemonSet"}:                    daemonSet,
		{apiVersion: "apps/v1", kind: "Deployment"}:                   deploymentAppsV1,
		{apiVersion: "apps/v1", kind: "ReplicaSet"}:                   replicaSetAppsV1,
		{apiVersion: "apps/v1", kind: "StatefulSet"}:                  statefulSet,
		{apiVersion: "batch/v1", kind: "Job"}:                         runJobStatus,
		{apiVersion: "v1", kind: "Pod"}:                               pod,
		{apiVersion: "v1", kind: "ReplicationController"}:             replicationController,
		{apiVersion: "v1", kind: "Service"}:                           service,
		{apiVersion: "v1", kind: "PersistentVolume"}:                  persistentVolume,
		{apiVersion: "networking.k8s.io/v1", kind: "Ingress"}:         runIngressStatus,
		{apiVersion: "apiregistration.k8s.io/v1", kind: "APIService"}: apiService,
	}
)

type ObjectStatus struct {
	nodeStatus component.NodeStatus
	Details    []component.Component
	Properties []component.Property
}

func (os *ObjectStatus) AddDetail(detail string) {
	os.Details = append(os.Details, component.NewText(detail))
}

func (os *ObjectStatus) AddProperty(label string, value component.Component) {
	os.Properties = append(os.Properties, component.Property{
		Label: label,
		Value: value,
	})
}

func (os *ObjectStatus) InsertProperty(label string, value component.Component) {
	os.Properties = append([]component.Property{{
		Label: label,
		Value: value,
	}}, os.Properties...)
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
func Status(ctx context.Context, object runtime.Object, o store.Store, link link.Interface) (ObjectStatus, error) {
	return status(ctx, object, o, defaultStatusLookup, link)
}

func status(ctx context.Context, object runtime.Object, o store.Store, lookup statusLookup, link link.Interface) (ObjectStatus, error) {
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

	var oStatus ObjectStatus
	fn, ok := lookup[statusKey{apiVersion: apiVersion, kind: kind}]

	if ok {
		oStatus, err = fn(ctx, object, o, link)
	} else {
		oStatus = ObjectStatus{
			nodeStatus: component.NodeStatusOK,
			Details:    []component.Component{component.NewText(fmt.Sprintf("%s %s is OK", apiVersion, kind))},
			Properties: []component.Property{},
		}
	}

	if labels := accessor.GetLabels(); len(labels) > 0 {
		oStatus.InsertProperty("Labels", component.NewLabels(labels))
	}
	oStatus.InsertProperty("Created", component.NewTimestamp(accessor.GetCreationTimestamp().Time))
	oStatus.InsertProperty("Namespace", component.NewText(accessor.GetNamespace()))

	return oStatus, nil
}
