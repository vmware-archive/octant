/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package testutil

import (
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// ToUnstructured converts an object to an unstructured.
func ToUnstructured(t *testing.T, object runtime.Object) *unstructured.Unstructured {
	m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(object)
	require.NoError(t, err)

	u := &unstructured.Unstructured{Object: m}

	return u
}

// ToUnstructuredList converts a list of objects to a list of unstructured.
func ToUnstructuredList(t *testing.T, objects ...runtime.Object) *unstructured.UnstructuredList {
	list := &unstructured.UnstructuredList{}

	for _, object := range objects {
		list.Items = append(list.Items, *ToUnstructured(t, object))
	}

	return list
}
