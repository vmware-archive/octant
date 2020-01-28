/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/wait"

	oerrors "github.com/vmware-tanzu/octant/internal/errors"
	"github.com/vmware-tanzu/octant/internal/log"
)

const (
	pollerWorkerCount = 2
)

// PollerFunc is a function run by the poller.
type PollerFunc func(context.Context) (bool, error)

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

func (a SingleRunPoller) ResetBackoff() {}

// InterruptiblePoller is a poller than runs an action and allows for interrupts.
type InterruptiblePoller struct {
	name    string
	backoff wait.Backoff
}

var _ Poller = (*InterruptiblePoller)(nil)

// NewInterruptiblePoller creates an instance of InterruptiblePoller.
func NewInterruptiblePoller(name string) *InterruptiblePoller {
	ip := &InterruptiblePoller{name: name}
	ip.resetBackoff()
	return ip
}

func (ip *InterruptiblePoller) resetBackoff() {
	ip.backoff = wait.Backoff{
		Duration: 1 * time.Second,
		Factor:   2.0,
		Jitter:   0.1,
		Steps:    50,
		Cap:      10 * time.Minute,
	}
}

// Run runs the poller.
func (ip *InterruptiblePoller) Run(ctx context.Context, ch <-chan struct{}, action PollerFunc, resetDuration time.Duration) {
	logger := log.From(ctx).With("poller-name", ip.name)
	ctx = log.WithLoggerContext(ctx, logger)

	jt := initJobTracker(ctx, action)
	defer jt.clear()

	resetCh := make(chan bool, 1)
	pollerQueue := make(chan job, 10)

	worker := func() {
		for j := range pollerQueue {
			select {
			case <-j.ctx.Done():
				// Job's context was canceled. Nothing else to do here.
			case <-j.run():
				backoffDuration := resetDuration
				if ip.isBackoffError(jt.logger, j.err) {
					backoffDuration = ip.backoff.Step()
					jt.logger.Debugf("poller backing off, next run %s)", backoffDuration)
				}

				select {
				case <-resetCh:
					break
				case <-time.After(backoffDuration):
					break
				}
				pollerQueue <- jt.create()
			}
		}
	}

	for i := 0; i < pollerWorkerCount; i++ {
		go worker()
	}

	go func() {
		for range ch {
			// Cancel all existing jobs before creating a new job.
			ip.resetBackoff()
			jt.clear()
			pollerQueue <- jt.create()
			resetCh <- true
		}
	}()

	pollerQueue <- jt.create()

	<-ctx.Done()
}

// isBackoffError checks if the error is one that should result in a backoff for the next call.
func (ip *InterruptiblePoller) isBackoffError(logger log.Logger, err error) bool {
	if err == nil {
		return false
	}

	if kerrors.IsUnauthorized(err) {
		return true
	}

	var ae *oerrors.AccessError
	if errors.As(err, &ae) {
		if ae.Name() == oerrors.OctantAccessError {
			return true
		}
	}

	es := err.Error()
	logger.Debugf("checking error string: %s", es)
	if strings.Contains(es, "Unauthorized") {
		return true
	}
	ip.resetBackoff()
	return false
}

type job struct {
	id         uuid.UUID
	ctx        context.Context
	cancelFunc context.CancelFunc
	action     PollerFunc
	err        error
}

func (j *job) run() <-chan bool {
	done := make(chan bool, 1)
	go func() {
		_, j.err = j.action(j.ctx)
		done <- true
	}()

	ch := make(chan bool, 1)
	go func() {
		select {
		case <-j.ctx.Done():
			ch <- true
		case <-done:
			ch <- true
		}
		defer close(ch)
	}()

	return ch
}

func createJob(ctx context.Context, action PollerFunc) job {
	ctx, cancel := context.WithCancel(ctx)

	return job{
		id:         uuid.New(),
		cancelFunc: cancel,
		ctx:        ctx,
		action:     action,
	}
}

func (j *job) cancel() {
	j.cancelFunc()
}

type jobTracker struct {
	jobs   map[uuid.UUID]job
	action PollerFunc
	mu     sync.Mutex
	ctx    context.Context
	logger log.Logger
}

func initJobTracker(ctx context.Context, action PollerFunc) *jobTracker {
	return &jobTracker{
		ctx:    ctx,
		action: action,
		jobs:   make(map[uuid.UUID]job),
		mu:     sync.Mutex{},
		logger: log.From(ctx),
	}
}

func (jt *jobTracker) create() job {
	jt.mu.Lock()
	defer jt.mu.Unlock()

	j := createJob(jt.ctx, jt.action)
	jt.jobs[j.id] = j

	return j
}

func (jt *jobTracker) clear() {
	jt.mu.Lock()
	defer jt.mu.Unlock()

	for id, j := range jt.jobs {
		j.cancel()
		delete(jt.jobs, id)
	}
}
