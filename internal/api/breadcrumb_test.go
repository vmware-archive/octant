/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/api"
	"github.com/vmware-tanzu/octant/internal/util/json"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_NavigationFromPathNamespace(t *testing.T) {
	tests := []struct {
		name              string
		path              string
		parentTitle       string
		parentUrl         string
		expectedTitle     string
		expectedSelection string
		expectedItems     int
	}{
		{
			name:          "Namespace Overview (root level)",
			path:          "overview/namespace/milan",
			expectedItems: 0,
		},
		{
			name:              "Config and Storage (1st level)",
			path:              "overview/namespace/milan/config-and-storage",
			parentTitle:       "Namespace Overview",
			parentUrl:         "overview/namespace/milan",
			expectedTitle:     "Namespace Overview",
			expectedSelection: "Config and Storage",
			expectedItems:     7,
		},
		{
			name:              "Config Maps (2nd level)",
			path:              "overview/namespace/milan/config-and-storage/config-maps",
			parentTitle:       "Config and Storage",
			parentUrl:         "overview/namespace/milan/config-and-storage",
			expectedTitle:     "Config Maps",
			expectedSelection: "Config Maps",
			expectedItems:     4,
		},
		{
			name:          "Specific Config Map (3rd level)",
			path:          "overview/namespace/milan/config-and-storage/config-maps/kafka-config",
			expectedItems: 0,
		},
		{
			name:              "Custom Resources (1st level)",
			path:              "overview/namespace/milan/custom-resources",
			parentTitle:       "Namespace Overview",
			parentUrl:         "overview/namespace/milan",
			expectedTitle:     "Namespace Overview",
			expectedSelection: "Custom Resources",
			expectedItems:     7,
		},
		{
			name:              "Custom Resources (2nd level)",
			path:              "overview/namespace/milan/custom-resources/brokers.eventing.knative.dev",
			parentTitle:       "Custom Resources",
			parentUrl:         "overview/namespace/milan/custom-resources",
			expectedTitle:     "brokers.eventing.knative.dev",
			expectedSelection: "brokers.eventing.knative.dev",
			expectedItems:     2,
		},
		{
			name:          "Custom Resources (3rd level)",
			path:          "overview/namespace/milan/custom-resources/brokers.eventing.knative.dev/default",
			expectedItems: 0,
		},
		{
			name:              "Events (1st level no children)",
			path:              "overview/namespace/milan/events",
			parentTitle:       "Namespace Overview",
			parentUrl:         "overview/namespace/milan",
			expectedTitle:     "Namespace Overview",
			expectedSelection: "Events",
			expectedItems:     7,
		},
		{
			name:          "Specific Event (2nd level)",
			path:          "overview/namespace/milan/events/kafka-0.16501e3e90730f40",
			expectedItems: 0,
		},
	}
	data, err := ioutil.ReadFile(filepath.Join("testdata", "namespace_navigation.json"))
	require.NoError(t, err)
	var namespaceNavigation []navigation.Navigation

	err = json.Unmarshal([]byte(data), &namespaceNavigation)
	require.NoError(t, err)
	require.Equal(t, 7, len(namespaceNavigation))

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nav, parent, selection := api.NavigationFromPath(namespaceNavigation, test.path)
			require.NotNil(t, nav)
			require.Equal(t, test.expectedItems, len(nav))
			require.Equal(t, test.parentTitle, parent.Title)
			require.Equal(t, test.parentUrl, parent.Url)
			if test.expectedItems > 0 {
				require.Equal(t, test.expectedTitle, nav[0].Title)
				require.Equal(t, test.expectedSelection, selection.Title)
			}
		})
	}
}

