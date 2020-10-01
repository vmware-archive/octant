/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package log

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/pkg/event"
)

func TestStreamer_Close(t *testing.T) {
	f := newStubMessageListenerFactory()
	s := NewStreamer(f)

	s.Close()
	require.True(t, f.isCanceled)
}

func TestStreamer_Stream(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "in general",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := Message{
				ID: "1",
			}

			f := newStubMessageListenerFactory()
			s := NewStreamer(f)

			readyCh := make(chan struct{}, 1)
			close(readyCh)

			ch, cancel := s.Stream(readyCh)

			done := make(chan struct{}, 1)
			go func() {
				e := <-ch
				want := event.Event{
					Type: event.EventTypeAppLogs,
					Data: []Message{message},
					Err:  nil,
				}
				require.Equal(t, want, e)
				close(done)
			}()

			f.ch <- message

			<-done
			cancel()

		})
	}
}

type stubMessageListenerFactory struct {
	ch         chan Message
	cancel     func()
	isCanceled bool
}

func newStubMessageListenerFactory() *stubMessageListenerFactory {
	f := &stubMessageListenerFactory{
		ch: make(chan Message, 1),
	}

	f.cancel = func() {
		if f.isCanceled {
			panic("already canceled")
		}
		f.isCanceled = true
	}

	return f
}

func (f *stubMessageListenerFactory) Listen() (<-chan Message, ListenCancelFunc) {
	return f.ch, f.cancel
}

var _ MessageListenerFactory = &stubMessageListenerFactory{}
