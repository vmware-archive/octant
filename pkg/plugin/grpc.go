/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"context"
	"encoding/json"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	ocontext "github.com/vmware-tanzu/octant/internal/context"

	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/plugin/dashboard"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// GRPCClient is the dashboard GRPC client.
type GRPCClient struct {
	broker Broker
	client dashboard.PluginClient
}

var _ Service = (*GRPCClient)(nil)
var _ ModuleService = (*GRPCClient)(nil)

// NewGRPCClient creates an instance of GRPCClient.
func NewGRPCClient(broker Broker, client dashboard.PluginClient) *GRPCClient {
	return &GRPCClient{
		client: client,
		broker: broker,
	}
}

func (c *GRPCClient) run(fn func() error) error {
	if fn == nil {
		return errors.New("client function is nil")
	}

	return fn()
}

// Content returns content from a plugin.
func (c *GRPCClient) Content(ctx context.Context, contentPath string) (component.ContentResponse, error) {
	var contentResponse component.ContentResponse

	err := c.run(func() error {
		req := &dashboard.ContentRequest{
			Path: contentPath,
		}

		resp, err := c.client.Content(ctx, req, grpc.WaitForReady(true))
		if err != nil {
			return errors.Wrap(err, "grpc client content")
		}

		if err := json.Unmarshal(resp.ContentResponse, &contentResponse); err != nil {
			return errors.Wrap(err, "unmarshal content response")
		}

		return nil
	})

	if err != nil {
		return component.ContentResponse{}, err
	}

	return contentResponse, nil
}

// HandleAction runs an action on a plugin.
func (c *GRPCClient) HandleAction(ctx context.Context, actionName string, payload action.Payload) error {
	err := c.run(func() error {
		clientID := ocontext.WebsocketClientIDFrom(ctx)
		data, err := json.Marshal(&payload)
		if err != nil {
			return err
		}

		req := &dashboard.HandleActionRequest{
			ActionName: actionName,
			Payload:    data,
			ClientID:   clientID,
		}

		_, err = c.client.HandleAction(ctx, req, grpc.WaitForReady(true))
		if err != nil {
			if s, isStatus := status.FromError(err); isStatus {
				return errors.Errorf("grpc error: %s", s.Message())
			}
			return err
		}

		return nil
	})

	return err
}

// Navigation returns navigation entries from a plugin.
func (c *GRPCClient) Navigation(ctx context.Context) (navigation.Navigation, error) {
	var entries navigation.Navigation

	if ctx.Err() == context.Canceled {
		return navigation.Navigation{}, nil
	}

	err := c.run(func() error {
		clientID := ocontext.WebsocketClientIDFrom(ctx)
		req := &dashboard.NavigationRequest{
			ClientID: clientID,
		}

		resp, err := c.client.Navigation(ctx, req, grpc.WaitForReady(true))
		if err != nil {
			return errors.Wrap(err, "grpc client response")
		}

		entries = convertToNavigation(resp.Navigation)

		return nil
	})

	if err != nil {
		return navigation.Navigation{}, err
	}

	return entries, nil
}

// Register register a plugin.
func (c *GRPCClient) Register(ctx context.Context, dashboardAPIAddress string) (Metadata, error) {
	var m Metadata

	err := c.run(func() error {
		registerRequest := &dashboard.RegisterRequest{
			DashboardAPIAddress: dashboardAPIAddress,
		}

		resp, err := c.client.Register(ctx, registerRequest, grpc.WaitForReady(true))
		if err != nil {
			spew.Dump(err)
			return errors.WithMessage(err, "unable to call register function")
		}

		capabilities := convertToCapabilities(resp.Capabilities)

		m = Metadata{
			Name:         resp.PluginName,
			Description:  resp.Description,
			Capabilities: capabilities,
		}

		return nil
	})

	if err != nil {
		return Metadata{}, err
	}

	return m, nil
}

// ObjectStatus gets an object status
func (c *GRPCClient) ObjectStatus(ctx context.Context, object runtime.Object) (ObjectStatusResponse, error) {
	var osr ObjectStatusResponse

	clientID := ocontext.WebsocketClientIDFrom(ctx)

	err := c.run(func() error {
		in, err := createObjectRequest(object, clientID)
		if err != nil {
			return err
		}

		resp, err := c.client.ObjectStatus(ctx, in, grpc.WaitForReady(true))
		if err != nil {
			return errors.Wrap(err, "grpc client object status")
		}

		var objectStatus component.PodSummary
		if err := json.Unmarshal(resp.ObjectStatus, &objectStatus); err != nil {
			return errors.Wrap(err, "convert object status")
		}

		osr = ObjectStatusResponse{
			ObjectStatus: objectStatus,
		}

		return nil
	})

	if err != nil {
		return ObjectStatusResponse{}, err
	}

	return osr, nil
}

