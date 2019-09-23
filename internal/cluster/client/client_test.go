/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package client

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vmware/octant/internal/cluster"
)

func TestClusterClientManager_Get(t *testing.T) {
	options := cluster.RESTConfigOptions{}
	kubeConfig := filepath.Join("../testdata", "kubeconfig.yaml")

	ccm, err := NewClusterClientManager(context.TODO(), kubeConfig, options)
	require.NoError(t, err)

	// contextName take from ../testdata/kubeconfig.yaml file.
	contextName := "my-cluster"
	_, err = ccm.Get(context.TODO(), contextName)
	require.NoError(t, err)

	_, err = ccm.Get(context.TODO(), "missing")
	require.Error(t, err)
}
