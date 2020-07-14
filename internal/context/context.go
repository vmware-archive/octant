/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package context

import "context"

type OctantContextKey string

const KubeConfigKey = OctantContextKey("kubeConfigPath")

func KubeConfigChFrom(ctx context.Context) chan string {
	return ctx.Value(KubeConfigKey).(chan string)
}

func WithKubeConfigCh(ctx context.Context) context.Context {
	return context.WithValue(ctx, KubeConfigKey, make(chan string))
}
