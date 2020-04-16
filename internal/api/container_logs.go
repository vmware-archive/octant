/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"context"
	"fmt"
	"github.com/vmware-tanzu/octant/internal/gvk"
	"github.com/vmware-tanzu/octant/pkg/store"
	"strings"
	"sync"
	"time"

	"github.com/vmware-tanzu/octant/internal/config"
	"github.com/vmware-tanzu/octant/internal/modules/overview/container"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/action"
)

type logEntry struct {
	Timestamp *time.Time `json:"timestamp,omitempty"`
	Container string     `json:"container,omitempty"`
	Message   string     `json:"message,omitempty"`
}

type logResponse struct {
	Entries []logEntry `json:"entries,omitempty"`
}

const (
	RequestPodLogsSubscribe   = "action.octant.dev/podLogs/subscribe"
	RequestPodLogsUnsubscribe = "action.octant.dev/podLogs/unsubscribe"
)

type podLogsStateManager struct {
	client OctantClient
	config config.Dash
	ctx    context.Context

	podLogSubscriptions sync.Map
}

var _ StateManager = (*podLogsStateManager)(nil)

// NewPodLogsStateManager returns a terminal state manager.
func NewPodLogsStateManager(dashConfig config.Dash) *podLogsStateManager {
	return &podLogsStateManager{
		config:              dashConfig,
		podLogSubscriptions: sync.Map{},
	}
}

// Handlers returns a slice of handlers.
func (s *podLogsStateManager) Handlers() []octant.ClientRequestHandler {
	return []octant.ClientRequestHandler{
		{
			RequestType: RequestPodLogsSubscribe,
			Handler:     s.StreamPodLogsSubscribe,
		},
		{
			RequestType: RequestPodLogsUnsubscribe,
			Handler:     s.StreamPodLogsUnsubscribe,
		},
	}
}

func (s *podLogsStateManager) StreamPodLogsSubscribe(state octant.State, payload action.Payload) error {
	namespace, err := payload.String("namespace")
	if err != nil {
		return fmt.Errorf("getting namespace from payload: %w", err)
	}
	podName, err := payload.String("podName")
	if err != nil {
		return fmt.Errorf("getting podName from payload: %w", err)
	}

	containerName, err := payload.String("containerName")
	if err != nil {
		return fmt.Errorf("getting containerName from payload: %w", err)
	}

	eventType := octant.NewLoggingEventType(namespace, podName)
	val, ok := s.podLogSubscriptions.Load(eventType)
	if ok {
		cancelFn, ok := val.(context.CancelFunc)
		if !ok {
			return fmt.Errorf("bad cancelFn conversion for %s", eventType)
		}
		cancelFn()
	}

	key := store.KeyFromGroupVersionKind(gvk.Pod)
	key.Name = podName
	key.Namespace = namespace

	logStreamer, err := container.NewLogStreamer(s.ctx, s.config, key, containerName)
	if err != nil {
		return fmt.Errorf("creating log streamer: %w", err)
	}

	cancelFn := s.startStream(key, logStreamer)
	s.podLogSubscriptions.Store(eventType, cancelFn)

	return nil
}

func (s *podLogsStateManager) StreamPodLogsUnsubscribe(state octant.State, payload action.Payload) error {
	namespace, err := payload.String("namespace")
	if err != nil {
		return fmt.Errorf("getting namespace from payload: %w", err)
	}

	podName, err := payload.String("podName")
	if err != nil {
		return fmt.Errorf("getting podName from payload: %w", err)
	}

	eventType := octant.NewLoggingEventType(namespace, podName)
	val, ok := s.podLogSubscriptions.Load(eventType)
	if ok {
		cancelFn, ok := val.(context.CancelFunc)
		if !ok {
			return fmt.Errorf("bad cancelFn conversion for %s", eventType)
		}
		s.podLogSubscriptions.Delete(eventType)
		cancelFn()
	}
	return nil
}

func (s *podLogsStateManager) Start(ctx context.Context, state octant.State, client OctantClient) {
	s.client = client
	s.ctx = ctx
}

func (s *podLogsStateManager) sendLogEvents(ctx context.Context, logEventType octant.EventType, logCh <-chan container.LogEntry) {
	ctx, cancelFn := context.WithCancel(s.ctx)
	for {
		select {
		case <-ctx.Done():
			cancelFn()
			return
		case entry, ok := <-logCh:
			if ok {
				le := newLogEntry(entry.Line(), entry.Container())
				logEvent := octant.Event{
					Type: logEventType,
					Data: le,
					Err:  nil,
				}
				s.client.Send(logEvent)
			} else {
				cancelFn()
				return
			}
		}
	}
}

func (s *podLogsStateManager) startStream(key store.Key, logStreamer container.LogStreamer) context.CancelFunc {
	ctx, cancelFn := context.WithCancel(s.ctx)

	eventType := octant.NewLoggingEventType(key.Namespace, key.Name)
	logCh := make(chan container.LogEntry)
	go s.sendLogEvents(ctx, eventType, logCh)

	logStreamer.Stream(ctx, logCh)

	return cancelFn
}

func newLogEntry(message, container string) logEntry {
	le := logEntry{
		Container: container,
		Message:   message,
		Timestamp: nil,
	}
	if message, ts, err := formatTimestamp(le.Message); err == nil {
		le.Message = message
		le.Timestamp = &ts
	}
	return le
}

func formatTimestamp(line string) (string, time.Time, error) {
	parts := strings.SplitN(line, " ", 2)
	ts, err := time.Parse(time.RFC3339, parts[0])
	if err != nil {
		return "", ts, err
	}
	return parts[1], ts, nil
}
