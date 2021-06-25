/*
   Copyright (c) 2019 the Octant contributors. All Rights Reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"context"

	"github.com/vmware-tanzu/octant/internal/octant"
)

//go:generate mockgen -destination=./fake/mock_state_manager.go -package=fake github.com/vmware-tanzu/octant/pkg/api StateManager
// StateManager manages states for WebsocketState.
type StateManager interface {
	Handlers() []octant.ClientRequestHandler
	Start(ctx context.Context, state octant.State, s OctantClient)
}
