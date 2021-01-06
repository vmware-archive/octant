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
	describerMap map[string]Describer
	path         string
	title        string
	mu           sync.RWMutex
}

var _ Describer = (*CRDSection)(nil)

func NewCRDSection(p, title string) *CRDSection {
	return &CRDSection{
		describerMap: make(map[string]Describer),
		path:         p,
		title:        title,
	}
}

func (csd *CRDSection) Add(name string, describer Describer) {
	csd.mu.Lock()
	defer csd.mu.Unlock()

	csd.describerMap[name] = describer
}

func (csd *CRDSection) Remove(name string) {
	csd.mu.Lock()
	defer csd.mu.Unlock()

	delete(csd.describerMap, name)
}

func (csd *CRDSection) Describe(ctx context.Context, namespace string, options Options) (component.ContentResponse, error) {
	title := component.Title(component.NewText(csd.title))
	list := component.NewList(title, nil)

	for _, d := range csd.describers() {
		cr, err := d.Describe(ctx, namespace, options)
		if err != nil {
			return component.EmptyContentResponse, err
		}

		for i := range cr.Components {
			switch c := cr.Components[i].(type) {
			case *component.List:
				for _, item := range c.Config.Items {
					if !item.IsEmpty() {
						list.Add(item)
					}
				}
			default:
				if !c.IsEmpty() {
					list.Add(c)
				}
			}
		}
	}

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

	for name := range csd.describerMap {
		logger.With("describer-name", name, "crd-section-path", csd.path).
			Debugf("removing crd from section")
		delete(csd.describerMap, name)
	}

	return nil
}

func (csd *CRDSection) describers() []Describer {
	csd.mu.RLock()
	defer csd.mu.RUnlock()

	var names []string
	for name := range csd.describerMap {
		names = append(names, name)
	}

	sort.Strings(names)

	var out []Describer

	for _, name := range names {
		out = append(out, csd.describerMap[name])
	}

	return out
}

func (csd *CRDSection) crdTable(ctx context.Context, namespace string, options Options) (*component.Table, error) {
	tableCols := component.NewTableCols("Name", "Labels", "Age")
	table := component.NewTable("Custom Resources", "", tableCols)

	describers := csd.describers()

	for _, describer := range describers {
		switch d := describer.(type) {
		case *crdList:
			key := store.KeyFromGroupVersionKind(gvk.CustomResourceDefinition)
			key.Name = d.name
			crd, err := options.ObjectStore().Get(ctx, key)
			if err != nil {
				return nil, err
			}

			crdTool, err := octant.NewCustomResourceDefinitionTool(crd)
			if err != nil {
				return nil, err
			}

			versions, err := crdTool.Versions()
			if err != nil {
				return nil, err
			}

			count := 0
			for _, version := range versions {
				crGVK, err := gvk.CustomResource(crd, version)
				if err != nil {
					return nil, err
				}
				key2 := store.KeyFromGroupVersionKind(crGVK)
				key2.Namespace = namespace
				list, _, err := options.ObjectStore().List(ctx, key2)
				if err != nil {
					return nil, err
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

	return table, nil
}
