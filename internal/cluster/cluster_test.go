/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package cluster

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FromKubeConfig(t *testing.T) {
	kubeConfig := filepath.Join("testdata", "kubeconfig.yaml")
	config := RESTConfigOptions{}

	_, err := FromKubeConfig(context.TODO(), kubeConfig, "", "", config)
	require.NoError(t, err)
}
