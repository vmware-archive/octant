/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package overview

import (
	"context"

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
		loading.IsObjectLoading(ctx, namespace, workloadsCronJobs, objectStore))
	neh.Add("Daemon Sets", "daemon-sets", icon.OverviewDaemonSet,
		loading.IsObjectLoading(ctx, namespace, workloadsDaemonSets, objectStore))
	neh.Add("Deployments", "deployments", icon.OverviewDeployment,
		loading.IsObjectLoading(ctx, namespace, workloadsDeployments, objectStore))
	neh.Add("Jobs", "jobs", icon.OverviewJob,
		loading.IsObjectLoading(ctx, namespace, workloadsJobs, objectStore))
	neh.Add("Pods", "pods", icon.OverviewPod,
		loading.IsObjectLoading(ctx, namespace, workloadsPods, objectStore))
	neh.Add("Replica Sets", "replica-sets", icon.OverviewReplicaSet,
		loading.IsObjectLoading(ctx, namespace, workloadsReplicaSets, objectStore))
	neh.Add("Replication Controllers", "replication-controllers", icon.OverviewReplicationController,
		loading.IsObjectLoading(ctx, namespace, workloadsReplicationControllers, objectStore))
	neh.Add("Stateful Sets", "stateful-sets", icon.OverviewStatefulSet,
		loading.IsObjectLoading(ctx, namespace, workloadsStatefulSets, objectStore))

	children, err := neh.Generate(prefix)
	if err != nil {
		return nil, false, err
	}

	return children, false, nil
}

func discoAndLBEntries(ctx context.Context, prefix, namespace string, objectStore store.Store, _ bool) ([]navigation.Navigation, bool, error) {
	neh := navigation.EntriesHelper{}
	neh.Add("Ingresses", "ingresses", icon.OverviewIngress,
		loading.IsObjectLoading(ctx, namespace, dlbIngresses, objectStore))
	neh.Add("Services", "services", icon.OverviewService,
		loading.IsObjectLoading(ctx, namespace, dlbServices, objectStore))

	children, err := neh.Generate(prefix)
	if err != nil {
		return nil, false, err
	}

	return children, false, nil
}

func configAndStorageEntries(ctx context.Context, prefix, namespace string, objectStore store.Store, _ bool) ([]navigation.Navigation, bool, error) {
	neh := navigation.EntriesHelper{}
	neh.Add("Config Maps", "config-maps", icon.OverviewConfigMap,
		loading.IsObjectLoading(ctx, namespace, csConfigMaps, objectStore))
	neh.Add("Persistent Volume Claims", "persistent-volume-claims", icon.OverviewPersistentVolumeClaim,
		loading.IsObjectLoading(ctx, namespace, csPVCs, objectStore))
	neh.Add("Secrets", "secrets", icon.OverviewSecret,
		loading.IsObjectLoading(ctx, namespace, csSecrets, objectStore))
	neh.Add("Service Accounts", "service-accounts", icon.OverviewServiceAccount,
		loading.IsObjectLoading(ctx, namespace, csServiceAccounts, objectStore))

	children, err := neh.Generate(prefix)
	if err != nil {
		return nil, false, err
	}

	return children, false, nil
}

func rbacEntries(ctx context.Context, prefix, namespace string, objectStore store.Store, _ bool) ([]navigation.Navigation, bool, error) {
	neh := navigation.EntriesHelper{}

	neh.Add("Roles", "roles", icon.OverviewRole,
		loading.IsObjectLoading(ctx, namespace, rbacRoles, objectStore))
	neh.Add("Role Bindings", "role-bindings", icon.OverviewRoleBinding,
		loading.IsObjectLoading(ctx, namespace, rbacRoleBindings, objectStore))

	children, err := neh.Generate(prefix)
	if err != nil {
		return nil, false, err
	}

	return children, false, nil
}
