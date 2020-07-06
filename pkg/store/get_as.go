/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package store

import (
	"context"

	"github.com/pkg/errors"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/internal/util/kubernetes"
)

// GetAs gets an object from the object store by key. If the object is not found,
// return false and a nil error.
func GetAs(ctx context.Context, o Store, key Key, as runtime.Object) (bool, error) {
	u, err := o.Get(ctx, key)
	if kerrors.IsNotFound(err) || u == nil {
		return false, nil
	}

	if err != nil {
		return false, errors.Wrap(err, "get object from object store")
	}

	if err := kubernetes.FromUnstructured(u, as); err != nil {
		return false, errors.Wrap(err, "unable to convert object to unstructured")
	}

	return true, nil
}
