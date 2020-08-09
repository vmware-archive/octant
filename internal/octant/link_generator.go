/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package octant

//go:generate mockgen -destination=./fake/mock_link_generator.go -package=fake . LinkGenerator

// LinkGenerator is an interface containing object path items.
type LinkGenerator interface {
	// ObjectPath returns the path of a reference.
	ObjectPath(namespace, apiVersion, kind, name string) (string, error)
}
