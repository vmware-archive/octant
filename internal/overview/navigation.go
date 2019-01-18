package overview

import (
	"path"
	"strings"

	"github.com/heptio/developer-dash/internal/hcli"
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

// NavigationFactory generates navigation entries.
type NavigationFactory struct {
	root string
}

// NewNavigationFactory creates an instance of NewNavigationFactory.
func NewNavigationFactory(root string) *NavigationFactory {
	rootPath := root
	if !strings.HasSuffix(rootPath, "/") {
		rootPath = rootPath + "/"
	}

	return &NavigationFactory{root: rootPath}
}

// Root returns the root of the navigation tree.
func (nf *NavigationFactory) Root() string {
	return nf.root
}

func (nf *NavigationFactory) pathFor(elements ...string) string {
	return path.Join(append([]string{nf.root}, elements...)...)
}

// Entries returns navigation entries.
func (nf *NavigationFactory) Entries() (*hcli.Navigation, error) {

	n := &hcli.Navigation{
		Title: "Overview",
		Path:  nf.root,
		Children: []*hcli.Navigation{

			nf.genNode("Workloads", nf.workloadEntries),

			// TODO: re-enable as functionality returns
			// nf.genNode("Discovery and Load Balancing", nf.discoAndLBEntries),
			// nf.genNode("Config and Storage", nf.configAndStorageEntries),
			// nf.genNode("Custom Resources", nil),
			// nf.genNode("RBAC", nf.rbacEntries),
			// nf.genNode("Events", nil),
		},
	}

	return n, nil
}

func (nf *NavigationFactory) genNode(name string, childFn func(string) []*hcli.Navigation) *hcli.Navigation {
	leaf := hcli.NewNavigation(name, nf.pathFor(navPathLookup[name]))
	if childFn != nil {
		leaf.Children = childFn(leaf.Path)
	}

	return leaf
}

func (nf *NavigationFactory) workloadEntries(prefix string) []*hcli.Navigation {
	return []*hcli.Navigation{
		// TODO: re-enable as functionality returns
		// hcli.NewNavigation("Cron Jobs", nf.pathFor(prefix, "cron-jobs")),
		// hcli.NewNavigation("Daemon Sets", nf.pathFor(prefix, "daemon-sets")),
		hcli.NewNavigation("Deployments", path.Join(prefix, "deployments")),
		// hcli.NewNavigation("Jobs", nf.pathFor(prefix, "jobs")),
		// hcli.NewNavigation("Pods", nf.pathFor(prefix, "pods")),
		// hcli.NewNavigation("Replica Sets", nf.pathFor(prefix, "replica-sets")),
		// hcli.NewNavigation("Replication Controllers", nf.pathFor(prefix, "replication-controllers")),
		// hcli.NewNavigation("Stateful Sets", nf.pathFor(prefix, "stateful-sets")),
	}
}

func (nf *NavigationFactory) discoAndLBEntries(prefix string) []*hcli.Navigation {
	return []*hcli.Navigation{
		hcli.NewNavigation("Ingresses", nf.pathFor(prefix, "ingresses")),
		hcli.NewNavigation("Services", nf.pathFor(prefix, "services")),
	}
}

func (nf *NavigationFactory) configAndStorageEntries(prefix string) []*hcli.Navigation {
	return []*hcli.Navigation{
		hcli.NewNavigation("Config Maps", nf.pathFor(prefix, "config-maps")),
		hcli.NewNavigation("Persistent Volume Claims", nf.pathFor(prefix, "persistent-volume-claims")),
		hcli.NewNavigation("Secrets", nf.pathFor(prefix, "secrets")),
		hcli.NewNavigation("Service Accounts", nf.pathFor(prefix, "service-accounts")),
	}
}

func (nf *NavigationFactory) rbacEntries(prefix string) []*hcli.Navigation {
	return []*hcli.Navigation{
		hcli.NewNavigation("Roles", nf.pathFor(prefix, "roles")),
		hcli.NewNavigation("Role Bindings", nf.pathFor(prefix, "role-bindings")),
	}
}
