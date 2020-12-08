/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package clusteroverview

import (
	"fmt"
	"path"

	"github.com/pkg/errors"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware-tanzu/octant/internal/gvk"
)

var (
	supportedGVKs = []schema.GroupVersionKind{
		gvk.ClusterRoleBinding,
		gvk.ClusterRole,
		gvk.Node,
		gvk.PersistentVolume,
		gvk.Namespace,
		gvk.CustomResourceDefinition,
		gvk.APIService,
		gvk.MutatingWebhookConfiguration,
		gvk.ValidatingWebhookConfiguration,
	}
)

const rbacAPIVersion = "rbac.authorization.k8s.io/v1"

func crdPath(namespace, crdName, version, name string) (string, error) {
	return path.Join("/cluster-overview/custom-resources", crdName, version, name), nil
}

func gvkPath(namespace, apiVersion, kind, name string) (string, error) {
	var p string

	switch {
	case apiVersion == rbacAPIVersion && kind == "ClusterRole":
		p = "/rbac/cluster-roles"
	case apiVersion == rbacAPIVersion && kind == "ClusterRoleBinding":
		p = "/rbac/cluster-role-bindings"
	case apiVersion == "v1" && kind == "Node":
		p = "/nodes"
	case apiVersion == "v1" && kind == "PersistentVolume":
		p = "/storage/persistent-volumes"
	case apiVersion == "v1" && kind == "Namespace":
		p = "/namespaces"
	case apiVersion == gvk.CustomResourceDefinition.GroupVersion().String() &&
		kind == gvk.CustomResourceDefinition.Kind:
		p = "/custom-resource-definitions"
	case apiVersion == "apiregistration.k8s.io/v1" && kind == "APIService":
		p = "/api-server/api-services"
	case apiVersion == "admissionregistration.k8s.io/v1beta1" && kind == "MutatingWebhookConfiguration":
		p = "/webhooks/mutating-webhooks"
	case apiVersion == "admissionregistration.k8s.io/v1beta1" && kind == "ValidatingWebhookConfiguration":
		p = "/webhooks/validating-webhooks"
	default:
		return "", fmt.Errorf("unknown object %s %s", apiVersion, kind)
	}

	return path.Join("/cluster-overview", p, name), nil
}

func gvkReversePath(contentPath, _ string) (schema.GroupVersionKind, error) {

	switch {
	case contentPath == "cluster-overview/rbac/cluster-roles":
		return gvk.ClusterRole, nil
	case contentPath == "cluster-overview/rbac/cluster-role-bindings":
		return gvk.ClusterRoleBinding, nil
	case contentPath == "cluster-overview/nodes":
		return gvk.Node, nil
	case contentPath == "cluster-overview/storage/persistent-volumes":
		return gvk.PersistentVolume, nil
	case contentPath == "cluster-overview/namespaces":
		return gvk.Namespace, nil
	case contentPath == "cluster-overview/custom-resource-definitions":
		return gvk.CustomResourceDefinition, nil
	case contentPath == "cluster-overview/api-server/api-services":
		return gvk.APIService, nil
	case contentPath == "cluster-overview/webhooks/mutating-webhooks":
		return gvk.MutatingWebhookConfiguration, nil
	case contentPath == "cluster-overview/webhooks/validating-webhooks":
		return gvk.ValidatingWebhookConfiguration, nil
	default:
		return schema.GroupVersionKind{}, errors.Errorf("unknown gvk %s", contentPath)
	}
}
