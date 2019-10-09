/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/vmware/octant/internal/api"
	"github.com/vmware/octant/internal/api/fake"
	configFake "github.com/vmware/octant/internal/config/fake"
	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/internal/octant"
	octantFake "github.com/vmware/octant/internal/octant/fake"
)

func TestContextManager_Handlers(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashConfig := configFake.NewMockDash(controller)

	manager := api.NewContextManager(dashConfig)
	AssertHandlers(t, manager, []string{api.RequestSetContext})
}

func TestContext_GenerateContexts(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	state := octantFake.NewMockState(controller)
	octantClient := fake.NewMockOctantClient(controller)

	ev := octant.Event{
		Type: "eventType",
	}
	octantClient.EXPECT().Send(ev)

	logger := log.NopLogger()

	dashConfig := configFake.NewMockDash(controller)
	dashConfig.EXPECT().Logger().Return(logger).AnyTimes()

	poller := api.NewSingleRunPoller()
	generatorFunc := func(ctx context.Context, state octant.State) (octant.Event, error) {
		return ev, nil
	}
	manager := api.NewContextManager(dashConfig,
		api.WithContextGenerator(generatorFunc),
		api.WithContextGeneratorPoll(poller))

	ctx := context.Background()
	manager.Start(ctx, state, octantClient)
}
