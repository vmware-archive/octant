/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/vmware-tanzu/octant/internal/api"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// defaultObjectTabs are the default tabs for an object (that is not a custom resource).
func defaultObjectTabs() []Tab {
	return []Tab{
		{Name: "Summary", Factory: SummaryTab},
		{Name: "Metadata", Factory: MetadataTab},
		{Name: "Resource Viewer", Factory: ResourceViewerTab},
		{Name: "YAML", Factory: YAMLViewerTab},
		{Name: "Logs", Factory: LogsTab},
		{Name: "Terminal", Factory: TerminalTab},
	}
}

// ObjectConfig is configuration for Object.
type ObjectConfig struct {
	Path           string
	BaseTitle      string
	ObjectType     func() interface{}
	StoreKey       store.Key
	IconName       string
	IconSource     string
	TabsGenerator  TabsGenerator
	TabDescriptors []Tab
}

// Object describes an object.
type Object struct {
	*base

	path                  string
	baseTitle             string
	objectType            func() interface{}
	objectStoreKey        store.Key
	disableResourceViewer bool
	tabFuncDescriptors    []Tab
	iconName              string
	iconSource            string
	tabsGenerator         TabsGenerator
}

// NewObject creates an instance of Object.
func NewObject(c ObjectConfig) *Object {
	tg := c.TabsGenerator
	if tg == nil {
		tg = NewObjectTabsGenerator()
	}

	td := c.TabDescriptors
	if td == nil {
		td = defaultObjectTabs()
	}

	o := &Object{
		path:               c.Path,
		baseTitle:          c.BaseTitle,
		base:               newBaseDescriber(),
		objectStoreKey:     c.StoreKey,
		objectType:         c.ObjectType,
		iconName:           c.IconName,
		iconSource:         c.IconSource,
		tabsGenerator:      tg,
		tabFuncDescriptors: td,
	}

	return o
}

// Describe describes an object. An object description is comprised of multiple tabs of content.
// By default, there will be the following tabs: summary, metadata, resource viewer, and yaml.
// If the object is a pod, there will also be a log and terminal tab. If plugins can contribute
// tabs to this object, those tabs will be included as well.
//
// This function should always return a content response even if there is an error.
func (d *Object) Describe(ctx context.Context, namespace string, options Options) (component.ContentResponse, error) {
	object, err := options.LoadObject(ctx, namespace, options.Fields, d.objectStoreKey)
	if err != nil {
		return component.EmptyContentResponse, api.NewNotFoundError(d.path)
	} else if object == nil {
		cr := component.NewContentResponse(component.TitleFromString("LoadObject Error"))
		c := CreateErrorTab("Error", fmt.Errorf("unable to load object %s", d.objectStoreKey))
		cr.Add(c)
		return *cr, nil
	}

	item := d.objectType()

	if err := scheme.Scheme.Convert(object, item, nil); err != nil {
		cr := component.NewContentResponse(component.TitleFromString("Converting Dynamic Object Error"))
		c := CreateErrorTab("Error", fmt.Errorf("converting dynamic object to a type: %w", err))
		cr.Add(c)
		return *cr, nil
	}

	if err := copyObjectMeta(item, object); err != nil {
		cr := component.NewContentResponse(component.TitleFromString("Copying Object Metadata Error"))
		c := CreateErrorTab("Error", fmt.Errorf("copying object metadata: %w", err))
		cr.Add(c)
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
		c := CreateErrorTab("Error", fmt.Errorf("expected item to be a runtime object. It was a %T", item))
		cr.Add(c)
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

		confirmation, err := DeleteObjectConfirmation(currentObject)
		if err != nil {
			return component.EmptyContentResponse, err
		}

		cr.AddButton("Delete", action.CreatePayload(octant.ActionDeleteObject,
			key.ToActionPayload()), confirmation)
	}

	config := TabsGeneratorConfig{
		Object:      currentObject,
		TabsFactory: objectTabsFactory(ctx, currentObject, d.tabFuncDescriptors, options),
		Options:     options,
	}
	tabComponents, err := d.tabsGenerator.Generate(ctx, config)
	if err != nil {
		return component.EmptyContentResponse, fmt.Errorf("generate tabs: %w", err)
	}

	cr.Add(tabComponents...)

	return *cr, nil
}

// DeleteObjectConfirmation create a button option confirmation for deleting an object.
func DeleteObjectConfirmation(object runtime.Object) (component.ButtonOption, error) {
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

// PathFilters returns the path filters for this object.
func (d *Object) PathFilters() []PathFilter {
	return []PathFilter{
		*NewPathFilter(d.path, d),
	}
}
