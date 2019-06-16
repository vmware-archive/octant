package dash

import (
	"context"

	"github.com/heptio/developer-dash/pkg/store"
	"github.com/heptio/developer-dash/internal/portforward"
	"github.com/heptio/developer-dash/pkg/plugin"
	"github.com/heptio/developer-dash/pkg/plugin/api"
	"github.com/pkg/errors"
)

func initPlugin(ctx context.Context, portForwarder portforward.PortForwarder, appObjectStore store.Store) (*plugin.Manager, error) {
	service := &api.GRPCService{
		ObjectStore:   appObjectStore,
		PortForwarder: portForwarder,
	}

	apiService, err := api.New(service)
	if err != nil {
		return nil, errors.Wrap(err, "create dashboard api")
	}

	m := plugin.NewManager(apiService)

	pluginList, err := plugin.AvailablePlugins(plugin.DefaultConfig)
	if err != nil {
		return nil, errors.Wrap(err, "finding available plugins")
	}

	for _, pluginPath := range pluginList {
		if err := m.Load(pluginPath); err != nil {
			return nil, errors.Wrapf(err, "initialize plugin %q", pluginPath)
		}

	}

	if err := m.Start(ctx); err != nil {
		return nil, errors.Wrap(err, "start plugin manager")
	}

	return m, nil
}
