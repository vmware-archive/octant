/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package terminal

import (
	"context"

	"github.com/pkg/errors"
	"github.com/vmware-tanzu/octant/internal/cluster"
	"github.com/vmware-tanzu/octant/pkg/store"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

//go:generate mockgen -source=manager.go -destination=./fake/mock_interface.go -package=fake github.com/vmware-tanzu/octant/internal/terminal TerminalManager

type manager struct {
	restClient  rest.Interface
	config      *rest.Config
	objectStore store.Store
	instances   map[string]Instance
}

var _ Manager = (*manager)(nil)

// NewTerminalManager creates a concrete TerminalMananger
func NewTerminalManager(ctx context.Context, client cluster.ClientInterface, objectStore store.Store) (Manager, error) {
	restClient, err := client.RESTClient()
	if err != nil {
		return nil, errors.Wrap(err, "fetching RESTClient")
	}

	tm := &manager{
		restClient:  restClient,
		config:      client.RESTConfig(),
		objectStore: objectStore,
		instances:   map[string]Instance{},
	}
	return tm, nil
}

func (tm *manager) Create(ctx context.Context, gvk schema.GroupVersionKind, name string, namespace string, container string, command string) (Instance, error) {
	t := NewTerminalInstance(ctx)
	tm.instances[t.ID(ctx)] = t
	return t, nil
}

func (tm *manager) Get(ctx context.Context, id string) (Instance, bool) {
	v, ok := tm.instances[id]
	return v, ok
}

func (tm *manager) List(ctx context.Context) []Instance {
	instances := make([]Instance, len(tm.instances))
	for _, instance := range tm.instances {
		instances = append(instances, instance)
	}
	return instances
}

func (tm *manager) StopAll(ctx context.Context) error {
	for _, instance := range tm.instances {
		instance.Stop(ctx)
	}
	return nil
}
