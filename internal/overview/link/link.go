package link

import (
	"fmt"
	"net/url"
	"path"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// gvkPathFromObject composes a path given an object.
func gvkPathFromObject(object runtime.Object) (string, error) {
	if object == nil {
		return "", errors.New("object is nil")
	}

	gvk := object.GetObjectKind().GroupVersionKind()
	apiVersion, kind := gvk.ToAPIVersionAndKind()

	accessor := meta.NewAccessor()
	name, err := accessor.Name(object)
	if err != nil {
		return "", errors.Wrap(err, "retrieve name from object")
	}

	ns, err := accessor.Namespace(object)
	if err != nil {
		return "", errors.Wrap(err, "retrieve namespace from object")
	}
	return gvkPath(ns, apiVersion, kind, name), nil
}

// gvkPath composes a path given resource coordinates
func gvkPath(namespace, apiVersion, kind, name string) string {
	var p string

	switch {
	case apiVersion == "apps/v1" && kind == "DaemonSet":
		p = "/workloads/daemon-sets"
	case apiVersion == "extensions/v1beta1" && kind == "ReplicaSet":
		p = "/workloads/replica-sets"
	case apiVersion == "apps/v1" && kind == "ReplicaSet":
		p = "/workloads/replica-sets"
	case apiVersion == "apps/v1" && kind == "StatefulSet":
		p = "/workloads/stateful-sets"
	case apiVersion == "extensions/v1beta1" && kind == "Deployment":
		p = "/workloads/deployments"
	case apiVersion == "apps/v1" && kind == "Deployment":
		p = "/workloads/deployments"
	case apiVersion == "batch/v1beta1" && kind == "CronJob":
		p = "/workloads/cron-jobs"
	case (apiVersion == "batch/v1beta1" || apiVersion == "batch/v1") && kind == "Job":
		p = "/workloads/jobs"
	case apiVersion == "v1" && kind == "ReplicationController":
		p = "/workloads/replication-controllers"
	case apiVersion == "v1" && kind == "Secret":
		p = "/config-and-storage/secrets"
	case apiVersion == "v1" && kind == "ConfigMap":
		p = "/config-and-storage/config-maps"
	case apiVersion == "v1" && kind == "PersistentVolumeClaim":
		p = "/config-and-storage/persistent-volume-claims"
	case apiVersion == "v1" && kind == "ServiceAccount":
		p = "/config-and-storage/service-accounts"
	case apiVersion == "extensions/v1beta1" && kind == "Ingress":
		p = "/discovery-and-load-balancing/ingresses"
	case apiVersion == "v1" && kind == "Service":
		p = "/discovery-and-load-balancing/services"
	case apiVersion == "rbac.authorization.k8s.io/v1" && kind == "Role":
		p = "/rbac/roles"
	case apiVersion == "rbac.authorization.k8s.io/v1" && kind == "ClusterRole":
		p = "/rbac/cluster-roles"
	case apiVersion == "rbac.authorization.k8s.io/v1" && kind == "RoleBinding":
		p = "/rbac/role-bindings"
	case apiVersion == "rbac.authorization.k8s.io/v1" && kind == "ClusterRoleBinding":
		p = "/rbac/cluster-role-bindings"
	case apiVersion == "v1" && kind == "Event":
		p = "/events"
	case apiVersion == "v1" && kind == "Pod":
		p = "/workloads/pods"
	default:
		return fmt.Sprintf("/content/overview/%s", namespace)
	}

	prefix := "/content/overview"

	if namespace != "" {
		prefix = fmt.Sprintf("%s/namespace/%s", prefix, namespace)
	}

	return path.Join(prefix, p, name)
}

// ForCustomResourceName returns a link component referencing
// a custom resource definition.
func ForCustomResourceDefinition(crdName, namespace string) *component.Link {
	ref := path.Join("/content/overview/namespace", namespace,
		"custom-resources", crdName)
	return component.NewLink("", crdName, ref)
}

// ForCustomResource returns a link component referenceing
// a custom resource.
func ForCustomResource(crdName string, object *unstructured.Unstructured) *component.Link {
	ref := path.Join("/content/overview/namespace", object.GetNamespace(),
		"custom-resources", crdName, object.GetName())
	return component.NewLink("", object.GetName(), ref)
}

// ForObject returns a link component referencing an object
// Returns an empty link if an error occurs.
func ForObject(object runtime.Object, text string) *component.Link {
	path, _ := gvkPathFromObject(object)
	return component.NewLink("", text, path)
}

// ForObjectWithQuery returns a link component references an object with a query.
// Return an empty link if an error occurs.
func ForObjectWithQuery(object runtime.Object, text string, query url.Values) *component.Link {
	path, _ := gvkPathFromObject(object)
	u := url.URL{Path: path}
	u.RawQuery = query.Encode()
	return component.NewLink("", text, u.String())
}

// ForGVK returns a link component referencing an object
func ForGVK(namespace, apiVersion, kind, name, text string) *component.Link {
	path := gvkPath(namespace, apiVersion, kind, name)
	return component.NewLink("", text, path)
}

// ForOwner returns a link component for an owner.
func ForOwner(parent runtime.Object, controllerRef *metav1.OwnerReference) *component.Link {
	if controllerRef == nil || parent == nil {
		return component.NewLink("", "none", "")
	}

	accessor := meta.NewAccessor()
	ns, err := accessor.Namespace(parent)
	if err != nil {
		return component.NewLink("", "none", "")
	}

	return ForGVK(
		ns,
		controllerRef.APIVersion,
		controllerRef.Kind,
		controllerRef.Name,
		controllerRef.Name,
	)
}
