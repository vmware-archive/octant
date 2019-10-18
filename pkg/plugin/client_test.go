/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware-tanzu/octant/internal/gvk"
)

func TestCapabilities_HasPrinterSupport(t *testing.T) {
	cases := []struct {
		name         string
		in           schema.GroupVersionKind
		capabilities Capabilities
		hasSupport   bool
	}{
		{
			name: "with printer support",
			in:   gvk.Pod,
			capabilities: Capabilities{
				SupportsPrinterConfig: []schema.GroupVersionKind{gvk.Pod},
			},
			hasSupport: true,
		},
		{
			name: "with out printer support",
			in:   gvk.Deployment,
			capabilities: Capabilities{
				SupportsPrinterConfig: []schema.GroupVersionKind{gvk.Pod},
			},
			hasSupport: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.hasSupport, tc.capabilities.HasPrinterSupport(tc.in))
		})
	}
}

func TestCapabilities_HasTabSupport(t *testing.T) {
	cases := []struct {
		name         string
		in           schema.GroupVersionKind
		capabilities Capabilities
		hasSupport   bool
	}{
		{
			name: "with tab support",
			in:   gvk.Pod,
			capabilities: Capabilities{
				SupportsTab: []schema.GroupVersionKind{gvk.Pod},
			},
			hasSupport: true,
		},
		{
			name: "with out tab support",
			in:   gvk.Deployment,
			capabilities: Capabilities{
				SupportsTab: []schema.GroupVersionKind{gvk.Pod},
			},
			hasSupport: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.hasSupport, tc.capabilities.HasTabSupport(tc.in))
		})
	}
}
