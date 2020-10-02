/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"context"

	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/spf13/viper"

	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/config"
	"github.com/vmware-tanzu/octant/pkg/plugin/api/proto"
	"github.com/vmware-tanzu/octant/pkg/store"
)

//go:generate mockgen -destination=./fake/mock_dashboard_client.go -package=fake github.com/vmware-tanzu/octant/pkg/plugin/api/proto DashboardClient
//go:generate mockgen -destination=./fake/mock_dashboard_connection.go -package=fake github.com/vmware-tanzu/octant/pkg/plugin/api DashboardConnection

type DashboardConnection interface {
	Close() error
	Client() proto.DashboardClient
}

type defaultDashboardConnection struct {
	conn *grpc.ClientConn
}

var _ DashboardConnection = (*defaultDashboardConnection)(nil)

func (d *defaultDashboardConnection) Close() error {
	return d.conn.Close()
}

func (d *defaultDashboardConnection) Client() proto.DashboardClient {
	return proto.NewDashboardClient(d.conn)
}

type ClientOption func(c *Client)

// Client is a dashboard service API client.
type Client struct {
	DashboardConnection DashboardConnection
}

var _ Service = (*Client)(nil)

// NewClient creates an instance of the API client. It requires the
// address of the API.
func NewClient(address string, options ...ClientOption) (*Client, error) {
	client := &Client{}
	viper.SetDefault("client-max-recv-msg-size", config.MaxMessageSize)

	for _, option := range options {
		option(client)
	}

	if client.DashboardConnection == nil {
		// NOTE: is it possible to make this secure? Is it even important?
		conn, err := grpc.Dial(address,
			grpc.WithInsecure(),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(viper.GetInt("client-max-recv-msg-size"))),
		)
		if err != nil {
			return nil, err

		}

		client.DashboardConnection = &defaultDashboardConnection{conn: conn}

	}

	return client, nil
}

// Close closes the client's connection.
func (c *Client) Close() error {
	return c.DashboardConnection.Close()
}

// List lists objects in the dashboard's object store.
func (c *Client) List(ctx context.Context, key store.Key) (*unstructured.UnstructuredList, error) {
	client := c.DashboardConnection.Client()

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

// Get retrieves an object from the dashboard's objectStore.
func (c *Client) Get(ctx context.Context, key store.Key) (*unstructured.Unstructured, error) {
	client := c.DashboardConnection.Client()

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

// Update updates an object in the store.
func (c *Client) Update(ctx context.Context, object *unstructured.Unstructured) error {
	client := c.DashboardConnection.Client()

	data, err := convertFromObject(object)
	if err != nil {
		return err
	}

	req := &proto.UpdateRequest{
		Object: data,
	}

	_, err = client.Update(ctx, req)

	return err
}

func (c *Client) Create(ctx context.Context, object *unstructured.Unstructured) error {
	client := c.DashboardConnection.Client()

	data, err := convertFromObject(object)
	if err != nil {
		return err
	}

	req := &proto.CreateRequest{
		Object: data,
	}

	_, err = client.Create(ctx, req)

	return err

}

// PortForward creates a port forward.
func (c *Client) PortForward(ctx context.Context, req PortForwardRequest) (PortForwardResponse, error) {
	client := c.DashboardConnection.Client()

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
	client := c.DashboardConnection.Client()

	req := &proto.CancelPortForwardRequest{
		PortForwardID: id,
	}

	_, err := client.CancelPortForward(ctx, req)
	if err != nil {
		logger := log.From(ctx)
		logger.Errorf("unable to cancel port forward: %v", err)
	}
}

// ListNamespaces lists namespaces.
func (c *Client) ListNamespaces(ctx context.Context) (NamespacesResponse, error) {
	client := c.DashboardConnection.Client()

	resp, err := client.ListNamespaces(ctx, &proto.Empty{})
	if err != nil {
		return NamespacesResponse{}, err
	}

	return NamespacesResponse{
		Namespaces: resp.Namespaces,
	}, nil
}

// ForceFrontendUpdate forces the frontend to update itself.
func (c *Client) ForceFrontendUpdate(ctx context.Context) error {
	client := c.DashboardConnection.Client()

	_, err := client.ForceFrontendUpdate(ctx, &proto.Empty{})
	return err
}

// SendAlert sends an alert
func (c *Client) SendAlert(ctx context.Context, clientID string, alert action.Alert) error {
	client := c.DashboardConnection.Client()

	alertRequest, err := convertFromAlert(alert)
	if err != nil {
		return err
	}

	alertRequest.ClientID = clientID

	_, err = client.SendAlert(ctx, alertRequest)
	return err
}
