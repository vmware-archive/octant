/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package testutil

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/pointer"
)

// LoadObjectFromFile loads a file from the `testdata` directory. It will
// assign it a `default` namespace if one is not set.
func LoadObjectFromFile(t *testing.T, objectFile string) runtime.Object {
	data, err := ioutil.ReadFile(filepath.Join("testdata", objectFile))
	require.NoError(t, err)

	object, _, err := scheme.Codecs.UniversalDeserializer().Decode(data, nil, nil)
	require.NoError(t, err, "unable to decode serialized data")

	accessor := meta.NewAccessor()
	namespace, err := accessor.Namespace(object)
	require.NoError(t, err)
	if namespace == "" {
		require.NoError(t, accessor.SetNamespace(object, "default"))
	}

	return object
}

// LoadTestData loads a file in the testdata directory.
func LoadTestData(t *testing.T, fileName string) []byte {
	data, err := ioutil.ReadFile(filepath.Join("testdata", fileName))
	require.NoError(t, err)

	// strip carriage returns for Widows
	data = bytes.Replace(data, []byte{13, 10}, []byte{10}, -1)

	return data
}

// LoadUnstructuredFromFile loads an object from a file in the in `testdata` directory.
// It will assign a `default` namespace if one is not set. This helper does not support
// multiple objects in a YAML file.
func LoadUnstructuredFromFile(t *testing.T, objectFile string) *unstructured.Unstructured {
	f, err := os.Open(filepath.Join("testdata", objectFile))
	require.NoError(t, err)
	defer func() {
		require.NoError(t, f.Close())
	}()

	d := yaml.NewYAMLOrJSONDecoder(f, 4096)

	ext := runtime.RawExtension{}
	require.NoError(t, d.Decode(&ext))

	obj, _, err := unstructured.UnstructuredJSONScheme.Decode(ext.Raw, nil, nil)
	require.NoError(t, err)

	u, ok := obj.(*unstructured.Unstructured)
	require.True(t, ok, "object is not an unstructured object")

	return u
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
func ToOwnerReferences(t *testing.T, objects ...runtime.Object) []metav1.OwnerReference {
	var list []metav1.OwnerReference

	for _, object := range objects {
		objectKind := object.GetObjectKind()
		apiVersion, kind := objectKind.GroupVersionKind().ToAPIVersionAndKind()

		accessor := meta.NewAccessor()
		name, err := accessor.Name(object)
		require.NoError(t, err)

		uid, err := accessor.UID(object)
		require.NoError(t, err)

		list = append(list, metav1.OwnerReference{
			APIVersion: apiVersion,
			Kind:       kind,
			Name:       name,
			UID:        uid,
			Controller: pointer.BoolPtr(true),
		})
	}

	return list
}

// Time generates a test time
func Time() time.Time {
	return time.Unix(1547211430, 0)
}

// RequireErrorOrNot or not is a helper that requires an error or not.
func RequireErrorOrNot(t *testing.T, wantErr bool, err error) {
	if wantErr {
		require.Error(t, err)
		return
	}
	require.NoError(t, err)
}
