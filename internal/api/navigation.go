package api

import (
	"encoding/json"
	"net/http"
)

type navigationResponse struct {
	Name     string               `json:"name,omitempty"`
	Path     string               `json:"path,omitempty"`
	Children []navigationResponse `json:"children,omitempty"`
}

type navigationsResponse struct {
	Navigation []navigationResponse `json:"navigation,omitempty"`
}

type navigation struct{}

var _ http.Handler = (*navigation)(nil)

func (n *navigation) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	nr := &navigationsResponse{
		Navigation: []navigationResponse{
			{
				Name: "Workloads",
				Path: "/overview/workloads",
				Children: []navigationResponse{
					{
						Name: "Cron Jobs",
						Path: "/overview/workloads/cron-jobs",
					},
					{
						Name: "Daemon Sets",
						Path: "/overview/workloads/daemon-sets",
					},
					{
						Name: "Deployments",
						Path: "/overview/workloads/deployments",
					},
					{
						Name: "Jobs",
						Path: "/overview/workloads/jobs",
					},
					{
						Name: "Pods",
						Path: "/overview/workloads/pods",
					},
					{
						Name: "Replica Sets",
						Path: "/overview/workloads/replica-sets",
					},
					{
						Name: "Replication Controllers",
						Path: "/overview/workloads/replication-controllers",
					},
					{
						Name: "Stateful Sets",
						Path: "/overview/workloads/stateful-sets",
					},
				},
			},
			{
				Name: "Discovery and Load Balancing",
				Path: "/overview/discovery-and-load-balancing",
				Children: []navigationResponse{
					{
						Name: "Ingresses",
						Path: "/overview/discovery-and-load-balancing/ingresses",
					},
					{
						Name: "Services",
						Path: "/overview/discovery-and-load-balancing/services",
					},
				},
			},
			{
				Name: "Config and Storage",
				Path: "/overview/config-and-storage",
				Children: []navigationResponse{
					{
						Name: "Config Maps",
						Path: "/overview/config-and-storage/config-maps",
					},
					{
						Name: "Persistent Volume Claims",
						Path: "/overview/config-and-storage/persistent-volume-claims",
					},
					{
						Name: "Secrets",
						Path: "/overview/config-and-storage/secrets",
					},
				},
			},
			{
				Name: "Custom Resources",
				Path: "/overview/custom-resources",
			},
			{
				Name: "RBAC",
				Path: "/overview/rbac",
				Children: []navigationResponse{
					{
						Name: "Roles",
						Path: "/overview/rbac/roles",
					},
					{
						Name: "Role Bindings",
						Path: "/overview/rbac/role-bindings",
					},
				},
			},
			{
				Name: "Events",
				Path: "/overview/events",
			},
		},
	}

	json.NewEncoder(w).Encode(nr)
}
