/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package clusteroverview

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/vmware/octant/internal/api"
	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/internal/mime"
	"github.com/vmware/octant/internal/portforward"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type portForwardCreateRequest struct {
	APIVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Type       string `json:"type,omitempty"`
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

type portForwardError struct {
	code     int
	message  string
	extraErr error
}

var _ error = (*portForwardError)(nil)

func (e *portForwardError) Error() string {
	return e.message
}

type portForwardsHandler struct {
	svc    portforward.PortForwarder
	logger log.Logger
}

var _ http.Handler = (*portForwardsHandler)(nil)

func newPortForwardsHandler(logger log.Logger, svc portforward.PortForwarder) (*portForwardsHandler, error) {
	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	if svc == nil {
		return nil, errors.New("port forward service is nil")
	}

	return &portForwardsHandler{
		svc: svc,
	}, nil
}

func (h *portForwardsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := log.WithLoggerContext(r.Context(), h.logger)

	defer r.Body.Close()

	switch r.Method {
	case http.MethodPost:
		err := createPortForward(ctx, r.Body, h.svc, w)
		handlePortForwardError(w, err, h.logger)
	default:
		api.RespondWithError(
			w,
			http.StatusNotFound,
			fmt.Sprintf("unhandled HTTP method %s for port forwards", r.Method),
			h.logger,
		)
	}
}

func createPortForward(ctx context.Context, body io.Reader, pfs portforward.PortForwarder, w http.ResponseWriter) error {
	if pfs == nil {
		return errors.New("port forward service is nil")
	}
	logger := log.From(ctx)

	req := portForwardCreateRequest{}
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		return &portForwardError{code: http.StatusBadRequest, message: "unable to decode request"}
	}

	if err := req.Validate(); err != nil {
		return &portForwardError{
			code:     http.StatusBadRequest,
			message:  "request is invalid",
			extraErr: err,
		}
	}

	resp, err := pfs.Create(ctx, req.gvk(), req.Name, req.Namespace, req.Port)
	if err != nil {
		return &portForwardError{
			code:     http.StatusInternalServerError,
			message:  "create port forward",
			extraErr: err,
		}
	}

	w.Header().Set("Content-Type", mime.JSONContentType)
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.With("err", err.Error()).Errorf("encoding JSON response")
	}

	return nil
}

func handlePortForwardError(w http.ResponseWriter, err error, logger log.Logger) {
	if err == nil {
		return
	}

	code := http.StatusInternalServerError
	message := err.Error()

	if cause, ok := errors.Cause(err).(*portForwardError); ok {
		code = cause.code
		message = cause.message

	}

	api.RespondWithError(w, code, message, logger)
}
