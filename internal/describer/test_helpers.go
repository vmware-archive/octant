/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"

	"github.com/vmware/octant/pkg/view/component"
)

type StubDescriber struct {
	path       string
	components []component.Component
}

var _ Describer = (*StubDescriber)(nil)

func NewStubDescriber(p string, components ...component.Component) *StubDescriber {
	return &StubDescriber{
		path:       p,
		components: components,
	}
}
func (d *StubDescriber) Describe(context.Context, string, string, Options) (component.ContentResponse, error) {
	return component.ContentResponse{
		Components: d.components,
	}, nil
}

func (d *StubDescriber) PathFilters() []PathFilter {
	return []PathFilter{
		*NewPathFilter(d.path, d),
	}
}

func (d *StubDescriber) Reset(ctx context.Context) error {
	panic("implement me")
}

func NewEmptyDescriber(p string) *StubDescriber {
	return &StubDescriber{
		path: p,
	}
}
