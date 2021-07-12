/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package websockets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebsocketClient_SetContentPath(t *testing.T) {
}

func Test_updateContentPathNamespace(t *testing.T) {
	tests := []struct {
		name      string
		in        string
		namespace string
		expected  string
	}{
		{
			name:      "content path with namespace",
			in:        "overview/namespace/default/workloads",
			namespace: "other",
			expected:  "overview/namespace/other/workloads",
		},
		{
			name:      "root content path with namespace",
			in:        "overview/namespace/default",
			namespace: "other",
			expected:  "overview/namespace/other",
		},
		{
			name:      "cluster scoped path",
			in:        "cluster-overview/path",
			namespace: "other",
			expected:  "cluster-overview/path",
		},
		{
			name:      "empty content path",
			in:        "",
			namespace: "other",
			expected:  "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := updateContentPathNamespace(test.in, test.namespace)
			assert.Equal(t, test.expected, got)
		})
	}
}
