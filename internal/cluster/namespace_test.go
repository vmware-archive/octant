/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package cluster

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	dynamicfake "k8s.io/client-go/dynamic/fake"
)

func Test_namespaceClient_Names(t *testing.T) {
	scheme := runtime.NewScheme()

	// NOTE: this should be reverted to the k8s.io/client-go/dynamic/fake when bug fix is
	// merged upstream
	dc := dynamicfake.NewSimpleDynamicClient(scheme,
		newUnstructured("v1", "Namespace", "", "default"),
		newUnstructured("v1", "Namespace", "", "app-1"),
	)

	nc := newNamespaceClient(dc, "default")

	got, err := nc.Names()
	require.NoError(t, err)

	expected := []string{"default", "app-1"}
	assert.Equal(t, expected, got)
}

func Test_namespaceClient_InitialNamespace(t *testing.T) {
	expected := "inital-namespace"
	nc := newNamespaceClient(nil, expected)
	assert.Equal(t, expected, nc.InitialNamespace())
}

func newUnstructured(apiVersion, kind, namespace, name string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": apiVersion,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"namespace": namespace,
				"name":      name,
			},
		},
	}
}
