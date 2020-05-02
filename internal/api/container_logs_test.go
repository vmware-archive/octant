package api

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	configFake "github.com/vmware-tanzu/octant/internal/config/fake"
	"github.com/vmware-tanzu/octant/internal/modules/overview/container"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/store"
)

func TestContainerLogs_NewLogEntry(t *testing.T) {
	le := newLogEntry("line", "container-name")

	assert.Equal(t, "container-name", le.Container)
	assert.Equal(t, "line", le.Message)
	assert.Nil(t, le.Timestamp)

	le = newLogEntry("1985-04-12T23:20:50.52Z line", "container-name")
	assert.Equal(t, "container-name", le.Container)
	assert.Equal(t, "line", le.Message)

	ts, err := time.Parse(time.RFC3339, "1985-04-12T23:20:50.52Z")
	assert.NoError(t, err)
	assert.Equal(t, ts.String(), le.Timestamp.String())
}

func TestContainerLogs_SendLogEventsStops(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashConfig := configFake.NewMockDash(controller)

	key := store.Key{
		Namespace: "test-ns",
		Name:      "test-pod",
	}

	eventType := octant.NewLoggingEventType(key.Namespace, key.Name)
	logCh := make(chan container.LogEntry)

	s := NewPodLogsStateManager(dashConfig)
	s.Start(context.Background(), nil, nil)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		s.streamEventsToClient(s.ctx, eventType, logCh)
		wg.Done()
	}()

	close(logCh)

	wg.Wait()
	_, ok := <-logCh
	assert.False(t, ok)
}

func TestContainerLogs_SendLogEventsClientSend(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashConfig := configFake.NewMockDash(controller)
	client := newOctantClient()

	key := store.Key{
		Namespace: "test-ns",
		Name:      "test-pod",
	}

	eventType := octant.NewLoggingEventType(key.Namespace, key.Name)
	logCh := make(chan container.LogEntry, 1)

	s := NewPodLogsStateManager(dashConfig)
	ctx, cancel := context.WithCancel(context.Background())
	s.Start(ctx, nil, client)

	go func() {
		s.streamEventsToClient(ctx, eventType, logCh)
	}()

	le := container.NewLogEntry("container-a", "testing log line")
	logCh <- le

	<-client.ch
	cancel()
	close(logCh)

	if assert.NotNil(t, client.sendCalledWith) {
		assert.Equal(t, eventType, client.sendCalledWith.Type)

		clientLe, ok := client.sendCalledWith.Data.(logEntry)
		assert.True(t, ok)
		assert.Equal(t, "container-a", clientLe.Container)
		assert.Equal(t, "testing log line", clientLe.Message)
		assert.Nil(t, clientLe.Timestamp)
	}
}

type octantClient struct {
	sendCalledWith octant.Event
	ch             chan bool
}

func newOctantClient() *octantClient {
	return &octantClient{
		ch: make(chan bool, 1),
	}
}

func (oc *octantClient) Send(event octant.Event) {
	oc.sendCalledWith = event
	oc.ch <- true
}
func (oc *octantClient) ID() string { return "" }
