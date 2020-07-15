/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package dash

import (
	"fmt"

	"github.com/vmware-tanzu/octant/internal/module"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/api"
)

func initPlugin(moduleManager module.ManagerInterface, actionManager *action.Manager, service api.Service) (*plugin.Manager, error) {
	apiService, err := api.New(service)
	if err != nil {
		return nil, fmt.Errorf("create dashboard api: %w", err)
	}

	m := plugin.NewManager(apiService, moduleManager, actionManager)

	pluginList, err := plugin.AvailablePlugins(plugin.DefaultConfig)
	if err != nil {
		return nil, fmt.Errorf("finding available plugins: %w", err)
	}

	for _, pluginPath := range pluginList {
		if plugin.IsJavaScriptPlugin(pluginPath) {
			continue
		}
		if err := m.Load(pluginPath); err != nil {
			return nil, fmt.Errorf("initialize plugin %q: %w", pluginPath, err)
		}

	}

	return m, nil
}
