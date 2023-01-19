/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package describer

import (
	"sync"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"

	"github.com/vmware-tanzu/octant/pkg/store"
)

var (
	namespacedOverviewOnce = sync.Once{}
	namespacedOverview     *Section
	namespacedCRDOnce      = sync.Once{}
	namespacedCRD          *CRDSection
)

// NamespacedObjects returns a describer for a namespaced overview.
func NamespacedOverview() *Section {
	namespacedOverviewOnce.Do(func() {
		namespacedOverview = initNamespacedOverview()
	})

	return namespacedOverview
}

// NamespacedCRD returns a describer for namespaces CRDs.
func NamespacedCRD() *CRDSection {
	namespacedCRDOnce.Do(func() {
		namespacedCRD = initNamespacedCRD()
	})

	return namespacedCRD
}

func initNamespacedCRD() *CRDSection {
	return NewCRDSection(
		"/custom-resources",
		"Custom Resources",
	)
}

func initNamespacedOverview() *Section {
	workloadsCronJobs := NewResource(ResourceOptions{
		Path:           "/workloads/cron-jobs",
		ObjectStoreKey: store.Key{APIVersion: "batch/v1", Kind: "CronJob"},
		ListType:       &batchv1.CronJobList{},
		ObjectType:     &batchv1.CronJob{},
		Titles:         ResourceTitle{List: "Cron Jobs", Object: "Cron Jobs"},
	})

	workloadsDaemonSets := NewResource(ResourceOptions{
		Path:           "/workloads/daemon-sets",
		ObjectStoreKey: store.Key{APIVersion: "apps/v1", Kind: "DaemonSet"},
		ListType:       &appsv1.DaemonSetList{},
		ObjectType:     &appsv1.DaemonSet{},
		Titles:         ResourceTitle{List: "Daemon Sets", Object: "Daemon Sets"},
	})

	workloadsDeployments := NewResource(ResourceOptions{
		Path:           "/workloads/deployments",
		ObjectStoreKey: store.Key{APIVersion: "apps/v1", Kind: "Deployment"},
		ListType:       &appsv1.DeploymentList{},
		ObjectType:     &appsv1.Deployment{},
		Titles:         ResourceTitle{List: "Deployments", Object: "Deployments"},
	})

	workloadsJobs := NewResource(ResourceOptions{
		Path:           "/workloads/jobs",
		ObjectStoreKey: store.Key{APIVersion: "batch/v1", Kind: "Job"},
		ListType:       &batchv1.JobList{},
		ObjectType:     &batchv1.Job{},
		Titles:         ResourceTitle{List: "Jobs", Object: "Jobs"},
	})

	workloadsPods := NewResource(ResourceOptions{
		Path:           "/workloads/pods",
		ObjectStoreKey: store.Key{APIVersion: "v1", Kind: "Pod"},
		ListType:       &corev1.PodList{},
		ObjectType:     &corev1.Pod{},
		Titles:         ResourceTitle{List: "Pods", Object: "Pods"},
	})

	workloadsReplicaSets := NewResource(ResourceOptions{
		Path:           "/workloads/replica-sets",
		ObjectStoreKey: store.Key{APIVersion: "apps/v1", Kind: "ReplicaSet"},
		ListType:       &appsv1.ReplicaSetList{},
		ObjectType:     &appsv1.ReplicaSet{},
		Titles:         ResourceTitle{List: "Replica Sets", Object: "Replica Sets"},
	})

	workloadsReplicationControllers := NewResource(ResourceOptions{
		Path:           "/workloads/replication-controllers",
		ObjectStoreKey: store.Key{APIVersion: "v1", Kind: "ReplicationController"},
		ListType:       &corev1.ReplicationControllerList{},
		ObjectType:     &corev1.ReplicationController{},
		Titles:         ResourceTitle{List: "Replication Controllers", Object: "Replication Controllers"},
	})
	workloadsStatefulSets := NewResource(ResourceOptions{
		Path:           "/workloads/stateful-sets",
		ObjectStoreKey: store.Key{APIVersion: "apps/v1", Kind: "StatefulSet"},
		ListType:       &appsv1.StatefulSetList{},
		ObjectType:     &appsv1.StatefulSet{},
		Titles:         ResourceTitle{List: "Stateful Sets", Object: "Stateful Sets"},
	})

	workloadsDescriber := NewSection(
		"/workloads",
		"Workloads",
		workloadsCronJobs,
		workloadsDaemonSets,
		workloadsDeployments,
		workloadsJobs,
		workloadsPods,
		workloadsReplicaSets,
		workloadsReplicationControllers,
		workloadsStatefulSets,
	)

	dlbHorizontalPodAutoscalers := NewResource(ResourceOptions{
		Path:           "/discovery-and-load-balancing/horizontal-pod-autoscalers",
		ObjectStoreKey: store.Key{APIVersion: "autoscaling/v1", Kind: "HorizontalPodAutoscaler"},
		ListType:       &autoscalingv1.HorizontalPodAutoscalerList{},
		ObjectType:     &autoscalingv1.HorizontalPodAutoscaler{},
		Titles:         ResourceTitle{List: "Horizontal Pod Autoscalers", Object: "Horizontal Pod Autoscalers"},
	})

	dlbIngresses := NewResource(ResourceOptions{
		Path:           "/discovery-and-load-balancing/ingresses",
		ObjectStoreKey: store.Key{APIVersion: "networking.k8s.io/v1", Kind: "Ingress"},
		ListType:       &networkingv1.IngressList{},
		ObjectType:     &networkingv1.Ingress{},
		Titles:         ResourceTitle{List: "Ingresses", Object: "Ingresses"},
	})

	dlbServices := NewResource(ResourceOptions{
		Path:           "/discovery-and-load-balancing/services",
		ObjectStoreKey: store.Key{APIVersion: "v1", Kind: "Service"},
		ListType:       &corev1.ServiceList{},
		ObjectType:     &corev1.Service{},
		Titles:         ResourceTitle{List: "Services", Object: "Services"},
	})

	dlbNetworkPolicies := NewResource(ResourceOptions{
		Path:           "/discovery-and-load-balancing/network-policies",
		ObjectStoreKey: store.Key{APIVersion: "networking.k8s.io/v1", Kind: "NetworkPolicy"},
		ListType:       &networkingv1.NetworkPolicyList{},
		ObjectType:     &networkingv1.NetworkPolicy{},
		Titles:         ResourceTitle{List: "Network Policies", Object: "Network Policy"},
	})

	discoveryAndLoadBalancingDescriber := NewSection(
		"/discovery-and-load-balancing",
		"Discovery and Load Balancing",
		dlbHorizontalPodAutoscalers,
		dlbIngresses,
		dlbServices,
		dlbNetworkPolicies,
	)

	csConfigMaps := NewResource(ResourceOptions{
		Path:           "/config-and-storage/config-maps",
		ObjectStoreKey: store.Key{APIVersion: "v1", Kind: "ConfigMap"},
		ListType:       &corev1.ConfigMapList{},
		ObjectType:     &corev1.ConfigMap{},
		Titles:         ResourceTitle{List: "Config Maps", Object: "Config Maps"},
	})

	csPVCs := NewResource(ResourceOptions{
		Path:           "/config-and-storage/persistent-volume-claims",
		ObjectStoreKey: store.Key{APIVersion: "v1", Kind: "PersistentVolumeClaim"},
		ListType:       &corev1.PersistentVolumeClaimList{},
		ObjectType:     &corev1.PersistentVolumeClaim{},
		Titles:         ResourceTitle{List: "Persistent Volume Claims", Object: "Persistent Volume Claims"},
	})

	csSecrets := NewResource(ResourceOptions{
		Path:           "/config-and-storage/secrets",
		ObjectStoreKey: store.Key{APIVersion: "v1", Kind: "Secret"},
		ListType:       &corev1.SecretList{},
		ObjectType:     &corev1.Secret{},
		Titles:         ResourceTitle{List: "Secrets", Object: "Secrets"},
	})

	csServiceAccounts := NewResource(ResourceOptions{
		Path:           "/config-and-storage/service-accounts",
		ObjectStoreKey: store.Key{APIVersion: "v1", Kind: "ServiceAccount"},
		ListType:       &corev1.ServiceAccountList{},
		ObjectType:     &corev1.ServiceAccount{},
		Titles:         ResourceTitle{List: "Service Accounts", Object: "Service Accounts"},
	})

	configAndStorageDescriber := NewSection(
		"/config-and-storage",
		"Config and Storage",
		csConfigMaps,
		csPVCs,
		csSecrets,
		csServiceAccounts,
	)

	rbacRoles := NewResource(ResourceOptions{
		Path:           "/rbac/roles",
		ObjectStoreKey: store.Key{APIVersion: "rbac.authorization.k8s.io/v1", Kind: "Role"},
		ListType:       &rbacv1.RoleList{},
		ObjectType:     &rbacv1.Role{},
		Titles:         ResourceTitle{List: "Roles", Object: "Roles"},
	})

	rbacRoleBindings := NewResource(ResourceOptions{
		Path:           "/rbac/role-bindings",
		ObjectStoreKey: store.Key{APIVersion: "rbac.authorization.k8s.io/v1", Kind: "RoleBinding"},
		ListType:       &rbacv1.RoleBindingList{},
		ObjectType:     &rbacv1.RoleBinding{},
		Titles:         ResourceTitle{List: "Role Bindings", Object: "Role Bindings"},
	})

	rbacDescriber := NewSection(
		"/rbac",
		"RBAC",
		rbacRoles,
		rbacRoleBindings,
	)

	eventsDescriber := NewResource(ResourceOptions{
		Path:                  "/events",
		ObjectStoreKey:        store.Key{APIVersion: "v1", Kind: "Event"},
		ListType:              &corev1.EventList{},
		ObjectType:            &corev1.Event{},
		Titles:                ResourceTitle{List: "Events", Object: "Events"},
		DisableResourceViewer: true,
	})

	rootDescriber := NewSection(
		"/",
		"Overview",
		workloadsDescriber,
		discoveryAndLoadBalancingDescriber,
		configAndStorageDescriber,
		NamespacedCRD(),
		rbacDescriber,
		eventsDescriber,
	)

	return rootDescriber
}
