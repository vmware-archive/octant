/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"github.com/vmware/octant/pkg/navigation"
	"github.com/vmware/octant/pkg/view/component"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// PluginName is the name of the dashboard plugin.
	PluginName = "plugin"
)

// Capabilities are plugin capabilities.
type Capabilities struct {
	// SupportsPrinterConfig are the GVKs the plugin will print configuration for.
	SupportsPrinterConfig []schema.GroupVersionKind
	// SupportsPrinterStatus are the GVKs the plugin will print status for.
	SupportsPrinterStatus []schema.GroupVersionKind
	// SupportsPrinterItems are the GVKs the plugin will print additional items for.
	SupportsPrinterItems []schema.GroupVersionKind
	// SupportsObjectStatus are the GVKs the plugin will generate object status for.
	SupportsObjectStatus []schema.GroupVersionKind
	// SupportsTab are the GVKs the plugin will create an additional tab for.
	SupportsTab []schema.GroupVersionKind
	// IsModule is true this plugin is a module.
	IsModule bool
}

// HasPrinterSupport returns true if this plugin supports the supplied GVK.
func (c Capabilities) HasPrinterSupport(gvk schema.GroupVersionKind) bool {
	return includesGVK(gvk, c.SupportsPrinterConfig) ||
		includesGVK(gvk, c.SupportsPrinterStatus) ||
		includesGVK(gvk, c.SupportsPrinterItems)
}

// HasTabSupport returns true if this plugins supports creating a tab for
// the supplied GVK.
func (c Capabilities) HasTabSupport(gvk schema.GroupVersionKind) bool {
	return includesGVK(gvk, c.SupportsTab)
}

// PrintResponse is a printer response from the plugin. The dashboard
// will use this to the add the plugin's output to a summary view.
type PrintResponse struct {
	// Config is additional summary sections for configuration.
	Config []component.SummarySection
	// Status is additional summary sections for status.
	Status []component.SummarySection
	// Items are additional view components.
	Items []component.FlexLayoutItem
}

// ObjectStatusResponse is an object status response from plugin.
type ObjectStatusResponse struct {
	// ObjectStatus is status of an object.
	ObjectStatus component.PodSummary
}

// Metadata is plugin metadata.
type Metadata struct {
	Name         string
	Description  string
	Capabilities Capabilities
}

// Service is the interface that is exposed as a plugin. The plugin is required to implement this
// interface.
type Service interface {
	Register(dashboardAPIAddress string) (Metadata, error)
	Print(object runtime.Object) (PrintResponse, error)
	PrintTab(object runtime.Object) (*component.Tab, error)
	ObjectStatus(object runtime.Object) (ObjectStatusResponse, error)
}

// ModuleService is the interface that is exposed as a plugin as a module. The plugin is required to implement this
// interface.
type ModuleService interface {
	Service

	Navigation() (navigation.Navigation, error)
	Content(contentPath string) (component.ContentResponse, error)
}

func includesGVK(gvk schema.GroupVersionKind, list []schema.GroupVersionKind) bool {
	for i := range list {
		if gvk.Group == list[i].Group &&
			gvk.Version == list[i].Version &&
			gvk.Kind == list[i].Kind {
			return true
		}
	}

	return false
}
