/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package errors

import (
	"sort"

	lru "github.com/hashicorp/golang-lru"

	"github.com/vmware-tanzu/octant/pkg/errors"
)

const maxErrors = 50

type errorStore struct {
	recentErrors *lru.Cache
}

var _ errors.ErrorStore = (*errorStore)(nil)

// NewErrorStore creates a new error store.
func NewErrorStore() (errors.ErrorStore, error) {
	cache, err := lru.New(maxErrors)
	if err != nil {
		return nil, err
	}

	return &errorStore{
		recentErrors: cache,
	}, nil
}

// Get returns an InternalError and if it was found in the store.
func (e *errorStore) Get(id string) (err errors.InternalError, found bool) {
	v, ok := e.recentErrors.Peek(id)
	if !ok {
		return nil, ok
	}
	return v.(errors.InternalError), ok
}

// AddError coverts a standard Go error to an InternalError and adds it to the error store.
func (e *errorStore) AddError(err error) errors.InternalError {
	intErr := convertError(err)
	e.Add(intErr)
	return intErr
}

// Add adds an InternalError directly to the error store.
func (e *errorStore) Add(intErr errors.InternalError) (found bool) {
	ok, _ := e.recentErrors.ContainsOrAdd(intErr.ID(), intErr)
	return ok
}

// List returns a list of all the error objects in the store from newest to oldest.
func (e *errorStore) List() []errors.InternalError {
	var intErrList []errors.InternalError
	for _, key := range e.recentErrors.Keys() {
		v, ok := e.recentErrors.Peek(key)
		if !ok {
			continue
		}
		intErrList = append(intErrList, v.(errors.InternalError))
	}
	sort.Slice(intErrList, func(i, j int) bool {
		return intErrList[i].Timestamp().After(intErrList[j].Timestamp())
	})
	return intErrList
}

func convertError(err error) errors.InternalError {
	return errors.NewGenericError(err)
}
