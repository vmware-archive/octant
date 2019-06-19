/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"context"

	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/pkg/store"
	"github.com/vmware/octant/pkg/plugin/api/proto"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Client is a dashboard service API client.
type Client struct {
	conn *grpc.ClientConn
}

var _ Service = (*Client)(nil)

// NewClient creates an instance of the API client. It requires the
// address of the API.
func NewClient(address string) (*Client, error) {
	// NOTE: is it possible to make this secure? Is it even important?
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err

	}

	client := &Client{
		conn: conn,
	}

	return client, nil
}

// Close closes the client's connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

// List lists objects in the dashboard's objectstore.
func (c *Client) List(ctx context.Context, key store.Key) ([]*unstructured.Unstructured, error) {
	client := proto.NewDashboardClient(c.conn)

	keyRequest, err := convertFromKey(key)
	if err != nil {
		return nil, err
	}

	resp, err := client.List(ctx, keyRequest)
	if err != nil {
		return nil, err
	}

	objects, err := convertToObjects(resp.Objects)
	if err != nil {
		return nil, err
	}

	return objects, nil
}

// Get retrieves an object from the dashboard's objectstore.
func (c *Client) Get(ctx context.Context, key store.Key) (*unstructured.Unstructured, error) {
	client := proto.NewDashboardClient(c.conn)

	keyRequest, err := convertFromKey(key)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, keyRequest)
	if err != nil {
		return nil, err
	}

	object, err := convertToObject(resp.Object)
	if err != nil {
		return nil, err
	}

	return object, nil
}

// PortForward creates a port forward.
func (c *Client) PortForward(ctx context.Context, req PortForwardRequest) (PortForwardResponse, error) {
	client := proto.NewDashboardClient(c.conn)

	pfRequest := &proto.PortForwardRequest{
		Namespace:  req.Namespace,
		PodName:    req.PodName,
		PortNumber: uint32(req.Port),
	}
	resp, err := client.PortForward(ctx, pfRequest)
	if err != nil {
		return PortForwardResponse{}, err
	}

	return PortForwardResponse{
		ID:   resp.PortForwardID,
		Port: uint16(resp.PortNumber),
	}, nil

}

// CancelPortForward cancels a port forward.
func (c *Client) CancelPortForward(ctx context.Context, id string) {
	client := proto.NewDashboardClient(c.conn)

	req := &proto.CancelPortForwardRequest{
		PortForwardID: id,
	}

	_, err := client.CancelPortForward(ctx, req)
	if err != nil {
		logger := log.From(ctx)
		logger.Errorf("unable to cancel port forward: %v", err)
	}
}
