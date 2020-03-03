/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/vmware-tanzu/octant/internal/api"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/modules/overview/logviewer"
	"github.com/vmware-tanzu/octant/internal/modules/overview/yamlviewer"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/internal/printer"
	"github.com/vmware-tanzu/octant/internal/resourceviewer"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
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
		{name: "Summary", tabFunc: o.addSummaryTab},
		{name: "Metadata", tabFunc: o.addMetadataTab},
		{name: "Resource Viewer", tabFunc: o.addResourceViewerTab},
		{name: "YAML", tabFunc: o.addYAMLViewerTab},
		{name: "Logs", tabFunc: o.addLogsTab},
	}

	return o
}

type tabFunc func(ctx context.Context, object runtime.Object, cr *component.ContentResponse, options Options) error

type tabFuncDescriptor struct {
	name    string
	tabFunc tabFunc
}

// Describe describes an object.
func (d *Object) Describe(ctx context.Context, namespace string, options Options) (component.ContentResponse, error) {
	logger := log.From(ctx)

	object, err := options.LoadObject(ctx, namespace, options.Fields, d.objectStoreKey)
	if err != nil {
		return component.EmptyContentResponse, api.NewNotFoundError(d.path)
	} else if object == nil {
		cr := component.NewContentResponse(component.TitleFromString("LoadObject Error"))
		addErrorTab(ctx, "Error", fmt.Errorf("unable to load object %s", d.objectStoreKey), cr)
		return *cr, nil
	}

	item := d.objectType()

	if err := scheme.Scheme.Convert(object, item, nil); err != nil {
		cr := component.NewContentResponse(component.TitleFromString("Converting Dynamic Object Error"))
		addErrorTab(ctx, "Error", fmt.Errorf("converting dynamic object to a type: %w", err), cr)
		return *cr, nil
	}

	if err := copyObjectMeta(item, object); err != nil {
		cr := component.NewContentResponse(component.TitleFromString("Copying Object Metadata Error"))
		addErrorTab(ctx, "Error", fmt.Errorf("copying object metadata: %w", err), cr)
		return *cr, nil
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
		addErrorTab(ctx, "Error", fmt.Errorf("expected item to be a runtime object. It was a %T", item), cr)
		return *cr, nil
	}

	objAccessor, err := meta.Accessor(currentObject)
	if err != nil {
		return component.EmptyContentResponse, err
	}

	if objAccessor.GetDeletionTimestamp() == nil {
		key, err := store.KeyFromObject(currentObject)
		if err != nil {
			return component.EmptyContentResponse, err
		}

		confirmation, err := deleteObjectConfirmation(currentObject)
		if err != nil {
			return component.EmptyContentResponse, err
		}

		cr.AddButton("Delete", action.CreatePayload(octant.ActionDeleteObject,
			key.ToActionPayload()), confirmation)
	}

	hasTabError := false
	for _, tfd := range d.tabFuncDescriptors {
		if err := tfd.tabFunc(ctx, currentObject, cr, options); err != nil {
			if ctx.Err() == context.Canceled {
				continue
			}
			hasTabError = true
			addErrorTab(ctx, tfd.name, err, cr)
		}
	}

	if hasTabError {
		logger.With("tab-object", object).Errorf("generated tabs with errors")
	}

	tabs, err := options.PluginManager().Tabs(ctx, object)
	if err != nil {
		addErrorTab(ctx, "Plugin Error", fmt.Errorf("getting tabs from plugins: %w", err), cr)
		return *cr, nil
	}

	for i := range tabs {
		tabs[i].Contents.SetAccessor(tabs[i].Name)
		cr.Add(&tabs[i].Contents)
	}

	return *cr, nil
}

func deleteObjectConfirmation(object runtime.Object) (component.ButtonOption, error) {
	if object == nil {
		return nil, fmt.Errorf("object is nil")
	}
	_, kind := object.GetObjectKind().GroupVersionKind().ToAPIVersionAndKind()

	accessor, err := meta.Accessor(object)
	if err != nil {
		return nil, err
	}

	confirmationTitle := fmt.Sprintf("Delete %s", kind)
	confirmationBody := fmt.Sprintf("Are you sure you want to delete *%s* **%s**? This action is permanent and cannot be recovered.", kind, accessor.GetName())
	return component.WithButtonConfirmation(confirmationTitle, confirmationBody), nil
}

func addErrorTab(ctx context.Context, name string, err error, cr *component.ContentResponse) {
	errComponent := component.NewError(component.TitleFromString(name), err)

	accessor := name
	strings.ReplaceAll(name, " ", "")

	errComponent.SetAccessor(accessor)
	cr.Add(errComponent)
}

func (d *Object) PathFilters() []PathFilter {
	return []PathFilter{
		*NewPathFilter(d.path, d),
	}
}

func (d *Object) addSummaryTab(ctx context.Context, object runtime.Object, cr *component.ContentResponse, options Options) error {
	vc, err := options.Printer.Print(ctx, object, options.PluginManager())
	if vc == nil {
		return fmt.Errorf("unable to print a nil object: %w", err)
	}

	if err != nil {
		return err
	}

	vc.SetAccessor("summary")
	cr.Add(vc)

	return nil
}

func (d *Object) addResourceViewerTab(ctx context.Context, object runtime.Object, cr *component.ContentResponse, options Options) error {
	if !d.disableResourceViewer {
		m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(object)
		if err != nil {
			component.NewError(component.TitleFromString("Show resource viewer for object"), err)
		}

		u := &unstructured.Unstructured{Object: m}

		resourceViewerComponent, err := resourceviewer.Create(ctx, options.Dash, options.Queryer, u)
		if err != nil {
			return err
		}

		resourceViewerComponent.SetAccessor("resourceViewer")
		cr.Add(resourceViewerComponent)
	}

	return nil
}

func (d *Object) addMetadataTab(ctx context.Context, object runtime.Object, cr *component.ContentResponse, options Options) error {
	metadataComponent, err := printer.MetadataHandler(object, options.Link)
	if err != nil {
		return err
	}

	metadataComponent.SetAccessor("metadata")
	cr.Add(metadataComponent)

	return nil
}

func (d *Object) addYAMLViewerTab(ctx context.Context, object runtime.Object, cr *component.ContentResponse, _ Options) error {
	yvComponent, err := yamlviewer.ToComponent(object)
	if err != nil {
		return err
	}

	yvComponent.SetAccessor("yaml")
	cr.Add(yvComponent)
	return nil
}

func (d *Object) addLogsTab(ctx context.Context, object runtime.Object, cr *component.ContentResponse, _ Options) error {
	if isPod(object) {
		logsComponent, err := logviewer.ToComponent(object)
		if err != nil {
			return err
		}

		logsComponent.SetAccessor("logs")
		cr.Add(logsComponent)
	}

	return nil
}
