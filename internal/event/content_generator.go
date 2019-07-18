/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package event

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/labels"

	"github.com/vmware/octant/internal/module"
	"github.com/vmware/octant/internal/octant"
	"github.com/vmware/octant/pkg/view/component"
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

	isRunning bool
	mu        sync.Mutex
}

var _ octant.Generator = (*ContentGenerator)(nil)

type dashResponse struct {
	Content component.ContentResponse `json:"content,omitempty"`
}

// IsRunning returns true if content is being generated.
func (g *ContentGenerator) IsRunning() bool {
	return false
}

// Event generates content events.
func (g *ContentGenerator) Event(ctx context.Context) (octant.Event, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.isRunning {
		return octant.Event{}, errNotReady
	}

	g.isRunning = true
	defer func() {
		g.isRunning = false
	}()

	return g.generateContent(ctx)
}

func (g *ContentGenerator) generateContent(ctx context.Context) (octant.Event, error) {
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
