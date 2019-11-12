/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package terminal

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/mime"
)

type terminalOutput struct {
	Scrollback []string `json:"scrollback,omitempty"`
}

func Handler(ctx context.Context, tm Manager) http.HandlerFunc {
	logger := log.From(ctx)

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uuid := vars["uuid"]

		t, ok := tm.Get(ctx, uuid)
		if !ok {
			respondWithError(w, http.StatusInternalServerError, "terminal not found", logger)
			return
		}

		output := terminalOutput{
			Scrollback: t.Scrollback(ctx),
		}

		if err := json.NewEncoder(w).Encode(&output); err != nil {
			logger := log.From(ctx)
			logger.With("err", err.Error()).Errorf("unable to encode log entries")
		}
	}
}

type errorMessage struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type errorResponse struct {
	Error errorMessage `json:"error,omitempty"`
}

// RespondWithError responds with an error message.
func respondWithError(w http.ResponseWriter, code int, message string, logger log.Logger) {
	r := &errorResponse{
		Error: errorMessage{
			Code:    code,
			Message: message,
		},
	}

	logger.With(
		"code", code,
		"message", message,
	).Infof("unable to serve")

	w.Header().Set("Content-Type", mime.JSONContentType)

	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(r); err != nil {
		logger.Errorf("encoding JSON response: %v", err)
	}
}
