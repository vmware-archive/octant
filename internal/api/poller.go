/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/vmware-tanzu/octant/internal/log"
)

const (
	pollerWorkerCount = 2
)

// PollerFunc is a function run by the poller.
type PollerFunc func(context.Context) bool

// Poller is a poller. It runs an action.
type Poller interface {
	// Run runs `action` and delays `resetDuration` before starting again. If a message
	// is sent to `ch`, it will cancel current work and restart `action`.
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
func (a SingleRunPoller) Run(ctx context.Context, _ <-chan struct{}, action PollerFunc, resetDuration time.Duration) {
	action(ctx)
}

// InterruptiblePoller is a poller than runs an action and allows for interrupts.
type InterruptiblePoller struct {
	name string
}

var _ Poller = (*InterruptiblePoller)(nil)

// NewInterruptiblePoller creates an instance of InterruptiblePoller.
func NewInterruptiblePoller(name string) *InterruptiblePoller {
	return &InterruptiblePoller{
		name: name,
	}
}

// Run runs the poller.
func (a *InterruptiblePoller) Run(ctx context.Context, ch <-chan struct{}, action PollerFunc, resetDuration time.Duration) {
	logger := log.From(ctx).With(
		"poller-name", a.name,
		"poller-instance", uuid.New().String())
	logger.Debugf("starting poller")

	timer := time.NewTimer(0)
	done := false
	for !done {
		func() {
			cur, cancel := context.WithCancel(log.WithLoggerContext(ctx, logger))
			defer cancel()
			canceled := false

			select {
			case <-ctx.Done():
				logger.Debugf("poller has been canceled")
				done = true
			case <-ch:
				canceled = true
				logger.Debugf("poller was interrupted")
			case <-timer.C:
				logger.Debugf("poller is running action")
				now := time.Now()
				action(cur)
				logger.With("elapsed", fmt.Sprintf("%s", time.Since(now))).
					Debugf("poller ran action")
			}

			if !canceled {
				timer.Reset(resetDuration)
			}

		}()
	}
}
