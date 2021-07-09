/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package errors

import (
	"errors"
	"fmt"
	"hash/fnv"
	"strings"
	"time"

	kerrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/vmware-tanzu/octant/pkg/store"
)

// The ID is using non-cryptographic hash to make sure ID is uniquely identify,
// GoLang provides many algorithms https://golang.org/pkg/hash, fnv was chosen,
// because of the information here
// https://softwareengineering.stackexchange.com/questions/49550/which-hashing-algorithm-is-best-for-uniqueness-and-speed
// fvn provides a good compromise between performance and number of collisions

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
	data := fnv.New64().Sum([]byte(OctantAccessError + ": " + err.Error()))

	return &AccessError{
		verb:      verb,
		err:       err,
		timestamp: time.Now(),
		id:        fmt.Sprintf("%x", data),
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
