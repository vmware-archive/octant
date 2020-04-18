/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"

	"github.com/pkg/errors"

	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// Section is a wrapper to combine content from multiple describers.
type Section struct {
	path       string
	title      string
	describers []Describer
}

var _ Describer = (*Section)(nil)

// NewSection creates a Section.
func NewSection(p, title string, describers ...Describer) *Section {
	return &Section{
		path:       p,
		title:      title,
		describers: describers,
	}
}

// Describe generates content.
func (d *Section) Describe(ctx context.Context, namespace string, options Options) (component.ContentResponse, error) {
	list, err := d.Component(ctx, namespace, options)
	if err != nil {
		return component.EmptyContentResponse, err
	}

	cr := component.ContentResponse{
		Components: []component.Component{list},
		Title:      component.Title(component.NewText(d.title)),
	}

	return cr, nil
}

func (d *Section) Component(ctx context.Context, namespace string, options Options) (*component.List, error) {
	title := append([]component.TitleComponent{}, component.NewText(d.title))
	list := component.NewList(title, nil)

	for describerIndex := range d.describers {
		cResponse, err := d.describers[describerIndex].Describe(ctx, namespace, options)
		if err != nil {
			return nil, err
		}

		for componentIndex := range cResponse.Components {
			if nestedList, ok := cResponse.Components[componentIndex].(*component.List); ok {
				for itemIndex := range nestedList.Config.Items {
					item := nestedList.Config.Items[itemIndex]
					if !item.IsEmpty() {
						list.Add(item)
					}
				}
			}
		}
	}

	return list, nil
}

// PathFilters returns path filters for the section.
func (d *Section) PathFilters() []PathFilter {
	PathFilters := []PathFilter{
		*NewPathFilter(d.path, d),
	}

	for _, child := range d.describers {
		PathFilters = append(PathFilters, child.PathFilters()...)
	}

	return PathFilters
}

func (d *Section) Reset(ctx context.Context) error {
	for i := range d.describers {
		if err := d.describers[i].Reset(ctx); err != nil {
			return errors.Wrapf(err, "reset describer in section %s", d.path)
		}
	}

	return nil
}
