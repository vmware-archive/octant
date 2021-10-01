/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/view/component"

	"google.golang.org/grpc/metadata"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	ocontext "github.com/vmware-tanzu/octant/internal/context"
	"github.com/vmware-tanzu/octant/pkg/event"

	"github.com/vmware-tanzu/octant/internal/gvk"
	"github.com/vmware-tanzu/octant/internal/portforward"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/cluster"
	"github.com/vmware-tanzu/octant/pkg/plugin/api/proto"
	"github.com/vmware-tanzu/octant/pkg/store"
)

// DashboardMetadataKey is a type used for metadata keys passed by plugins
type DashboardMetadataKey string

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

// NamespacesResponse is a response from listing namespaces
type NamespacesResponse struct {
	Namespaces []string
}

type LinkResponse struct {
	Link component.Link
}

// Service is the dashboard service.
type Service interface {
	List(ctx context.Context, key store.Key) (*unstructured.UnstructuredList, error)
	Get(ctx context.Context, key store.Key) (*unstructured.Unstructured, error)
	PortForward(ctx context.Context, req PortForwardRequest) (PortForwardResponse, error)
	CancelPortForward(ctx context.Context, id string)
	ListNamespaces(ctx context.Context) (NamespacesResponse, error)
	Update(ctx context.Context, object *unstructured.Unstructured) error
	Create(ctx context.Context, object *unstructured.Unstructured) error
	ApplyYAML(ctx context.Context, namespace, yaml string) ([]string, error)
	Delete(ctx context.Context, key store.Key) error
	ForceFrontendUpdate(ctx context.Context) error
	SendAlert(ctx context.Context, clientID string, alert action.Alert) error
	CreateLink(ctx context.Context, key store.Key) (LinkResponse, error)
	SendEvent(ctx context.Context, clientID string, eventName event.EventType, payload action.Payload) error
}

// FrontendUpdateController can control the frontend. ie. the web gui
type FrontendUpdateController interface {
	ForceUpdate() error
}

// FrontendProxy is a proxy for messaging the frontend.
type FrontendProxy struct {
	FrontendUpdateController FrontendUpdateController
}

// ForceFrontendUpdate forces the frontend to update
func (proxy *FrontendProxy) ForceFrontendUpdate() error {
	if proxy.FrontendUpdateController == nil {
		return nil
	}

	return proxy.FrontendUpdateController.ForceUpdate()
}

// GRPCService is an implementation of the dashboard service based on GRPC.
type GRPCService struct {
	ObjectStore            store.Store
	PortForwarder          portforward.PortForwarder
	FrontendProxy          FrontendProxy
	NamespaceInterface     cluster.NamespaceInterface
	WebsocketClientManager event.WSClientGetter
	LinkGenerator          octant.LinkGenerator
}

var _ Service = (*GRPCService)(nil)

// List lists objects.
func (s *GRPCService) List(ctx context.Context, key store.Key) (*unstructured.UnstructuredList, error) {
	// TODO: support hasSynced
	ctx = extractObjectStoreMetadata(ctx)
	list, _, err := s.ObjectStore.List(ctx, key)
	return list, err
}

// Get retrieves an object.
func (s *GRPCService) Get(ctx context.Context, key store.Key) (*unstructured.Unstructured, error) {
	ctx = extractObjectStoreMetadata(ctx)
	return s.ObjectStore.Get(ctx, key)
}

func (s *GRPCService) Update(ctx context.Context, object *unstructured.Unstructured) error {
	key, err := store.KeyFromObject(object)
	if err != nil {
		return err
	}

	ctx = extractObjectStoreMetadata(ctx)
	return s.ObjectStore.Update(ctx, key, func(u *unstructured.Unstructured) error {
		u.Object = object.Object
		return nil
	})
}

func (s *GRPCService) Create(ctx context.Context, object *unstructured.Unstructured) error {
	ctx = extractObjectStoreMetadata(ctx)
	return s.ObjectStore.Create(ctx, object)
}

