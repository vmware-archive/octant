/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package testutil

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/pointer"
)

// LoadObjectFromFile loads a file from the `testdata` directory. It will
// assign it a `default` namespace if one is not set.
func LoadObjectFromFile(t *testing.T, objectFile string) runtime.Object {
	data, err := ioutil.ReadFile(filepath.Join("testdata", objectFile))
	require.NoError(t, err)

	object, _, err := scheme.Codecs.UniversalDeserializer().Decode(data, nil, nil)
	require.NoError(t, err)

	accessor := meta.NewAccessor()
	namespace, err := accessor.Namespace(object)
	require.NoError(t, err)
	if namespace == "" {
		require.NoError(t, accessor.SetNamespace(object, "default"))
	}

	return object
}

// LoadTypedObjectFromFile loads a file from the `testdata` directory. It will
// assign it a `default` namespace if one is not set.
func LoadTypedObjectFromFile(t *testing.T, objectFile string, into runtime.Object) {
	data, err := ioutil.ReadFile(filepath.Join("testdata", objectFile))
	require.NoError(t, err)

	gvk := into.GetObjectKind().GroupVersionKind()
	object, _, err := scheme.Codecs.UniversalDeserializer().Decode(data, &gvk, into)
	require.NoError(t, err)

	accessor := meta.NewAccessor()
	namespace, err := accessor.Namespace(object)
	require.NoError(t, err)
	if namespace == "" {
		require.NoError(t, accessor.SetNamespace(object, "default"))
	}
}

// ToOwnerReferences converts an object to owner references.
func ToOwnerReferences(t *testing.T, object runtime.Object) []metav1.OwnerReference {
	objectKind := object.GetObjectKind()
	apiVersion, kind := objectKind.GroupVersionKind().ToAPIVersionAndKind()

	accessor := meta.NewAccessor()
	name, err := accessor.Name(object)
	require.NoError(t, err)

	uid, err := accessor.UID(object)
	require.NoError(t, err)

	return []metav1.OwnerReference{
		{
			APIVersion: apiVersion,
			Kind:       kind,
			Name:       name,
			UID:        uid,
			Controller: pointer.BoolPtr(true),
		},
	}
}
