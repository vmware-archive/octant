/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api_test

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware/octant/internal/api"
	"github.com/vmware/octant/internal/api/fake"
	configFake "github.com/vmware/octant/internal/config/fake"
	"github.com/vmware/octant/internal/log"
	moduleFake "github.com/vmware/octant/internal/module/fake"
	"github.com/vmware/octant/internal/octant"
)

func TestWebsocketState_Start(t *testing.T) {
	mocks := newWebsocketStateMocks(t, "default")
	defer mocks.finish()

	started := make(chan bool, 1)
	mocks.stateManager.EXPECT().Start(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, state octant.State, wsClient api.OctantClient) {
			started <- true
		})
	s := mocks.factory()

	ctx, cancel := context.WithCancel(context.Background())
	s.Start(ctx)

	<-started

	cancel()
}

func TestWebsocketState_SetContentPath(t *testing.T) {
	tests := []struct {
		name        string
		contentPath string
		namespace   string
		setup       func(mocks *websocketStateMocks)
		verify      func(t *testing.T, s *api.WebsocketState)
	}{
		{
			name:        "set content path without namespace change",
			contentPath: "overview/namespace/default",
			namespace:   "default",
			setup: func(mocks *websocketStateMocks) {
				contentPath := "overview/namespace/default"
				mocks.moduleManager.EXPECT().
					ModuleForContentPath(contentPath).
					Return(mocks.module, true)
			},
			verify: func(t *testing.T, s *api.WebsocketState) {
				contentPath := "overview/namespace/default"
				assert.Equal(t, "default", s.GetNamespace())
				assert.Equal(t, contentPath, s.GetContentPath())
			},
		},
		{
			name:        "non namespaced content path",
			contentPath: "overview/foo",
			namespace:   "default",
			setup: func(mocks *websocketStateMocks) {
				contentPath := "overview/foo"
				mocks.moduleManager.EXPECT().
					ModuleForContentPath(contentPath).
					Return(mocks.module, true)
			},
			verify: func(t *testing.T, s *api.WebsocketState) {
				contentPath := "overview/foo"
				assert.Equal(t, "default", s.GetNamespace())
				assert.Equal(t, contentPath, s.GetContentPath())
			},
		},
		{
			name:        "set content path with namespace change",
			contentPath: "overview/namespace/kube-system",
			namespace:   "default",
			setup: func(mocks *websocketStateMocks) {
				contentPath := "overview/namespace/kube-system"
				mocks.moduleManager.EXPECT().
					ModuleForContentPath(contentPath).
					Return(mocks.module, true)
			},
			verify: func(t *testing.T, s *api.WebsocketState) {
				contentPath := "overview/namespace/kube-system"
				assert.Equal(t, "kube-system", s.GetNamespace())
				assert.Equal(t, contentPath, s.GetContentPath())
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mocks := newWebsocketStateMocks(t, test.namespace)
			defer mocks.finish()

			require.NotNil(t, test.setup)
			test.setup(mocks)
			s := mocks.factory()
			s.SetContentPath(test.contentPath)

			require.NotNil(t, test.verify)
			test.verify(t, s)
		})
	}
}

func TestWebsocketState_OnContentPathUpdate(t *testing.T) {
	mocks := newWebsocketStateMocks(t, "default")
	defer mocks.finish()

	contentPath := "overview/foo"
	mocks.wsClient.EXPECT().Send(gomock.Any()).AnyTimes()
	mocks.moduleManager.EXPECT().
		ModuleForContentPath(contentPath).
		Return(mocks.module, true).AnyTimes()

	s := mocks.factory()

	called := make(chan bool, 1)
	cancelUpdate := s.OnContentPathUpdate(func(s string) {
		called <- true
	})

	ctx, cancel1 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel1()

	s.SetContentPath(contentPath)

	select {
	case <-ctx.Done():
		t.Error("should have been called")
	case <-called:
	}

	cancelUpdate()

	ctx, cancel2 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel2()

	s.SetContentPath(contentPath)

	select {
	case <-ctx.Done():
	case <-called:
		t.Error("should not have been called")
	}
}

