/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInterruptiblePoller_Run(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	ip := NewInterruptiblePoller("poller")

	resetDuration := 10 * time.Millisecond
	ch := make(chan struct{}, 1)

	ready := make(chan bool, 1)
	ran := false
	action := func(ctx context.Context) (bool, error) {
		ready <- true
		ran = true
		return false, nil
	}

	exited := make(chan bool, 1)
	go func() {
		ip.Run(ctx, ch, action, resetDuration)
		exited <- true
	}()

	<-ready
	close(ch)
	cancel()
	<-exited

	assert.True(t, ran)
}
