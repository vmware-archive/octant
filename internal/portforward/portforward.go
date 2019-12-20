/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package portforward

import (
	"io"
	"net/http"
	"net/url"

	"k8s.io/client-go/rest"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

// Options contains all the options for running a port-forward
// <snip> from pkg/kubectl/cmd/portforward/portforward.go <snip>
type Options struct {
	Config        *restclient.Config
	RESTClient    rest.Interface
	Address       []string
	Ports         []string
	PortForwarder portForwarder
	StopChannel   <-chan struct{}
	ReadyChannel  chan struct{}
	PortsChannel  chan []ForwardedPort
}

type portForwarder interface {
	ForwardPorts(method string, url *url.URL, opts Options) error
}

type DefaultPortForwarder struct {
	IOStreams
}

// IOStreams provides the standard names for iostreams.  This is useful for embedding and for unit testing.
// Inconsistent and different names make it hard to read and review code
type IOStreams struct {
	// In think, os.Stdin
	In io.Reader
	// Out think, os.Stdout
	Out io.Writer
	// ErrOut think, os.Stderr
	ErrOut io.Writer
}

type ForwardedPort struct {
	Local  uint16
	Remote uint16
}

func (f *DefaultPortForwarder) ForwardPorts(method string, url *url.URL, opts Options) error {
	transport, upgrader, err := spdy.RoundTripperFor(opts.Config)
	if err != nil {
		return err
	}
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, method, url)
	fw, err := portforward.NewOnAddresses(dialer, opts.Address, opts.Ports, opts.StopChannel, opts.ReadyChannel, f.Out, f.ErrOut)
	if err != nil {
		return err
	}

	// Wait for resolved ports to become available and send them on up
	localPortsHandler(fw, opts)

	// Forward and block
	return fw.ForwardPorts()
}

// </snip> from pkg/kubectl/cmd/portforward/portforward.go </snip>

// localPortsHandler manages passing up the resolved local ports from the forwarder when
// they become available via the PortsChannel.
func localPortsHandler(fw *portforward.PortForwarder, opts Options) {
	if fw == nil {
		return
	}

	if opts.ReadyChannel != nil && opts.PortsChannel != nil {
		go func() {
			select {
			case <-opts.ReadyChannel:
				ports, err := fw.GetPorts()
				if err != nil {
					// TODO: alert user (GH#497)
					return
				}
				fp := make([]ForwardedPort, len(ports))
				for i := range ports {
					fp[i].Local = ports[i].Local
					fp[i].Remote = ports[i].Remote
				}
				opts.PortsChannel <- fp

			case <-opts.StopChannel:
				return
			}
		}()
	}
}
