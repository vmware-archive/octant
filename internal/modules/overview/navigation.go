package overview

import (
	"context"
	"path"

	"github.com/heptio/developer-dash/internal/clustereye"
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

func workloadEntries(_ context.Context, prefix, _ string, _ store.Store) ([]clustereye.Navigation, error) {
	return []clustereye.Navigation{
		*clustereye.NewNavigation("Cron Jobs", path.Join(prefix, "cron-jobs")),
		*clustereye.NewNavigation("Daemon Sets", path.Join(prefix, "daemon-sets")),
		*clustereye.NewNavigation("Deployments", path.Join(prefix, "deployments")),
		*clustereye.NewNavigation("Jobs", path.Join(prefix, "jobs")),
		*clustereye.NewNavigation("Pods", path.Join(prefix, "pods")),
		*clustereye.NewNavigation("Replica Sets", path.Join(prefix, "replica-sets")),
		*clustereye.NewNavigation("Replication Controllers", path.Join(prefix, "replication-controllers")),
		*clustereye.NewNavigation("Stateful Sets", path.Join(prefix, "stateful-sets")),
	}, nil
}

func discoAndLBEntries(_ context.Context, prefix, _ string, _ store.Store) ([]clustereye.Navigation, error) {
	return []clustereye.Navigation{
		*clustereye.NewNavigation("Ingresses", path.Join(prefix, "ingresses")),
		*clustereye.NewNavigation("Services", path.Join(prefix, "services")),
	}, nil
}

func configAndStorageEntries(_ context.Context, prefix, _ string, _ store.Store) ([]clustereye.Navigation, error) {
	return []clustereye.Navigation{
		*clustereye.NewNavigation("Config Maps", path.Join(prefix, "config-maps")),
		*clustereye.NewNavigation("Persistent Volume Claims", path.Join(prefix, "persistent-volume-claims")),
		*clustereye.NewNavigation("Secrets", path.Join(prefix, "secrets")),
		*clustereye.NewNavigation("Service Accounts", path.Join(prefix, "service-accounts")),
	}, nil
}

func rbacEntries(_ context.Context, prefix, _ string, _ store.Store) ([]clustereye.Navigation, error) {
	return []clustereye.Navigation{
		*clustereye.NewNavigation("Roles", path.Join(prefix, "roles")),
		*clustereye.NewNavigation("Role Bindings", path.Join(prefix, "role-bindings")),
	}, nil
}


