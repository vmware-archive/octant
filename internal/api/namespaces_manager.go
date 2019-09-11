/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api

import (
	"context"
	"time"

	"github.com/vmware/octant/internal/cluster"
	"github.com/vmware/octant/internal/event"
	"github.com/vmware/octant/internal/octant"
)

// NamespacesManager manages namespaces.
type NamespacesManager struct {
	clusterClient cluster.ClientInterface
}

var _ StateManager = (*NamespacesManager)(nil)

// NewNamespacesManager creates an instance of NamespacesManager.
func NewNamespacesManager(clusterClient cluster.ClientInterface) *NamespacesManager {
	return &NamespacesManager{
		clusterClient: clusterClient,
	}
}

// Handlers returns nil.
func (n NamespacesManager) Handlers() []octant.ClientRequestHandler {
	return nil
}

// Start starts the manager. It periodically generates a list of namespaces.
func (n *NamespacesManager) Start(ctx context.Context, state octant.State, s OctantClient) {
	timer := time.NewTimer(0)

	done := false
	for !done {
		select {
		case <-ctx.Done():
			done = true
			break
		case <-timer.C:
			generator, err := n.initGenerator(state)
			if err != nil {
				// do something with this error?
				break
			}
			ev, err := generator.Event(ctx)
			if err != nil {
				// do something with this error?
				break
			}

			s.Send(ev)
			timer.Reset(generator.ScheduleDelay())
		}
	}
	timer.Stop()
}

func (n *NamespacesManager) initGenerator(state octant.State) (*event.NamespacesGenerator, error) {
	nsClient, err := n.clusterClient.NamespaceClient()
	if err != nil {
		return nil, err
	}
	return &event.NamespacesGenerator{
		NamespaceClient: nsClient,
	}, nil
}
