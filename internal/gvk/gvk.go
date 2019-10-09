/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package gvk

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	AppReplicaSet            = schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "ReplicaSet"}
	ClusterRoleBinding       = schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRoleBinding"}
	ClusterRole              = schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"}
	ConfigMap                = schema.GroupVersionKind{Version: "v1", Kind: "ConfigMap"}
	CronJob                  = schema.GroupVersionKind{Group: "batch", Version: "v1beta1", Kind: "CronJob"}
	CustomResourceDefinition = schema.GroupVersionKind{Group: "apiextensions.k8s.io", Version: "v1beta1", Kind: "CustomResourceDefinition"}
	DaemonSet                = schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "DaemonSet"}
	Deployment               = schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}
	ExtDeployment            = schema.GroupVersionKind{Group: "extensions", Version: "v1beta1", Kind: "Deployment"}
	ExtReplicaSet            = schema.GroupVersionKind{Group: "extensions", Version: "v1beta1", Kind: "ReplicaSet"}
	Event                    = schema.GroupVersionKind{Version: "v1", Kind: "Event"}
	HorizontalPodAutoscaler  = schema.GroupVersionKind{Group: "autoscaling", Version: "v1", Kind: "HorizontalPodAutoscaler"}
	Ingress                  = schema.GroupVersionKind{Group: "extensions", Version: "v1beta1", Kind: "Ingress"}
	Job                      = schema.GroupVersionKind{Group: "batch", Version: "v1", Kind: "Job"}
	Node                     = schema.GroupVersionKind{Version: "v1", Kind: "Node"}
	ServiceAccount           = schema.GroupVersionKind{Version: "v1", Kind: "ServiceAccount"}
	Secret                   = schema.GroupVersionKind{Version: "v1", Kind: "Secret"}
	Service                  = schema.GroupVersionKind{Version: "v1", Kind: "Service"}
	Pod                      = schema.GroupVersionKind{Version: "v1", Kind: "Pod"}
	PersistentVolumeClaim    = schema.GroupVersionKind{Version: "v1", Kind: "PersistentVolumeClaim"}
	ReplicationController    = schema.GroupVersionKind{Version: "v1", Kind: "ReplicationController"}
	StatefulSet              = schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "StatefulSet"}
	RoleBinding              = schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "RoleBinding"}
	Role                     = schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "Role"}
)
