package overview

import (
	"context"
	"path"

	"github.com/heptio/developer-dash/internal/octant"
	"github.com/heptio/developer-dash/pkg/store"
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
	return []octant.Navigation{
		*octant.NewNavigation("Cron Jobs", path.Join(prefix, "cron-jobs")),
		*octant.NewNavigation("Daemon Sets", path.Join(prefix, "daemon-sets")),
		*octant.NewNavigation("Deployments", path.Join(prefix, "deployments")),
		*octant.NewNavigation("Jobs", path.Join(prefix, "jobs")),
		*octant.NewNavigation("Pods", path.Join(prefix, "pods")),
		*octant.NewNavigation("Replica Sets", path.Join(prefix, "replica-sets")),
		*octant.NewNavigation("Replication Controllers", path.Join(prefix, "replication-controllers")),
		*octant.NewNavigation("Stateful Sets", path.Join(prefix, "stateful-sets")),
	}, nil
}

func discoAndLBEntries(_ context.Context, prefix, _ string, _ store.Store) ([]octant.Navigation, error) {
	return []octant.Navigation{
		*octant.NewNavigation("Ingresses", path.Join(prefix, "ingresses")),
		*octant.NewNavigation("Services", path.Join(prefix, "services")),
	}, nil
}

func configAndStorageEntries(_ context.Context, prefix, _ string, _ store.Store) ([]octant.Navigation, error) {
	return []octant.Navigation{
		*octant.NewNavigation("Config Maps", path.Join(prefix, "config-maps")),
		*octant.NewNavigation("Persistent Volume Claims", path.Join(prefix, "persistent-volume-claims")),
		*octant.NewNavigation("Secrets", path.Join(prefix, "secrets")),
		*octant.NewNavigation("Service Accounts", path.Join(prefix, "service-accounts")),
	}, nil
}

func rbacEntries(_ context.Context, prefix, _ string, _ store.Store) ([]octant.Navigation, error) {
	return []octant.Navigation{
		*octant.NewNavigation("Roles", path.Join(prefix, "roles")),
		*octant.NewNavigation("Role Bindings", path.Join(prefix, "role-bindings")),
	}, nil
}


