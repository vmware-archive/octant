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
		"Port Forwarding":              "portforward",
	}
)

// NavigationFactory generates navigation entries.
type NavigationFactory struct {
	root      string
	namespace string
}

// NewNavigationFactory creates an instance of NewNavigationFactory.
func NewNavigationFactory(namespace string, root string) *NavigationFactory {
	var rootPath = root
	if namespace != "" {
		rootPath = path.Join(root, "namespace", namespace, "")
	}
	if !strings.HasSuffix(rootPath, "/") {
		rootPath = rootPath + "/"
	}

	return &NavigationFactory{root: rootPath, namespace: namespace}
}

// Root returns the root of the navigation tree.
func (nf *NavigationFactory) Root() string {
	return nf.root
}

func (nf *NavigationFactory) pathFor(elements ...string) string {
	// return path.Join(append([]string{nf.root, "namespace", nf.namespace}, elements...)...)
	return path.Join(append([]string{nf.root}, elements...)...)
}

// Entries returns navigation entries.
func (nf *NavigationFactory) Entries() (*hcli.Navigation, error) {

	n := &hcli.Navigation{
		Title: "Overview",
		Path:  nf.root,
		Children: []*hcli.Navigation{

			nf.genNode("Workloads", nf.workloadEntries),
			nf.genNode("Discovery and Load Balancing", nf.discoAndLBEntries),
			nf.genNode("Config and Storage", nf.configAndStorageEntries),
			nf.genNode("Custom Resources", nil),
			nf.genNode("RBAC", nf.rbacEntries),
			nf.genNode("Events", nil),
			nf.genNode("Port Forwarding", nil),
		},
	}

	return n, nil
}

func (nf *NavigationFactory) genNode(name string, childFn func(string) []*hcli.Navigation) *hcli.Navigation {
	node := hcli.NewNavigation(name, nf.pathFor(navPathLookup[name]))
	if childFn != nil {
		node.Children = childFn(node.Path)
	}

	return node
}

func (nf *NavigationFactory) workloadEntries(prefix string) []*hcli.Navigation {
	return []*hcli.Navigation{
		hcli.NewNavigation("Cron Jobs", path.Join(prefix, "cron-jobs")),
		hcli.NewNavigation("Daemon Sets", path.Join(prefix, "daemon-sets")),
		hcli.NewNavigation("Deployments", path.Join(prefix, "deployments")),
		hcli.NewNavigation("Jobs", path.Join(prefix, "jobs")),
		hcli.NewNavigation("Pods", path.Join(prefix, "pods")),
		hcli.NewNavigation("Replica Sets", path.Join(prefix, "replica-sets")),
		hcli.NewNavigation("Replication Controllers", path.Join(prefix, "replication-controllers")),
		hcli.NewNavigation("Stateful Sets", path.Join(prefix, "stateful-sets")),
	}
}

func (nf *NavigationFactory) discoAndLBEntries(prefix string) []*hcli.Navigation {
	return []*hcli.Navigation{
		hcli.NewNavigation("Ingresses", path.Join(prefix, "ingresses")),
		hcli.NewNavigation("Services", path.Join(prefix, "services")),
	}
}

func (nf *NavigationFactory) configAndStorageEntries(prefix string) []*hcli.Navigation {
	return []*hcli.Navigation{
		hcli.NewNavigation("Config Maps", path.Join(prefix, "config-maps")),
		hcli.NewNavigation("Persistent Volume Claims", path.Join(prefix, "persistent-volume-claims")),
		hcli.NewNavigation("Secrets", path.Join(prefix, "secrets")),
		hcli.NewNavigation("Service Accounts", path.Join(prefix, "service-accounts")),
	}
}

func (nf *NavigationFactory) rbacEntries(prefix string) []*hcli.Navigation {
	return []*hcli.Navigation{
		hcli.NewNavigation("Roles", path.Join(prefix, "roles")),
		hcli.NewNavigation("Role Bindings", path.Join(prefix, "role-bindings")),
	}
}
