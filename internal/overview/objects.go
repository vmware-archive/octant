package overview

import (
	"github.com/heptio/developer-dash/pkg/cacheutil"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
)

var (
	workloadsCronJobs = NewResource(ResourceOptions{
		Path:       "/workloads/cron-jobs",
		CacheKey:   cacheutil.Key{APIVersion: "batch/v1beta1", Kind: "CronJob"},
		ListType:   &batchv1beta1.CronJobList{},
		ObjectType: &batchv1beta1.CronJob{},
		Titles:     ResourceTitle{List: "Cron Jobs", Object: "Cron Job"},
	})

	workloadsDaemonSets = NewResource(ResourceOptions{
		Path:       "/workloads/daemon-sets",
		CacheKey:   cacheutil.Key{APIVersion: "apps/v1", Kind: "DaemonSet"},
		ListType:   &appsv1.DaemonSetList{},
		ObjectType: &appsv1.DaemonSet{},
		Titles:     ResourceTitle{List: "Daemon Sets", Object: "Daemon Set"},
	})

	workloadsDeployments = NewResource(ResourceOptions{
		Path:       "/workloads/deployments",
		CacheKey:   cacheutil.Key{APIVersion: "apps/v1", Kind: "Deployment"},
		ListType:   &appsv1.DeploymentList{},
		ObjectType: &appsv1.Deployment{},
		Titles:     ResourceTitle{List: "Deployments", Object: "Deployment"},
	})

	workloadsJobs = NewResource(ResourceOptions{
		Path:       "/workloads/jobs",
		CacheKey:   cacheutil.Key{APIVersion: "batch/v1", Kind: "Job"},
		ListType:   &batchv1.JobList{},
		ObjectType: &batchv1.Job{},

		Titles: ResourceTitle{List: "Jobs", Object: "Job"},
	})

	workloadsPods = NewResource(ResourceOptions{
		Path:       "/workloads/pods",
		CacheKey:   cacheutil.Key{APIVersion: "v1", Kind: "Pod"},
		ListType:   &corev1.PodList{},
		ObjectType: &corev1.Pod{},
		Titles:     ResourceTitle{List: "Pods", Object: "Pod"},
	})

	workloadsReplicaSets = NewResource(ResourceOptions{
		Path:       "/workloads/replica-sets",
		CacheKey:   cacheutil.Key{APIVersion: "apps/v1", Kind: "ReplicaSet"},
		ListType:   &appsv1.ReplicaSetList{},
		ObjectType: &appsv1.ReplicaSet{},
		Titles:     ResourceTitle{List: "Replica Sets", Object: "Replica Set"},
	})

	workloadsReplicationControllers = NewResource(ResourceOptions{
		Path:       "/workloads/replication-controllers",
		CacheKey:   cacheutil.Key{APIVersion: "v1", Kind: "ReplicationController"},
		ListType:   &corev1.ReplicationControllerList{},
		ObjectType: &corev1.ReplicationController{},
		Titles:     ResourceTitle{List: "Replication Controllers", Object: "Replication Controller"},
	})
	workloadsStatefulSets = NewResource(ResourceOptions{
		Path:       "/workloads/stateful-sets",
		CacheKey:   cacheutil.Key{APIVersion: "apps/v1", Kind: "StatefulSet"},
		ListType:   &appsv1.StatefulSetList{},
		ObjectType: &appsv1.StatefulSet{},
		Titles:     ResourceTitle{List: "Stateful Sets", Object: "Stateful Set"},
	})

	workloadsDescriber = NewSectionDescriber(
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

	dlbIngresses = NewResource(ResourceOptions{
		Path:       "/discovery-and-load-balancing/ingresses",
		CacheKey:   cacheutil.Key{APIVersion: "extensions/v1beta1", Kind: "Ingress"},
		ListType:   &v1beta1.IngressList{},
		ObjectType: &v1beta1.Ingress{},
		Titles:     ResourceTitle{List: "Ingresses", Object: "Ingress"},
	})

	dlbServices = NewResource(ResourceOptions{
		Path:       "/discovery-and-load-balancing/services",
		CacheKey:   cacheutil.Key{APIVersion: "v1", Kind: "Service"},
		ListType:   &corev1.ServiceList{},
		ObjectType: &corev1.Service{},
		Titles:     ResourceTitle{List: "Services", Object: "Service"},
	})

	discoveryAndLoadBalancingDescriber = NewSectionDescriber(
		"/discovery-and-load-balancing",
		"Discovery and Load Balancing",
		dlbIngresses,
		dlbServices,
	)

	csConfigMaps = NewResource(ResourceOptions{
		Path:       "/config-and-storage/config-maps",
		CacheKey:   cacheutil.Key{APIVersion: "v1", Kind: "ConfigMap"},
		ListType:   &corev1.ConfigMapList{},
		ObjectType: &corev1.ConfigMap{},
		Titles:     ResourceTitle{List: "Config Maps", Object: "Config Map"},
	})

	csPVCs = NewResource(ResourceOptions{
		Path:       "/config-and-storage/persistent-volume-claims",
		CacheKey:   cacheutil.Key{APIVersion: "v1", Kind: "PersistentVolumeClaim"},
		ListType:   &corev1.PersistentVolumeClaimList{},
		ObjectType: &corev1.PersistentVolumeClaim{},
		Titles:     ResourceTitle{List: "Persistent Volume Claims", Object: "Persistent Volume Claim"},
	})

	csSecrets = NewResource(ResourceOptions{
		Path:       "/config-and-storage/secrets",
		CacheKey:   cacheutil.Key{APIVersion: "v1", Kind: "Secret"},
		ListType:   &corev1.SecretList{},
		ObjectType: &corev1.Secret{},
		Titles:     ResourceTitle{List: "Secrets", Object: "Secret"},
	})

	csServiceAccounts = NewResource(ResourceOptions{
		Path:       "/config-and-storage/service-accounts",
		CacheKey:   cacheutil.Key{APIVersion: "v1", Kind: "ServiceAccount"},
		ListType:   &corev1.ServiceAccountList{},
		ObjectType: &corev1.ServiceAccount{},
		Titles:     ResourceTitle{List: "Service Accounts", Object: "Service Account"},
	})

	configAndStorageDescriber = NewSectionDescriber(
		"/config-and-storage",
		"Config and Storage",
		csConfigMaps,
		csPVCs,
		csSecrets,
		csServiceAccounts,
	)

	customResourcesDescriber = newCRDSectionDescriber(
		"/custom-resources",
		"Custom Resources",
	)

	rbacRoles = NewResource(ResourceOptions{
		Path:       "/rbac/roles",
		CacheKey:   cacheutil.Key{APIVersion: "rbac.authorization.k8s.io/v1", Kind: "Role"},
		ListType:   &rbacv1.RoleList{},
		ObjectType: &rbacv1.Role{},
		Titles:     ResourceTitle{List: "Roles", Object: "Role"},
	})

	rbacClusterRoles = NewResource(ResourceOptions{
		Path:        "/rbac/cluster-roles",
		CacheKey:    cacheutil.Key{APIVersion: "rbac.authorization.k8s.io/v1", Kind: "ClusterRole"},
		ListType:    &rbacv1.ClusterRoleList{},
		ObjectType:  &rbacv1.ClusterRole{},
		Titles:      ResourceTitle{List: "Cluster Roles", Object: "Cluster Role"},
		ClusterWide: true,
	})

	rbacRoleBindings = NewResource(ResourceOptions{
		Path:       "/rbac/role-bindings",
		CacheKey:   cacheutil.Key{APIVersion: "rbac.authorization.k8s.io/v1", Kind: "RoleBinding"},
		ListType:   &rbacv1.RoleBindingList{},
		ObjectType: &rbacv1.RoleBinding{},
		Titles:     ResourceTitle{List: "Role Bindings", Object: "Role Binding"},
	})

	rbacClusterRoleBindings = NewResource(ResourceOptions{
		Path:        "/rbac/cluster-role-bindings",
		CacheKey:    cacheutil.Key{APIVersion: "rbac.authorization.k8s.io/v1", Kind: "ClusterRoleBinding"},
		ListType:    &rbacv1.ClusterRoleBindingList{},
		ObjectType:  &rbacv1.ClusterRoleBinding{},
		Titles:      ResourceTitle{List: "Cluster Role Bindings", Object: "Cluster Role Binding"},
		ClusterWide: true,
	})

	rbacDescriber = NewSectionDescriber(
		"/rbac",
		"RBAC",
		rbacClusterRoles,
		rbacClusterRoleBindings,
		rbacRoles,
		rbacRoleBindings,
	)

	rootDescriber = NewSectionDescriber(
		"/",
		"Overview",
		workloadsDescriber,
		discoveryAndLoadBalancingDescriber,
		configAndStorageDescriber,
		customResourcesDescriber,
		rbacDescriber,
		portForwardDescriber,
	)

	eventsDescriber = NewResource(ResourceOptions{
		Path:                  "/events",
		CacheKey:              cacheutil.Key{APIVersion: "v1", Kind: "Event"},
		ListType:              &corev1.EventList{},
		ObjectType:            &corev1.Event{},
		Titles:                ResourceTitle{List: "Events", Object: "Event"},
		DisableResourceViewer: true,
	})

	portForwardDescriber = NewPortForwardListDescriber()
)
