/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package clusteroverview

import (
	"path"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware/octant/internal/gvk"
)

var (
	supportedGVKs = []schema.GroupVersionKind{
		gvk.ClusterRoleBinding,
		gvk.ClusterRole,
		gvk.Node,
	}
)

const rbacAPIVersion = "rbac.authorization.k8s.io/v1"

func crdPath(namespace, crdName, name string) (string, error) {
	return path.Join("/content/cluster-overview/custom-resources", crdName, name), nil
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
	default:
		return "", errors.Errorf("unknown object %s %s", apiVersion, kind)
	}

	return path.Join("/content/cluster-overview", p, name), nil
}
