package clusteroverview

import (
	"path"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/heptio/developer-dash/internal/gvk"
)

type objectPath struct{}

func (op *objectPath) SupportedGroupVersionKind() []schema.GroupVersionKind{
	return []schema.GroupVersionKind{
		gvk.ClusterRoleBindingGVK,
		gvk.ClusterRoleGVK,
	}
}

func (op *objectPath) GroupVersionKindPath(namespace, apiVersion, kind, name string) (string, error) {
	return gvkPath(apiVersion, kind, name)
}

func gvkPath(apiVersion, kind, name string) (string, error) {
	var p string

	switch {
	case apiVersion == "rbac.authorization.k8s.io/v1" && kind == "ClusterRole":
		p = "/rbac/cluster-roles"
	case apiVersion == "rbac.authorization.k8s.io/v1" && kind == "ClusterRoleBinding":
		p = "/rbac/cluster-role-bindings"
	default:
		return "", errors.Errorf("unknown object %s %s", apiVersion, kind)
	}

	return path.Join("/content/cluster-overview", p, name), nil

}
