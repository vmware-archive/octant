/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package kubeconfig

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	internalCluster "github.com/vmware-tanzu/octant/internal/cluster"
	"github.com/vmware-tanzu/octant/pkg/cluster"
)

func Test_NewKubeConfigs(t *testing.T) {
	kubeConfig := filepath.Join("testdata", "kubeconfig.yaml")
	config := cluster.RESTConfigOptions{}

	_, err := NewKubeConfigContextManager(
		context.TODO(),
		WithKubeConfigList(kubeConfig),
		FromClusterOption(internalCluster.WithRESTConfigOptions(config)),
	)
	require.NoError(t, err)
}

func Test_SwitchContextUpdatesCurrentContext(t *testing.T) {
	kubeConfigs, err := NewKubeConfigContextManager(
		context.TODO(),
		WithKubeConfigList(filepath.Join("testdata", "kubeconfig.yaml")),
	)
	require.NoError(t, err)

	kubeConfigs.SwitchContext(context.TODO(), "other-context")

	require.Equal(t, "other-context", kubeConfigs.CurrentContext())
}

func Test_SwitchContextToEmptyUpdatesCurrentContextFromFileSystem(t *testing.T) {
	kubeConfigs, err := NewKubeConfigContextManager(
		context.TODO(),
		WithKubeConfigList(filepath.Join("testdata", "kubeconfig.yaml")),
	)
	require.NoError(t, err)

	kubeConfigs.SwitchContext(context.TODO(), "")

	require.Equal(t, "my-cluster", kubeConfigs.CurrentContext())
}

func Test_SwitchContextUpdatesClientNamespace(t *testing.T) {
	kubeConfigs, err := NewKubeConfigContextManager(
		context.TODO(),
		WithKubeConfigList(filepath.Join("testdata", "kubeconfig.yaml")),
	)
	require.NoError(t, err)

	kubeConfigs.SwitchContext(context.TODO(), "other-context")

	require.Equal(t, "non-default", kubeConfigs.ClusterClient().DefaultNamespace())
}

func TestFSLoader_Load(t *testing.T) {
	dir, err := ioutil.TempDir("", "loader-test")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(dir))
	}()

	inputs := []string{"kubeconfig-1.yaml", "kubeconfig-2.yaml"}
	var paths []string
	for i := range inputs {
		data, err := ioutil.ReadFile(filepath.Join("testdata", inputs[i]))
		require.NoError(t, err)
		kubeConfigPath := filepath.Join(dir, inputs[i])
		require.NoError(t, ioutil.WriteFile(kubeConfigPath, data, 0644))
		paths = append(paths, kubeConfigPath)
	}

	kc, err := NewKubeConfigContextManager(
		context.TODO(),
		WithKubeConfigList(strings.Join(paths, string(os.PathListSeparator))),
	)
	require.NoError(t, err)

	assert.Equal(t, "dev-frontend", kc.CurrentContext())
	assert.Equal(t, []Context{
		{Name: "dev-frontend"},
		{Name: "dev-storage"},
		{Name: "exp-scratch"},
	}, kc.Contexts())
}

func Test_NewKubeConfigNoCluster(t *testing.T) {
	noClusterOptions := KubeConfigOption{nil, nil}

	_, err := NewKubeConfigContextManager(
		context.TODO(),
		WithKubeConfigList(filepath.Join("testdata", "kubeconfig.yaml")),
		FromClusterOption(internalCluster.WithRESTConfigOptions(cluster.RESTConfigOptions{})),
		noClusterOptions,
	)
	require.NoError(t, err)
}
