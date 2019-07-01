/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware/octant/internal/gvk"
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
			in:   gvk.PodGVK,
			capabilities: Capabilities{
				SupportsPrinterConfig: []schema.GroupVersionKind{gvk.PodGVK},
			},
			hasSupport: true,
		},
		{
			name: "with out printer support",
			in:   gvk.DeploymentGVK,
			capabilities: Capabilities{
				SupportsPrinterConfig: []schema.GroupVersionKind{gvk.PodGVK},
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
			in:   gvk.PodGVK,
			capabilities: Capabilities{
				SupportsTab: []schema.GroupVersionKind{gvk.PodGVK},
			},
			hasSupport: true,
		},
		{
			name: "with out tab support",
			in:   gvk.DeploymentGVK,
			capabilities: Capabilities{
				SupportsTab: []schema.GroupVersionKind{gvk.PodGVK},
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
