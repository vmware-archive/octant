/*
Copyright (c) 2021 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package path_util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NamespacedPath(t *testing.T) {
	cases := []struct {
		name      string
		base      string
		namespace string
		args      []string
		expected  string
	}{
		{
			name:      "base path, no args",
			base:      "overview",
			namespace: "default",
			args:      []string{},
			expected:  "overview/namespace/default",
		},
		{
			name:      "base path, single arg",
			base:      "overview",
			namespace: "default",
			args:      []string{"arg"},
			expected:  "overview/namespace/default/arg",
		},
		{
			name:      "base path, two args",
			base:      "overview",
			namespace: "default",
			args:      []string{"arg1", "arg2"},
			expected:  "overview/namespace/default/arg1/arg2",
		},
		{
			name:      "no base path, two args",
			base:      "",
			namespace: "default",
			args:      []string{"arg1", "arg2"},
			expected:  "namespace/default/arg1/arg2",
		},
	}

	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			got := NamespacedPath(tc.base, tc.namespace, tc.args...)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func Test_PrefixedPath(t *testing.T) {
	cases := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "base path no slash",
			path:     "path",
			expected: "/path",
		},
		{
			name:     "base path with slash",
			path:     "/path",
			expected: "/path",
		},
		{
			name:     "empty string",
			path:     "",
			expected: "/",
		},
	}

	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			got := PrefixedPath(tc.path)
			assert.Equal(t, tc.expected, got)
		})
	}
}
