/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package overview

import (
	"context"

	"github.com/vmware/octant/internal/gvk"
	"github.com/vmware/octant/internal/loading"
	"github.com/vmware/octant/pkg/icon"
	"github.com/vmware/octant/pkg/navigation"
	"github.com/vmware/octant/pkg/store"
)

var (
	navPathLookup = map[string]string{
		"Workloads":                    "workloads",
		"Discovery and Load Balancing": "discovery-and-load-balancing",
		"Config and Storage":           "config-and-storage",
		"Custom Resources":             "custom-resources",
		"RBAC":                         "rbac",
		"Events":                       "events",
	}
)

func workloadEntries(ctx context.Context, prefix, namespace string, objectStore store.Store, _ bool) ([]navigation.Navigation, bool, error) {
	neh := navigation.EntriesHelper{}

	neh.Add("Cron Jobs", "cron-jobs", icon.OverviewCronJob,
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.CronJob), objectStore))
	neh.Add("Daemon Sets", "daemon-sets", icon.OverviewDaemonSet,
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.DaemonSet), objectStore))
	neh.Add("Deployments", "deployments", icon.OverviewDeployment,
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.Deployment), objectStore))
	neh.Add("Jobs", "jobs", icon.OverviewJob,
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.Job), objectStore))
	neh.Add("Pods", "pods", icon.OverviewPod,
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.Pod), objectStore))
	neh.Add("Replica Sets", "replica-sets", icon.OverviewReplicaSet,
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.ExtReplicaSet), objectStore))
	neh.Add("Replication Controllers", "replication-controllers", icon.OverviewReplicationController,
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.ReplicationController), objectStore))
	neh.Add("Stateful Sets", "stateful-sets", icon.OverviewStatefulSet,
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.StatefulSet), objectStore))

	children, err := neh.Generate(prefix)
	if err != nil {
		return nil, false, err
	}

	return children, false, nil
}

func discoAndLBEntries(ctx context.Context, prefix, namespace string, objectStore store.Store, _ bool) ([]navigation.Navigation, bool, error) {
	neh := navigation.EntriesHelper{}
	neh.Add("Horizontal Pod Autoscalers", "horizontal-pod-autoscalers", icon.OverviewHorizontalPodAutoscaler,
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.HorizontalPodAutoscaler), objectStore))
	neh.Add("Ingresses", "ingresses", icon.OverviewIngress,
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.Ingress), objectStore))
	neh.Add("Services", "services", icon.OverviewService,
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.Service), objectStore))

	children, err := neh.Generate(prefix)
	if err != nil {
		return nil, false, err
	}

	return children, false, nil
}

func configAndStorageEntries(ctx context.Context, prefix, namespace string, objectStore store.Store, _ bool) ([]navigation.Navigation, bool, error) {
	neh := navigation.EntriesHelper{}
	neh.Add("Config Maps", "config-maps", icon.OverviewConfigMap,
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.ConfigMap), objectStore))
	neh.Add("Persistent Volume Claims", "persistent-volume-claims", icon.OverviewPersistentVolumeClaim,
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.PersistentVolumeClaim), objectStore))
	neh.Add("Secrets", "secrets", icon.OverviewSecret,
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.Secret), objectStore))
	neh.Add("Service Accounts", "service-accounts", icon.OverviewServiceAccount,
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.ServiceAccount), objectStore))

	children, err := neh.Generate(prefix)
	if err != nil {
		return nil, false, err
	}

	return children, false, nil
}

func rbacEntries(ctx context.Context, prefix, namespace string, objectStore store.Store, _ bool) ([]navigation.Navigation, bool, error) {
	neh := navigation.EntriesHelper{}

	neh.Add("Roles", "roles", icon.OverviewRole,
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.Role), objectStore))
	neh.Add("Role Bindings", "role-bindings", icon.OverviewRoleBinding,
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.RoleBinding), objectStore))

	children, err := neh.Generate(prefix)
	if err != nil {
		return nil, false, err
	}

	return children, false, nil
}
