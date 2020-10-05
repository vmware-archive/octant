/*
 *  Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 *  SPDX-License-Identifier: Apache-2.0
 *
 */

package describer

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

//go:generate mockgen  -destination=./fake/mock_tabs_generator.go -package=fake github.com/vmware-tanzu/octant/internal/describer TabsGenerator

// TabsFactory creates a list of tab descriptors.
type TabsFactory func() ([]Tab, error)

// TabsGeneratorConfig is configuration for TabsGenerator.
type TabsGeneratorConfig struct {
	Object      runtime.Object
	TabsFactory TabsFactory
	Options     Options
}

// TabsGenerator generates tabs for an object.
type TabsGenerator interface {
	// Generate generates tabs given a configuration and returns a list of components.
	Generate(ctx context.Context, config TabsGeneratorConfig) ([]component.Component, error)
}

// ObjectTabsGenerator generates tabs for an object.
type ObjectTabsGenerator struct {
}

var _ TabsGenerator = &ObjectTabsGenerator{}

// NewObjectTabsGenerator creates an instance of ObjectTabsGenerator.
func NewObjectTabsGenerator() *ObjectTabsGenerator {
	tg := ObjectTabsGenerator{}

	return &tg
}

// Generate generates tabs for an object.
func (t ObjectTabsGenerator) Generate(ctx context.Context, config TabsGeneratorConfig) ([]component.Component, error) {
	if config.Object == nil {
		return nil, fmt.Errorf("can't generate tabs for nil object")
	}

	if config.TabsFactory == nil {
		return nil, fmt.Errorf("tabs factory is nil")
	}

	logger := log.From(ctx)

	var indexedComponents []indexedComponent
	var mu sync.Mutex

	var g errgroup.Group

	descriptors, err := config.TabsFactory()
	if err != nil {
		return []component.Component{
			CreateErrorTab("Error", fmt.Errorf("generate tabs: %w", err)),
		}, nil
	}

	for i := range descriptors {
		i := i
		descriptor := descriptors[i]
		g.Go(func() error {
			c, err := descriptor.Factory(ctx, config.Object, config.Options)
			if err != nil {
				c = CreateErrorTab(descriptor.Name, err)
			}

			mu.Lock()
			indexedComponents = append(indexedComponents, indexedComponent{
				c:     c,
				index: i,
			})
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		logger.WithErr(err).Errorf("create tabs")
	}

	sort.Slice(indexedComponents, func(i, j int) bool {
		return indexedComponents[i].index < indexedComponents[j].index
	})

	var list []component.Component
	for _, ci := range indexedComponents {
		if ci.c != nil {
			list = append(list, ci.c)
		}
	}

	return list, nil
}

// CreateErrorTab creates an error tab given a name and an error.
func CreateErrorTab(name string, err error) component.Component {
	errComponent := component.NewError(component.TitleFromString(name), err)

	accessor := name
	strings.ReplaceAll(name, " ", "")

	errComponent.SetAccessor(accessor)

	return errComponent
}

// TabFactory is a function that generates a component which describes an object as a component.
type TabFactory func(
	ctx context.Context,
	object runtime.Object,
	options Options) (component.Component, error)

// Tab describes a tab. It contains the name and a factory function to generate the content for the tab.
type Tab struct {
	// Name is the name of the tab.
	Name string
	// Factory is a function that generates the contents for a tab (as a component).
	Factory TabFactory
}

type indexedComponent struct {
	c     component.Component
	index int
}

// objectTabsFactory generates tabs for an object. It includes plugin tabs.
func objectTabsFactory(
	ctx context.Context,
	object runtime.Object,
	descriptors []Tab, options Options) TabsFactory {
	return func() ([]Tab, error) {
		list := append(descriptors)
		pluginList, err := pluginTabsFactory(ctx, object, options)
		if err != nil {
			return nil, fmt.Errorf("generate plugin tabs: %w", err)
		}

		return append(list, pluginList...), nil
	}
}

// pluginTabsFactory generates plugin tabs for an object.
func pluginTabsFactory(
	ctx context.Context,
	object runtime.Object,
	options Options) ([]Tab, error) {
	var list []Tab

	tabs, err := options.PluginManager().Tabs(ctx, object)
	if err != nil {
		list = append(list, Tab{
			Name: "Plugin",
			Factory: func(ctx context.Context, object runtime.Object, options Options) (component.Component, error) {
				tab := CreateErrorTab("Plugin Error", fmt.Errorf("getting tabs from plugins: %w", err))
				return tab, nil
			},
		})

		return list, nil
	}

	for _, tab := range tabs {
		tab := tab
		list = append(list, Tab{
			Name: tab.Name,
			Factory: func(ctx context.Context, object runtime.Object, options Options) (component.Component, error) {
				return &tab.Contents, nil
			},
		})
	}

	return list, nil
}
