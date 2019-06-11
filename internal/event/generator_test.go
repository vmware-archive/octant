package event

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/heptio/developer-dash/internal/clustereye"
	clustereyeFake "github.com/heptio/developer-dash/internal/clustereye/fake"
	eventFake "github.com/heptio/developer-dash/internal/event/fake"
)

func TestStream(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	controller := gomock.NewController(t)
	defer controller.Finish()

	event := clustereye.Event{
		Type: clustereye.EventType("test"),
		Data: []byte("data"),
	}
	generator := clustereyeFake.NewMockGenerator(controller)
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
		Do(func(ctx context.Context, ch <-chan clustereye.Event) {
			<-ch
			done <- true

		})

	go func() {
		err := Stream(ctx, streamer, []clustereye.Generator{generator}, "/request-path", "/content-path")
		require.NoError(t, err)
	}()

	<-done
	cancel()
}
