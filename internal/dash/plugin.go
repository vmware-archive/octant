/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package dash

import (
	"github.com/pkg/errors"

	"github.com/vmware/octant/internal/module"
	"github.com/vmware/octant/pkg/action"
	"github.com/vmware/octant/pkg/plugin"
	"github.com/vmware/octant/pkg/plugin/api"
)

func initPlugin(moduleManager module.ManagerInterface, actionManager *action.Manager, service api.Service) (*plugin.Manager, error) {
	apiService, err := api.New(service)
	if err != nil {
		return nil, errors.Wrap(err, "create dashboard api")
	}

	m := plugin.NewManager(apiService, moduleManager, actionManager)

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
