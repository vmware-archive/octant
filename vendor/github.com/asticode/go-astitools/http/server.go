package astihttp

import (
	"context"
	"net"
	"net/http"

	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

func Serve(ctx context.Context, h http.Handler, addr string, fn func(a net.Addr)) (err error) {
	// Create listener
	var l net.Listener
	if l, err = net.Listen("tcp", addr); err != nil {
		err = errors.Wrap(err, "astihttp: net.Listen failed")
		return
	}
	defer l.Close()

	// Create server
	astilog.Debugf("astihttp: serving on %s", l.Addr())
	srv := &http.Server{Handler: h}
	defer func() {
		astilog.Debugf("astihttp: shutting down server on %s", l.Addr())
		if err := srv.Shutdown(ctx); err != nil {
			astilog.Error(errors.Wrapf(err, "astihttp: shutting down server on %s failed", l.Addr()))
		}
	}()

	// Serve
	var chanDone = make(chan error)
	go func() {
		if err := srv.Serve(l); err != nil {
			chanDone <- err
		}
	}()

	// Execute custom callback
	if fn != nil {
		fn(l.Addr())
	}

	// Wait for context or chanDone to be done
	select {
	case <-ctx.Done():
		if ctx.Err() != context.Canceled {
			err = errors.Wrap(err, "astihttp: context error")
		}
		return
	case err = <-chanDone:
		if err != nil {
			err = errors.Wrap(err, "astihttp: serving failed")
		}
		return
	}
}
