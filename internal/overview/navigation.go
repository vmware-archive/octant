package overview

// Navigation is a set of navigation entries.
type Navigation struct {
	Title    string        `json:"title,omitempty"`
	Path     string        `json:"path,omitempty"`
	Children []*Navigation `json:"children,omitempty"`
}

func navigationEntries() (*Navigation, error) {
	n := &Navigation{
		Title: "Overview",
		Path:  "/overview",
		Children: []*Navigation{
			{
				Title: "Workloads",
				Path:  "/overview/workloads",
				Children: []*Navigation{
					{
						Title: "Cron Jobs",
						Path:  "/overview/workloads/cron-jobs",
					},
					{
						Title: "Daemon Sets",
						Path:  "/overview/workloads/daemon-sets",
					},
					{
						Title: "Deployments",
						Path:  "/overview/workloads/deployments",
					},
					{
						Title: "Jobs",
						Path:  "/overview/workloads/jobs",
					},
					{
						Title: "Pods",
						Path:  "/overview/workloads/pods",
					},
					{
						Title: "Replica Sets",
						Path:  "/overview/workloads/replica-sets",
					},
					{
						Title: "Replication Controllers",
						Path:  "/overview/workloads/replication-controllers",
					},
					{
						Title: "Stateful Sets",
						Path:  "/overview/workloads/stateful-sets",
					},
				},
			},
			{
				Title: "Discovery and Load Balancing",
				Path:  "/overview/discovery-and-load-balancing",
				Children: []*Navigation{
					{
						Title: "Ingresses",
						Path:  "/overview/discovery-and-load-balancing/ingresses",
					},
					{
						Title: "Services",
						Path:  "/overview/discovery-and-load-balancing/services",
					},
				},
			},
			{
				Title: "Config and Storage",
				Path:  "/overview/config-and-storage",
				Children: []*Navigation{
					{
						Title: "Config Maps",
						Path:  "/overview/config-and-storage/config-maps",
					},
					{
						Title: "Persistent Volume Claims",
						Path:  "/overview/config-and-storage/persistent-volume-claims",
					},
					{
						Title: "Secrets",
						Path:  "/overview/config-and-storage/secrets",
					},
				},
			},
			{
				Title: "Custom Resources",
				Path:  "/overview/custom-resources",
			},
			{
				Title: "RBAC",
				Path:  "/overview/rbac",
				Children: []*Navigation{
					{
						Title: "Roles",
						Path:  "/overview/rbac/roles",
					},
					{
						Title: "Role Bindings",
						Path:  "/overview/rbac/role-bindings",
					},
				},
			},
			{
				Title: "Events",
				Path:  "/overview/events",
			},
		},
	}

	return n, nil
}
