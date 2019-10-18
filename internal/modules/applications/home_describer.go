/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package applications

import (
	"context"

	"github.com/pkg/errors"

	"github.com/vmware-tanzu/octant/internal/describer"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// HomeDescriberOption is an option for configuring HomeDescriber.
type HomeDescriberOption func(d *HomeDescriber)

// WithHomeDescriberSummarizer configures the Summarizer for HomeDescriber.
func WithHomeDescriberSummarizer(s Summarizer) HomeDescriberOption {
	return func(d *HomeDescriber) {
		d.summarizer = s
	}
}

// HomeDescriber describes content for applications.
type HomeDescriber struct {
	summarizer Summarizer
}

var _ describer.Describer = (*HomeDescriber)(nil)

// NewHomeDescriber creates an instance of HomeDescriber.
func NewHomeDescriber(options ...HomeDescriberOption) *HomeDescriber {
	d := &HomeDescriber{}

	for _, option := range options {
		option(d)
	}

	if d.summarizer == nil {
		d.summarizer = &summarizer{}
	}

	return d
}

// Describe prints a summary of applications.
func (l *HomeDescriber) Describe(ctx context.Context, namespace string, options describer.Options) (component.ContentResponse, error) {
	table, err := l.summarizer.Summarize(ctx, namespace, options)
	if err != nil {
		return component.EmptyContentResponse, errors.Wrap(err, "summarize applications")
	}

	contentResponse := component.ContentResponse{
		Title:      component.TitleFromString("Applications"),
		Components: []component.Component{table},
		IconName:   "",
		IconSource: "",
	}

	return contentResponse, nil
}

// PathFilters returns PathFilters for this describer. It is the root of a the module.
func (l *HomeDescriber) PathFilters() []describer.PathFilter {
	return []describer.PathFilter{
		*describer.NewPathFilter("/", l),
	}
}

// Reset does nothing.
func (l HomeDescriber) Reset(ctx context.Context) error {
	return nil
}
