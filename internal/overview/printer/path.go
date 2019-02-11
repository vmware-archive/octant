package printer

import (
	"path"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
)

type objectReferenceKey struct {
	apiVersion string
	kind       string
}

var (
	objectReferenceLookup = map[objectReferenceKey]string{
		objectReferenceKey{apiVersion: "batch/v1beta1", kind: "CronJob"}:      "workloads/cron-jobs",
		objectReferenceKey{apiVersion: "apps/v1", kind: "DaemonSet"}:          "workloads/daemon-sets",
		objectReferenceKey{apiVersion: "apps/v1", kind: "Deployment"}:         "workloads/deployments",
		objectReferenceKey{apiVersion: "batch/v1", kind: "Job"}:               "workloads/jobs",
		objectReferenceKey{apiVersion: "v1", kind: "Pod"}:                     "workloads/pods",
		objectReferenceKey{apiVersion: "apps/v1", kind: "ReplicaSet"}:         "workloads/replica-sets",
		objectReferenceKey{apiVersion: "v1", kind: "ReplicationController"}:   "workloads/replication-controllers",
		objectReferenceKey{apiVersion: "apps/v1", kind: "StatefulSet"}:        "workloads/stateful-sets",
		objectReferenceKey{apiVersion: "extensions/v1beta1", kind: "Ingress"}: "discovery-and-load-balancing/ingresses",
		objectReferenceKey{apiVersion: "v1", kind: "Service"}:                 "discovery-and-load-balancing/services",
		objectReferenceKey{apiVersion: "v1", kind: "ConfigMap"}:               "config-and-storage/config-maps",
		objectReferenceKey{apiVersion: "v1", kind: "PersistentVolumeClaim"}:   "config-and-storage/persistent-volume-claims",
		objectReferenceKey{apiVersion: "v1", kind: "Secret"}:                  "config-and-storage/secrets",
		objectReferenceKey{apiVersion: "v1", kind: "ServiceAccount"}:          "config-and-storage/service-accounts",
		objectReferenceKey{apiVersion: "v1", kind: "Role"}:                    "rbac/roles",
		objectReferenceKey{apiVersion: "v1", kind: "RoleBinding"}:             "rbac/role-bindings",
		objectReferenceKey{apiVersion: "v1", kind: "Event"}:                   "events",
	}
)

// ObjectReferencePath returns the overview path for an object reference.
// Currently, this does not support custom resources.
func ObjectReferencePath(or corev1.ObjectReference) (string, error) {
	key := objectReferenceKey{
		apiVersion: or.APIVersion,
		kind:       or.Kind,
	}

	section, ok := objectReferenceLookup[key]
	if !ok {
		return "", errors.Errorf("unable to locate overview section for %s/%s",
			key.apiVersion, key.kind)
	}

	var objectPath string
	if or.Namespace != "" {
		objectPath = path.Join("/content/overview/namespace", or.Namespace, section, or.Name)
	} else {
		objectPath = path.Join("/content/overview", section, or.Name)
	}
	return objectPath, nil
}
