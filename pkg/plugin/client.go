/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

var (
	// PluginName is the name of the dashboard plugin.
	Name = "plugin"
)

// Capabilities are plugin capabilities.
type Capabilities struct {
	// SupportsPrinterConfig are the GVKs the plugin will print configuration for.
	SupportsPrinterConfig []schema.GroupVersionKind `json:",omitempty"`
	// SupportsPrinterStatus are the GVKs the plugin will print status for.
	SupportsPrinterStatus []schema.GroupVersionKind `json:",omitempty"`
	// SupportsPrinterItems are the GVKs the plugin will print additional items for.
	SupportsPrinterItems []schema.GroupVersionKind `json:",omitempty"`
	// SupportsObjectStatus are the GVKs the plugin will generate object status for.
	SupportsObjectStatus []schema.GroupVersionKind `json:",omitempty"`
	// SupportsTab are the GVKs the plugin will create an additional tab for.
	SupportsTab []schema.GroupVersionKind `json:",omitempty"`
	// IsModule is true this plugin is a module.
	IsModule bool `json:",omitempty"`
	// ActionNames is a list of action names this plugin handles
	ActionNames []string `json:",omitempty"`
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

// HasObjectStatusSupport returns true if this plugins supports creating object status for
// the supplied GVK.
func (c Capabilities) HasObjectStatusSupport(gvk schema.GroupVersionKind) bool {
	return includesGVK(gvk, c.SupportsObjectStatus)
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

// TabResponse is a tab printer response from the plugin. The
// dashboard will use this to create an additional tab for
// an object.
type TabResponse struct {
	Tab *component.Tab `json:"tab"`
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
	Register(ctx context.Context, dashboardAPIAddress string) (Metadata, error)
	Print(ctx context.Context, object runtime.Object) (PrintResponse, error)
	PrintTab(ctx context.Context, object runtime.Object) (TabResponse, error)
	ObjectStatus(ctx context.Context, object runtime.Object) (ObjectStatusResponse, error)
	HandleAction(ctx context.Context, actionName string, payload action.Payload) error
}

// ModuleService is the interface that is exposed as a plugin as a module. The plugin is required to implement this
// interface.
type ModuleService interface {
	Service

	Navigation(ctx context.Context) (navigation.Navigation, error)
	Content(ctx context.Context, contentPath string) (component.ContentResponse, error)
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
