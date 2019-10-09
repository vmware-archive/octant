/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package portforward

import (
	"context"
	"os"

	"github.com/vmware/octant/internal/cluster"
	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/pkg/store"

	"github.com/pkg/errors"
)

// Default create a port forward instance.
func Default(ctx context.Context, client cluster.ClientInterface, objectStore store.Store) (PortForwarder, error) {
	logger := log.From(ctx)
	restClient, err := client.RESTClient()
	if err != nil {
		return nil, errors.Wrap(err, "fetching RESTClient")
	}

	go func() {

	}()

	pfOpts := ServiceOptions{
		RESTClient:  restClient,
		Config:      client.RESTConfig(),
		ObjectStore: objectStore,
		PortForwarder: &DefaultPortForwarder{
			IOStreams: IOStreams{
				In:     os.Stdin,
				Out:    os.Stdout,
				ErrOut: os.Stderr,
			},
		},
	}

	// FIXME: logger is in context
	svc := New(ctx, pfOpts, logger)

	return svc, nil
}
