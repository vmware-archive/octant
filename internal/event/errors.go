/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package event

import "github.com/pkg/errors"

type notFound interface {
	NotFound() bool
	Path() string
}

var (
	errNotReady = errors.New("not ready")
)
