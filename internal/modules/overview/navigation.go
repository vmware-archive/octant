/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package overview

import (
	"context"

	"github.com/vmware-tanzu/octant/internal/gvk"
	"github.com/vmware-tanzu/octant/internal/loading"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/store"
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

	neh.Add("Cron Jobs", "cron-jobs",
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.CronJob), objectStore))
	neh.Add("Daemon Sets", "daemon-sets",
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.DaemonSet), objectStore))
	neh.Add("Deployments", "deployments",
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.Deployment), objectStore))
	neh.Add("Jobs", "jobs",
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.Job), objectStore))
	neh.Add("Pods", "pods",
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.Pod), objectStore))
	neh.Add("Replica Sets", "replica-sets",
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.ExtReplicaSet), objectStore))
	neh.Add("Replication Controllers", "replication-controllers",
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.ReplicationController), objectStore))
	neh.Add("Stateful Sets", "stateful-sets",
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.StatefulSet), objectStore))

	children, err := neh.Generate(prefix, namespace, "")

	if err != nil {
		return nil, false, err
	}

	return children, false, nil
}

func discoAndLBEntries(ctx context.Context, prefix, namespace string, objectStore store.Store, _ bool) ([]navigation.Navigation, bool, error) {
	neh := navigation.EntriesHelper{}

	neh.Add("Horizontal Pod Autoscalers", "horizontal-pod-autoscalers",
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.HorizontalPodAutoscaler), objectStore))
	neh.Add("Ingresses", "ingresses",
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.Ingress), objectStore))
	neh.Add("Services", "services",
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.Service), objectStore))
	neh.Add("Network Policies", "network-policies",
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.NetworkPolicy), objectStore))

	children, err := neh.Generate(prefix, namespace, "")
	if err != nil {
		return nil, false, err
	}

	return children, false, nil
}

func configAndStorageEntries(ctx context.Context, prefix, namespace string, objectStore store.Store, _ bool) ([]navigation.Navigation, bool, error) {
	neh := navigation.EntriesHelper{}

	neh.Add("Config Maps", "config-maps",
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.ConfigMap), objectStore))
	neh.Add("Persistent Volume Claims", "persistent-volume-claims",
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.PersistentVolumeClaim), objectStore))
	neh.Add("Secrets", "secrets",
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.Secret), objectStore))
	neh.Add("Service Accounts", "service-accounts",
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.ServiceAccount), objectStore))

	children, err := neh.Generate(prefix, namespace, "")
	if err != nil {
		return nil, false, err
	}

	return children, false, nil
}

func rbacEntries(ctx context.Context, prefix, namespace string, objectStore store.Store, _ bool) ([]navigation.Navigation, bool, error) {
	neh := navigation.EntriesHelper{}

	neh.Add("Roles", "roles",
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.Role), objectStore))
	neh.Add("Role Bindings", "role-bindings",
		loading.IsObjectLoading(ctx, namespace, store.KeyFromGroupVersionKind(gvk.RoleBinding), objectStore))

	children, err := neh.Generate(prefix, namespace, "")
	if err != nil {
		return nil, false, err
	}

	return children, false, nil
}
