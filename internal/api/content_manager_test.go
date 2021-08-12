/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api_test

import (
	"context"
	"sort"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/api"
	configFake "github.com/vmware-tanzu/octant/internal/config/fake"
	ocontext "github.com/vmware-tanzu/octant/internal/context"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/module"
	moduleFake "github.com/vmware-tanzu/octant/internal/module/fake"
	"github.com/vmware-tanzu/octant/internal/octant"
	octantFake "github.com/vmware-tanzu/octant/internal/octant/fake"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/api/fake"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func TestContentManager_Handlers(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashConfig := configFake.NewMockDash(controller)
	moduleManager := moduleFake.NewMockManagerInterface(controller)

	logger := log.NopLogger()

	manager := api.NewContentManager(moduleManager, dashConfig, logger)
	AssertHandlers(t, manager, []string{
		api.RequestSetContentPath,
		action.RequestSetNamespace,
		api.CheckLoading,
	})
}

func TestContentManager_GenerateContent(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	params := map[string][]string{}
	filters := []octant.Filter{{Key: "foo", Value: "bar"}}

	dashConfig := configFake.NewMockDash(controller)
	moduleManager := moduleFake.NewMockManagerInterface(controller)
	fakeModule := moduleFake.NewMockModule(controller)
	state := octantFake.NewMockState(controller)

	dashConfig.EXPECT().CurrentContext().Return("foo-context")
	state.EXPECT().GetClientID().Return("foo-client")
	state.EXPECT().GetFilters().Return(filters).AnyTimes()
	state.EXPECT().GetNamespace().Return("foo-namespace").AnyTimes()
	state.EXPECT().GetQueryParams().Return(params)
	state.EXPECT().GetContentPath().Return(".").AnyTimes()
	state.EXPECT().OnContentPathUpdate(gomock.Any()).DoAndReturn(func(fn octant.ContentPathUpdateFunc) octant.UpdateCancelFunc {
		fn("foo")
		return func() {}
	})
	octantClient := fake.NewMockOctantClient(controller)

	stopCh := make(chan struct{}, 1)

	contentResponse := component.ContentResponse{}
	contentEvent := api.CreateContentEvent(contentResponse, "foo-namespace", ".", params)
	octantClient.EXPECT().Send(contentEvent).AnyTimes()
	octantClient.EXPECT().StopCh().Return(stopCh).AnyTimes()

	moduleManager.EXPECT().ModuleForContentPath(gomock.Any()).Return(fakeModule, true).Times(2)
	moduleManager.EXPECT().Navigation(gomock.Any(), "foo-namespace", "foo-module").Return([]navigation.Navigation{}, nil)
	fakeModule.EXPECT().Name().Return("foo-module").AnyTimes()
	fakeModule.EXPECT().Content(gomock.Any(), ".", gomock.Any()).
		Do(func(ctx context.Context, _ string, _ module.ContentOptions) {
			clientState := ocontext.ClientStateFrom(ctx)
			require.Equal(t, "foo-namespace", clientState.Namespace)
			require.Equal(t, "foo", clientState.Filters[0].Key)
			require.Equal(t, "bar", clientState.Filters[0].Value)
			require.Equal(t, "foo-client", clientState.ClientID)
			require.Equal(t, "foo-context", clientState.ContextName)
		}).
		Return(contentResponse, nil)

	logger := log.NopLogger()

	poller := api.NewSingleRunPoller()

	manager := api.NewContentManager(moduleManager, dashConfig, logger,
		api.WithContentGeneratorPoller(poller))

	ctx := context.Background()
	manager.Start(ctx, state, octantClient)
}

