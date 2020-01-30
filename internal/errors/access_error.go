/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package errors

import (
	"errors"
	"fmt"
	"strings"
	"time"

	kerrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/vmware-tanzu/octant/pkg/store"
)

const OctantAccessError = "AccessError"

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

func (o *AccessError) Name() string {
	return OctantAccessError
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
	e := fmt.Sprintf("%s: %s", o.verb, o.key)
	if o.err != nil {
		e = fmt.Sprintf("%s (error: %s)", e, o.err.Error())
	}
	return e
}

// Key returns the key for the error.
func (o *AccessError) Key() store.Key {
	return o.key
}

// Verb returns the verb for the error.
func (o *AccessError) Verb() string {
	return o.verb
}

func IsBackoffError(err error) bool {
	if err == nil {
		return false
	}

	if kerrors.IsUnauthorized(err) {
		return true
	}

	var ae *AccessError
	if errors.As(err, &ae) {
		if ae.Name() == OctantAccessError {
			return true
		}
	}

	es := err.Error()
	if strings.Contains(es, "Unauthorized") {
		return true
	}

	return false
}
