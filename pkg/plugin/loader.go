/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

const configDir = "octant"

// Config is configuration for the plugin manager.
type Config interface {
	// PluginDirs returns the location of the plugin directories.
	PluginDirs() ([]string, error)
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

// PluginDirs returns the plugin directories. Current only works on macOS and Linux
// and not in a container.
func (c *defaultConfig) PluginDirs() ([]string, error) {
	home := c.Home()

	if home == "" {
		// home could be blank if running in a container, so bail out...
		return []string{}, errors.Errorf("running dash in a container is not yet supported: No $HOME env var")
	}

	defaultDir := filepath.Join(home, ".config", configDir, "plugins")

	if runtime.GOOS == "windows" || os.Getenv("XDG_CONFIG_HOME") != "" {
		defaultDir = filepath.Join(home, configDir, "plugins")
	}

	if path := os.Getenv("OCTANT_PLUGIN_PATH"); path != "" {
		path = strings.Trim(path, string(filepath.ListSeparator))
		return append(filepath.SplitList(path), defaultDir), nil
	}

	return []string{defaultDir}, nil
}

func (c *defaultConfig) Home() string {
	if c.homeFn == nil {
		c.homeFn = func() string {
			switch runtime.GOOS {
			case "windows":
				return os.Getenv("LOCALAPPDATA")

			case "darwin":
				return os.Getenv("HOME")

			default: // Unix
				if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
					return dir
				}
			}
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

	dirs, err := config.PluginDirs()
	if err != nil {
		return nil, errors.Wrap(err, "get plugin directory")
	}

	var list []string

	for _, dir := range dirs {
		_, err = config.Fs().Stat(dir)
		if err != nil {
			if os.IsNotExist(err) {
				// no-op
				continue
			}
			return nil, errors.Wrap(err, "check plugin directory")
		}

		fis, err := afero.ReadDir(config.Fs(), dir)
		if err != nil {
			return nil, errors.Wrap(err, "read files in plugin directory")
		}

		for _, fi := range fis {
			mode := fi.Mode()

			if mode|64 == mode {
				pluginPath := filepath.Join(dir, fi.Name())
				list = append(list, pluginPath)
			}
		}
	}

	sort.Strings(list)

	return list, nil
}