func TestWebsocketState_SetNamespace(t *testing.T) {
	tests := []struct {
		name             string
		initialNamespace string
		newNamespace     string
		setup            func(mocks *websocketStateMocks)
	}{
		{
			name:             "set to existing namespace",
			initialNamespace: "default",
			newNamespace:     "default",
			setup: func(mocks *websocketStateMocks) {
			},
		},
		{
			name:             "set to new namespace",
			initialNamespace: "default",
			newNamespace:     "other",
			setup: func(mocks *websocketStateMocks) {
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mocks := newWebsocketStateMocks(t, test.initialNamespace)
			defer mocks.finish()

			require.NotNil(t, test.setup)
			test.setup(mocks)

			s := mocks.factory()
			s.SetNamespace(test.newNamespace)

		})
	}
}

func TestWebsocketState_OnNamespaceUpdate(t *testing.T) {
	mocks := newWebsocketStateMocks(t, "default")
	defer mocks.finish()

	mocks.wsClient.EXPECT().Send(gomock.Any()).AnyTimes()

	s := mocks.factory()

	called := make(chan bool, 1)
	cancelUpdate := s.OnNamespaceUpdate(func(s string) {
		called <- true
	})

	ctx, cancel1 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel1()

	s.SetNamespace("new-namespace")

	select {
	case <-ctx.Done():
		t.Error("should have been called")
	case <-called:
	}

	cancelUpdate()

	ctx, cancel2 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel2()

	s.SetNamespace("new-namespace")

	select {
	case <-ctx.Done():
	case <-called:
		t.Error("should not have been called")
	}
}

func TestWebsocketState_AddFilter(t *testing.T) {
	mocks := newWebsocketStateMocks(t, "default")
	defer mocks.finish()
	s := mocks.factory()

	s.AddFilter(octant.Filter{
		Key:   "key",
		Value: "value",
	})

	got := s.GetFilters()

	expected := []octant.Filter{
		{Key: "key", Value: "value"},
	}
	assert.Equal(t, expected, got)
}

type websocketStateMocks struct {
	controller       *gomock.Controller
	module           *moduleFake.MockModule
	moduleManager    *moduleFake.MockManagerInterface
	dashConfig       *configFake.MockDash
	wsClient         *fake.MockOctantClient
	stateManager     *fake.MockStateManager
	actionDispatcher *fake.MockActionDispatcher
}

func newWebsocketStateMocks(t *testing.T, namespace string) *websocketStateMocks {
	controller := gomock.NewController(t)
	m := moduleFake.NewMockModule(controller)
	m.EXPECT().Name().Return("overview").AnyTimes()
	moduleManager := moduleFake.NewMockManagerInterface(controller)
	dashConfig := configFake.NewMockDash(controller)
	dashConfig.EXPECT().DefaultNamespace().Return(namespace)
	dashConfig.EXPECT().ModuleManager().Return(moduleManager).AnyTimes()
	dashConfig.EXPECT().Logger().Return(log.NopLogger()).AnyTimes()
	octantClient := fake.NewMockOctantClient(controller)
	stateManager := fake.NewMockStateManager(controller)
	actionDispatcher := fake.NewMockActionDispatcher(controller)

	return &websocketStateMocks{
		controller:       controller,
		module:           m,
		moduleManager:    moduleManager,
		dashConfig:       dashConfig,
		wsClient:         octantClient,
		stateManager:     stateManager,
		actionDispatcher: actionDispatcher,
	}
}

func (w *websocketStateMocks) options() []api.WebsocketStateOption {
	return []api.WebsocketStateOption{
		api.WebsocketStateManagers([]api.StateManager{w.stateManager}),
	}
}

func (w *websocketStateMocks) finish() {
	w.controller.Finish()
}

func (w *websocketStateMocks) factory() *api.WebsocketState {
	return api.NewWebsocketState(w.dashConfig, w.actionDispatcher, w.wsClient, w.options()...)

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