// Print prints an object.
func (c *GRPCClient) Print(ctx context.Context, object runtime.Object) (PrintResponse, error) {
	var pr PrintResponse

	clientID := ocontext.WebsocketClientIDFrom(ctx)

	err := c.run(func() error {
		in, err := createObjectRequest(object, clientID)
		if err != nil {
			return err
		}

		resp, err := c.client.Print(ctx, in, grpc.WaitForReady(true))
		if err != nil {
			return errors.Wrap(err, "grpc client print")
		}

		var items []component.FlexLayoutItem
		if len(resp.Items) > 0 {
			if err := json.Unmarshal(resp.Items, &items); err != nil {
				return err
			}
		}

		configSection, err := convertToSummarySections(resp.Config)
		if err != nil {
			return errors.Wrap(err, "convert config sections")
		}

		summarySection, err := convertToSummarySections(resp.Status)
		if err != nil {
			return errors.Wrap(err, "convert summary sections")
		}

		pr = PrintResponse{
			Config: configSection,
			Status: summarySection,
			Items:  items,
		}

		return nil
	})

	if err != nil {
		return PrintResponse{}, err
	}

	return pr, nil
}

func createObjectRequest(object runtime.Object, clientID string) (*dashboard.ObjectRequest, error) {
	data, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}

	or := &dashboard.ObjectRequest{
		Object:   data,
		ClientID: clientID,
	}

	return or, err
}

// PrintTab creates a tab for an object.
func (c *GRPCClient) PrintTab(ctx context.Context, object runtime.Object) (TabResponse, error) {
	var tab component.Tab

	clientID := ocontext.WebsocketClientIDFrom(ctx)

	err := c.run(func() error {
		in, err := createObjectRequest(object, clientID)
		if err != nil {
			return err
		}

		resp, err := c.client.PrintTab(ctx, in, grpc.WaitForReady(true))
		if err != nil {
			return errors.Wrap(err, "grpc client print tab")
		}

		var to component.TypedObject
		if err := json.Unmarshal(resp.Layout, &to); err != nil {
			return err
		}

		c, err := to.ToComponent()
		if err != nil {
			return err
		}

		layout, ok := c.(*component.FlexLayout)
		if !ok {
			return errors.Errorf("expected to be flex layout was: %T", c)
		}

		tab.Name = resp.Name
		tab.Contents = *layout

		return nil
	})

	if err != nil {
		return TabResponse{}, err
	}

	return TabResponse{Tab: &tab}, nil
}

// GRPCServer is the grpc server the dashboard will use to communicate with the
// the plugin.
type GRPCServer struct {
	Impl   Service
	broker Broker
}

var _ dashboard.PluginServer = (*GRPCServer)(nil)

// Content returns content from a plugin.
func (s *GRPCServer) Content(ctx context.Context, req *dashboard.ContentRequest) (*dashboard.ContentResponse, error) {
	service, ok := s.Impl.(ModuleService)
	if !ok {
		return nil, errors.Errorf("plugin is not a module, it's a %T", s.Impl)
	}

	ctx = ocontext.WithWebsocketClientID(ctx, req.ClientID)
	contentResponse, err := service.Content(ctx, req.Path)
	if err != nil {
		return nil, err
	}

	contentResponseBytes, err := json.Marshal(&contentResponse)
	if err != nil {
		return nil, err
	}

	return &dashboard.ContentResponse{
		ContentResponse: contentResponseBytes,
	}, nil
}

// HandleAction runs an action in a plugin.
func (s *GRPCServer) HandleAction(ctx context.Context, handleActionRequest *dashboard.HandleActionRequest) (*dashboard.HandleActionResponse, error) {
	var payload action.Payload
	if err := json.Unmarshal(handleActionRequest.Payload, &payload); err != nil {
		return nil, err
	}

	ctx = ocontext.WithWebsocketClientID(ctx, handleActionRequest.ClientID)
	if err := s.Impl.HandleAction(ctx, handleActionRequest.ActionName, payload); err != nil {
		return nil, err
	}

	return &dashboard.HandleActionResponse{}, nil
}

