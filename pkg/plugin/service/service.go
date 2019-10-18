/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package service

import (
	"context"
	"path"
	"strings"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/plugin"
)

func defaultServerFactory(service plugin.Service) {
	plugin.Serve(service)
}

// PluginOption is an option for configuring Plugin.
type PluginOption func(p *Plugin)

// WithPrinter configures the plugin to have a printer.
func WithPrinter(fn HandlerPrinterFunc) PluginOption {
	return func(p *Plugin) {
		p.pluginHandler.HandlerFuncs.Print = fn
	}
}

// WithTabPrinter configures the plugin to have a tab printer.
func WithTabPrinter(fn HandlerTabPrintFunc) PluginOption {
	return func(p *Plugin) {
		p.pluginHandler.HandlerFuncs.PrintTab = fn
	}
}

// WithObjectStatus configures the plugin to supply object status.
func WithObjectStatus(fn HandlerObjectStatusFunc) PluginOption {
	return func(p *Plugin) {
		p.pluginHandler.HandlerFuncs.ObjectStatus = fn
	}
}

// WithActionHandler configures the plugin to handle actions.
func WithActionHandler(fn HandlerActionFunc) PluginOption {
	return func(p *Plugin) {
		p.pluginHandler.HandlerFuncs.HandleAction = fn
	}
}

// WithNavigation configures the plugin to handle navigation and routes.
func WithNavigation(fn HandlerNavigationFunc, routerInit HandlerInitRoutesFunc) PluginOption {
	return func(p *Plugin) {
		p.pluginHandler.HandlerFuncs.Navigation = fn
		p.pluginHandler.HandlerFuncs.InitRoutes = routerInit
	}
}

// Plugin is a plugin service helper.
type Plugin struct {
	pluginHandler *Handler
	serverFactory func(service plugin.Service)
}

// Register registers a plugin with Octant.
func Register(name, description string, capabilities *plugin.Capabilities, options ...PluginOption) (*Plugin, error) {
	router := NewRouter()

	p := &Plugin{
		pluginHandler: &Handler{
			name:             name,
			description:      description,
			capabilities:     capabilities,
			dashboardFactory: NewDashboardClient,
			router:           router,
		},

		serverFactory: defaultServerFactory,
	}

	for _, option := range options {
		option(p)
	}

	if p.pluginHandler.InitRoutes != nil {
		p.pluginHandler.InitRoutes(router)
	}

	if err := p.Validate(); err != nil {
		return nil, err
	}

	return p, nil
}

// Validate validates this helper.
func (p *Plugin) Validate() error {
	var list []string

	if p.pluginHandler.name == "" {
		list = append(list, "requires name")
	}
	if p.pluginHandler.description == "" {
		list = append(list, "requires description")
	}

	if p.pluginHandler.capabilities == nil {
		list = append(list, "requires capabilities")
	}

	if err := p.pluginHandler.Validate(); err != nil {
		list = append(list, err.Error())
	}

	if len(list) == 0 {
		return nil
	}

	return errors.Errorf("validation errors: %s", strings.Join(list, ", "))
}

// Serve serves a plugin.
func (p *Plugin) Serve() {
	p.serverFactory(p.pluginHandler)
}

type baseRequest struct {
	ctx        context.Context
	pluginName string
}

func newBaseRequest(ctx context.Context, pluginName string) baseRequest {
	return baseRequest{
		ctx:        ctx,
		pluginName: pluginName,
	}
}

func (r *baseRequest) Context() context.Context {
	return r.ctx
}

func (r *baseRequest) GeneratePath(pathParts ...string) string {
	return path.Join(append([]string{r.pluginName}, pathParts...)...)
}

// PrintRequest is a request for printing.
type PrintRequest struct {
	baseRequest

	DashboardClient Dashboard
	Object          runtime.Object
}

// ActionRequest is a request for actions.
type ActionRequest struct {
	baseRequest

	DashboardClient Dashboard
	Payload         action.Payload
}

// NavigationRequest is a request for navigation.
type NavigationRequest struct {
	baseRequest

	DashboardClient Dashboard
}

type HandlerPrinterFunc func(request *PrintRequest) (plugin.PrintResponse, error)
type HandlerTabPrintFunc func(request *PrintRequest) (plugin.TabResponse, error)
type HandlerObjectStatusFunc func(request *PrintRequest) (plugin.ObjectStatusResponse, error)
type HandlerActionFunc func(request *ActionRequest) error
type HandlerNavigationFunc func(request *NavigationRequest) (navigation.Navigation, error)
type HandlerInitRoutesFunc func(router *Router)

// HandlerFuncs are functions for configuring a plugin.
type HandlerFuncs struct {
	Print        HandlerPrinterFunc
	PrintTab     HandlerTabPrintFunc
	ObjectStatus HandlerObjectStatusFunc
	HandleAction HandlerActionFunc
	Navigation   HandlerNavigationFunc
	InitRoutes   HandlerInitRoutesFunc
}
