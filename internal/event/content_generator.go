/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package event

import (
	"context"
	"encoding/json"
	"time"

	"k8s.io/apimachinery/pkg/labels"

	"github.com/heptio/developer-dash/internal/octant"
	"github.com/heptio/developer-dash/internal/module"
	"github.com/heptio/developer-dash/pkg/view/component"
)

// ContentGenerator generates content events.
type ContentGenerator struct {
	// ResponseFactory is a function that generates a content response.
	ResponseFactory func(ctx context.Context, path, prefix, namespace string, opts module.ContentOptions) (component.ContentResponse, error)

	// Path is the path to the content.
	Path string

	// Prefix the API path prefix. It could be prepended to the path to create
	// a resolvable path.
	Prefix string

	// Namespace is the current namespace.
	Namespace string

	// LabelSet is a label set to filter any content.
	LabelSet *labels.Set

	// RunEvery is how often the event generator should be run.
	RunEvery time.Duration
}

var _ octant.Generator = (*ContentGenerator)(nil)

type dashResponse struct {
	Content component.ContentResponse `json:"content,omitempty"`
}

// Event generates content events.
func (g *ContentGenerator) Event(ctx context.Context) (octant.Event, error) {
	resp, err := g.ResponseFactory(ctx, g.Path, g.Prefix, g.Namespace, module.ContentOptions{LabelSet: g.LabelSet})
	if err != nil {
		return octant.Event{}, err
	}

	dr := dashResponse{
		Content: resp,
	}

	data, err := json.Marshal(dr)
	if err != nil {
		return octant.Event{}, err
	}

	return octant.Event{
		Type: octant.EventTypeContent,
		Data: data,
	}, nil
}

// ScheduleDelay returns how long to delay before running this generator again.
func (g *ContentGenerator) ScheduleDelay() time.Duration {
	return g.RunEvery
}

// Name returns the name of this generator.
func (*ContentGenerator) Name() string {
	return "content"
}
