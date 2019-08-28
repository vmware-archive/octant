/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware/octant/pkg/navigation"
	"github.com/vmware/octant/pkg/plugin"
	"github.com/vmware/octant/pkg/plugin/service"
	"github.com/vmware/octant/pkg/store"
	"github.com/vmware/octant/pkg/view/component"
	"github.com/vmware/octant/pkg/view/flexlayout"
)

var pluginName = "plugin-name"

// This is a sample plugin showing the features of Octant's plugin API.
func main() {
	// Remove the prefix from the go logger since Octant will print logs with timestamps.
	log.SetPrefix("")

	// This plugin is interested in Pods
	podGVK := schema.GroupVersionKind{Version: "v1", Kind: "Pod"}

	// Tell Octant to call this plugin when printing configuration or tabs for Pods
	capabilities := &plugin.Capabilities{
		SupportsPrinterConfig: []schema.GroupVersionKind{podGVK},
		SupportsTab:           []schema.GroupVersionKind{podGVK},
		IsModule:              true,
	}

	// Set up what should happen when Octant calls this plugin.
	options := []service.PluginOption{
		service.WithPrinter(handlePrint),
		service.WithTabPrinter(handleTab),
		service.WithNavigation(handleNavigation, initRoutes),
	}

	// Use the plugin service helper to register this plugin.
	p, err := service.Register(pluginName, "a description", capabilities, options...)
	if err != nil {
		log.Fatal(err)
	}

	// The plugin can log and the log messages will show up in Octant.
	log.Printf("octant-sample-plugin is starting")
	p.Serve()
}

// handleTab is called when Octant wants to print a tab for an object.
func handleTab(request *service.PrintRequest) (plugin.TabResponse, error) {
	if request.Object == nil {
		return plugin.TabResponse{}, errors.New("object is nil")
	}

	// Octant uses flex layouts to display information. It's a flexible
	// grid. A flex layout is composed of multiple section. Each section
	// can contain multiple components. Components are displayed given
	// a width. In the case below, the width is half of the visible space.
	// Create sections to separate your components as each section will
	// start a new row.
	layout := flexlayout.New()
	section := layout.AddSection()

	// Octant contain's a library of components that can be used to display content.
	// This example uses markdown text.
	contents := component.NewMarkdownText("content from a *plugin*")

	err := section.Add(contents, component.WidthHalf)
	if err != nil {
		return plugin.TabResponse{}, err
	}

	// In this example, this plugin will tell Octant to create a new
	// tab when showing pods. This tab's name will be "Extra Pod Details".
	tab := component.NewTabWithContents(*layout.ToComponent("Extra Pod Details"))

	return plugin.TabResponse{Tab: tab}, nil
}

// handlePrint is called when Octant wants to print an object.
func handlePrint(request *service.PrintRequest) (plugin.PrintResponse, error) {
	if request.Object == nil {
		return plugin.PrintResponse{}, errors.Errorf("object is nil")
	}

	// load an object from the cluster and use that object to create a response.

	// Octant has a helper function to generate a key from an object. The key
	// is used to find the object in the cluster.
	key, err := store.KeyFromObject(request.Object)
	if err != nil {
		return plugin.PrintResponse{}, err
	}
	u, found, err := request.DashboardClient.Get(request.Context(), key)
	if err != nil {
		return plugin.PrintResponse{}, err
	}

	// The plugin can check if the object it requested exists.
	if !found {
		return plugin.PrintResponse{}, errors.New("object doesn't exist")
	}

	// Octant has a component library that can be used to build content for a plugin.
	// In this case, the plugin is creating a card.
	podCard := component.NewCard(fmt.Sprintf("Extra Output for %s", u.GetName()))
	podCard.SetBody(component.NewMarkdownText("This output was generated from _octant-sample-plugin_"))

	msg := fmt.Sprintf("update from plugin at %s", time.Now().Format(time.RFC3339))

	// When printing an object, you can create multiple types of content. In this
	// example, the plugin is:
	//
	// * adding a field to the configuration section for this object.
	// * adding a field to the status section for this object.
	// * create a new piece of content that will be embedded in the
	//   summary section for the component.
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

// handlePrint creates a navigation tree for this plugin. Navigation is dynamic and will
// be called frequently from Octant. Navigation is a tree of `Navigation` structs.
// The plugin can use whatever paths it likes since these paths can be namespaced to the
// the plugin.
func handleNavigation(request *service.NavigationRequest) (navigation.Navigation, error) {
	return navigation.Navigation{
		Title: "Sample Plugin",
		Path:  request.GeneratePath(),
		Children: []navigation.Navigation{
			{
				Title:    "Nested Once",
				Path:     request.GeneratePath("nested-once"),
				IconName: "folder",
				Children: []navigation.Navigation{
					{
						Title:    "Nested Twice",
						Path:     request.GeneratePath("nested-once", "nested-twice"),
						IconName: "folder",
					},
				},
			},
		},
		IconName: "cloud",
	}, nil
}

// initRoutes routes for this plugin. In this example, there is a global catch all route
// that will return the content for every single path.
func initRoutes(router *service.Router) {
	gen := func(name, accessor, requestPath string) component.Component {
		cardBody := component.NewText(fmt.Sprintf("hello from plugin: path %s", requestPath))
		card := component.NewCard(fmt.Sprintf("My Card - %s", name))
		card.SetBody(cardBody)
		cardList := component.NewCardList(name)
		cardList.AddCard(*card)
		cardList.SetAccessor(accessor)

		return cardList
	}

	router.HandleFunc("/*", func(request *service.Request) (component.ContentResponse, error) {
		// For each page, generate two tabs with a some content.
		component1 := gen("Tab 1", "tab1", request.Path)
		component2 := gen("Tab 2", "tab2", request.Path)

		contentResponse := component.NewContentResponse(component.TitleFromString("Example"))
		contentResponse.Add(component1, component2)

		return *contentResponse, nil
	})
}
