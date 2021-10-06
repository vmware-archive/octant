/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package cluster

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	dynamicfake "k8s.io/client-go/dynamic/fake"
)

func Test_namespaceClient_Names(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()

	dc := dynamicfake.NewSimpleDynamicClientWithCustomListKinds(scheme,
		map[schema.GroupVersionResource]string{
			{Group: "", Version: "v1", Resource: "namespaces"}: "NamespacesList",
		},
		newUnstructured("v1", "Namespace", "default", "default"),
		newUnstructured("v1", "Namespace", "app-1", "app-1"),
	)

	nc := newNamespaceClient(dc, nil, nil, "default", []string{})

	got, err := nc.Names(ctx)
	require.NoError(t, err)

	expected := []string{"app-1", "default"}
	assert.Equal(t, expected, got)
}

func Test_namespaceClient_providedNamespaces(t *testing.T) {
	ctx := context.Background()
	providedNamespaces := []string{"default", "user-1"}

	scheme := runtime.NewScheme()
	dc := dynamicfake.NewSimpleDynamicClient(scheme)
	nc := newNamespaceClient(dc, nil, nil, "default", providedNamespaces)

	assert.Equal(t, providedNamespaces, nc.ProvidedNamespaces(ctx))

	nc = newNamespaceClient(dc, nil, nil, "default", []string{})
	assert.Equal(t, nc.ProvidedNamespaces(ctx), []string{"default"})
}

func Test_namespaceClient_InitialNamespace(t *testing.T) {
	expected := "initial-namespace"
	nc := newNamespaceClient(nil, nil, nil, expected, []string{})
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
