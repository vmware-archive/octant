/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_AvailablePlugins(t *testing.T) {
	tests := []struct {
		homePath string
		envVar   string
	}{
		{
			homePath: filepath.Join("/home", "user"),
			envVar:   "OCTANT_PLUGIN_PATH",
		},
		{
			homePath: filepath.Join("/home", "xdg_config_path"),
			envVar:   "XDG_CONFIG_HOME",
		},
	}

	for _, test := range tests {
		defer os.Unsetenv(test.envVar)
		fs := afero.NewMemMapFs()

		c := &defaultConfig{
			fs: fs,
			homeFn: func() string {
				return test.homePath
			},
		}

		switch test.envVar {
		case "OCTANT_PLUGIN_PATH":
			customPath := "/example/test"
			envPaths := customPath + ":/another/one"
			os.Setenv(test.envVar, envPaths)

			configPath := filepath.Join(test.homePath, ".config", configDir, "plugins")

			err := fs.MkdirAll(configPath, 0700)
			require.NoError(t, err, "unable to create test home directory")

			if os.Getenv(test.envVar) != "" {
				for _, path := range filepath.SplitList(envPaths) {
					err := fs.MkdirAll(path, 0700)
					require.NoError(t, err, "unable to create directory from environment variable")
				}
			}

			stagePlugin := func(t *testing.T, path string, name string, mode os.FileMode) {
				p := filepath.Join(path, name)
				err = afero.WriteFile(fs, p, []byte("guts"), mode)
				require.NoError(t, err)
			}

			stagePlugin(t, configPath, "z-plugin", 0755)
			stagePlugin(t, configPath, "a-plugin", 0755)
			stagePlugin(t, configPath, "not-a-plugin", 0600)
			stagePlugin(t, customPath, "e-plugin", 0755)

			got, err := AvailablePlugins(c)
			require.NoError(t, err)

			expected := []string{
				"/example/test/e-plugin",
				"/home/user/.config/octant/plugins/a-plugin",
				"/home/user/.config/octant/plugins/z-plugin",
			}

			assert.Equal(t, expected, got)

		case "XDG_CONFIG_HOME":
			xdgPath := "/home/xdg_config_path"
			os.Setenv(test.envVar, xdgPath)

			configPath := filepath.Join(test.homePath, configDir, "plugins")

			err := fs.MkdirAll(configPath, 0700)
			require.NoError(t, err, "unable to create test home directory")

			stagePlugin := func(t *testing.T, path string, name string, mode os.FileMode) {
				p := filepath.Join(path, name)
				err = afero.WriteFile(fs, p, []byte("guts"), mode)
				require.NoError(t, err)
			}

			stagePlugin(t, configPath, "a-plugin", 0755)

			got, err := AvailablePlugins(c)
			require.NoError(t, err)

			expected := []string{
				"/home/xdg_config_path/octant/plugins/a-plugin",
			}

			assert.Equal(t, expected, got)
		}
	}
}

func Test_AvailablePlugins_no_plugin_dir(t *testing.T) {
	fs := afero.NewMemMapFs()

	homePath := filepath.Join("/home", "user")

	c := &defaultConfig{
		fs: fs,
		homeFn: func() string {
			return homePath
		},
	}

	got, err := AvailablePlugins(c)
	require.NoError(t, err)

	expected := []string(nil)

	assert.Equal(t, expected, got)
}
