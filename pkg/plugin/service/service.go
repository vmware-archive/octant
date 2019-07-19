/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package service

import (
	"strings"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware/octant/pkg/action"
	"github.com/vmware/octant/pkg/navigation"
	"github.com/vmware/octant/pkg/plugin"
	"github.com/vmware/octant/pkg/view/component"
)

func defaultServerFactory(service plugin.Service) {
	plugin.Serve(service)
}

// PluginOption is an option for configuring Plugin.
type PluginOption func(p *Plugin)

// Plugin is a plugin service helper.
type Plugin struct {
	pluginHandler *Handler
	serverFactory func(service plugin.Service)
}

// Register registers a plugin with Octant.
func Register(name, description string, capabilities *plugin.Capabilities, handlers HandlerFuncs, options ...PluginOption) (*Plugin, error) {
	p := &Plugin{
		pluginHandler: &Handler{
			name:             name,
			description:      description,
			capabilities:     capabilities,
			HandlerFuncs:     handlers,
			dashboardFactory: NewDashboardClient,
		},

		serverFactory: defaultServerFactory,
	}

	for _, option := range options {
		option(p)
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

// HandlerFuncs are functions for configuring a plugin.
type HandlerFuncs struct {
	Print        func(dashboardClient Dashboard, object runtime.Object) (plugin.PrintResponse, error)
	PrintTab     func(dashboardClient Dashboard, object runtime.Object) (*component.Tab, error)
	ObjectStatus func(dashboardClient Dashboard, object runtime.Object) (plugin.ObjectStatusResponse, error)
	HandleAction func(dashboardClient Dashboard, payload action.Payload) error
	Navigation   func(dashboardClient Dashboard) (navigation.Navigation, error)
	Content      func(dashboardClient Dashboard, contentPath string) (component.ContentResponse, error)
}