// Navigation returns navigation entries from a plugin.
func (s *GRPCServer) Navigation(ctx context.Context, req *dashboard.NavigationRequest) (*dashboard.NavigationResponse, error) {
	service, ok := s.Impl.(ModuleService)
	if !ok {
		return nil, errors.Errorf("plugin is not a module, it's a %T", s.Impl)
	}

	ctx = ocontext.WithWebsocketClientID(ctx, req.ClientID)
	entry, err := service.Navigation(ctx)
	if err != nil {
		return nil, err
	}

	converted := convertFromNavigation(entry)

	return &dashboard.NavigationResponse{
		Navigation: converted,
	}, nil

}

// Register register a plugin.
func (s *GRPCServer) Register(ctx context.Context, registerRequest *dashboard.RegisterRequest) (*dashboard.RegisterResponse, error) {
	m, err := s.Impl.Register(ctx, registerRequest.DashboardAPIAddress)
	if err != nil {
		return nil, err
	}

	capabilities := convertFromCapabilities(m.Capabilities)

	return &dashboard.RegisterResponse{
		PluginName:   m.Name,
		Description:  m.Description,
		Capabilities: capabilities,
	}, nil
}

// Print prints an object.
func (s *GRPCServer) Print(ctx context.Context, objectRequest *dashboard.ObjectRequest) (*dashboard.PrintResponse, error) {
	u, err := decodeObjectRequest(objectRequest)
	if err != nil {
		return nil, err
	}

	ctx = ocontext.WithWebsocketClientID(ctx, objectRequest.ClientID)
	pr, err := s.Impl.Print(ctx, u)
	if err != nil {
		return nil, errors.Wrap(err, "grpc server print")
	}

	itemBytes, err := json.Marshal(pr.Items)
	if err != nil {
		return nil, err
	}

	configSection, err := convertFromSummarySections(pr.Config)
	if err != nil {
		return nil, err
	}

	statusSection, err := convertFromSummarySections(pr.Status)
	if err != nil {
		return nil, err
	}

	out := &dashboard.PrintResponse{
		Config: configSection,
		Status: statusSection,
		Items:  itemBytes,
	}

	return out, nil
}

// ObjectStatus generates status for an object.
func (s *GRPCServer) ObjectStatus(ctx context.Context, objectRequest *dashboard.ObjectRequest) (*dashboard.ObjectStatusResponse, error) {
	u, err := decodeObjectRequest(objectRequest)
	if err != nil {
		return nil, err
	}

	ctx = ocontext.WithWebsocketClientID(ctx, objectRequest.ClientID)
	osr, err := s.Impl.ObjectStatus(ctx, u)
	if err != nil {
		return nil, errors.Wrap(err, "grpc server object status")
	}

	objectStatusBytes, err := json.Marshal(osr.ObjectStatus)
	if err != nil {
		return nil, err
	}

	out := &dashboard.ObjectStatusResponse{
		ObjectStatus: objectStatusBytes,
	}

	return out, nil
}

func decodeObjectRequest(req *dashboard.ObjectRequest) (*unstructured.Unstructured, error) {
	m := map[string]interface{}{}

	err := json.Unmarshal(req.Object, &m)
	if err != nil {
		return nil, err
	}

	u := &unstructured.Unstructured{Object: m}
	return u, nil
}

// PrintTab prints a tab for an object.
func (s *GRPCServer) PrintTab(ctx context.Context, objectRequest *dashboard.ObjectRequest) (*dashboard.PrintTabResponse, error) {
	u, err := decodeObjectRequest(objectRequest)
	if err != nil {
		return nil, err
	}

	ctx = ocontext.WithWebsocketClientID(ctx, objectRequest.ClientID)
	tabResponse, err := s.Impl.PrintTab(ctx, u)
	if err != nil {
		return nil, errors.Wrap(err, "grpc server print tab")
	}

	if tabResponse.Tab == nil {
		return nil, errors.New("tab is nil")
	}

	layoutBytes, err := json.Marshal(tabResponse.Tab.Contents)
	if err != nil {
		return nil, err
	}

	out := &dashboard.PrintTabResponse{
		Name:   tabResponse.Tab.Name,
		Layout: layoutBytes,
	}

	return out, nil
}

// WatchAdd is called when a watched GVK has a new object added.
func (s *GRPCServer) WatchAdd(context.Context, *dashboard.WatchRequest) (*dashboard.Empty, error) {
	panic("not implemented")
}

// WatchUpdate is called when a watched GVK has an object updated.
func (s *GRPCServer) WatchUpdate(context.Context, *dashboard.WatchRequest) (*dashboard.Empty, error) {
	panic("not implemented")
}

// WatchDelete is called when a watched GVK has an object deleted.
func (s *GRPCServer) WatchDelete(context.Context, *dashboard.WatchRequest) (*dashboard.Empty, error) {
	panic("not implemented")
}
