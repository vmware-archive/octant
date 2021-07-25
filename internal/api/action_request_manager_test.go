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

	configFake "github.com/vmware-tanzu/octant/internal/config/fake"

	"github.com/vmware-tanzu/octant/internal/api"
	ocontext "github.com/vmware-tanzu/octant/internal/context"
	"github.com/vmware-tanzu/octant/internal/octant"
	octantFake "github.com/vmware-tanzu/octant/internal/octant/fake"
	"github.com/vmware-tanzu/octant/pkg/action"
)

func TestActionRequestManager_Handlers(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashConfig := configFake.NewMockDash(controller)
	manager := api.NewActionRequestManager(dashConfig)
	AssertHandlers(t, manager, []string{api.RequestPerformAction, api.TerminatingThreshold})
}

func TestActionRequestManager_PerformAction(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashConfig := configFake.NewMockDash(controller)
	dashConfig.EXPECT().CurrentContext().Return("foo-context")
	state := octantFake.NewMockState(controller)
	state.EXPECT().GetFilters().Return([]octant.Filter{{Key: "foo", Value: "bar"}})
	state.EXPECT().GetNamespace().Return("foo-namespace")
	state.EXPECT().GetClientID().Return("foo-client")

	manager := api.NewActionRequestManager(dashConfig)

	payload := action.CreatePayload(api.RequestPerformAction, map[string]interface{}{
		"foo": "bar",
	})

	state.EXPECT().
		Dispatch(gomock.Any(), api.RequestPerformAction, payload).
		Do(func(ctx context.Context, _ string, _ action.Payload) {
			clientState := ocontext.ClientStateFrom(ctx)
			require.Equal(t, "foo-namespace", clientState.Namespace)
			require.Equal(t, "foo", clientState.Filters[0].Key)
			require.Equal(t, "bar", clientState.Filters[0].Value)
			require.Equal(t, "foo-client", clientState.ClientID)
			require.Equal(t, "foo-context", clientState.ContextName)
		}).
		Return(nil)

	require.NoError(t, manager.PerformAction(state, payload))
}

func TestActionRequestManager_SetTerminateThreshold(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashConfig := configFake.NewMockDash(controller)
	dashConfig.EXPECT().SetTerminateThreshold(int64(5))
	state := octantFake.NewMockState(controller)

	manager := api.NewActionRequestManager(dashConfig)

	payload := action.CreatePayload(api.TerminatingThreshold, map[string]interface{}{
		"threshold": "5",
	})

	require.NoError(t, manager.SetTerminateThreshold(state, payload))
}
