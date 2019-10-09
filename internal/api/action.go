/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"context"

	"github.com/vmware/octant/pkg/action"
)

//go:generate mockgen -destination=./fake/mock_action_dispatcher.go -package=fake github.com/vmware/octant/internal/api ActionDispatcher

// ActionDispatcher dispatches actions.
type ActionDispatcher interface {
	Dispatch(ctx context.Context, alerter action.Alerter, actionName string, payload action.Payload) error
}
