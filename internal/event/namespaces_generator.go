/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package event

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/vmware/octant/internal/cluster"
	"github.com/vmware/octant/internal/octant"
)

type namespacesResponse struct {
	Namespaces []string `json:"namespaces,omitempty"`
}

// NamespacesGenerator generates namespaces events.
type NamespacesGenerator struct {
	// NamespaceClient is a namespaces client.
	NamespaceClient cluster.NamespaceInterface
}

var _ octant.Generator = (*NamespacesGenerator)(nil)

// Event generates namespaces events
func (g *NamespacesGenerator) Event(ctx context.Context) (octant.Event, error) {
	if g.NamespaceClient == nil {
		return octant.Event{}, errors.New("unable to query namespaces, client is nil")
	}

	names, err := g.NamespaceClient.Names()
	if err != nil {
		initialNamespace := g.NamespaceClient.InitialNamespace()
		names = []string{initialNamespace}
	}

	nr := &namespacesResponse{Namespaces: names}

	return octant.Event{
		Type: octant.EventTypeNamespaces,
		Data: nr,
	}, nil
}

// ScheduleDelay returns how long to delay before running this generator again.
func (NamespacesGenerator) ScheduleDelay() time.Duration {
	return DefaultScheduleDelay
}

// Name returns the generator's name.
func (NamespacesGenerator) Name() string {
	return "namespaces"
}
