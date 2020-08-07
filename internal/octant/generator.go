/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package octant

import (
	"context"
	"time"

	"github.com/vmware-tanzu/octant/pkg/event"
)

//go:generate mockgen -destination=./fake/mock_generator.go -package=fake github.com/vmware-tanzu/octant/internal/octant Generator

// Generator generates events.
type Generator interface {
	// Event generates events using the returned channel.
	Event(ctx context.Context) (event.Event, error)

	// ScheduleDelay is how long to wait before scheduling this generator again.
	ScheduleDelay() time.Duration

	// Name is the generator name.
	Name() string
}
