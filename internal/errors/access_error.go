/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package errors

import (
	"fmt"
	"time"

	"github.com/vmware/octant/pkg/store"
)

type AccessError struct {
	id        string
	key       store.Key
	timestamp time.Time
	verb      string
	err       error
}

var _ InternalError = (*AccessError)(nil)

func NewAccessError(key store.Key, verb string, err error) *AccessError {
	return &AccessError{
		verb:      verb,
		err:       err,
		timestamp: time.Now(),
		id:        key.String(),
		key:       key,
	}
}

// ID returns the error unique ID.
func (o *AccessError) ID() string {
	return o.id
}

// Timestamp returns the error timestamp.
func (o *AccessError) Timestamp() time.Time {
	return o.timestamp
}

// Error returns an error string.
func (o *AccessError) Error() string {
	return fmt.Sprintf("%s: %s: %s", o.verb, o.key, o.err)
}

// Key returns the key for the error.
func (o *AccessError) Key() store.Key {
	return o.key
}

// Verb returns the verb for the error.
func (o *AccessError) Verb() string {
	return o.verb
}
