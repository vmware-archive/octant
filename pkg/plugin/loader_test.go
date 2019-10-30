/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_AvailablePlugins(t *testing.T) {
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("OCTANT")
	viper.AutomaticEnv()

	tests := []struct {
		homePath string
		key      string
	}{
		{
			homePath: filepath.Join("/home", "user"),
			key:      "plugin-path",
		},
		{
			homePath: filepath.Join("/home", "xdg_config_path"),
			key:      "xdg-config-home",
		},
	}

	for _, test := range tests {
		viper.Set(test.key, test.homePath)

		fs := afero.NewMemMapFs()

		c := &defaultConfig{
			fs: fs,
			homeFn: func() string {
				return test.homePath
			},
		}

		switch test.key {
		case "plugin-path":
			customPath := "/example/test"
			envPaths := customPath + ":/another/one"
			viper.Set(test.key, envPaths)

			configPath := filepath.Join(test.homePath, ".config", configDir, "plugins")

			err := fs.MkdirAll(configPath, 0700)
			require.NoError(t, err, "unable to create test home directory")

			if viper.GetString(test.key) != "" {
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

		case "xdg-config-home":
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
