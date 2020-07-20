/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/vmware-tanzu/octant/internal/api"
	"github.com/vmware-tanzu/octant/internal/api/fake"
	configFake "github.com/vmware-tanzu/octant/internal/config/fake"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/octant"
	octantFake "github.com/vmware-tanzu/octant/internal/octant/fake"
)

func TestHelperManager_GenerateContent(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	state := octantFake.NewMockState(controller)
	octantClient := fake.NewMockOctantClient(controller)

	ev := octant.Event{
		Type: "event.octant.dev/buildInfo",
	}
	octantClient.EXPECT().Send(ev)

	logger := log.NopLogger()

	dashConfig := configFake.NewMockDash(controller)
	dashConfig.EXPECT().Logger().Return(logger).AnyTimes()

	poller := api.NewSingleRunPoller()
	generatorFunc := func(ctx context.Context, state octant.State) (octant.Event, error) {
		return ev, nil
	}

	manager := api.NewHelperStateManager(dashConfig,
		api.WithHelperGenerator(generatorFunc),
		api.WithHelperGeneratorPoll(poller))
	ctx := context.Background()
	manager.Start(ctx, state, octantClient)
}