func Test_NavigationFromPathCluster(t *testing.T) {
	tests := []struct {
		name              string
		path              string
		parentTitle       string
		parentUrl         string
		expectedTitle     string
		expectedSelection string
		expectedItems     int
	}{
		{
			name:          "Cluster Overview (Root)",
			path:          "cluster-overview",
			expectedTitle: "",
			expectedItems: 0,
		},
		{
			name:              "Webhooks (1st level)",
			path:              "cluster-overview/webhooks",
			parentTitle:       "Cluster Overview",
			parentUrl:         "cluster-overview",
			expectedTitle:     "Cluster Overview",
			expectedSelection: "Webhooks",
			expectedItems:     9,
		},
		{
			name:              "Mutating Webhooks (2nd level)",
			path:              "cluster-overview/webhooks/mutating-webhooks",
			parentTitle:       "Webhooks",
			parentUrl:         "cluster-overview/webhooks",
			expectedTitle:     "Mutating Webhooks",
			expectedSelection: "Mutating Webhooks",
			expectedItems:     2,
		},
		{
			name:          "Specific Webhook (3rd level)",
			path:          "cluster-overview/webhooks/mutating-webhooks/cert-manager-webhook",
			expectedTitle: "",
			expectedItems: 0,
		},
		{
			name:              "Custom Resources (1st level)",
			path:              "cluster-overview/custom-resources",
			parentTitle:       "Cluster Overview",
			parentUrl:         "cluster-overview",
			expectedTitle:     "Cluster Overview",
			expectedSelection: "Custom Resources",
			expectedItems:     9,
		},
		{
			name:              "Custom Resources (2nd level)",
			path:              "cluster-overview/custom-resources/storagestates.migration.k8s.io",
			parentTitle:       "Custom Resources",
			parentUrl:         "cluster-overview/custom-resources",
			expectedTitle:     "clusterissuers.cert-manager.io",
			expectedSelection: "storagestates.migration.k8s.io",
			expectedItems:     7,
		},
		{
			name:              "Namespaces (1st level no children)",
			path:              "cluster-overview/namespaces",
			parentTitle:       "Cluster Overview",
			parentUrl:         "cluster-overview",
			expectedTitle:     "Cluster Overview",
			expectedSelection: "Namespaces",
			expectedItems:     9,
		},
		{
			name:          "Specific Namespace (2nd level)",
			path:          "cluster-overview/namespaces/milan",
			expectedItems: 0,
		},
	}
	data, err := ioutil.ReadFile(filepath.Join("testdata", "cluster_navigation.json"))
	require.NoError(t, err)
	var namespaceNavigation []navigation.Navigation

	err = json.Unmarshal([]byte(data), &namespaceNavigation)
	require.NoError(t, err)
	require.Equal(t, 9, len(namespaceNavigation))

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nav, parent, selection := api.NavigationFromPath(namespaceNavigation, test.path)
			require.NotNil(t, nav)
			require.Equal(t, test.expectedItems, len(nav))
			require.Equal(t, test.parentTitle, parent.Title)
			require.Equal(t, test.parentUrl, parent.Url)
			if test.expectedItems > 0 {
				require.Equal(t, test.expectedTitle, nav[0].Title)
				require.Equal(t, test.expectedSelection, selection.Title)
			}
		})
	}
}

func Test_CreateNavigationBreadcrumbNamespace(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		expectedTitle string
		lastTitle     string
		lastUrl       string
		expectedItems int
	}{
		{
			name:          "Namespace Overview (root level)",
			path:          "overview/namespace/milan",
			expectedItems: 0,
		},
		{
			name:          "Config and Storage (1st level)",
			path:          "overview/namespace/milan/config-and-storage",
			expectedTitle: "Namespace Overview",
			lastTitle:     "Config and Storage",
			lastUrl:       "overview/namespace/milan/config-and-storage",
			expectedItems: 1,
		},
		{
			name:          "Config Maps (2nd level)",
			path:          "overview/namespace/milan/config-and-storage/config-maps",
			expectedTitle: "Config Maps",
			lastTitle:     "Config Maps",
			lastUrl:       "overview/namespace/milan/config-and-storage/config-maps",
			expectedItems: 2,
		},
		{
			name:          "Specific Config Map (3rd level)",
			path:          "overview/namespace/milan/config-and-storage/config-maps/kafka-config",
			expectedItems: 0,
		},
		{
			name:          "Custom Resources (1st level)",
			path:          "overview/namespace/milan/custom-resources",
			expectedTitle: "Namespace Overview",
			lastTitle:     "Custom Resources",
			lastUrl:       "overview/namespace/milan/custom-resources",
			expectedItems: 1,
		},
		{
			name:          "Custom Resources (2nd level)",
			path:          "overview/namespace/milan/custom-resources/brokers.eventing.knative.dev",
			expectedTitle: "brokers.eventing.knative.dev",
			lastTitle:     "brokers.eventing.knative.dev",
			lastUrl:       "overview/namespace/milan/custom-resources/brokers.eventing.knative.dev",
			expectedItems: 2,
		},
		{
			name:          "Custom Resources (3rd level)",
			path:          "overview/namespace/milan/custom-resources/brokers.eventing.knative.dev/default",
			expectedItems: 0,
		},
		{
			name:          "Events (1st level no children)",
			path:          "overview/namespace/milan/events",
			expectedTitle: "Namespace Overview",
			lastTitle:     "Events",
			lastUrl:       "overview/namespace/milan/events",
			expectedItems: 1,
		},
		{
			name:          "Specific Event (2nd level)",
			path:          "overview/namespace/milan/events/kafka-0.16501e3e90730f40",
			expectedItems: 0,
		},
	}
	data, err := ioutil.ReadFile(filepath.Join("testdata", "namespace_navigation.json"))
	require.NoError(t, err)
	var namespaceNavigation []navigation.Navigation

	err = json.Unmarshal([]byte(data), &namespaceNavigation)
	require.NoError(t, err)
	require.Equal(t, 7, len(namespaceNavigation))

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			last, title := api.CreateNavigationBreadcrumb(namespaceNavigation, test.path)
			if test.expectedItems > 0 {
				require.NotNil(t, title)
				require.Equal(t, test.expectedItems, len(title))
				require.Equal(t, test.lastTitle, last.Title)
				require.Equal(t, test.lastUrl, last.Url)
			}
		})
	}
}

