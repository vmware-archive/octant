/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/vmware/octant/internal/octant"
)

// NotFoundError is a not found error.
type NotFoundError struct {
	path string
}

// NewNotFoundError creates an instance of NotFoundError
func NewNotFoundError(path string) *NotFoundError {
	return &NotFoundError{path: path}
}

// Path is the path of the error.
func (e *NotFoundError) Path() string {
	return e.path
}

// NotFound returns true to signify this is a not found error.
func (e *NotFoundError) NotFound() bool { return true }

// Error returns the error string.
func (e *NotFoundError) Error() string {
	return "Not found"
}

type eventSourceStreamer struct {
	w http.ResponseWriter
}

func (s *eventSourceStreamer) Stream(ctx context.Context, ch <-chan octant.Event) {
	flusher, ok := s.w.(http.Flusher)
	if !ok {
		http.Error(s.w, "server sent events are unsupported", http.StatusInternalServerError)
		return
	}

	s.w.Header().Set("Content-Type", "text/event-stream")
	s.w.Header().Set("Cache-Control", "no-cache")
	s.w.Header().Set("Connection", "keep-alive")
	s.w.Header().Set("Access-Control-Allow-Origin", "*")

	isStreaming := true

	for isStreaming {
		select {
		case <-ctx.Done():
			isStreaming = false
		case e := <-ch:
			if e.Type != "" {
				_, _ = fmt.Fprintf(s.w, "event: %s\n", e.Type)
			}
			_, _ = fmt.Fprintf(s.w, "data: %s\n\n", string(e.Data))
			flusher.Flush()
		}
	}
}

func notFoundRedirectPath(requestPath string) string {
	parts := strings.Split(requestPath, "/")
	if len(parts) < 5 {
		return ""
	}
	return path.Join(append([]string{"/"}, parts[3:len(parts)-2]...)...)
}
