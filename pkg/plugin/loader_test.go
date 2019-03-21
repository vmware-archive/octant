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
	fs := afero.NewMemMapFs()

	homePath := filepath.Join("/home", "user")

	c := &defaultConfig{
		fs: fs,
		homeFn: func() string {
			return homePath
		},
	}

	configPath := filepath.Join(homePath, ".config", configDir, "plugins")

	err := fs.MkdirAll(configPath, 0700)
	require.NoError(t, err, "unable to create test home directory")

	stagePlugin := func(t *testing.T, name string, mode os.FileMode) {
		p := filepath.Join(configPath, name)
		err = afero.WriteFile(fs, p, []byte("guts"), mode)
		require.NoError(t, err)
	}

	stagePlugin(t, "z-plugin", 0755)
	stagePlugin(t, "a-plugin", 0755)
	stagePlugin(t, "not-a-plugin", 0600)

	got, err := AvailablePlugins(c)
	require.NoError(t, err)

	expected := []string{
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

	expected := []string{}

	assert.Equal(t, expected, got)
}
