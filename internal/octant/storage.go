/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package octant

import "github.com/vmware-tanzu/octant/pkg/store"

//go:generate mockgen -destination=./fake/mock_storage.go -package=fake . Storage

// Storage is an interface containing storage items.
type Storage interface {
	// ObjectStore returns the object store.
	ObjectStore() store.Store
}
