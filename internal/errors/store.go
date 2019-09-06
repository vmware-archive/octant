package errors

import (
	"sort"

	"github.com/google/uuid"
	lru "github.com/hashicorp/golang-lru"
	"github.com/vmware/octant/pkg/action"
)

const maxErrors = 50

type ErrorStore interface {
	List() []*InternalError
	Get(uuid.UUID) (*InternalError, bool)
	Add(*InternalError)
	NewError(string, action.Payload, error) *InternalError
}

func NewErrorStore() (ErrorStore, error) {
	cache, err := lru.New(maxErrors)
	if err != nil {
		return nil, err
	}

	return &errorStore{
		recentErrors: cache,
	}, nil
}

type errorStore struct {
	recentErrors *lru.Cache
}

// NewError creates a new InternalError and adds it to the store.
func (e *errorStore) NewError(requestType string, payload action.Payload, err error) *InternalError {
	intErr := NewInternalError(requestType, payload, err)
	e.Add(intErr)
	return intErr
}

// Get returns an InternalError and if it was found in the store.
func (e *errorStore) Get(id uuid.UUID) (err *InternalError, found bool) {
	v, ok := e.recentErrors.Peek(id)
	if !ok {
		return nil, ok
	}
	return v.(*InternalError), ok
}

// Add adds a new InternalError directly in to the store.
func (e *errorStore) Add(intErr *InternalError) {
	e.recentErrors.ContainsOrAdd(intErr.id, intErr)
}

// List returns a list of all the IntrenalError objects in the store from newest to oldest.
func (e *errorStore) List() []*InternalError {
	var intErrList []*InternalError
	for _, key := range e.recentErrors.Keys() {
		v, ok := e.recentErrors.Peek(key)
		if !ok {
			continue
		}
		intErrList = append(intErrList, v.(*InternalError))
	}
	sort.Slice(intErrList, func(i, j int) bool {
		return intErrList[i].timestamp.After(intErrList[j].timestamp)
	})
	return intErrList
}
