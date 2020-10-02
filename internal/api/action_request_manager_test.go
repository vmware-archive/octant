/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/api"
	octantFake "github.com/vmware-tanzu/octant/internal/octant/fake"
	"github.com/vmware-tanzu/octant/pkg/action"
)

func TestActionRequestManager_Handlers(t *testing.T) {
	manager := api.NewActionRequestManager()
	AssertHandlers(t, manager, []string{api.RequestPerformAction})
}

func TestActionRequestManager_PerformAction(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	state := octantFake.NewMockState(controller)

	manager := api.NewActionRequestManager()

	payload := action.CreatePayload(api.RequestPerformAction, map[string]interface{}{
		"foo": "bar",
	})

	state.EXPECT().
		Dispatch(gomock.Any(), api.RequestPerformAction, payload).
		Return(nil)
	state.EXPECT().GetClientID()

	require.NoError(t, manager.PerformAction(state, payload))
}
