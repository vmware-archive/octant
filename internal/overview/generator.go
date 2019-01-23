package overview

import (
	"context"
	"fmt"
	"regexp"
	"sync"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/overview/printer"
	"github.com/heptio/developer-dash/internal/view"
	"github.com/heptio/developer-dash/internal/view/component"

	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
)

type pathFilter struct {
	path      string
	describer Describer

	re *regexp.Regexp
}

func newPathFilter(path string, describer Describer) *pathFilter {
	re := regexp.MustCompile(fmt.Sprintf("^%s/?$", path))

	return &pathFilter{
		re:        re,
		path:      path,
		describer: describer,
	}
}

func (pf *pathFilter) Match(path string) bool {
	return pf.re.MatchString(path)
}

func (pf *pathFilter) Fields(path string) map[string]string {
	out := make(map[string]string)

	match := pf.re.FindStringSubmatch(path)
	for i, name := range pf.re.SubexpNames() {
		if i != 0 && name != "" {
			out[name] = match[i]
		}
	}

	return out
}

var (
	workloadsCronJobs = NewResource(ResourceOptions{
		Path:       "/workloads/cron-jobs",
		CacheKey:   cache.Key{APIVersion: "batch/v1beta1", Kind: "CronJob"},
		ListType:   &batchv1beta1.CronJobList{},
		ObjectType: &batchv1beta1.CronJob{},
		Titles:     ResourceTitle{List: "Cron Jobs", Object: "Cron Job"},
		Transforms: cronJobTransforms,
		Sections: []ContentSection{
			{
				Views: []view.ViewFactory{
					NewCronJobSummary,
					NewCronJobJobs,
					NewEventList,
				},
				Title: "Summary",
			},
			{
				Views: []view.ViewFactory{
					NewResourceViewerStub,
				},
				Title: "Resource Viewer",
			},
		},
	})

	workloadsDaemonSets = NewResource(ResourceOptions{
		Path:       "/workloads/daemon-sets",
		CacheKey:   cache.Key{APIVersion: "apps/v1", Kind: "DaemonSet"},
		ListType:   &appsv1.DaemonSetList{},
		ObjectType: &appsv1.DaemonSet{},
		Titles:     ResourceTitle{List: "Daemon Sets", Object: "Daemon Set"},
		Transforms: daemonSetTransforms,
		Sections: []ContentSection{
			{
				Title: "Summary",
				Views: []view.ViewFactory{
					NewDaemonSetSummary,
					NewContainerSummary,
					NewPodList,
					NewEventList,
				},
			},
			{
				Title: "Resource Viewer",
				Views: []view.ViewFactory{
					workloadViewFactory,
				},
			},
		},
	})

	workloadsDeployments = NewResource(ResourceOptions{
		Path:       "/workloads/deployments",
		CacheKey:   cache.Key{APIVersion: "apps/v1", Kind: "Deployment"},
		ListType:   &appsv1.DeploymentList{},
		ObjectType: &appsv1.Deployment{},
		Titles:     ResourceTitle{List: "Deployments", Object: "Deployment"},
		Transforms: deploymentTransforms,
		Sections: []ContentSection{
			{
				Title: "Summary",
				Views: []view.ViewFactory{
					NewDeploymentSummary,
					NewContainerSummary,
					NewDeploymentReplicaSets,
					NewEventList,
				},
			},
			{
				Title: "Resource Viewer",
				Views: []view.ViewFactory{
					workloadViewFactory,
				},
			},
		},
	})

	workloadsJobs = NewResource(ResourceOptions{
		Path:       "/workloads/jobs",
		CacheKey:   cache.Key{APIVersion: "batch/v1", Kind: "Job"},
		ListType:   &batchv1.JobList{},
		ObjectType: &batchv1.Job{},

		Titles:     ResourceTitle{List: "Jobs", Object: "Job"},
		Transforms: jobTransforms,
		Sections: []ContentSection{
			{
				Views: []view.ViewFactory{
					NewJobSummary,
					NewContainerSummary,
					NewPodList,
					NewEventList,
				},
			},
		},
	})

	workloadsPods = NewResource(ResourceOptions{
		Path:       "/workloads/pods",
		CacheKey:   cache.Key{APIVersion: "v1", Kind: "Pod"},
		ListType:   &corev1.PodList{},
		ObjectType: &corev1.Pod{},
		Titles:     ResourceTitle{List: "Pods", Object: "Pod"},
		Transforms: podTransforms,
		Sections: []ContentSection{
			{
				Title: "Summary",
				Views: []view.ViewFactory{
					NewPodSummary,
					NewPodContainer,
					NewPodCondition,
					NewPodVolume,
					NewEventList,
				},
			},
			{
				Title: "Resource Viewer",
				Views: []view.ViewFactory{
					workloadViewFactory,
				},
			},
		},
	})

	workloadsReplicaSets = NewResource(ResourceOptions{
		Path:       "/workloads/replica-sets",
		CacheKey:   cache.Key{APIVersion: "apps/v1", Kind: "ReplicaSet"},
		ListType:   &appsv1.ReplicaSetList{},
		ObjectType: &appsv1.ReplicaSet{},
		Titles:     ResourceTitle{List: "Replica Sets", Object: "Replica Set"},
		Transforms: replicaSetTransforms,
		Sections: []ContentSection{
			{
				Title: "Summary",
				Views: []view.ViewFactory{
					NewReplicaSetSummary,
					NewContainerSummary,
					NewPodList,
					NewEventList,
				},
			},
			{
				Title: "Resource Viewer",
				Views: []view.ViewFactory{
					workloadViewFactory,
				},
			},
		},
	})

	workloadsReplicationControllers = NewResource(ResourceOptions{
		Path:       "/workloads/replication-controllers",
		CacheKey:   cache.Key{APIVersion: "v1", Kind: "ReplicationController"},
		ListType:   &corev1.ReplicationControllerList{},
		ObjectType: &corev1.ReplicationController{},
		Titles:     ResourceTitle{List: "Replication Controllers", Object: "Replication Controller"},
		Transforms: replicationControllerTransforms,
		Sections: []ContentSection{
			{
				Title: "Summary",
				Views: []view.ViewFactory{
					NewReplicationControllerSummary,
					NewContainerSummary,
					NewPodList,
					NewEventList,
				},
			},
			{
				Title: "Resource Viewer",
				Views: []view.ViewFactory{
					workloadViewFactory,
				},
			},
		},
	})
	workloadsStatefulSets = NewResource(ResourceOptions{
		Path:       "/workloads/stateful-sets",
		CacheKey:   cache.Key{APIVersion: "apps/v1", Kind: "StatefulSet"},
		ListType:   &appsv1.StatefulSetList{},
		ObjectType: &appsv1.StatefulSet{},
		Titles:     ResourceTitle{List: "Stateful Sets", Object: "Stateful Set"},
		Transforms: statefulSetTransforms,
		Sections: []ContentSection{
			{
				Title: "Summary",
				Views: []view.ViewFactory{
					NewStatefulSetSummary,
					NewContainerSummary,
					NewPodList,
					NewEventList,
				},
			},
			{
				Title: "Resource Viewer",
				Views: []view.ViewFactory{
					workloadViewFactory,
				},
			},
		},
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
		CacheKey:   cache.Key{APIVersion: "extensions/v1beta1", Kind: "Ingress"},
		ListType:   &v1beta1.IngressList{},
		ObjectType: &v1beta1.Ingress{},
		Titles:     ResourceTitle{List: "Ingresses", Object: "Ingress"},
		Transforms: ingressTransforms,
		Sections: []ContentSection{
			{
				Title: "Summary",
				Views: []view.ViewFactory{
					NewIngressSummary,
					NewIngressDetails,
					NewEventList,
				},
			},
			{
				Title: "Resource Viewer",
				Views: []view.ViewFactory{
					workloadViewFactory,
				},
			},
		},
	})

	dlbServices = NewResource(ResourceOptions{
		Path:       "/discovery-and-load-balancing/services",
		CacheKey:   cache.Key{APIVersion: "v1", Kind: "Service"},
		ListType:   &corev1.ServiceList{},
		ObjectType: &corev1.Service{},
		Titles:     ResourceTitle{List: "Services", Object: "Service"},
		Transforms: serviceTransforms,
		Sections: []ContentSection{
			{
				Title: "Summary",
				Views: []view.ViewFactory{
					NewServiceSummary,
					NewServicePort,
					NewServiceEndpoints,
					NewEventList,
				},
			},
			{
				Title: "Resource Viewer",
				Views: []view.ViewFactory{
					workloadViewFactory,
				},
			},
		},
	})

	discoveryAndLoadBalancingDescriber = NewSectionDescriber(
		"/discovery-and-load-balancing",
		"Discovery and Load Balancing",
		dlbIngresses,
		dlbServices,
	)

	csConfigMaps = NewResource(ResourceOptions{
		Path:       "/config-and-storage/config-maps",
		CacheKey:   cache.Key{APIVersion: "v1", Kind: "ConfigMap"},
		ListType:   &corev1.ConfigMapList{},
		ObjectType: &corev1.ConfigMap{},
		Titles:     ResourceTitle{List: "Config Maps", Object: "Config Map"},
		Transforms: configMapTransforms,
		Sections: []ContentSection{
			{
				Views: []view.ViewFactory{
					NewConfigMapSummary,
					NewConfigMapDetails,
					NewEventList,
				},
			},
		},
	})

	csPVCs = NewResource(ResourceOptions{
		Path:       "/config-and-storage/persistent-volume-claims",
		CacheKey:   cache.Key{APIVersion: "v1", Kind: "PersistentVolumeClaim"},
		ListType:   &corev1.PersistentVolumeClaimList{},
		ObjectType: &corev1.PersistentVolumeClaim{},
		Titles:     ResourceTitle{List: "Persistent Volume Claims", Object: "Persistent Volume Claim"},
		Transforms: pvcTransforms,
		Sections: []ContentSection{
			{
				Views: []view.ViewFactory{
					NewPersistentVolumeClaimSummary,
					NewEventList,
				},
			},
		},
	})

	csSecrets = NewResource(ResourceOptions{
		Path:       "/config-and-storage/secrets",
		CacheKey:   cache.Key{APIVersion: "v1", Kind: "Secret"},
		ListType:   &corev1.SecretList{},
		ObjectType: &corev1.Secret{},
		Titles:     ResourceTitle{List: "Secrets", Object: "Secret"},
		Transforms: secretTransforms,
		Sections: []ContentSection{
			{
				Views: []view.ViewFactory{
					NewSecretSummary,
					NewSecretData,
					NewEventList,
				},
			},
		},
	})

	csServiceAccounts = NewResource(ResourceOptions{
		Path:       "/config-and-storage/service-accounts",
		CacheKey:   cache.Key{APIVersion: "v1", Kind: "ServiceAccount"},
		ListType:   &corev1.ServiceAccountList{},
		ObjectType: &corev1.ServiceAccount{},
		Titles:     ResourceTitle{List: "Service Accounts", Object: "Service Account"},
		Transforms: serviceAccountTransforms,
		Sections: []ContentSection{
			{
				Views: []view.ViewFactory{
					NewServiceAccountSummary,
					NewEventList,
				},
			},
		},
	})

	configAndStorageDescriber = NewSectionDescriber(
		"/config-and-storage",
		"Config and Storage",
		csConfigMaps,
		csPVCs,
		csSecrets,
		csServiceAccounts,
	)

	customResourcesDescriber = NewSectionDescriber(
		"/custom-resources",
		"Custom Resources",
	)

	rbacRoles = NewResource(ResourceOptions{
		Path:       "/rbac/roles",
		CacheKey:   cache.Key{APIVersion: "rbac.authorization.k8s.io/v1", Kind: "Role"},
		ListType:   &rbacv1.RoleList{},
		ObjectType: &rbacv1.Role{},
		Titles:     ResourceTitle{List: "Roles", Object: "Role"},
		Transforms: roleTransforms,
		Sections: []ContentSection{
			{
				Views: []view.ViewFactory{
					NewRoleSummary,
					NewRoleRule,
					NewEventList,
				},
			},
		},
	})

	rbacRoleBindings = NewResource(ResourceOptions{
		Path:       "/rbac/role-bindings",
		CacheKey:   cache.Key{APIVersion: "rbac.authorization.k8s.io/v1", Kind: "RoleBinding"},
		ListType:   &rbacv1.RoleBindingList{},
		ObjectType: &rbacv1.RoleBinding{},
		Titles:     ResourceTitle{List: "Role Bindings", Object: "Role Binding"},
		Transforms: roleBindingTransforms,
		Sections: []ContentSection{
			{
				Views: []view.ViewFactory{
					NewRoleBindingSummary,
					NewRoleBindingSubjects,
					NewEventList,
				},
			},
		},
	})

	rbacDescriber = NewSectionDescriber(
		"/rbac",
		"RBAC",
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
	)

	eventsDescriber = NewResource(ResourceOptions{
		Path:       "/events",
		CacheKey:   cache.Key{APIVersion: "v1", Kind: "Event"},
		ListType:   &corev1.EventList{},
		ObjectType: &corev1.Event{},
		Titles:     ResourceTitle{List: "Events", Object: "Event"},
		Transforms: roleBindingTransforms,
	})
)

