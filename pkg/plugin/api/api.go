/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"context"
	"net"

	"github.com/spf13/viper"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/pkg/plugin/api/proto"
)

// API controls the dashboard API service.
type API interface {
	// Addr is the address of the API service.
	Addr() string
	// Start starts the API. To stop the API, cancel the context.
	Start(context.Context) error
}

// grpcAPI is in implementation of API backed by GRPC.
type grpcAPI struct {
	Service  Service
	listener net.Listener
}

const dashServiceAddress = "127.0.0.1:0"

var _ API = (*grpcAPI)(nil)

// New creates a new API instance for DashService.
func New(service Service) (API, error) {
	listener, err := net.Listen("tcp", dashServiceAddress)
	if err != nil {
		return nil, errors.Wrap(err, "create listener")
	}

	return &grpcAPI{
		Service:  service,
		listener: listener,
	}, nil
}

// Start starts the API.
func (a *grpcAPI) Start(ctx context.Context) error {
	logger := log.From(ctx)

	dashboardServer := &grpcServer{
		service: a.Service,
	}

	s := grpc.NewServer(
		grpc.MaxRecvMsgSize(viper.GetInt("client-max-recv-msg-size")),
	)

	proto.RegisterDashboardServer(s, dashboardServer)

	logger.Debugf("dashboard plugin api is starting")
	go func() {
		if err := s.Serve(a.listener); err != nil {
			logger.WithErr(err).Errorf("unable to serve GRPC")
			return
		}
	}()

	go func() {
		<-ctx.Done()
		logger.Debugf("dashboard plugin api is stopping")
		s.Stop()
	}()

	return nil
}

func (a *grpcAPI) Addr() string {
	return a.listener.Addr().String()
}
