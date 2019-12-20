/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package portforward

import (
	"context"
	"os"

	"github.com/vmware-tanzu/octant/internal/cluster"
	"github.com/vmware-tanzu/octant/pkg/store"

	"github.com/pkg/errors"
)

// Default create a port forward instance.
func Default(ctx context.Context, client cluster.ClientInterface, objectStore store.Store) (PortForwarder, error) {
	restClient, err := client.RESTClient()
	if err != nil {
		return nil, errors.Wrap(err, "fetching RESTClient")
	}

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

	svc := New(ctx, pfOpts)

	return svc, nil
}
