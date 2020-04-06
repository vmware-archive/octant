/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package octant

import (
	"context"
	"path"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/store"
)

// EntriesFunc is a function that can create navigation entries.
type EntriesFunc func(ctx context.Context, prefix, namespace string, objectStore store.Store, wantsClusterScoped bool) ([]navigation.Navigation, bool, error)

// NavigationEntries help construct navigation entries.
type NavigationEntries struct {
	Lookup       map[string]string
	EntriesFuncs map[string]EntriesFunc
	IconMap      map[string]string
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
func (nf *NavigationFactory) Generate(ctx context.Context, module string, wantsClusterScoped bool) ([]navigation.Navigation, error) {
	n := []navigation.Navigation{}

	var mu sync.Mutex
	var g errgroup.Group

	for index, name := range nf.entries.Order {
		g.Go(func() error {
			child, err := nf.genNode(ctx, name, nf.entries.EntriesFuncs[name], wantsClusterScoped)
			if err != nil {
				return errors.Wrapf(err, "generate entries for %s", name)
			}

			if iconName, ok := nf.entries.IconMap[name]; ok {
				child.IconName = iconName
			}

			// Setting module creates a divider in navigation
			if (module != "") && (index == 0) {
				child.Module = module
			}

			mu.Lock()
			n = append(n, *child)
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

func (nf *NavigationFactory) genNode(ctx context.Context, name string, childFn EntriesFunc, wantsClusterScoped bool) (*navigation.Navigation, error) {
	node, err := navigation.New(name, nf.pathFor(nf.entries.Lookup[name]))
	if err != nil {
		return nil, err
	}

	if childFn != nil {
		children, loading, err := childFn(ctx, node.Path, nf.namespace, nf.objectStore, wantsClusterScoped)
		if err != nil {
			return nil, err
		}
		node.Children = children
		node.Loading = loading
	}

	return node, nil
}
