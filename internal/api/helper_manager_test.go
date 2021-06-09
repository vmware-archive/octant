/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api_test

import (
	"context"
	"testing"

	"github.com/vmware-tanzu/octant/pkg/event"

	"github.com/golang/mock/gomock"

	"github.com/vmware-tanzu/octant/internal/api"
	configFake "github.com/vmware-tanzu/octant/internal/config/fake"
	"github.com/vmware-tanzu/octant/internal/log"
	octantFake "github.com/vmware-tanzu/octant/internal/octant/fake"
	"github.com/vmware-tanzu/octant/pkg/api/fake"
)

func TestHelperManager_GenerateContent(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	state := octantFake.NewMockState(controller)
	octantClient := fake.NewMockOctantClient(controller)

	bev := event.Event{
		Type: event.EventTypeBuildInfo,
	}
	octantClient.EXPECT().Send(bev)

	kEv := event.Event{
		Type: event.EventTypeKubeConfigPath,
	}
	octantClient.EXPECT().Send(kEv)

	logger := log.NopLogger()

	dashConfig := configFake.NewMockDash(controller)
	dashConfig.EXPECT().Logger().Return(logger).AnyTimes()

	poller := api.NewSingleRunPoller()
	generatorFunc := func(ctx context.Context) ([]event.Event, error) {
		return []event.Event{bev, kEv}, nil
	}

	manager := api.NewHelperStateManager(dashConfig,
		api.WithHelperGenerator(generatorFunc),
		api.WithHelperGeneratorPoll(poller))
	ctx := context.Background()
	manager.Start(ctx, state, octantClient)
}
