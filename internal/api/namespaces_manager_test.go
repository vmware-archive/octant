/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/api"
	clusterFake "github.com/vmware-tanzu/octant/internal/cluster/fake"
	configFake "github.com/vmware-tanzu/octant/internal/config/fake"
	octantFake "github.com/vmware-tanzu/octant/internal/octant/fake"
	"github.com/vmware-tanzu/octant/pkg/api/fake"
)

func TestNamespacesManager_GenerateNamespaces(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashConfig := configFake.NewMockDash(controller)

	state := octantFake.NewMockState(controller)
	octantClient := fake.NewMockOctantClient(controller)

	namespaces := []string{"default"}

	octantClient.EXPECT().
		Send(api.CreateNamespacesEvent(namespaces))

	poller := api.NewSingleRunPoller()
	manager := api.NewNamespacesManager(dashConfig,
		api.WithNamespacesGeneratorPoller(poller),
		api.WithNamespacesGenerator(func(ctx context.Context, config api.NamespaceManagerConfig) (strings []string, e error) {
			return namespaces, nil
		}))

	ctx := context.Background()
	manager.Start(ctx, state, octantClient)
}

func TestNamespacesGenerator(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		setup    func(controller *gomock.Controller) *configFake.MockDash
		isErr    bool
		expected []string
	}{
		{
			name: "in general",
			setup: func(controller *gomock.Controller) *configFake.MockDash {
				namespaces := []string{"ns-1"}

				namespaceClient := clusterFake.NewMockNamespaceInterface(controller)
				namespaceClient.EXPECT().Names(ctx).Return(namespaces, nil)
				namespaceClient.EXPECT().ProvidedNamespaces(ctx).Return([]string{})

				clusterClient := clusterFake.NewMockClientInterface(controller)
				clusterClient.EXPECT().NamespaceClient().Return(namespaceClient, nil)

				dashConfig := configFake.NewMockDash(controller)
				dashConfig.EXPECT().ClusterClient().Return(clusterClient)

				return dashConfig
			},
			expected: []string{"ns-1"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			require.NotNil(t, test.setup)
			dashConfig := test.setup(controller)

			ctx := context.Background()
			got, err := api.NamespacesGenerator(ctx, dashConfig)
			require.NoError(t, err)

			require.Equal(t, test.expected, got)
		})
	}
}
