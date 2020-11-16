/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package event

import (
	"context"
	"sort"
	"time"

	"github.com/vmware-tanzu/octant/pkg/event"

	"github.com/vmware-tanzu/octant/internal/kubeconfig"
	"github.com/vmware-tanzu/octant/internal/octant"
)

// kubeContextsResponse is a response for current kube contexts.
type kubeContextsResponse struct {
	Contexts       []kubeconfig.Context `json:"contexts"`
	CurrentContext string               `json:"currentContext"`
}

type ContextGeneratorOption func(generator *ContextsGenerator)

type KubeContextStore interface {
	CurrentContext() string
	Contexts() []kubeconfig.Context
}

// ContextsGenerator generates kube contexts for the front end.
type ContextsGenerator struct {
	KubeContextStore KubeContextStore
}

var _ octant.Generator = (*ContextsGenerator)(nil)

func NewContextsGenerator(kubeContextStore KubeContextStore, options ...ContextGeneratorOption) *ContextsGenerator {
	kcg := &ContextsGenerator{
		KubeContextStore: kubeContextStore,
	}

	for _, option := range options {
		option(kcg)
	}

	return kcg
}

func (g *ContextsGenerator) Event(ctx context.Context) (event.Event, error) {
	resp := kubeContextsResponse{
		CurrentContext: g.KubeContextStore.CurrentContext(),
		Contexts:       g.KubeContextStore.Contexts(),
	}

	sort.Slice(resp.Contexts, func(i, j int) bool {
		return resp.Contexts[i].Name < resp.Contexts[j].Name
	})

	e := event.Event{
		Type: event.EventTypeKubeConfig,
		Data: resp,
	}

	return e, nil
}

func (ContextsGenerator) ScheduleDelay() time.Duration {
	return DefaultScheduleDelay
}

func (ContextsGenerator) Name() string {
	return "kubeConfig"
}
