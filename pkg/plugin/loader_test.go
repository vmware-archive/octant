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

func Test_PluginDirs(t *testing.T) {
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("OCTANT")
	viper.AutomaticEnv()

	tests := []struct {
		expectedPaths []string
		homePath      string
		customPath    string
		key           string
	}{
		{
			homePath:      filepath.Join("/home", "userA"),
			expectedPaths: []string{filepath.Join("/home", "userA", ".config", "octant", "plugins")},
		},
		{
			homePath:      filepath.Join("/home", "userB"),
			expectedPaths: []string{filepath.Join("/home", "userB", ".config", "octant", "plugins")},
		},
		{
			homePath:   filepath.Join("/home", "userC"),
			customPath: filepath.Join("/my", "custom", "path"),
			expectedPaths: []string{
				filepath.Join("/my", "custom", "path"),
				filepath.Join("/home", "userC", ".config", "octant", "plugins"),
			},
		},
	}

	for _, test := range tests {
		viper.Set("home", test.homePath)

		fs := afero.NewMemMapFs()
		for _, expectedPath := range test.expectedPaths {
			err := fs.MkdirAll(expectedPath, 0700)
			require.NoError(t, err, "unable to create test home directory")
		}

		if test.customPath != "" {
			err := fs.MkdirAll(test.customPath, 0700)
			require.NoError(t, err, "unable to create test home directory")
			viper.Set("plugin-path", test.customPath)
		}

		c := &defaultConfig{
			fs: fs,
			os: "unix",
		}

		results, err := c.PluginDirs(test.homePath)
		require.NoError(t, err)
		assert.Equal(t, test.expectedPaths, results)
	}
}

func Test_AvailablePlugins(t *testing.T) {
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
		replacer := strings.NewReplacer("-", "_")
		viper.SetEnvKeyReplacer(replacer)
		viper.SetEnvPrefix("OCTANT")
		viper.AutomaticEnv()

		fs := afero.NewMemMapFs()

		c := &defaultConfig{
			fs: fs,
			homeFn: func() string {
				return test.homePath
			},
		}

		switch test.key {
		case "plugin-path":
			c.os = "unix"
			customPath := filepath.Join("/example", "test")
			envPaths := customPath + string(filepath.ListSeparator) + filepath.Join("/another", "one")
			viper.Set("plugin-path", envPaths)

			configPath := filepath.Join(test.homePath, ".config", configDir, "plugins")

			err := fs.MkdirAll(configPath, 0700)
			require.NoError(t, err, "unable to create test home directory")

			for _, path := range filepath.SplitList(envPaths) {
				err := fs.MkdirAll(path, 0700)
				require.NoError(t, err, "unable to create directory from environment variable")
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
				filepath.Join("/example", "test", "e-plugin"),
				filepath.Join(configPath, "a-plugin"),
				filepath.Join(configPath, "z-plugin"),
			}

			assert.Equal(t, expected, got)

		case "xdg-config-home":
			configPath := filepath.Join(test.homePath, configDir, "plugins")
			err := fs.MkdirAll(configPath, 0700)
			require.NoError(t, err, "unable to create test home directory")
			viper.Set("xdg-config-home", test.homePath)
			viper.Set("plugin-path", "")

			stagePlugin := func(t *testing.T, path string, name string, mode os.FileMode) {
				p := filepath.Join(path, name)
				err = afero.WriteFile(fs, p, []byte("guts"), mode)
				require.NoError(t, err)
			}

			stagePlugin(t, configPath, "a-plugin", 0755)

			got, err := AvailablePlugins(c)
			require.NoError(t, err)

			expected := []string{
				filepath.Join("/home", "xdg_config_path", "octant", "plugins", "a-plugin"),
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