func TestContentManager_GenerateContent_ClusterOverviewNamespace(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	params := map[string][]string{}
	filters := []octant.Filter{{Key: "foo", Value: "bar"}}

	dashConfig := configFake.NewMockDash(controller)
	moduleManager := moduleFake.NewMockManagerInterface(controller)
	fakeModule := moduleFake.NewMockModule(controller)
	state := octantFake.NewMockState(controller)

	dashConfig.EXPECT().CurrentContext().Return("foo-context")
	state.EXPECT().GetClientID().Return("foo-client")
	state.EXPECT().GetFilters().Return(filters).AnyTimes()
	state.EXPECT().GetNamespace().Return("foo-namespace").AnyTimes()
	state.EXPECT().GetQueryParams().Return(params).AnyTimes()
	state.EXPECT().GetContentPath().Return(".").AnyTimes()
	state.EXPECT().OnContentPathUpdate(gomock.Any()).DoAndReturn(func(fn octant.ContentPathUpdateFunc) octant.UpdateCancelFunc {
		fn("foo")
		return func() {}
	})
	octantClient := fake.NewMockOctantClient(controller)

	stopCh := make(chan struct{}, 1)

	contentResponse := component.ContentResponse{}
	contentEvent := api.CreateContentEvent(contentResponse, "", ".", params)
	octantClient.EXPECT().Send(contentEvent).AnyTimes()
	octantClient.EXPECT().StopCh().Return(stopCh).AnyTimes()

	moduleManager.EXPECT().ModuleForContentPath(gomock.Any()).Return(fakeModule, true).Times(2)
	moduleManager.EXPECT().Navigation(gomock.Any(), "foo-namespace", "cluster-overview").Return([]navigation.Navigation{}, nil)
	fakeModule.EXPECT().Name().Return("cluster-overview").AnyTimes()
	fakeModule.EXPECT().Content(gomock.Any(), ".", gomock.Any()).Return(contentResponse, nil)

	logger := log.NopLogger()

	poller := api.NewSingleRunPoller()

	manager := api.NewContentManager(moduleManager, dashConfig, logger,
		api.WithContentGeneratorPoller(poller))

	ctx := context.Background()
	manager.Start(ctx, state, octantClient)
}

func TestContentManager_SetContentPath(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	m := moduleFake.NewMockModule(controller)
	m.EXPECT().Name().Return("name").AnyTimes()

	moduleManager := moduleFake.NewMockManagerInterface(controller)
	dashConfig := configFake.NewMockDash(controller)

	state := octantFake.NewMockState(controller)
	state.EXPECT().SetContentPath("/path")

	logger := log.NopLogger()

	manager := api.NewContentManager(moduleManager, dashConfig, logger,
		api.WithContentGeneratorPoller(api.NewSingleRunPoller()))

	payload := action.Payload{
		"contentPath": "/path",
	}

	require.NoError(t, manager.SetContentPath(state, payload))
}

func TestContentManager_SetNamespace(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	m := moduleFake.NewMockModule(controller)
	m.EXPECT().Name().Return("name").AnyTimes()

	moduleManager := moduleFake.NewMockManagerInterface(controller)
	dashConfig := configFake.NewMockDash(controller)

	state := octantFake.NewMockState(controller)
	state.EXPECT().SetNamespace("kube-system")
	logger := log.NopLogger()

	manager := api.NewContentManager(moduleManager, dashConfig, logger,
		api.WithContentGeneratorPoller(api.NewSingleRunPoller()))

	payload := action.Payload{
		"namespace": "kube-system",
	}
	state.EXPECT().Dispatch(nil, action.RequestSetNamespace, payload)
	require.NoError(t, manager.SetNamespace(state, payload))
}

func TestContentManager_SetQueryParams(t *testing.T) {
	tests := []struct {
		name    string
		payload action.Payload
		setup   func(state *octantFake.MockState)
	}{
		{
			name: "single filter",
			payload: action.Payload{
				"params": map[string]interface{}{
					"filters": "foo:bar",
				},
			},
			setup: func(state *octantFake.MockState) {
				state.EXPECT().SetFilters([]octant.Filter{
					{Key: "foo", Value: "bar"},
				})
			},
		},
		{
			name: "multiple filters",
			payload: action.Payload{
				"params": map[string]interface{}{
					"filters": []interface{}{
						"foo:bar",
						"baz:qux",
					},
				},
			},
			setup: func(state *octantFake.MockState) {
				state.EXPECT().SetFilters([]octant.Filter{
					{Key: "foo", Value: "bar"},
					{Key: "baz", Value: "qux"},
				})
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			m := moduleFake.NewMockModule(controller)
			m.EXPECT().Name().Return("name").AnyTimes()

			moduleManager := moduleFake.NewMockManagerInterface(controller)
			dashConfig := configFake.NewMockDash(controller)

			state := octantFake.NewMockState(controller)
			require.NotNil(t, test.setup)
			test.setup(state)

			logger := log.NopLogger()

			manager := api.NewContentManager(moduleManager, dashConfig, logger,
				api.WithContentGeneratorPoller(api.NewSingleRunPoller()))
			require.NoError(t, manager.SetQueryParams(state, test.payload))
		})
	}
}

func AssertHandlers(t *testing.T, manager api.StateManager, expected []string) {
	handlers := manager.Handlers()
	var got []string
	for _, h := range handlers {
		got = append(got, h.RequestType)
	}
	sort.Strings(got)
	sort.Strings(expected)

	assert.Equal(t, expected, got)
}
