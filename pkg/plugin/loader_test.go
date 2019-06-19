/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
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
	pluginEnv := "OCTANT_PLUGIN_PATH"
	defer os.Unsetenv(pluginEnv)
	fs := afero.NewMemMapFs()

	homePath := filepath.Join("/home", "user")

	c := &defaultConfig{
		fs: fs,
		homeFn: func() string {
			return homePath
		},
	}

	customPath := "/example/test"
	envPaths := customPath + ":/another/one"
	os.Setenv(pluginEnv, envPaths)

	configPath := filepath.Join(homePath, ".config", configDir, "plugins")

	err := fs.MkdirAll(configPath, 0700)
	require.NoError(t, err, "unable to create test home directory")

	if os.Getenv(pluginEnv) != "" {
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
		"/home/user/.config/vmdash/plugins/a-plugin",
		"/home/user/.config/vmdash/plugins/z-plugin",
	}

	assert.Equal(t, expected, got)
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