var contentNotFound = errors.Errorf("content not found")

type realGenerator struct {
	cache         cache.Cache
	pathFilters   []pathFilter
	clusterClient cluster.ClientInterface
	printer       printer.Printer

	mu sync.Mutex
}

func newGenerator(cache cache.Cache, pathFilters []pathFilter, clusterClient cluster.ClientInterface) (*realGenerator, error) {
	p := printer.NewResource(cache)

	if err := AddPrintHandlers(p); err != nil {
		return nil, errors.Wrap(err, "add print handlers")
	}

	return &realGenerator{
		cache:         cache,
		pathFilters:   pathFilters,
		clusterClient: clusterClient,
		printer:       p,
	}, nil
}

func (g *realGenerator) Generate(ctx context.Context, path, prefix, namespace string) (component.ContentResponse, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, pf := range g.pathFilters {
		if !pf.Match(path) {
			continue
		}

		fields := pf.Fields(path)
		options := DescriberOptions{
			Cache:   g.cache,
			Fields:  fields,
			Printer: g.printer,
		}

		cResponse, err := pf.describer.Describe(ctx, prefix, namespace, g.clusterClient, options)
		if err != nil {
			return emptyContentResponse, err
		}

		return cResponse, nil
	}

	fmt.Println("content not found for", path)
	return emptyContentResponse, contentNotFound
}

// PrinterHandler configures handlers for a printer.
type PrinterHandler interface {
	Handler(printFunc interface{}) error
}

// AddPrintHandlers adds print handlers to a printer.
func AddPrintHandlers(p PrinterHandler) error {
	handlers := []interface{}{
		printer.DeploymentHandler,
		printer.DeploymentListHandler,
		printer.ReplicaSetHandler,
		printer.ReplicaSetListHandler,
		printer.PodHandler,
		printer.PodListHandler,
	}

	for _, handler := range handlers {
		if err := p.Handler(handler); err != nil {
			return err
		}
	}

	return nil
}
