/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"
	"sort"
	"sync"

	"github.com/vmware-tanzu/octant/internal/gvk"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

type CRDSection struct {
	describers map[string]Describer
	path       string
	title      string
	rootPath   ResourceLink
	mu sync.Mutex
}

var _ Describer = (*CRDSection)(nil)

func NewCRDSection(p, title string, rootPath ResourceLink) *CRDSection {
	return &CRDSection{
		describers: make(map[string]Describer),
		path:       p,
		title:      title,
		rootPath:   rootPath,
	}
}

func (csd *CRDSection) Add(name string, describer Describer) {
	csd.mu.Lock()
	defer csd.mu.Unlock()

	csd.describers[name] = describer
}

func (csd *CRDSection) Remove(name string) {
	csd.mu.Lock()
	defer csd.mu.Unlock()

	delete(csd.describers, name)
}

func (csd *CRDSection) Describe(ctx context.Context, namespace string, options Options) (component.ContentResponse, error) {
	csd.mu.Lock()
	defer csd.mu.Unlock()

	var names []string
	for name := range csd.describers {
		names = append(names, name)
	}

	sort.Strings(names)

	tableCols := component.NewTableCols("Name", "Labels", "Age")
	table := component.NewTable("Custom Resources", "", tableCols)

	for _, name := range names {
		switch d := csd.describers[name].(type) {
		case *crdList:
			key := store.KeyFromGroupVersionKind(gvk.CustomResourceDefinition)
			key.Name = d.name
			crd, err := options.ObjectStore().Get(ctx, key)
			if err != nil {
				return component.EmptyContentResponse, err
			}

			crdObject, err := octant.NewCustomResourceDefinition(crd)
			if err != nil {
				return component.EmptyContentResponse, err
			}

			versions, err := crdObject.Versions()
			if err != nil {
				return component.EmptyContentResponse, err
			}

			count := 0
			for _, version := range versions {
				crGVK, err := gvk.CustomResource(crd, version)
				if err != nil {
					return component.EmptyContentResponse, err
				}
				key2 := store.KeyFromGroupVersionKind(crGVK)
				key2.Namespace = namespace
				list, _, err := options.ObjectStore().List(ctx, key2)
				if err != nil {
					return component.EmptyContentResponse, err
				}
				count += len(list.Items)
			}

			if count > 0 {
				row := component.TableRow{}

				row["Name"] = component.NewLink("", crd.GetName(), getCrdUrl(namespace, crd))
				row["Labels"] = component.NewLabels(crd.GetLabels())
				row["Age"] = component.NewTimestamp(crd.GetCreationTimestamp().Time)

				table.Add(row)
			}

		}
	}

	title := getBreadcrumb(csd.rootPath, csd.title, "", namespace)
	list := component.NewList(title, nil)
	list.Add(table)
	return component.ContentResponse{
		Components: []component.Component{list},
	}, nil
}

func (csd *CRDSection) PathFilters() []PathFilter {
	return []PathFilter{
		*NewPathFilter(csd.path, csd),
	}
}

func (csd *CRDSection) Reset(ctx context.Context) error {
	csd.mu.Lock()
	defer csd.mu.Unlock()

	logger := log.From(ctx)

	for name := range csd.describers {
		logger.With("describer-name", name, "crd-section-path", csd.path).
			Debugf("removing crd from section")
		delete(csd.describers, name)
	}

	return nil
}
