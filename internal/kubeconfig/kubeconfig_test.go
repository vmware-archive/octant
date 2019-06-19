/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package kubeconfig

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFSLoader_Load(t *testing.T) {
	fs := afero.NewMemMapFs()

	data, err := ioutil.ReadFile(filepath.Join("testdata", "kubeconfig.yaml"))
	require.NoError(t, err)
	require.NoError(t, afero.WriteFile(fs, "/path", data, 0644))

	opt := func(l *FSLoader) {
		l.AppFS = fs
	}

	l := NewFSLoader(opt)

	kc, err := l.Load("/path")
	require.NoError(t, err)

	expected := &KubeConfig{
		Contexts: []Context{
			{Name: "dev-frontend"},
			{Name: "dev-storage"},
			{Name: "exp-scratch"},
		},
		CurrentContext: "dev-frontend",
	}

	assert.Equal(t, expected, kc)
}
