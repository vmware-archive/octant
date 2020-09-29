/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package log

import (
	"container/ring"
	"sort"
	"sync"

	"k8s.io/apimachinery/pkg/util/rand"

	"github.com/vmware-tanzu/octant/pkg/event"
)

// MessageListenerFactory is a factory which generates message listeners.
type MessageListenerFactory interface {
	// Listen creates a listener.
	Listen() (<-chan Message, ListenCancelFunc)
}

// EventStreamer is an interface for streaming.
type EventStreamer interface {
	// Stream streams events.
	Stream(ready <-chan struct{}) (<-chan event.Event, func())
	// Close closes the streamer. No more events will be generated
	// one the streamer is closed.
	Close()
}

// Streamer streams events using a ring buffer.
type Streamer struct {
	messageCh     <-chan Message
	messageCancel ListenCancelFunc
	r             *ring.Ring
	eventCh       chan event.Event

	mu        sync.RWMutex
	listeners map[string]chan event.Event
}

var _ EventStreamer = &Streamer{}

// NewStreamer creates an instance of Streamer.
func NewStreamer(messageListenerFactory MessageListenerFactory) *Streamer {
	messageCh, messageCancel := messageListenerFactory.Listen()

	lm := &Streamer{
		messageCh:     messageCh,
		r:             ring.New(1500),
		eventCh:       make(chan event.Event),
		listeners:     make(map[string]chan event.Event),
		messageCancel: messageCancel,
	}

	go func() {
		for m := range messageCh {
			lm.handle(m)
		}
	}()

	return lm
}

// Close closes the streamer.
func (lm *Streamer) Close() {
	lm.messageCancel()
}

func (lm *Streamer) handle(m Message) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	lm.r.Value = m
	lm.r = lm.r.Next()

	e := lm.createEvent()

	for _, ch := range lm.listeners {
		ch <- e
	}

}

func (lm *Streamer) createEvent() event.Event {
	var messages []Message
	lm.r.Do(func(i interface{}) {
		if i != nil {
			m := i.(Message)
			messages = append(messages, m)
		}
	})

	sort.Slice(messages, func(i, j int) bool {
		return messages[j].Date < messages[i].Date
	})

	e := event.Event{
		Type: event.EventTypeAppLogs,
		Data: messages,
	}

	return e
}

// Stream creates a channel which streams messages and a cancel function.
func (lm *Streamer) Stream(readyCh <-chan struct{}) (<-chan event.Event, func()) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	id := rand.String(6)
	ch := make(chan event.Event, 10)

	go func() {
		lm.mu.RLock()
		defer lm.mu.RUnlock()

		if _, ok := lm.listeners[id]; !ok {
			return
		}

		<-readyCh
		ch <- lm.createEvent()
	}()

	lm.listeners[id] = ch

	return ch, func() {
		lm.mu.Lock()
		defer lm.mu.Unlock()

		close(ch)
		delete(lm.listeners, id)
	}
}
