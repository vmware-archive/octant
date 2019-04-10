package portforward

import (
	"context"
	"os"

	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/objectstore"
	"github.com/pkg/errors"
)

// Default create a portforward instance.
func Default(ctx context.Context, client cluster.ClientInterface, objectStore objectstore.ObjectStore) (PortForwarder, error) {
	logger := log.From(ctx)
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

	// FIXME: logger is in context
	svc := New(ctx, pfOpts, logger)

	return svc, nil
}
