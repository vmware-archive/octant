package octant

import (
	"context"
	"path"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/heptio/developer-dash/pkg/store"
)

// EntriesFunc is a function that can create navigation entries.
type EntriesFunc func(ctx context.Context, prefix, namespace string, objectStore store.Store) ([]Navigation, error)

// NavigationEntries help construct navigation entries.
type NavigationEntries struct {
	Lookup       map[string]string
	EntriesFuncs map[string]EntriesFunc
	Order        []string
}

// NavigationFactory generates navigation entries.
type NavigationFactory struct {
	rootPath    string
	namespace   string
	entries     NavigationEntries
	objectStore store.Store
}

// NewNavigationFactory creates an instance of NewNavigationFactory.
func NewNavigationFactory(namespace string, root string, objectStore store.Store, entries NavigationEntries) *NavigationFactory {
	var rootPath = root
	if namespace != "" {
		rootPath = path.Join(root, "namespace", namespace, "")
	}
	if !strings.HasSuffix(rootPath, "/") {
		rootPath = rootPath + "/"
	}

	return &NavigationFactory{
		rootPath:    rootPath,
		namespace:   namespace,
		objectStore: objectStore,
		entries:     entries,
	}
}

// Root returns the rootPath of the navigation tree.
func (nf *NavigationFactory) Root() string {
	return nf.rootPath
}

// Generate returns navigation entries.
func (nf *NavigationFactory) Generate(ctx context.Context, title string) (*Navigation, error) {
	n := &Navigation{
		Title:    title,
		Path:     nf.rootPath,
		Children: []Navigation{},
	}

	var mu sync.Mutex
	var g errgroup.Group

	for _, name := range nf.entries.Order {
		g.Go(func() error {
			children, err := nf.genNode(ctx, name, nf.entries.EntriesFuncs[name])
			if err != nil {
				return errors.Wrapf(err, "generate entries for %s", name)
			}

			mu.Lock()
			n.Children = append(n.Children, *children)
			mu.Unlock()

			return nil
		})

		if err := g.Wait(); err != nil {
			return nil, err
		}

	}

	return n, nil
}

func (nf *NavigationFactory) pathFor(elements ...string) string {
	return path.Join(append([]string{nf.rootPath}, elements...)...)
}

func (nf *NavigationFactory) genNode(ctx context.Context, name string, childFn EntriesFunc) (*Navigation, error) {
	node := NewNavigation(name, nf.pathFor(nf.entries.Lookup[name]))
	if childFn != nil {
		children, err := childFn(ctx, node.Path, nf.namespace, nf.objectStore)
		if err != nil {
			return nil, err
		}
		node.Children = children
	}

	return node, nil
}
