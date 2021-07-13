/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api

import (
	"context"

	ocontext "github.com/vmware-tanzu/octant/internal/context"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/api"
	"github.com/vmware-tanzu/octant/pkg/config"
)

const (
	RequestPerformAction = "action.octant.dev/performAction"
)

// ActionRequestManager manages action requests. Action requests allow a generic interface
// for supporting dynamic requests from clients.
type ActionRequestManager struct {
	dashConfig config.Dash
}

var _ StateManager = (*ActionRequestManager)(nil)

// NewActionRequestManager creates an instance of ActionRequestManager.
func NewActionRequestManager(dashConfig config.Dash) *ActionRequestManager {
	return &ActionRequestManager{
		dashConfig: dashConfig,
	}
}

func (a ActionRequestManager) Start(ctx context.Context, state octant.State, s api.OctantClient) {
}

// Handlers returns the handlers this manager supports.
func (a *ActionRequestManager) Handlers() []octant.ClientRequestHandler {
	return []octant.ClientRequestHandler{
		{
			RequestType: RequestPerformAction,
			Handler:     a.PerformAction,
		},
	}
}

// PerformAction is a handler than runs an action.
func (a *ActionRequestManager) PerformAction(state octant.State, payload action.Payload) error {
	clientState := ocontext.ClientState{
		ClientID:    state.GetClientID(),
		Filters:     state.GetFilters(),
		Namespace:   state.GetNamespace(),
		ContextName: a.dashConfig.CurrentContext(),
	}
	ctx := ocontext.WithClientState(context.TODO(), clientState)

	actionName, err := payload.String("action")
	if err != nil {
		// TODO: alert the user this action doesn't exist (GH#493)
		return nil
	}

	if err := state.Dispatch(ctx, actionName, payload); err != nil {
		return err
	}

	return nil
}
