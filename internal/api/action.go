/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/pkg/action"
)

//go:generate mockgen -destination=./fake/mock_action_dispatcher.go -package=fake github.com/vmware/octant/internal/api ActionDispatcher

// ActionDispatcher dispatches actions.
type ActionDispatcher interface {
	Dispatch(ctx context.Context, actionName string, payload action.Payload) error
}

// UpdateRequest is an update request. It contains the action payload.
type UpdateRequest struct {
	Update action.Payload `json:"update"`
}

// ActionHandler is a handler that responds to action messages.
type ActionHandler struct {
	logger           log.Logger
	actionDispatcher ActionDispatcher
}

var _ http.Handler = (*ActionHandler)(nil)

// NewActionHandler creates an instance of ActionHandler.
func NewActionHandler(logger log.Logger, actionDispatcher ActionDispatcher) *ActionHandler {
	return &ActionHandler{
		logger:           logger,
		actionDispatcher: actionDispatcher,
	}
}

// ServeHTTP serves the handler.
func (a *ActionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req UpdateRequest

	defer func() {
		if cErr := r.Body.Close(); cErr != nil {
			a.logger.WithErr(cErr).Errorf("unable to close action request body")
		}
	}()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error(), a.logger)
		return
	}

	actionName, err := req.Update.String("action")
	if err != nil {
		RespondWithError(w, http.StatusNotFound, fmt.Sprintf("unknown action %v", req.Update), a.logger)
		return
	}

	if err := a.actionDispatcher.Dispatch(r.Context(), actionName, req.Update); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error(), a.logger)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
