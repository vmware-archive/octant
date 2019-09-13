/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
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
	Dispatch(ctx context.Context, actionName string, payload action.Payload) error
}
