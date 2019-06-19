/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"encoding/json"
	"net/http"

	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/internal/module"
)

type namespace struct {
	moduleManager module.ManagerInterface
	logger        log.Logger
}

func newNamespace(moduleManager module.ManagerInterface, logger log.Logger) *namespace {
	return &namespace{
		moduleManager: moduleManager,
		logger:        logger,
	}
}

type namespaceRequest struct {
	Namespace string `json:"namespace,omitempty"`
}

func (n *namespace) update(w http.ResponseWriter, r *http.Request) {
	var nr namespaceRequest

	err := json.NewDecoder(r.Body).Decode(&nr)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "unable to decode request", n.logger)
		return
	}

	if nr.Namespace == "" {
		RespondWithError(w, http.StatusBadRequest, "unable to decode request", n.logger)
		return
	}

	n.moduleManager.SetNamespace(nr.Namespace)

	w.WriteHeader(http.StatusNoContent)
}

type namespaceResponse struct {
	Namespace string `json:"namespace,omitempty"`
}

func (n *namespace) read(w http.ResponseWriter, r *http.Request) {
	ns := n.moduleManager.GetNamespace()
	nr := &namespaceResponse{
		Namespace: ns,
	}

	if err := json.NewEncoder(w).Encode(nr); err != nil {
		n.logger.Errorf("encoding namespace error: %v", err)
	}
}
