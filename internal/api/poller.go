/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api

import (
	"context"
	"sync"
	"time"
)

// PollerFunc is a function run by the poller.
type PollerFunc func(context.Context) bool

// Poller is a poller. It runs an action.
type Poller interface {
	Run(ctx context.Context, ch <-chan struct{}, action PollerFunc, resetDuration time.Duration)
}

// SingleRunPoller is a a poller runs the supplied action once. It is useful for testing.
type SingleRunPoller struct{}

var _ Poller = (*SingleRunPoller)(nil)

// NewSingleRunPoller creates an instance of SingleRunPoller.
func NewSingleRunPoller() *SingleRunPoller {
	return &SingleRunPoller{}
}

// Run runs the poller.
func (a SingleRunPoller) Run(ctx context.Context, ch <-chan struct{}, action PollerFunc, resetDuration time.Duration) {
	action(ctx)
}

// InterruptiblePoller is a poller than runs an action and allows for interrupts.
type InterruptiblePoller struct {
}

var _ Poller = (*InterruptiblePoller)(nil)

// NewInterruptiblePoller creates an instance of InterruptiblePoller.
func NewInterruptiblePoller() *InterruptiblePoller {
	return &InterruptiblePoller{}
}

// Run runs the poller.
func (ip *InterruptiblePoller) Run(ctx context.Context, ch <-chan struct{}, action PollerFunc, resetDuration time.Duration) {
	timer := time.NewTimer(0)
	var cancel context.CancelFunc
	var mu sync.Mutex

	go func() {
		for _ = range ch {
			mu.Lock()
			if cancel != nil {
				cancel()
			}
			mu.Unlock()

			timer.Reset(0)
		}
	}()

	done := false
	for !done {
		select {
		case <-ctx.Done():
			done = true
			break
		case <-timer.C:
			var actionContext context.Context
			go func() {
				mu.Lock()
				actionContext, cancel = context.WithCancel(ctx)
				mu.Unlock()

				rerun := action(actionContext)
				if actionContext.Err() == nil {
					dur := resetDuration
					if rerun {
						dur = 0
					}
					timer.Reset(dur)
				}
				cancel()
			}()
		}
	}
}