func Test_CreateNavigationBreadcrumbCluster(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		expectedTitle string
		lastTitle     string
		lastUrl       string
		expectedItems int
	}{
		{
			name:          "Cluster Overview (Root)",
			path:          "cluster-overview",
			expectedTitle: "",
			expectedItems: 0,
		},
		{
			name:          "Webhooks (1st level)",
			path:          "cluster-overview/webhooks",
			expectedTitle: "Cluster Overview",
			lastTitle:     "Webhooks",
			lastUrl:       "cluster-overview/webhooks",
			expectedItems: 1,
		},
		{
			name:          "Mutating Webhooks (2nd level)",
			path:          "cluster-overview/webhooks/mutating-webhooks",
			expectedTitle: "Mutating Webhooks",
			lastTitle:     "Mutating Webhooks",
			lastUrl:       "cluster-overview/webhooks/mutating-webhooks",
			expectedItems: 2,
		},
		{
			name:          "Specific Webhook (3rd level)",
			path:          "cluster-overview/webhooks/mutating-webhooks/cert-manager-webhook",
			expectedTitle: "",
			expectedItems: 0,
		},
		{
			name:          "Custom Resources (1st level)",
			path:          "cluster-overview/custom-resources",
			expectedTitle: "Cluster Overview",
			lastTitle:     "Custom Resources",
			lastUrl:       "cluster-overview/custom-resources",
			expectedItems: 1,
		},
		{
			name:          "Custom Resources (2nd level)",
			path:          "cluster-overview/custom-resources/storagestates.migration.k8s.io",
			expectedTitle: "clusterissuers.cert-manager.io",
			lastTitle:     "storagestates.migration.k8s.io",
			lastUrl:       "cluster-overview/custom-resources/storagestates.migration.k8s.io",
			expectedItems: 2,
		},
		{
			name:          "Namespaces (1st level no children)",
			path:          "cluster-overview/namespaces",
			expectedTitle: "Cluster Overview",
			lastTitle:     "Namespaces",
			lastUrl:       "cluster-overview/namespaces",
			expectedItems: 1,
		},
		{
			name:          "Specific Namespace (2nd level)",
			path:          "cluster-overview/namespaces/milan",
			expectedItems: 0,
		},
	}
	data, err := ioutil.ReadFile(filepath.Join("testdata", "cluster_navigation.json"))
	require.NoError(t, err)
	var namespaceNavigation []navigation.Navigation

	err = json.Unmarshal([]byte(data), &namespaceNavigation)
	require.NoError(t, err)
	require.Equal(t, 9, len(namespaceNavigation))

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			last, title := api.CreateNavigationBreadcrumb(namespaceNavigation, test.path)
			if test.expectedItems > 0 {
				require.NotNil(t, title)
				require.Equal(t, test.expectedItems, len(title))
				require.Equal(t, test.lastTitle, last.Title)
				require.Equal(t, test.lastUrl, last.Url)
			}
		})
	}
}

func Test_CreateNavigationBreadcrumbApplications(t *testing.T) {
	data, err := ioutil.ReadFile(filepath.Join("testdata", "application_navigation.json"))
	require.NoError(t, err)
	var namespaceNavigation []navigation.Navigation

	err = json.Unmarshal([]byte(data), &namespaceNavigation)
	require.NoError(t, err)
	require.Equal(t, 1, len(namespaceNavigation))

	tests := []struct {
		name          string
		path          string
		lastTitle     string
		lastUrl       string
		expectedTitle component.TitleComponent
		expectedItems int
	}{
		{
			name:          "Applications Detail Breadcumb",
			path:          "workloads/namespace/milan/detail/simple-app",
			expectedTitle: component.NewLink("", "Applications", "/workloads/namespace/milan"),
			expectedItems: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			last, title := api.CreateNavigationBreadcrumb(namespaceNavigation, test.path)
			require.Equal(t, test.expectedItems, len(title))
			require.Equal(t, test.expectedTitle, title[0])
			require.NotNil(t, title)
			require.Equal(t, "", last.Title)
			require.Equal(t, "", last.Url)
		})
	}
}
