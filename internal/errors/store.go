/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package errors

import (
	"sort"

	lru "github.com/hashicorp/golang-lru"
)

const maxErrors = 50

// Observer Initial implementation of observer pattern
type Observer interface {
	Update()
}

type ErrorStore interface {
	List() []InternalError
	Get(string) (InternalError, bool)
	Add(InternalError) (found bool)
	AddError(error) InternalError
	Subscribe(observer Observer)
}

type errorStore struct {
	recentErrors *lru.Cache
	subscribers  []Observer
}

var _ ErrorStore = (*errorStore)(nil)

// NewErrorStore creates a new error store.
func NewErrorStore() (ErrorStore, error) {
	cache, err := lru.New(maxErrors)
	if err != nil {
		return nil, err
	}

	return &errorStore{
		recentErrors: cache,
	}, nil
}

// Get returns an InternalError and if it was found in the store.
func (e *errorStore) Get(id string) (err InternalError, found bool) {
	v, ok := e.recentErrors.Peek(id)
	if !ok {
		return nil, ok
	}
	return v.(InternalError), ok
}

// AddError coverts a standard Go error to an InternalError and adds it to the error store.
func (e *errorStore) AddError(err error) InternalError {
	intErr := convertError(err)
	e.Add(intErr)
	return intErr
}

// Add adds an InternalError directly to the error store.
func (e *errorStore) Add(intErr InternalError) (found bool) {
	ok, _ := e.recentErrors.ContainsOrAdd(intErr.ID(), intErr)

	if !ok {
		for i := 0; i < len(e.subscribers); i++ {
			e.subscribers[i].Update()
		}
	}

	return ok
}

// List returns a list of all the error objects in the store from newest to oldest.
func (e *errorStore) List() []InternalError {
	var intErrList []InternalError
	for _, key := range e.recentErrors.Keys() {
		v, ok := e.recentErrors.Peek(key)
		if !ok {
			continue
		}
		intErrList = append(intErrList, v.(InternalError))
	}
	sort.Slice(intErrList, func(i, j int) bool {
		return intErrList[i].Timestamp().After(intErrList[j].Timestamp())
	})
	return intErrList
}

func (e *errorStore) Subscribe(manager Observer) {
	e.subscribers = append(e.subscribers, manager)
}

func convertError(err error) InternalError {
	return NewGenericError(err)
}
