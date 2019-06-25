/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package overview

import (
	"context"

	"github.com/vmware/octant/internal/icon"
	"github.com/vmware/octant/internal/octant"
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

func workloadEntries(_ context.Context, prefix, _ string, _ store.Store) ([]octant.Navigation, error) {
	neh := octant.NavigationEntriesHelper{}
	neh.Add("Cron Jobs", "cron-jobs", icon.OverviewCronJob)
	neh.Add("Daemon Sets", "daemon-sets", icon.OverviewDaemonSet)
	neh.Add("Deployments", "deployments", icon.OverviewDeployment)
	neh.Add("Jobs", "jobs", icon.OverviewJob)
	neh.Add("Pods", "pods", icon.OverviewPod)
	neh.Add("Replica Sets", "replica-sets", icon.OverviewReplicaSet)
	neh.Add("Replication Controllers", "replication-controllers", icon.OverviewReplicationController)
	neh.Add("Stateful Sets", "stateful-sets", icon.OverviewStatefulSet)

	return neh.Generate(prefix)
}

func discoAndLBEntries(_ context.Context, prefix, _ string, _ store.Store) ([]octant.Navigation, error) {
	neh := octant.NavigationEntriesHelper{}
	neh.Add("Ingresses", "ingresses", icon.OverviewIngress)
	neh.Add("Services", "services", icon.OverviewService)

	return neh.Generate(prefix)
}

func configAndStorageEntries(_ context.Context, prefix, _ string, _ store.Store) ([]octant.Navigation, error) {
	neh := octant.NavigationEntriesHelper{}
	neh.Add("Config Maps", "config-maps", icon.OverviewConfigMap)
	neh.Add("Persistent Volume Claims", "persistent-volume-claims", icon.OverviewPersistentVolumeClaim)
	neh.Add("Secrets", "secrets", icon.OverviewSecret)
	neh.Add("Service Accounts", "service-accounts", icon.OverviewServiceAccount)

	return neh.Generate(prefix)
}

func rbacEntries(_ context.Context, prefix, _ string, _ store.Store) ([]octant.Navigation, error) {
	neh := octant.NavigationEntriesHelper{}

	neh.Add("Roles", "roles", icon.OverviewRole)
	neh.Add("Role Bindings", "role-bindings", icon.OverviewRoleBinding)

	return neh.Generate(prefix)
}
