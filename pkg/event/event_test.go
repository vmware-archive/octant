package event_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vmware-tanzu/octant/pkg/event"
)

func TestFindEvent(t *testing.T) {
	assert := assert.New(t)
	ev := event.CreateEvent(event.EventTypeAlert, make(map[string]interface{}))
	events := []event.Event{ev}

	res, err := event.FindEvent(events, event.EventTypeAlert)
	assert.Equal(ev, res)
	assert.Nil(err)

	res, err = event.FindEvent(events, event.EventTypeAppLogs)
	assert.NotNil(err)
}
