/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package gvk

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	ClusterRoleBindingGVK       = schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRoleBinding"}
	ClusterRoleGVK              = schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"}
	ConfigMapGVK                = schema.GroupVersionKind{Version: "v1", Kind: "ConfigMap"}
	CronJobGVK                  = schema.GroupVersionKind{Group: "batch", Version: "v1beta1", Kind: "CronJob"}
	CustomResourceDefinitionGVK = schema.GroupVersionKind{Group: "apiextensions.k8s.io", Version: "v1beta1", Kind: "CustomResourceDefinition"}
	DaemonSetGVK                = schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "DaemonSet"}
	DeploymentGVK               = schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}
	Event                      = schema.GroupVersionKind{Version: "v1", Kind: "Event"}
	IngressGVK                  = schema.GroupVersionKind{Group: "extensions", Version: "v1beta1", Kind: "Ingress"}
	JobGVK                      = schema.GroupVersionKind{Group: "batch", Version: "v1", Kind: "Job"}
	ServiceAccountGVK           = schema.GroupVersionKind{Version: "v1", Kind: "ServiceAccount"}
	SecretGVK                   = schema.GroupVersionKind{Version: "v1", Kind: "Secret"}
	ServiceGVK                  = schema.GroupVersionKind{Version: "v1", Kind: "Service"}
	PodGVK                      = schema.GroupVersionKind{Version: "v1", Kind: "Pod"}
	ReplicaSetGVK               = schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "ReplicaSet"}
	PersistentVolumeClaimGVK    = schema.GroupVersionKind{Version: "v1", Kind: "PersistentVolumeClaim"}
	ReplicationControllerGVK    = schema.GroupVersionKind{Version: "v1", Kind: "ReplicationController"}
	StatefulSetGVK              = schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "StatefulSet"}
	RoleBindingGVK              = schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "RoleBinding"}
	RoleGVK                     = schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "Role"}
)
