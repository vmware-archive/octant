package overview

import (
	"path"

	"github.com/heptio/developer-dash/internal/hcli"
)

func navigationEntries(root string) (*hcli.Navigation, error) {
	n := &hcli.Navigation{
		Title: "Overview",
		Path:  path.Join(root, "/"),
		Children: []*hcli.Navigation{
			{
				Title: "Workloads",
				Path:  path.Join(root, "workloads"),
				Children: []*hcli.Navigation{
					{
						Title: "Cron Jobs",
						Path:  path.Join(root, "workloads/cron-jobs"),
					},
					{
						Title: "Daemon Sets",
						Path:  path.Join(root, "workloads/daemon-sets"),
					},
					{
						Title: "Deployments",
						Path:  path.Join(root, "workloads/deployments"),
					},
					{
						Title: "Jobs",
						Path:  path.Join(root, "workloads/jobs"),
					},
					{
						Title: "Pods",
						Path:  path.Join(root, "workloads/pods"),
					},
					{
						Title: "Replica Sets",
						Path:  path.Join(root, "workloads/replica-sets"),
					},
					{
						Title: "Replication Controllers",
						Path:  path.Join(root, "workloads/replication-controllers"),
					},
					{
						Title: "Stateful Sets",
						Path:  path.Join(root, "workloads/stateful-sets"),
					},
				},
			},
			{
				Title: "Discovery and Load Balancing",
				Path:  path.Join(root, "discovery-and-load-balancing"),
				Children: []*hcli.Navigation{
					{
						Title: "Ingresses",
						Path:  path.Join(root, "discovery-and-load-balancing/ingresses"),
					},
					{
						Title: "Services",
						Path:  path.Join(root, "discovery-and-load-balancing/services"),
					},
				},
			},
			{
				Title: "Config and Storage",
				Path:  path.Join(root, "config-and-storage"),
				Children: []*hcli.Navigation{
					{
						Title: "Config Maps",
						Path:  path.Join(root, "config-and-storage/config-maps"),
					},
					{
						Title: "Persistent Volume Claims",
						Path:  path.Join(root, "config-and-storage/persistent-volume-claims"),
					},
					{
						Title: "Secrets",
						Path:  path.Join(root, "config-and-storage/secrets"),
					},
				},
			},
			{
				Title: "Custom Resources",
				Path:  path.Join(root, "custom-resources"),
			},
			{
				Title: "RBAC",
				Path:  path.Join(root, "rbac"),
				Children: []*hcli.Navigation{
					{
						Title: "Roles",
						Path:  path.Join(root, "rbac/roles"),
					},
					{
						Title: "Role Bindings",
						Path:  path.Join(root, "rbac/role-bindings"),
					},
				},
			},
			{
				Title: "Events",
				Path:  path.Join(root, "events"),
			},
		},
	}

	return n, nil
}
