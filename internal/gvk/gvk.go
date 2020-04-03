/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package gvk

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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
	Namespace                = schema.GroupVersionKind{Version: "v1", Kind: "Namespace"}
	NetworkPolicy            = schema.GroupVersionKind{Group: "networking.k8s.io", Version: "v1", Kind: "NetworkPolicy"}
	ServiceAccount           = schema.GroupVersionKind{Version: "v1", Kind: "ServiceAccount"}
	Secret                   = schema.GroupVersionKind{Version: "v1", Kind: "Secret"}
	Service                  = schema.GroupVersionKind{Version: "v1", Kind: "Service"}
	Pod                      = schema.GroupVersionKind{Version: "v1", Kind: "Pod"}
	PodMetrics               = schema.GroupVersionKind{Group: "metrics.k8s.io", Version: "v1beta1", Kind: "PodMetrics"}
	PersistentVolume         = schema.GroupVersionKind{Version: "v1", Kind: "PersistentVolume"}
	PersistentVolumeClaim    = schema.GroupVersionKind{Version: "v1", Kind: "PersistentVolumeClaim"}
	ReplicationController    = schema.GroupVersionKind{Version: "v1", Kind: "ReplicationController"}
	StatefulSet              = schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "StatefulSet"}
	RoleBinding              = schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "RoleBinding"}
	Role                     = schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "Role"}
)

// CustomResource generates a `schema.GroupVersionKind` for a custom resource given a version.
func CustomResource(crd *unstructured.Unstructured, version string) (schema.GroupVersionKind, error) {
	if crd == nil {
		return schema.GroupVersionKind{}, fmt.Errorf("custom resource definition is nil")
	}

	crdGroup := crd.GroupVersionKind().Group
	crdKind := crd.GroupVersionKind().Kind

	if crdGroup != CustomResourceDefinition.Group || crdKind != CustomResourceDefinition.Kind {
		return schema.GroupVersionKind{}, fmt.Errorf("input was not a crd {group: %q, kind: %q}", crdGroup, crdKind)
	}

	group, _, err := unstructured.NestedString(crd.Object, "spec", "group")
	if err != nil {
		return schema.GroupVersionKind{}, fmt.Errorf("get crd group: %w", err)
	}

	kind, _, err := unstructured.NestedString(crd.Object, "spec", "names", "kind")
	if err != nil {
		return schema.GroupVersionKind{}, fmt.Errorf("get crd kind: %w", err)
	}

	return schema.GroupVersionKind{
		Group:   group,
		Version: version,
		Kind:    kind,
	}, nil
}
