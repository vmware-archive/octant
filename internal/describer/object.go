/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/internal/modules/overview/logviewer"
	"github.com/vmware/octant/internal/modules/overview/resourceviewer"
	"github.com/vmware/octant/internal/modules/overview/yamlviewer"
	"github.com/vmware/octant/pkg/store"
	"github.com/vmware/octant/pkg/view/component"
)

type ObjectConfig struct {
	Path                  string
	BaseTitle             string
	ObjectType            func() interface{}
	StoreKey              store.Key
	DisableResourceViewer bool
	IconName              string
	IconSource            string
}

// Object describes an object.
type Object struct {
	*base

	path                  string
	baseTitle             string
	objectType            func() interface{}
	objectStoreKey        store.Key
	disableResourceViewer bool
	tabFuncDescriptors    []tabFuncDescriptor
	iconName              string
	iconSource            string
}

// NewObject creates an instance of Object.
func NewObject(c ObjectConfig) *Object {
	o := &Object{
		path:                  c.Path,
		baseTitle:             c.BaseTitle,
		base:                  newBaseDescriber(),
		objectStoreKey:        c.StoreKey,
		objectType:            c.ObjectType,
		disableResourceViewer: c.DisableResourceViewer,
		iconName:              c.IconName,
		iconSource:            c.IconSource,
	}

	o.tabFuncDescriptors = []tabFuncDescriptor{
		{name: "summary", tabFunc: o.addSummaryTab},
		{name: "resource viewer", tabFunc: o.addResourceViewerTab},
		{name: "yaml", tabFunc: o.addYAMLViewerTab},
		{name: "logs", tabFunc: o.addLogsTab},
	}

	return o
}

type tabFunc func(ctx context.Context, object runtime.Object, cr *component.ContentResponse, options Options) error

type tabFuncDescriptor struct {
	name    string
	tabFunc tabFunc
}

// Describe describes an object.
func (d *Object) Describe(ctx context.Context, prefix, namespace string, options Options) (component.ContentResponse, error) {
	logger := log.From(ctx)

	object, err := options.LoadObject(ctx, namespace, options.Fields, d.objectStoreKey)
	if err != nil {
		return EmptyContentResponse, errors.Wrapf(err, "loading object with %s", d.objectStoreKey.String())
	} else if object == nil {
		return EmptyContentResponse, errors.Errorf("unable to load object %s", d.objectStoreKey)
	}

	item := d.objectType()

	if err := scheme.Scheme.Convert(object, item, nil); err != nil {
		return EmptyContentResponse, errors.Wrapf(err, "converting dynamic object to a type")
	}

	if err := copyObjectMeta(item, object); err != nil {
		return EmptyContentResponse, errors.Wrap(err, "copying object metadata")
	}

	accessor := meta.NewAccessor()
	objectName, _ := accessor.Name(object)

	title := append([]component.TitleComponent{}, component.NewText(d.baseTitle))
	if objectName != "" {
		title = append(title, component.NewText(objectName))
	}

	cr := component.NewContentResponse(title)
	cr.IconSource = d.iconSource
	cr.IconName = d.iconName

	currentObject, ok := item.(runtime.Object)
	if !ok {
		return EmptyContentResponse, errors.Errorf("expected item to be a runtime object. It was a %T",
			item)
	}

	hasTabError := false
	for _, tfd := range d.tabFuncDescriptors {
		if err := tfd.tabFunc(ctx, currentObject, cr, options); err != nil {
			hasTabError = true
			logger.With(
				"err", err,
				"tab-name", tfd.name,
			).Errorf("generating object Describer tab")
		}
	}

	if hasTabError {
		logger.With("tab-object", object).Errorf("unable to generate all tabs for object")
	}

	tabs, err := options.PluginManager().Tabs(object)
	if err != nil {
		return EmptyContentResponse, errors.Wrap(err, "getting tabs from plugins")
	}

	for _, tab := range tabs {
		tab.Contents.SetAccessor(tab.Name)
		cr.Add(&tab.Contents)
	}

	return *cr, nil
}

func (d *Object) PathFilters() []PathFilter {
	return []PathFilter{
		*NewPathFilter(d.path, d),
	}
}

func (d *Object) addSummaryTab(ctx context.Context, object runtime.Object, cr *component.ContentResponse, options Options) error {
	vc, err := options.Printer.Print(ctx, object, options.PluginManager())
	if vc == nil {
		return errors.Wrap(err, "unable to print a nil object")
	}

	if err != nil {
		errComponent := component.NewError(component.TitleFromString("Summary"), err)
		cr.Add(errComponent)

		logger := log.From(ctx)
		logger.Errorf("printing object: %s", err)

		return nil
	}

	vc.SetAccessor("summary")
	cr.Add(vc)

	return nil
}

func (d *Object) addResourceViewerTab(ctx context.Context, object runtime.Object, cr *component.ContentResponse, options Options) error {
	if !d.disableResourceViewer {

		cacheFn := resourceviewer.CachedResourceViewer(ctx, object, options.Dash, options.Queryer)
		resourceViewerComponent, err := options.Dash.ComponentCache().Update(ctx, cacheFn)
		if err != nil {
			return errors.Wrap(err, "retrieve resource viewer from component cache")
		}

		resourceViewerComponent.SetAccessor("resourceViewer")
		cr.Add(resourceViewerComponent)
	}

	return nil
}

func (d *Object) addYAMLViewerTab(ctx context.Context, object runtime.Object, cr *component.ContentResponse, options Options) error {
	yvComponent, err := yamlviewer.ToComponent(object)

	if err != nil {
		errComponent := component.NewError(component.TitleFromString("YAML"), err)
		cr.Add(errComponent)

		logger := log.From(ctx)
		logger.Errorf("converting object to YAML: %s", err)

		return nil
	}

	yvComponent.SetAccessor("yaml")
	cr.Add(yvComponent)
	return nil

}

func (d *Object) addLogsTab(ctx context.Context, object runtime.Object, cr *component.ContentResponse, options Options) error {
	if isPod(object) {
		logsComponent, err := logviewer.ToComponent(object)
		if err != nil {
			errComponent := component.NewError(component.TitleFromString("Logs"), err)
			cr.Add(errComponent)

			logger := log.From(ctx)
			logger.Errorf("retrieving logs for pod: %s", err)

			return nil
		}

		logsComponent.SetAccessor("logs")
		cr.Add(logsComponent)
	}

	return nil
}
