/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package dash

import (
	"github.com/pkg/errors"

	"github.com/vmware/octant/internal/module"
	"github.com/vmware/octant/internal/portforward"
	"github.com/vmware/octant/pkg/plugin"
	"github.com/vmware/octant/pkg/plugin/api"
	"github.com/vmware/octant/pkg/store"
)

func initPlugin(portForwarder portforward.PortForwarder, appObjectStore store.Store, moduleManager module.ManagerInterface) (*plugin.Manager, error) {
	service := &api.GRPCService{
		ObjectStore:   appObjectStore,
		PortForwarder: portForwarder,
	}

	apiService, err := api.New(service)
	if err != nil {
		return nil, errors.Wrap(err, "create dashboard api")
	}

	m := plugin.NewManager(apiService, moduleManager)

	pluginList, err := plugin.AvailablePlugins(plugin.DefaultConfig)
	if err != nil {
		return nil, errors.Wrap(err, "finding available plugins")
	}

	for _, pluginPath := range pluginList {
		if err := m.Load(pluginPath); err != nil {
			return nil, errors.Wrapf(err, "initialize plugin %q", pluginPath)
		}

	}

	return m, nil
}
