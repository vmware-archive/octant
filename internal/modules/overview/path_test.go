/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package overview

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_crdPath(t *testing.T) {
	got, err := crdPath("default", "crdName", "name")
	require.NoError(t, err)

	expected := path.Join("/overview", "namespace", "default", "custom-resources", "crdName", "name")
	assert.Equal(t, expected, got)
}

func Test_gvk_path(t *testing.T) {
	tests := []struct {
		name       string
		namespace  string
		apiVersion string
		kind       string
		objectName string
		expected   string
		isErr      bool
	}{
		{
			name:       "pod",
			namespace:  "default",
			apiVersion: "v1",
			kind:       "Pod",
			objectName: "pod",
			expected:   path.Join("/overview", "namespace", "default", "workloads", "pods", "pod"),
		},
		{
			name:       "no namespace",
			apiVersion: "v1",
			kind:       "Pod",
			objectName: "pod",
			isErr:      true,
		},
		{
			name:       "unknown",
			namespace:  "default",
			apiVersion: "unknown",
			kind:       "ClusterRoleBinding",
			objectName: "cluster-role-binding",
			isErr:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := gvkPath(test.namespace, test.apiVersion, test.kind, test.objectName)
			if test.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, test.expected, got)
		})
	}
}
