/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package overview

import (
	"path"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware-tanzu/octant/internal/gvk"
)

var (
	supportedGVKs = []schema.GroupVersionKind{
		gvk.AppReplicaSet,
		gvk.CronJob,
		gvk.DaemonSet,
		gvk.Deployment,
		gvk.ExtReplicaSet,
		gvk.Job,
		gvk.Pod,
		gvk.ReplicationController,
		gvk.StatefulSet,
		gvk.HorizontalPodAutoscaler,
		gvk.Ingress,
		gvk.Service,
		gvk.NetworkPolicy,
		gvk.ConfigMap,
		gvk.Secret,
		gvk.PersistentVolumeClaim,
		gvk.ServiceAccount,
		gvk.RoleBinding,
		gvk.Role,
		gvk.Event,
	}
)

func crdPath(namespace, crdName, version, name string) (string, error) {
	if namespace == "" {
		return "", errors.Errorf("unable to create CRD path for %s due to missing namespace", crdName)
	}

	return path.Join("/overview/namespace", namespace, "custom-resources", crdName, version, name), nil
}

func gvkPath(namespace, apiVersion, kind, name string) (string, error) {
	if namespace == "" {
		return "", errors.Errorf("unable to create  path for %s %s due to missing namespace", apiVersion, kind)
	}

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
	case (apiVersion == "autoscaling/v1" || apiVersion == "autoscaling/v2beta2") && kind == "HorizontalPodAutoscaler":
		p = "/discovery-and-load-balancing/horizontal-pod-autoscalers"
	case apiVersion == "extensions/v1beta1" && kind == "Ingress":
		p = "/discovery-and-load-balancing/ingresses"
	case apiVersion == "v1" && kind == "Service":
		p = "/discovery-and-load-balancing/services"
	case apiVersion == "networking.k8s.io/v1" && kind == "NetworkPolicy":
		p = "/discovery-and-load-balancing/network-policies"
	case apiVersion == "rbac.authorization.k8s.io/v1" && kind == "Role":
		p = "/rbac/roles"
	case apiVersion == "rbac.authorization.k8s.io/v1" && kind == "RoleBinding":
		p = "/rbac/role-bindings"
	case apiVersion == "v1" && kind == "Event":
		p = "/events"
	case apiVersion == "v1" && kind == "Pod":
		p = "/workloads/pods"
	default:
		return "", errors.Errorf("unknown object %s %s", apiVersion, kind)
	}

	return path.Join("/overview/namespace", namespace, p, name), nil
}
