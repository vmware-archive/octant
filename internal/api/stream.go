/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"fmt"
	"path"
	"strings"
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
	return fmt.Sprintf("Not found: %s", e.path)
}

func notFoundRedirectPath(requestPath string) string {
	if strings.HasSuffix(requestPath, "/") {
		// ignore trailing slash
		requestPath = requestPath[0 : len(requestPath)-1]
	}
	parts := strings.Split(requestPath, "/")
	if len(parts) < 2 {
		return ""
	}
	return path.Join(parts[0 : len(parts)-1]...)
}
