package dash

import (
	"context"

	"github.com/heptio/developer-dash/pkg/plugin"
	"github.com/pkg/errors"
)

func initPlugin(ctx context.Context) (*plugin.Manager, error) {
	m := plugin.NewManager()

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
