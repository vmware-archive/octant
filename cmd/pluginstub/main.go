/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"context"
	"fmt"
	"log"
	"path"
	"time"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware/octant/pkg/navigation"
	"github.com/vmware/octant/pkg/plugin"
	"github.com/vmware/octant/pkg/plugin/service"
	"github.com/vmware/octant/pkg/store"
	"github.com/vmware/octant/pkg/view/component"
	"github.com/vmware/octant/pkg/view/flexlayout"
)

var pluginName = "plugin-name"
var pluginPath = path.Join("content", pluginName)

func init() {
	// Remove the prefix from the go logger since Octant will print logs with timestamps.
	log.SetPrefix("")
}

func main() {
	// This plugin is interested in Pods
	podGVK := schema.GroupVersionKind{Version: "v1", Kind: "Pod"}

	// Tell Octant to call this plugin when printing configuration or tabs for Pods
	capabilities := &plugin.Capabilities{
		SupportsPrinterConfig: []schema.GroupVersionKind{podGVK},
		SupportsTab:           []schema.GroupVersionKind{podGVK},
		IsModule:              true,
	}

	// Set up what should happen when Octant calls this plugin.
	handlers := service.HandlerFuncs{
		Print:      handlePrint,
		PrintTab:   handleTab,
		Navigation: handleNavigation,
		Content:    handleContent,
	}

	// Use the plugin service helper to register this plugin.
	p, err := service.Register(pluginName, "a description", capabilities, handlers)
	if err != nil {
		log.Fatal(err)
	}

	// The plugin can log and the log messages will show up in Octant.
	log.Printf("octant-sample-plugin is starting")
	p.Serve()
}

// handleTab is called when Octant wants to print a tab for an object.
func handleTab(dashboardClient service.Dashboard, object runtime.Object) (*component.Tab, error) {
	if object == nil {
		return nil, errors.New("object is nil")
	}

	layout := flexlayout.New()
	section := layout.AddSection()
	err := section.Add(component.NewText("content from a plugin"), component.WidthHalf)
	if err != nil {
		return nil, err
	}

	tab := component.Tab{
		Name:     "PluginStub",
		Contents: *layout.ToComponent("Plugin"),
	}

	return &tab, nil
}

// handlePrint is called when Octant wants to print an object.
func handlePrint(dashboardClient service.Dashboard, object runtime.Object) (plugin.PrintResponse, error) {
	if object == nil {
		return plugin.PrintResponse{}, errors.Errorf("object is nil")
	}

	ctx := context.Background()
	key, err := store.KeyFromObject(object)
	if err != nil {
		return plugin.PrintResponse{}, err
	}
	u, err := dashboardClient.Get(ctx, key)
	if err != nil {
		return plugin.PrintResponse{}, err
	}

	podCard := component.NewCard(fmt.Sprintf("Extra Outout for %s", u.GetName()))
	podCard.SetBody(component.NewMarkdownText("This output was generated from _octant-sample-plugin_"))

	msg := fmt.Sprintf("update from plugin at %s", time.Now().Format(time.RFC3339))

	return plugin.PrintResponse{
		Config: []component.SummarySection{
			{Header: "from-plugin", Content: component.NewText(msg)},
		},
		Status: []component.SummarySection{
			{Header: "from-plugin", Content: component.NewText(msg)},
		},
		Items: []component.FlexLayoutItem{
			{
				Width: component.WidthHalf,
				View:  podCard,
			},
		},
	}, nil
}

func handleNavigation(dashboardClient service.Dashboard) (navigation.Navigation, error) {
	return navigation.Navigation{
		Title: "Module Plugin",
		Path:  pluginPath,
		Children: []navigation.Navigation{
			{
				Title:    "Nested Once",
				Path:     path.Join(pluginPath, "nested-once"),
				IconName: "folder",
				Children: []navigation.Navigation{
					{
						Title:    "Nested Twice",
						Path:     path.Join(pluginPath, "nested-once", "nested-twice"),
						IconName: "folder",
					},
				},
			},
		},
		IconName: "cloud",
	}, nil

}

func handleContent(dashboardClient service.Dashboard, contentPath string) (component.ContentResponse, error) {
	return component.ContentResponse{
		Components: []component.Component{
			component.NewText(fmt.Sprintf("hello from plugin: path %s", contentPath)),
		},
	}, nil
}
