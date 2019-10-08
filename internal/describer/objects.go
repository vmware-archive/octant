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
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"

	"github.com/vmware/octant/pkg/icon"
	"github.com/vmware/octant/pkg/store"
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
		ObjectStoreKey: store.Key{APIVersion: "batch/v1beta1", Kind: "CronJob"},
		ListType:       &batchv1beta1.CronJobList{},
		ObjectType:     &batchv1beta1.CronJob{},
		Titles:         ResourceTitle{List: "Workloads / Cron Jobs", Object: "Cron Job"},
		IconName:       icon.OverviewCronJob,
	})

	workloadsDaemonSets := NewResource(ResourceOptions{
		Path:           "/workloads/daemon-sets",
		ObjectStoreKey: store.Key{APIVersion: "apps/v1", Kind: "DaemonSet"},
		ListType:       &appsv1.DaemonSetList{},
		ObjectType:     &appsv1.DaemonSet{},
		Titles:         ResourceTitle{List: "Workloads / Daemon Sets", Object: "Daemon Set"},
		IconName:       icon.OverviewDaemonSet,
	})

	workloadsDeployments := NewResource(ResourceOptions{
		Path:           "/workloads/deployments",
		ObjectStoreKey: store.Key{APIVersion: "apps/v1", Kind: "Deployment"},
		ListType:       &appsv1.DeploymentList{},
		ObjectType:     &appsv1.Deployment{},
		Titles:         ResourceTitle{List: "Workloads / Deployments", Object: "Deployment"},
		IconName:       icon.OverviewDeployment,
	})

	workloadsJobs := NewResource(ResourceOptions{
		Path:           "/workloads/jobs",
		ObjectStoreKey: store.Key{APIVersion: "batch/v1", Kind: "Job"},
		ListType:       &batchv1.JobList{},
		ObjectType:     &batchv1.Job{},
		Titles:         ResourceTitle{List: "Workloads / Jobs", Object: "Job"},
		IconName:       icon.OverviewJob,
	})

	workloadsPods := NewResource(ResourceOptions{
		Path:           "/workloads/pods",
		ObjectStoreKey: store.Key{APIVersion: "v1", Kind: "Pod"},
		ListType:       &corev1.PodList{},
		ObjectType:     &corev1.Pod{},
		Titles:         ResourceTitle{List: "Workloads / Pods", Object: "Pod"},
		IconName:       icon.OverviewPod,
	})

	workloadsReplicaSets := NewResource(ResourceOptions{
		Path:           "/workloads/replica-sets",
		ObjectStoreKey: store.Key{APIVersion: "apps/v1", Kind: "ReplicaSet"},
		ListType:       &appsv1.ReplicaSetList{},
		ObjectType:     &appsv1.ReplicaSet{},
		Titles:         ResourceTitle{List: "Workloads / Replica Sets", Object: "Replica Set"},
		IconName:       icon.OverviewReplicaSet,
	})

	workloadsReplicationControllers := NewResource(ResourceOptions{
		Path:           "/workloads/replication-controllers",
		ObjectStoreKey: store.Key{APIVersion: "v1", Kind: "ReplicationController"},
		ListType:       &corev1.ReplicationControllerList{},
		ObjectType:     &corev1.ReplicationController{},
		Titles:         ResourceTitle{List: "Workloads / Replication Controllers", Object: "Replication Controller"},
		IconName:       icon.OverviewReplicationController,
	})
	workloadsStatefulSets := NewResource(ResourceOptions{
		Path:           "/workloads/stateful-sets",
		ObjectStoreKey: store.Key{APIVersion: "apps/v1", Kind: "StatefulSet"},
		ListType:       &appsv1.StatefulSetList{},
		ObjectType:     &appsv1.StatefulSet{},
		Titles:         ResourceTitle{List: "Workloads / Stateful Sets", Object: "Stateful Set"},
		IconName:       icon.OverviewStatefulSet,
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
		Titles:         ResourceTitle{List: "Discovery & Load Balancing / Horizontal Pod Autoscaler", Object: "Horizontal Pod Autoscaler"},
		IconName:       icon.OverviewHorizontalPodAutoscaler,
	})

	dlbIngresses := NewResource(ResourceOptions{
		Path:           "/discovery-and-load-balancing/ingresses",
		ObjectStoreKey: store.Key{APIVersion: "extensions/v1beta1", Kind: "Ingress"},
		ListType:       &v1beta1.IngressList{},
		ObjectType:     &v1beta1.Ingress{},
		Titles:         ResourceTitle{List: "Discovery & Load Balancing / Ingresses", Object: "Ingress"},
		IconName:       icon.OverviewIngress,
	})

	dlbServices := NewResource(ResourceOptions{
		Path:           "/discovery-and-load-balancing/services",
		ObjectStoreKey: store.Key{APIVersion: "v1", Kind: "Service"},
		ListType:       &corev1.ServiceList{},
		ObjectType:     &corev1.Service{},
		Titles:         ResourceTitle{List: "Discovery & Load Balancing / Services", Object: "Service"},
		IconName:       icon.OverviewService,
	})

	discoveryAndLoadBalancingDescriber := NewSection(
		"/discovery-and-load-balancing",
		"Discovery and Load Balancing",
		dlbHorizontalPodAutoscalers,
		dlbIngresses,
		dlbServices,
	)

	csConfigMaps := NewResource(ResourceOptions{
		Path:           "/config-and-storage/config-maps",
		ObjectStoreKey: store.Key{APIVersion: "v1", Kind: "ConfigMap"},
		ListType:       &corev1.ConfigMapList{},
		ObjectType:     &corev1.ConfigMap{},
		Titles:         ResourceTitle{List: "Config & Storage / Config Maps", Object: "Config Map"},
		IconName:       icon.OverviewConfigMap,
	})

	csPVCs := NewResource(ResourceOptions{
		Path:           "/config-and-storage/persistent-volume-claims",
		ObjectStoreKey: store.Key{APIVersion: "v1", Kind: "PersistentVolumeClaim"},
		ListType:       &corev1.PersistentVolumeClaimList{},
		ObjectType:     &corev1.PersistentVolumeClaim{},
		Titles:         ResourceTitle{List: "Config & Storage / Persistent Volume Claims", Object: "Persistent Volume Claim"},
		IconName:       icon.OverviewPersistentVolumeClaim,
	})

	csSecrets := NewResource(ResourceOptions{
		Path:           "/config-and-storage/secrets",
		ObjectStoreKey: store.Key{APIVersion: "v1", Kind: "Secret"},
		ListType:       &corev1.SecretList{},
		ObjectType:     &corev1.Secret{},
		Titles:         ResourceTitle{List: "Config & Storage / Secrets", Object: "Secret"},
		IconName:       icon.OverviewSecret,
	})

	csServiceAccounts := NewResource(ResourceOptions{
		Path:           "/config-and-storage/service-accounts",
		ObjectStoreKey: store.Key{APIVersion: "v1", Kind: "ServiceAccount"},
		ListType:       &corev1.ServiceAccountList{},
		ObjectType:     &corev1.ServiceAccount{},
		Titles:         ResourceTitle{List: "Config & Storage / Service Accounts", Object: "Service Account"},
		IconName:       icon.OverviewServiceAccount,
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
		Titles:         ResourceTitle{List: "RBAC / Roles", Object: "Role"},
		IconName:       icon.OverviewRole,
	})

	rbacRoleBindings := NewResource(ResourceOptions{
		Path:           "/rbac/role-bindings",
		ObjectStoreKey: store.Key{APIVersion: "rbac.authorization.k8s.io/v1", Kind: "RoleBinding"},
		ListType:       &rbacv1.RoleBindingList{},
		ObjectType:     &rbacv1.RoleBinding{},
		Titles:         ResourceTitle{List: "RBAC / Role Bindings", Object: "Role Binding"},
		IconName:       icon.OverviewRoleBinding,
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
		Titles:                ResourceTitle{List: "Events", Object: "Event"},
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
