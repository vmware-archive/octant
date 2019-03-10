package overview

import (
	"context"
	"path"
	"sort"
	"strings"
	"sync"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/sugarloaf"
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
		"Port Forwards":                "portforward",
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
func (nf *NavigationFactory) Entries(ctx context.Context) (*sugarloaf.Navigation, error) {
	m := map[string]entriesFunc{
		"Workloads":                    nf.workloadEntries,
		"Discovery and Load Balancing": nf.discoAndLBEntries,
		"Config and Storage":           nf.configAndStorageEntries,
		"Custom Resources":             nf.crdEntries,
		"RBAC":                         nf.rbacEntries,
		"Events":                       nil,
		"Port Forwards":                nil,
	}

	navOrder := []string{
		"Workloads",
		"Discovery and Load Balancing",
		"Config and Storage",
		"Custom Resources",
		"RBAC",
		"Events",
		"Port Forwards",
	}

	n := &sugarloaf.Navigation{
		Title:    "Overview",
		Path:     nf.root,
		Children: []*sugarloaf.Navigation{},
	}

	var mu sync.Mutex
	var g errgroup.Group

	for _, name := range navOrder {
		g.Go(func() error {
			children, err := nf.genNode(ctx, name, m[name])
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

type entriesFunc func(context.Context, string) ([]*sugarloaf.Navigation, error)

func (nf *NavigationFactory) genNode(ctx context.Context, name string, childFn entriesFunc) (*sugarloaf.Navigation, error) {
	node := sugarloaf.NewNavigation(name, nf.pathFor(navPathLookup[name]))
	if childFn != nil {
		children, err := childFn(ctx, node.Path)
		if err != nil {
			return nil, err
		}
		node.Children = children
	}

	return node, nil
}

func (nf *NavigationFactory) workloadEntries(ctx context.Context, prefix string) ([]*sugarloaf.Navigation, error) {
	return []*sugarloaf.Navigation{
		sugarloaf.NewNavigation("Cron Jobs", path.Join(prefix, "cron-jobs")),
		sugarloaf.NewNavigation("Daemon Sets", path.Join(prefix, "daemon-sets")),
		sugarloaf.NewNavigation("Deployments", path.Join(prefix, "deployments")),
		sugarloaf.NewNavigation("Jobs", path.Join(prefix, "jobs")),
		sugarloaf.NewNavigation("Pods", path.Join(prefix, "pods")),
		sugarloaf.NewNavigation("Replica Sets", path.Join(prefix, "replica-sets")),
		sugarloaf.NewNavigation("Replication Controllers", path.Join(prefix, "replication-controllers")),
		sugarloaf.NewNavigation("Stateful Sets", path.Join(prefix, "stateful-sets")),
	}, nil
}

func (nf *NavigationFactory) discoAndLBEntries(ctx context.Context, prefix string) ([]*sugarloaf.Navigation, error) {
	return []*sugarloaf.Navigation{
		sugarloaf.NewNavigation("Ingresses", path.Join(prefix, "ingresses")),
		sugarloaf.NewNavigation("Services", path.Join(prefix, "services")),
	}, nil
}

func (nf *NavigationFactory) configAndStorageEntries(ctx context.Context, prefix string) ([]*sugarloaf.Navigation, error) {
	return []*sugarloaf.Navigation{
		sugarloaf.NewNavigation("Config Maps", path.Join(prefix, "config-maps")),
		sugarloaf.NewNavigation("Persistent Volume Claims", path.Join(prefix, "persistent-volume-claims")),
		sugarloaf.NewNavigation("Secrets", path.Join(prefix, "secrets")),
		sugarloaf.NewNavigation("Service Accounts", path.Join(prefix, "service-accounts")),
	}, nil
}

func (nf *NavigationFactory) rbacEntries(ctx context.Context, prefix string) ([]*sugarloaf.Navigation, error) {
	return []*sugarloaf.Navigation{
		sugarloaf.NewNavigation("Cluster Roles", path.Join(prefix, "cluster-roles")),
		sugarloaf.NewNavigation("Cluster Role Bindings", path.Join(prefix, "cluster-role-bindings")),
		sugarloaf.NewNavigation("Roles", path.Join(prefix, "roles")),
		sugarloaf.NewNavigation("Role Bindings", path.Join(prefix, "role-bindings")),
	}, nil
}

func (nf *NavigationFactory) crdEntries(ctx context.Context, prefix string) ([]*sugarloaf.Navigation, error) {
	var list []*sugarloaf.Navigation

	crdNames, err := customResourceDefinitionNames(ctx, nf.cache)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving CRD names")
	}

	sort.Strings(crdNames)

	for _, name := range crdNames {
		crd, err := customResourceDefinition(ctx, name, nf.cache)
		if err != nil {
			return nil, errors.Wrapf(err, "load %q custom resource definition", name)
		}

		objects, err := listCustomResources(ctx, crd, nf.namespace, nf.cache)
		if err != nil {
			return nil, err
		}

		if len(objects) > 0 {
			list = append(list, sugarloaf.NewNavigation(name, path.Join(prefix, name)))
		}
	}

	return list, nil
}
