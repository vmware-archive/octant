/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package cluster

// NamespaceInterface is an interface for querying namespace details.
type NamespaceInterface interface {
	Names() ([]string, error)
	InitialNamespace() string
	ProvidedNamespaces() []string
	HasNamespace(namespace string) bool
}
