/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package cluster

import "context"

// NamespaceInterface is an interface for querying namespace details.
type NamespaceInterface interface {
	Names(ctx context.Context) ([]string, error)
	InitialNamespace() string
	ProvidedNamespaces(ctx context.Context) []string
	HasNamespace(ctx context.Context, namespace string) bool
}
