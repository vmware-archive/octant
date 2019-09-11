/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package octant

//go:generate mockgen -destination=./fake/mock_state.go -package=fake github.com/vmware/octant/internal/octant State

// UpdateCancelFunc cancels the update.
type UpdateCancelFunc func()

// State represents Octant's view state.
type State interface {
	// SetContentPath sets the content path.
	SetContentPath(string)
	// GetContentPath returns the content path.
	GetContentPath() string
	// OnNamespaceUpdate registers a function to be called with the content path
	// is changed.
	OnContentPathUpdate(fn ContentPathUpdateFunc) UpdateCancelFunc
	// SetNamespace sets the namespace.
	SetNamespace(namespace string)
	// GetNamespace returns the namespace.
	GetNamespace() string
	// OnNamespaceUpdate returns a function to be called when the namespace
	// is changed.
	OnNamespaceUpdate(fun NamespaceUpdateFunc) UpdateCancelFunc
	// AddFilter adds a label to filtered.
	AddFilter(filter Filter)
	// RemoveFilter removes a filter.
	RemoveFilter(filter Filter)
	// GetFilters returns a slice of filters.
	GetFilters() []Filter
	// SetFilters replaces the current filters with a slice of filters.
	// The slice can be empty.
	SetFilters(filters []Filter)
	// SetContext sets the current context.
	SetContext(requestedContext string)
}

// ContentPathUpdateFunc is a function that is called when content path is updated.
type ContentPathUpdateFunc func(contentPath string)

// NamespaceUpdateFunc is a function that is called when namespace is updated.
type NamespaceUpdateFunc func(namespace string)
