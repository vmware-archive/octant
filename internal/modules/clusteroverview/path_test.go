/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package clusteroverview

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_crdPath(t *testing.T) {
	got, err := crdPath("namespace", "crdName", "name")
	require.NoError(t, err)

	expected := path.Join("/content", "cluster-overview", "custom-resources", "crdName", "name")
	assert.Equal(t, expected, got)
}

func Test_gvk_path(t *testing.T) {
	tests := []struct {
		name       string
		apiVersion string
		kind       string
		objectName string
		expected   string
		isErr      bool
	}{
		{
			name:       "ClusterRole",
			apiVersion: rbacAPIVersion,
			kind:       "ClusterRole",
			objectName: "cluster-role",
			expected:   path.Join("/content", "cluster-overview", "rbac", "cluster-roles", "cluster-role"),
		},
		{
			name:       "ClusterRoleBinding",
			apiVersion: rbacAPIVersion,
			kind:       "ClusterRoleBinding",
			objectName: "cluster-role-binding",
			expected:   path.Join("/content", "cluster-overview", "rbac", "cluster-role-bindings", "cluster-role-binding"),
		},
		{
			name:       "unknown",
			apiVersion: "unknown",
			kind:       "ClusterRoleBinding",
			objectName: "cluster-role-binding",
			isErr:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := gvkPath("", test.apiVersion, test.kind, test.objectName)
			if test.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, test.expected, got)
		})
	}
}
