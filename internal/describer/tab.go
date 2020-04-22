/*
 *  Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 *  SPDX-License-Identifier: Apache-2.0
 *
 */

package describer

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/internal/link"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/modules/overview/logviewer"
	"github.com/vmware-tanzu/octant/internal/modules/overview/terminalviewer"
	"github.com/vmware-tanzu/octant/internal/modules/overview/yamlviewer"
	"github.com/vmware-tanzu/octant/internal/printer"
	"github.com/vmware-tanzu/octant/internal/resourceviewer"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// CustomResourceSummaryTab generates a summary tab for a custom resource. This function
// returns a TabFactory since the crd name might not be available when this factory
// is invoked.
func CustomResourceSummaryTab(crdName string) TabFactory {
	return func(ctx context.Context, object runtime.Object, options Options) (component.Component, error) {
		crd, err := CustomResourceDefinition(ctx, crdName, options.ObjectStore())
		if err != nil {
			return nil, fmt.Errorf("unable to find custom resource definition: %w", err)
		}

		cr, ok := object.(*unstructured.Unstructured)
		if !ok {
			return nil, fmt.Errorf("invalid custom resource")
		}

		linkGenerator, err := link.NewFromDashConfig(options)
		if err != nil {
			return nil, fmt.Errorf("create link generator: %w", err)
		}

		printOptions := printer.Options{
			DashConfig: options,
			Link:       linkGenerator,
		}

		return printer.CustomResourceHandler(ctx, crd, cr, printOptions)
	}
}

// SummaryTab generates a summary tab for an object.
func SummaryTab(ctx context.Context, object runtime.Object, options Options) (component.Component, error) {
	vc, err := options.Printer.Print(ctx, object, options.PluginManager())
	if err != nil {
		return nil, fmt.Errorf("print summary tab: %w", err)
	} else if vc == nil {
		return nil, fmt.Errorf("printer generated a nil object")
	}

	vc.SetAccessor("summary")

	return vc, nil
}

// MetadataTab generates a metadata tab for an object.
func MetadataTab(_ context.Context, object runtime.Object, options Options) (component.Component, error) {
	metadataComponent, err := printer.MetadataHandler(object, options.Link)
	if err != nil {
		return nil, fmt.Errorf("print metadata: %w", err)
	}

	metadataComponent.SetAccessor("metadata")

	return metadataComponent, nil
}

// ResourceViewerTab generates a resource viewer tab for an object.
func ResourceViewerTab(ctx context.Context, object runtime.Object, options Options) (component.Component, error) {
	m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(object)
	if err != nil {
		component.NewError(component.TitleFromString("Show resource viewer for object"), err)
	}

	u := &unstructured.Unstructured{Object: m}

	resourceViewerComponent, err := resourceviewer.Create(ctx, options.Dash, options.Queryer, u)
	if err != nil {
		return nil, fmt.Errorf("create resource viewer: %w", err)
	}

	resourceViewerComponent.SetAccessor("resourceViewer")
	return resourceViewerComponent, nil
}

// YAMLViewerTab generates a yaml viewer for an object.
func YAMLViewerTab(_ context.Context, object runtime.Object, _ Options) (component.Component, error) {
	yvComponent, err := yamlviewer.ToComponent(object)
	if err != nil {
		return nil, fmt.Errorf("create yaml viewer: %w", err)
	}

	yvComponent.SetAccessor("yaml")
	return yvComponent, nil
}

// LogsTab generates a logs tab for a pod. If the object is not a pod, the
// returned component will be nil with a nil error.
func LogsTab(_ context.Context, object runtime.Object, _ Options) (component.Component, error) {
	if isPod(object) {
		logsComponent, err := logviewer.ToComponent(object)
		if err != nil {
			return nil, fmt.Errorf("create log viewer: %w", err)
		}

		logsComponent.SetAccessor("logs")
		return logsComponent, nil
	}

	return nil, nil
}

// TerminalTab generates a terminal tab for a pod. If the object is not a pod,
// the returned component will be nil with a nil error.
func TerminalTab(ctx context.Context, object runtime.Object, options Options) (component.Component, error) {
	if isPod(object) {
		logger := log.From(ctx)

		terminalComponent, err := terminalviewer.ToComponent(ctx, object, options.TerminalManager(), logger)
		if err != nil {
			return nil, fmt.Errorf("create terminal viewer: %w", err)
		}

		terminalComponent.SetAccessor("terminal")
		return terminalComponent, nil
	}

	return nil, nil
}
