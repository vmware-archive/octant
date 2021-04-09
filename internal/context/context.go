/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package context

import (
	"context"

	"github.com/vmware-tanzu/octant/internal/octant"
)

type OctantContextKey string

const KubeConfigKey = OctantContextKey("kubeConfigPath")

func KubeConfigChFrom(ctx context.Context) chan string {
	return ctx.Value(KubeConfigKey).(chan string)
}

func WithKubeConfigCh(ctx context.Context) context.Context {
	return context.WithValue(ctx, KubeConfigKey, make(chan string))
}

type OctantClientState string

const ClientStateKey = OctantClientState("clientState")

type ClientState struct {
	ClientID    string          `json:"clientID"`
	Filters     []octant.Filter `json:"filters"`
	Namespace   string          `json:"namespace"`
	ContextName string          `json:"contextName"`
}

func WithClientState(ctx context.Context, state ClientState) context.Context {
	return context.WithValue(ctx, ClientStateKey, state)
}

func ClientStateFrom(ctx context.Context) ClientState {
	if ctx.Value(ClientStateKey) == nil {
		return ClientState{}
	}
	return ctx.Value(ClientStateKey).(ClientState)
}
