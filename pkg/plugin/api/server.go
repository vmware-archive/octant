/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"context"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/heptio/developer-dash/internal/gvk"
	"github.com/heptio/developer-dash/internal/portforward"
	"github.com/heptio/developer-dash/pkg/store"
	"github.com/heptio/developer-dash/pkg/plugin/api/proto"
)

// PortForwardRequest describes a port forward request.
type PortForwardRequest struct {
	Namespace     string
	PodName       string
	ContainerName string
	Port          uint16
}

// PortForwardResponse is the response from a port forward request.
type PortForwardResponse struct {
	ID   string
	Port uint16
}

// Service is the dashboard service.
type Service interface {
	List(ctx context.Context, key store.Key) ([]*unstructured.Unstructured, error)
	Get(ctx context.Context, key store.Key) (*unstructured.Unstructured, error)
	PortForward(ctx context.Context, req PortForwardRequest) (PortForwardResponse, error)
	CancelPortForward(ctx context.Context, id string)
}

// GRPCService is an implementation of the dashboard service based on GRPC.
type GRPCService struct {
	ObjectStore   store.Store
	PortForwarder portforward.PortForwarder
}

var _ Service = (*GRPCService)(nil)

// List lists objects.
func (s *GRPCService) List(ctx context.Context, key store.Key) ([]*unstructured.Unstructured, error) {
	return s.ObjectStore.List(ctx, key)
}

// Get retrieves an object.
func (s *GRPCService) Get(ctx context.Context, key store.Key) (*unstructured.Unstructured, error) {
	return s.ObjectStore.Get(ctx, key)
}

// PortForward creates a port forward.
func (s *GRPCService) PortForward(ctx context.Context, req PortForwardRequest) (PortForwardResponse, error) {
	pfResponse, err := s.PortForwarder.Create(
		ctx,
		gvk.PodGVK,
		req.PodName,
		req.Namespace,
		req.Port)
	if err != nil {
		return PortForwardResponse{}, err
	}

	resp := PortForwardResponse{
		ID:   pfResponse.ID,
		Port: pfResponse.Ports[0].Local,
	}

	return resp, nil
}

// CancelPortForward cancels a port forward
func (s *GRPCService) CancelPortForward(ctx context.Context, id string) {
	s.PortForwarder.StopForwarder(id)
}

type grpcServer struct {
	service Service
}

var _ proto.DashboardServer = (*grpcServer)(nil)

func (c *grpcServer) List(ctx context.Context, in *proto.KeyRequest) (*proto.ListResponse, error) {
	key, err := convertToKey(in)
	if err != nil {
		return nil, err
	}

	objects, err := c.service.List(ctx, key)
	if err != nil {
		return nil, err
	}

	encodedObjects, err := convertFromObjects(objects)
	if err != nil {
		return nil, err
	}

	out := &proto.ListResponse{
		Objects: encodedObjects,
	}

	return out, nil
}

func (c *grpcServer) Get(ctx context.Context, in *proto.KeyRequest) (*proto.GetResponse, error) {
	key, err := convertToKey(in)
	if err != nil {
		return nil, err
	}

	object, err := c.service.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	encodedObject, err := convertFromObject(object)
	if err != nil {
		return nil, err
	}

	out := &proto.GetResponse{
		Object: encodedObject,
	}

	return out, nil
}

func (c *grpcServer) PortForward(ctx context.Context, in *proto.PortForwardRequest) (*proto.PortForwardResponse, error) {
	req, err := convertToPortForwardRequest(in)
	if err != nil {
		return nil, err
	}

	pfResp, err := c.service.PortForward(ctx, *req)
	if err != nil {
		return nil, err
	}

	resp := &proto.PortForwardResponse{
		PortForwardID: pfResp.ID,
		PortNumber:    uint32(pfResp.Port),
	}

	return resp, nil
}

func (c *grpcServer) CancelPortForward(ctx context.Context, in *proto.CancelPortForwardRequest) (*proto.Empty, error) {
	if in == nil {
		return nil, errors.New("request is nil")
	}

	c.service.CancelPortForward(ctx, in.PortForwardID)
	return &proto.Empty{}, nil
}
