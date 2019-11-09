/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package clusteroverview

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/vmware-tanzu/octant/internal/api"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/terminal"
	"github.com/vmware-tanzu/octant/pkg/action"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type terminalCreateRequest struct {
	APIVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Name       string `json:"name,omitempty"`
	Namespace  string `json:"namespace,omitempty"`
	Container  string `json:"container,omitempty"`
	Command    string `json:"command,omitempty"`
}

func (req *terminalCreateRequest) gvk() schema.GroupVersionKind {
	return schema.FromAPIVersionAndKind(req.APIVersion, req.Kind)
}

func (req *terminalCreateRequest) Validate() error {
	if req.APIVersion != "v1" && req.Kind == "Pod" {
		return errors.New("only supports terminals for v1 Pods")
	}

	if req.Name == "" {
		return errors.New("pod name is blank")
	}

	if req.Namespace == "" {
		return errors.New("pod namespace is blank")
	}

	if req.Container == "" {
		return errors.New("pod container is blank")
	}

	if req.Command == "" {
		return errors.New("terminal command is blank")
	}

	return nil
}

func terminalRequestFromPayload(payload action.Payload) (*terminalCreateRequest, error) {
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

	container, err := payload.String("container")
	if err != nil {
		return nil, err
	}

	command, err := payload.String("command")
	if err != nil {
		return nil, err
	}

	req := &terminalCreateRequest{
		APIVersion: apiVersion,
		Kind:       kind,
		Name:       name,
		Namespace:  namespace,
		Container:  container,
		Command:    command,
	}

	if err := req.Validate(); err != nil {
		return nil, err
	}

	return req, nil
}

type terminalError struct {
	code     int
	message  string
	extraErr error
}

var _ error = (*terminalError)(nil)

func (e *terminalError) Error() string {
	return e.message
}

type terminalHandler struct {
	manager terminal.Manager
	logger  log.Logger
}

var _ http.Handler = (*terminalHandler)(nil)

func newTerminalHandler(logger log.Logger, manager terminal.Manager) (*terminalHandler, error) {
	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	if manager == nil {
		return nil, errors.New("terminal service is nil")
	}

	return &terminalHandler{
		manager: manager,
	}, nil
}

func (h *terminalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := log.WithLoggerContext(r.Context(), h.logger)

	defer r.Body.Close()

	switch r.Method {
	case http.MethodPost:
		err := createTerminal(ctx, r.Body, h.manager, w)
		handleTerminalError(w, err, h.logger)
	default:
		api.RespondWithError(
			w,
			http.StatusNotFound,
			fmt.Sprintf("unhandled HTTP method %s for terminals", r.Method),
			h.logger,
		)
	}
}

func createTerminal(ctx context.Context, body io.Reader, terminalManager terminal.Manager, w http.ResponseWriter) error {
	return nil
}

func handleTerminalError(w http.ResponseWriter, err error, logger log.Logger) {}
