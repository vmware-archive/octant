/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package event

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	eventFake "github.com/vmware/octant/internal/event/fake"
	"github.com/vmware/octant/internal/octant"
	octantFake "github.com/vmware/octant/internal/octant/fake"
)

func TestStream(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	controller := gomock.NewController(t)
	defer controller.Finish()

	event := octant.Event{
		Type: octant.EventType("test"),
		Data: []byte("data"),
	}
	generator := octantFake.NewMockGenerator(controller)
	generator.EXPECT().
		Event(gomock.Any()).Return(event, nil)
	generator.EXPECT().
		Name().Return("test").AnyTimes()
	generator.EXPECT().
		ScheduleDelay().Return(DefaultScheduleDelay).AnyTimes()

	done := make(chan bool, 1)

	streamer := eventFake.NewMockStreamer(controller)
	streamer.EXPECT().
		Stream(gomock.Any(), gomock.Any()).
		Do(func(ctx context.Context, ch <-chan octant.Event) {
			<-ch
			done <- true

		})

	go func() {
		forceCh := make(chan bool, 1)
		defer close(forceCh)
		err := Stream(ctx, streamer, forceCh, []octant.Generator{generator}, "/request-path", "/content-path")
		require.NoError(t, err)
	}()

	<-done
	cancel()
}
