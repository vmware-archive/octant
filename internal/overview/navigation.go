package overview

import (
	"path"
	"sort"
	"strings"
	"sync"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/hcli"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
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
	cache     cache.Cache
}

// NewNavigationFactory creates an instance of NewNavigationFactory.
func NewNavigationFactory(namespace string, root string, c cache.Cache) *NavigationFactory {
	var rootPath = root
	if namespace != "" {
		rootPath = path.Join(root, "namespace", namespace, "")
	}
	if !strings.HasSuffix(rootPath, "/") {
		rootPath = rootPath + "/"
	}

	return &NavigationFactory{
		root:      rootPath,
		namespace: namespace,
		cache:     c,
	}
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
	m := map[string]entriesFunc{
		"Workloads":                    nf.workloadEntries,
		"Discovery and Load Balancing": nf.discoAndLBEntries,
		"Config and Storage":           nf.configAndStorageEntries,
		"Custom Resources":             nf.crdEntries,
		"RBAC":                         nf.rbacEntries,
		"Events":                       nil,
		"Port Forwarding":              nil,
	}

	navOrder := []string{
		"Workloads",
		"Discovery and Load Balancing",
		"Config and Storage",
		"Custom Resources",
		"RBAC",
		"Events",
		"Port Forwarding",
	}

	n := &hcli.Navigation{
		Title:    "Overview",
		Path:     nf.root,
		Children: []*hcli.Navigation{},
	}

	var mu sync.Mutex
	var g errgroup.Group

	for _, name := range navOrder {
		g.Go(func() error {
			children, err := nf.genNode(name, m[name])
			if err != nil {
				return errors.Wrapf(err, "generate entries for %s", name)
			}

			mu.Lock()
			n.Children = append(n.Children, children)
			mu.Unlock()

			return nil
		})

		if err := g.Wait(); err != nil {
			return nil, err
		}

	}

	return n, nil
}

type entriesFunc func(string) ([]*hcli.Navigation, error)

func (nf *NavigationFactory) genNode(name string, childFn entriesFunc) (*hcli.Navigation, error) {
	node := hcli.NewNavigation(name, nf.pathFor(navPathLookup[name]))
	if childFn != nil {
		children, err := childFn(node.Path)
		if err != nil {
			return nil, err
		}
		node.Children = children
	}

	return node, nil
}

func (nf *NavigationFactory) workloadEntries(prefix string) ([]*hcli.Navigation, error) {
	return []*hcli.Navigation{
		hcli.NewNavigation("Cron Jobs", path.Join(prefix, "cron-jobs")),
		hcli.NewNavigation("Daemon Sets", path.Join(prefix, "daemon-sets")),
		hcli.NewNavigation("Deployments", path.Join(prefix, "deployments")),
		hcli.NewNavigation("Jobs", path.Join(prefix, "jobs")),
		hcli.NewNavigation("Pods", path.Join(prefix, "pods")),
		hcli.NewNavigation("Replica Sets", path.Join(prefix, "replica-sets")),
		hcli.NewNavigation("Replication Controllers", path.Join(prefix, "replication-controllers")),
		hcli.NewNavigation("Stateful Sets", path.Join(prefix, "stateful-sets")),
	}, nil
}

func (nf *NavigationFactory) discoAndLBEntries(prefix string) ([]*hcli.Navigation, error) {
	return []*hcli.Navigation{
		hcli.NewNavigation("Ingresses", path.Join(prefix, "ingresses")),
		hcli.NewNavigation("Services", path.Join(prefix, "services")),
	}, nil
}

func (nf *NavigationFactory) configAndStorageEntries(prefix string) ([]*hcli.Navigation, error) {
	return []*hcli.Navigation{
		hcli.NewNavigation("Config Maps", path.Join(prefix, "config-maps")),
		hcli.NewNavigation("Persistent Volume Claims", path.Join(prefix, "persistent-volume-claims")),
		hcli.NewNavigation("Secrets", path.Join(prefix, "secrets")),
		hcli.NewNavigation("Service Accounts", path.Join(prefix, "service-accounts")),
	}, nil
}

func (nf *NavigationFactory) rbacEntries(prefix string) ([]*hcli.Navigation, error) {
	return []*hcli.Navigation{
		hcli.NewNavigation("Cluster Roles", path.Join(prefix, "cluster-roles")),
		hcli.NewNavigation("Cluster Role Bindings", path.Join(prefix, "cluster-role-bindings")),
		hcli.NewNavigation("Roles", path.Join(prefix, "roles")),
		hcli.NewNavigation("Role Bindings", path.Join(prefix, "role-bindings")),
	}, nil
}

func (nf *NavigationFactory) crdEntries(prefix string) ([]*hcli.Navigation, error) {
	var list []*hcli.Navigation

	crdNames, err := customResourceDefinitionNames(nf.cache)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving CRD names")
	}

	sort.Strings(crdNames)

	for _, name := range crdNames {
		crd, err := customResourceDefinition(name, nf.cache)
		if err != nil {
			return nil, errors.Wrapf(err, "load %q custom resource definition", name)
		}

		objects, err := listCustomResources(crd, nf.namespace, nf.cache)
		if err != nil {
			return nil, err
		}

		if len(objects) > 0 {
			list = append(list, hcli.NewNavigation(name, path.Join(prefix, name)))
		}
	}

	return list, nil
}
