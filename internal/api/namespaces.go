/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"encoding/json"
	"net/http"

	"github.com/vmware/octant/internal/cluster"
	"github.com/vmware/octant/internal/log"
)

type namespacesResponse struct {
	Namespaces []string `json:"namespaces,omitempty"`
}

type namespaces struct {
	nsClient cluster.NamespaceInterface
	logger   log.Logger
}

var _ http.Handler = (*namespaces)(nil)

func newNamespaces(nsClient cluster.NamespaceInterface, logger log.Logger) *namespaces {
	return &namespaces{
		nsClient: nsClient,
		logger:   logger,
	}
}

// ServeHTTP implements http.Handler and returns a list of namespace names for a cluster.
func (n *namespaces) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	names, err := n.nsClient.Names()
	if err != nil {
		// Fallback to initial namespace
		initialNamespace := n.nsClient.InitialNamespace()
		n.logger.Debugf("could not list namespaces, falling back to context namespace: %v (%v)", initialNamespace, err)
		names = []string{initialNamespace}
	}

	nr := &namespacesResponse{
		Namespaces: names,
	}

	if err := json.NewEncoder(w).Encode(nr); err != nil {
		n.logger.Errorf("encoding namespaces: %v", err)
	}
}
