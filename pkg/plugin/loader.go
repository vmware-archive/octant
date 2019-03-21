package plugin

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

const configDir = "vmdash"

// Config is configuration for the plugin manager.
type Config interface {
	// PluginDir returns the location of the plugin directory.
	PluginDir() (string, error)
	// Home returns the user's home directory.
	Home() string
	// Fs is the afero filesystem
	Fs() afero.Fs
}

type defaultConfig struct {
	fs     afero.Fs
	homeFn func() string
}

var (
	// DefaultConfig is the default plugin manager configuration.
	DefaultConfig = &defaultConfig{}
)

var _ Config = (*defaultConfig)(nil)

// PluginDir returns the plugin directory. Current only works on macOS and Linux
// and not in a container.
func (c *defaultConfig) PluginDir() (string, error) {
	home := c.Home()

	if home == "" {
		// home could be blank if running in a container, so bail out...
		return "", errors.Errorf("running dash in a container is not yet supported: No $HOME env var")
	}

	return filepath.Join(home, ".config", configDir, "plugins"), nil
}

func (c *defaultConfig) Home() string {
	if c.homeFn == nil {
		c.homeFn = func() string {
			// TODO: make me work in windows
			return os.Getenv("HOME")
		}
	}

	return c.homeFn()
}

func (c *defaultConfig) Fs() afero.Fs {
	if c.fs == nil {
		c.fs = afero.NewOsFs()
	}

	return c.fs
}

// AvailablePlugins returns a list of available plugins.
func AvailablePlugins(config Config) ([]string, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	dir, err := config.PluginDir()
	if err != nil {
		return nil, errors.Wrap(err, "get plugin directory")
	}

	_, err = config.Fs().Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, errors.Wrap(err, "check plugin directory")
	}

	fis, err := afero.ReadDir(config.Fs(), dir)
	if err != nil {
		return nil, errors.Wrap(err, "read files in plugin directory")
	}

	var list []string

	for _, fi := range fis {
		mode := fi.Mode()

		if mode|64 == mode {
			pluginPath := filepath.Join(dir, fi.Name())
			list = append(list, pluginPath)
		}
	}

	sort.Strings(list)

	return list, nil
}
