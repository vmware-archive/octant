package printer

import (
	"path"

	"github.com/heptio/developer-dash/internal/view/component"
)

// gvkPath composes a path given resource coordinates
func gvkPath(apiVersion, kind, name string) string {
	var p string

	switch {
	case apiVersion == "apps/v1" && kind == "DaemonSet":
		p = "/content/overview/workloads/daemon-sets"
	case apiVersion == "extensions/v1beta1" && kind == "ReplicaSet":
		p = "/content/overview/workloads/replica-sets"
	case apiVersion == "apps/v1" && kind == "StatefulSet":
		p = "/content/overview/workloads/stateful-sets"
	case apiVersion == "extensions/v1beta1" && kind == "Deployment":
		p = "/content/overview/workloads/deployments"
	case apiVersion == "apps/v1" && kind == "Deployment":
		p = "/content/overview/workloads/deployments"
	case apiVersion == "batch/v1beta1" && kind == "CronJob":
		p = "/content/overview/workloads/cron-jobs"
	case (apiVersion == "batch/v1beta1" || apiVersion == "batch/v1") && kind == "Job":
		p = "/content/overview/workloads/jobs"
	case apiVersion == "v1" && kind == "ReplicationController":
		p = "/content/overview/workloads/replication-controllers"
	case apiVersion == "v1" && kind == "Secret":
		p = "/content/overview/config-and-storage/secrets"
	case apiVersion == "v1" && kind == "ConfigMap":
		p = "/content/overview/config-and-storage/configmaps"
	case apiVersion == "v1" && kind == "PersistentVolumeClaim":
		p = "/content/overview/config-and-storage/persistent-volume-claims"
	case apiVersion == "v1" && kind == "ServiceAccount":
		p = "/content/overview/config-and-storage/service-accounts"
	case apiVersion == "v1" && kind == "Service":
		p = "/content/overview/discovery-and-load-balancing/services"
	case apiVersion == "rbac.authorization.k8s.io/v1" && kind == "Role":
		p = "/content/overview/rbac/roles"
	default:
		return "/content/overview"
	}

	return path.Join(p, name)
}

// linkForObject returns a link component referencing an object
func linkForObject(apiVersion, kind, name, text string) *component.Link {
	path := gvkPath(apiVersion, kind, name)
	return component.NewLink("", text, path)
}
