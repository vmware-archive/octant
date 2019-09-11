/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/vmware/octant/internal/log"
)

var (
	ErrPathNotFound = errors.New("path not found")
)

type PathMatcher struct {
	name    string
	filters map[string]PathFilter

	sync.Mutex
}

func NewPathMatcher(name string) *PathMatcher {
	return &PathMatcher{
		name:    name,
		filters: make(map[string]PathFilter),
	}
}

func (pm *PathMatcher) Register(ctx context.Context, pf PathFilter) {
	logger := log.From(ctx)

	pm.Lock()
	defer pm.Unlock()

	logger.With(
		"name", pm.name,
		"path", pf.filterPath,
	).Debugf("register path")
	pm.filters[pf.filterPath] = pf
}

func (pm *PathMatcher) Deregister(ctx context.Context, paths ...string) {
	logger := log.From(ctx)

	pm.Lock()
	defer pm.Unlock()

	for _, p := range paths {
		logger.With(
			"name", pm.name,
			"path", p,
		).Debugf("deregister path")
		delete(pm.filters, p)
	}

}

func (pm *PathMatcher) Find(path string) (PathFilter, error) {
	pm.Lock()
	defer pm.Unlock()

	for _, pf := range pm.filters {
		if pf.Match(path) {
			return pf, nil
		}
	}

	return PathFilter{}, ErrPathNotFound
}