func (s *GRPCService) ApplyYAML(ctx context.Context, namespace, yaml string) ([]string, error) {
	ctx = extractObjectStoreMetadata(ctx)
	return s.ObjectStore.CreateOrUpdateFromYAML(ctx, namespace, yaml)
}

func (s *GRPCService) Delete(ctx context.Context, key store.Key) error {
	ctx = extractObjectStoreMetadata(ctx)
	return s.ObjectStore.Delete(ctx, key)
}

// PortForward creates a port forward.
func (s *GRPCService) PortForward(ctx context.Context, req PortForwardRequest) (PortForwardResponse, error) {
	pfResponse, err := s.PortForwarder.Create(ctx, nil, gvk.Pod, req.PodName, req.Namespace, req.Port)
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

// ListNamespaces lists namespaces
func (s *GRPCService) ListNamespaces(ctx context.Context) (NamespacesResponse, error) {
	namespaces, err := s.NamespaceInterface.Names(ctx)
	if err != nil {
		return NamespacesResponse{}, err
	}

	resp := NamespacesResponse{
		Namespaces: namespaces,
	}
	return resp, nil
}

func (s *GRPCService) ForceFrontendUpdate(ctx context.Context) error {
	return s.FrontendProxy.ForceFrontendUpdate()
}

// SendAlert sends an alert. This method is deprecated and will be removed in a future release
func (s *GRPCService) SendAlert(ctx context.Context, clientID string, alert action.Alert) error {
	return nil
}

func (s *GRPCService) SendEvent(ctx context.Context, clientID string, eventName event.EventType, payload action.Payload) error {
	event := event.CreateEvent(eventName, payload)

	if s.WebsocketClientManager == nil {
		return fmt.Errorf("websocket client manager is nil")
	}

	if clientID == "" {
		return fmt.Errorf("no websocket client id")
	}

	sender := s.WebsocketClientManager.Get(clientID)
	if sender == nil {
		clientID := ocontext.ClientStateFrom(ctx).ClientID
		return fmt.Errorf("unable to find ws client %s", clientID)
	}
	sender.Send(event)
	return nil
}

func (s *GRPCService) CreateLink(_ context.Context, key store.Key) (LinkResponse, error) {
	ref, err := s.LinkGenerator.ObjectPath(key.Namespace, key.APIVersion, key.Kind, key.Name)
	if err != nil {
		return LinkResponse{}, err
	}

	linkComponent := component.NewLink("", "", ref)
	return LinkResponse{
		Link: *linkComponent,
	}, nil
}

func NewGRPCServer(service Service) *grpcServer {
	return &grpcServer{
		service: service,
	}
}

type grpcServer struct {
	service Service
	proto.UnimplementedDashboardServer
}

var _ proto.DashboardServer = (*grpcServer)(nil)

// List lists objects.
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

// Get gets an object.
func (c *grpcServer) Get(ctx context.Context, in *proto.KeyRequest) (*proto.GetResponse, error) {
	key, err := convertToKey(in)
	if err != nil {
		return nil, err
	}

	object, err := c.service.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var out *proto.GetResponse

	if object != nil {
		encodedObject, err := convertFromObject(object)
		if err != nil {
			return nil, err
		}

		out = &proto.GetResponse{
			Object: encodedObject,
		}
	} else {
		return &proto.GetResponse{}, nil
	}

	return out, nil
}

// Update updates an object.
func (c *grpcServer) Update(ctx context.Context, in *proto.UpdateRequest) (*proto.UpdateResponse, error) {
	object, err := convertToObject(in.Object)
	if err != nil {
		return nil, err
	}

	if object == nil {
		return &proto.UpdateResponse{}, fmt.Errorf("can't update an object that doesn't exist")
	}

	if err := c.service.Update(ctx, object); err != nil {
		return nil, err
	}

	return &proto.UpdateResponse{}, nil
}

// ApplyYAML updates an object.
func (c *grpcServer) ApplyYAML(ctx context.Context, in *proto.ApplyYAMLRequest) (*proto.ApplyYAMLResponse, error) {
	var res []string
	var err error

	if res, err = c.service.ApplyYAML(ctx, in.Namespace, in.Yaml); err != nil {
		return nil, err
	}

	return &proto.ApplyYAMLResponse{
		Resources: res,
	}, nil
}

// Create creates an object in the cluster.
func (c *grpcServer) Create(ctx context.Context, in *proto.CreateRequest) (*proto.CreateResponse, error) {
	object, err := convertToObject(in.Object)
	if err != nil {
		return nil, err
	}

	if object == nil {
		return nil, fmt.Errorf("unable to create a nil object")
	}

	if err := c.service.Create(ctx, object); err != nil {
		return nil, err
	}

	return &proto.CreateResponse{}, nil
}

// Delete deletes an object.
func (c *grpcServer) Delete(ctx context.Context, in *proto.KeyRequest) (*proto.DeleteResponse, error) {
	key, err := convertToKey(in)
	if err != nil {
		return nil, err
	}

	err = c.service.Delete(ctx, key)
	if err != nil {
		return nil, err
	}

	return &proto.DeleteResponse{}, nil
}

// PortForward creates a port forward.
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

// CancelPortForward cancels a port forward.
func (c *grpcServer) CancelPortForward(ctx context.Context, in *proto.CancelPortForwardRequest) (*proto.Empty, error) {
	if in == nil {
		return nil, fmt.Errorf("request is nil")
	}

	c.service.CancelPortForward(ctx, in.PortForwardID)
	return &proto.Empty{}, nil
}

// ListNamespaces lists namespaces.
func (c *grpcServer) ListNamespaces(ctx context.Context, _ *proto.Empty) (*proto.NamespacesResponse, error) {
	nsResp, err := c.service.ListNamespaces(ctx)
	if err != nil {
		return nil, err
	}

	resp := &proto.NamespacesResponse{
		Namespaces: nsResp.Namespaces,
	}
	return resp, nil
}

// ForceFrontendUpdate forces the front end to update.
func (c *grpcServer) ForceFrontendUpdate(ctx context.Context, _ *proto.Empty) (*proto.Empty, error) {
	if err := c.service.ForceFrontendUpdate(ctx); err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

// SendAlert sends an alert. This method is deprecated and will be removed in a future release
func (c *grpcServer) SendAlert(ctx context.Context, in *proto.AlertRequest) (*proto.Empty, error) {
	alert, err := convertToAlert(in)
	if err != nil {
		return nil, err
	}

	c.service.SendEvent(ctx, in.ClientID, event.EventTypeAlert, action.Payload{
		"type":       alert.Type,
		"message":    alert.Message,
		"expiration": alert.Expiration,
	})
	return &proto.Empty{}, nil
}

func (c *grpcServer) CreateLink(ctx context.Context, in *proto.KeyRequest) (*proto.LinkResponse, error) {
	key, err := convertToKey(in)
	if err != nil {
		return nil, err
	}

	resp, err := c.service.CreateLink(ctx, key)
	if err != nil {
		return nil, err
	}

	var ref string
	if &resp.Link != nil && &resp.Link.Config != nil {
		ref = resp.Link.Config.Ref
	}

	return &proto.LinkResponse{
		Ref: ref,
	}, nil
}

func extractObjectStoreMetadata(ctx context.Context) context.Context {
	// Second value is ignored as metadata is always set by grpc.
	md, _ := metadata.FromIncomingContext(ctx)
	for k, v := range md {
		if strings.HasPrefix(k, "x-octant-") {
			ctxKey := strings.Replace(k, "x-octant-", "", 1)
			ctx = context.WithValue(ctx, DashboardMetadataKey(ctxKey), v[0])
		}
	}

	return ctx
}

func (c *grpcServer) SendEvent(ctx context.Context, in *proto.EventRequest) (*proto.EventResponse, error) {
	clientID, eventName, payload, err := convertToEvent(in)
	if err != nil {
		return nil, err
	}

	err = c.service.SendEvent(ctx, clientID, event.EventType(eventName), payload)
	return &proto.EventResponse{}, err
}
