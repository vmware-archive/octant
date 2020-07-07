/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package octant

import (
	"context"

	"github.com/pkg/errors"
	"github.com/vmware-tanzu/octant/internal/portforward"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/log"
	"github.com/vmware-tanzu/octant/pkg/store"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

// PortForward creates a port forwarder
type PortForward struct {
	logger        log.Logger
	objectStore   store.Store
	portForwarder portforward.PortForwarder
}

var _ action.Dispatcher = (*PortForward)(nil)

// NewPortForward creates an instance of PortForward
func NewPortForward(logger log.Logger, objectStore store.Store, portForwarder portforward.PortForwarder) *PortForward {
	return &PortForward{
		logger:        logger,
		objectStore:   objectStore,
		portForwarder: portForwarder,
	}
}

// ActionName returns the name of this action
func (p *PortForward) ActionName() string {
	return "overview/startPortForward"
}

// Handle starts a port forward
func (p *PortForward) Handle(ctx context.Context, alerter action.Alerter, payload action.Payload) error {
	p.logger.With("payload", payload).Debugf("received action payload")
	request, err := portForwardRequestFromPayload(payload)
	if err != nil {
		return errors.Wrap(err, "convert payload to port forward request")
	}
	p.logger.Debugf("%s", request)

	_, err = p.portForwarder.Create(ctx, request.gvk(), request.Name, request.Namespace, request.Port)
	if err != nil {
		return errors.Wrap(err, "create port forwarder")
	}
	return nil
}

// PortForwardDelete stops a port forwarder
type PortForwardDelete struct {
	logger        log.Logger
	objectStore   store.Store
	portForwarder portforward.PortForwarder
}

var _ action.Dispatcher = (*PortForwardDelete)(nil)

// NewPortForwardDelete creates an instance of PortForwardDelete
func NewPortForwardDelete(logger log.Logger, objectStore store.Store, portForwarder portforward.PortForwarder) *PortForwardDelete {
	return &PortForwardDelete{
		logger:        logger,
		objectStore:   objectStore,
		portForwarder: portForwarder,
	}
}

// ActionName returns the name of this action
func (p *PortForwardDelete) ActionName() string {
	return "overview/stopPortForward"
}

// Handle stops a port forward
func (p *PortForwardDelete) Handle(ctx context.Context, alerter action.Alerter, payload action.Payload) error {
	p.logger.With("payload", payload).Debugf("received action payload")
	id, err := payload.String("id")
	if err != nil {
		return errors.Wrap(err, "convert payload to stop port forward request")
	}

	p.portForwarder.StopForwarder(id)
	return nil
}

type portForwardCreateRequest struct {
	APIVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Name       string `json:"name,omitempty"`
	Namespace  string `json:"namespace,omitempty"`
	Port       uint16 `json:"port,omitempty"`
}

func (req *portForwardCreateRequest) Validate() error {
	if req.APIVersion != "v1" && req.Kind == "Pod" {
		return errors.New("only supports forwards for v1 Pods")
	}

	if req.Name == "" {
		return errors.New("pod name is blank")
	}

	if req.Namespace == "" {
		return errors.New("pod namespace is blank")
	}

	if req.Port < 1 {
		return errors.New("port must be greater than 0")
	}

	return nil
}

func (req *portForwardCreateRequest) gvk() schema.GroupVersionKind {
	return schema.FromAPIVersionAndKind(req.APIVersion, req.Kind)
}

func portForwardRequestFromPayload(payload action.Payload) (*portForwardCreateRequest, error) {
	apiVersion, err := payload.String("apiVersion")
	if err != nil {
		return nil, err
	}

	kind, err := payload.String("kind")
	if err != nil {
		return nil, err
	}

	name, err := payload.String("name")
	if err != nil {
		return nil, err
	}

	namespace, err := payload.String("namespace")
	if err != nil {
		return nil, err
	}

	port, err := payload.Uint16("port")
	if err != nil {
		return nil, err
	}

	req := &portForwardCreateRequest{
		APIVersion: apiVersion,
		Kind:       kind,
		Name:       name,
		Namespace:  namespace,
		Port:       port,
	}

	if err := req.Validate(); err != nil {
		return nil, err
	}

	return req, nil
}
