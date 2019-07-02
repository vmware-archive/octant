/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"go.opencensus.io/trace"

	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/pkg/navigation"
)

type navSections interface {
	Sections(ctx context.Context, namespace string) ([]navigation.Navigation, error)
}

type navigationResponse struct {
	Sections []navigation.Navigation `json:"sections,omitempty"`
}

type navigationHandler struct {
	navSections navSections
	logger      log.Logger
}

var _ http.Handler = (*navigationHandler)(nil)

func newNavigationHandler(ns navSections, logger log.Logger) *navigationHandler {
	if logger == nil {
		logger = log.NopLogger()
	}

	return &navigationHandler{
		navSections: ns,
		logger:      logger,
	}
}

func (n *navigationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "api:navigationHandler")
	defer span.End()

	if n.navSections == nil {
		RespondWithError(w, http.StatusInternalServerError,
			"unable to generate navigationHandler sections", n.logger)
		return
	}

	vars := mux.Vars(r)
	namespace := vars["namespace"] // optional
	if namespace == "" {
		// Fallback to legacy query parameter
		namespace = "default"
	}

	n.logger.Debugf("navigationHandler for namespace %s", namespace)

	ns, err := n.navSections.Sections(ctx, namespace)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError,
			"unable to generate navigationHandler sections", n.logger)
		return
	}

	nr := navigationResponse{
		Sections: ns,
	}

	serveAsJSON(w, &nr, n.logger)
}
