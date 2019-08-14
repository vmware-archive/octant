/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package kubeconfig

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

	l := NewFSLoader()

	kc, err := l.Load(strings.Join(paths, ":"))
	require.NoError(t, err)

	expected := &KubeConfig{
		Contexts: []Context{
			{Name: "dev-frontend"},
			{Name: "dev-storage"},
			{Name: "exp-scratch"},
		},
		CurrentContext: "dev-frontend",
	}

	assert.Equal(t, fmt.Sprint(*expected), fmt.Sprint(*kc))
}
