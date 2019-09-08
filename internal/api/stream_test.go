/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_notFoundRedirectPath(t *testing.T) {
	cases := []struct {
		name     string
		expected string
	}{
		{
			name:     "overview/namespace/default/workloads/deployments/nginx-deployment/",
			expected: "overview/namespace/default/workloads/deployments",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := notFoundRedirectPath(tc.name)
			assert.Equal(t, tc.expected, got)
		})
	}
}
