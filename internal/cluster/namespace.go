/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package cluster

import (
	"context"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

//go:generate mockgen -source=namespace.go -destination=./fake/mock_namespace_interface.go -package=fake github.com/vmware-tanzu/octant/internal/cluster NamespaceInterface

// NamespaceInterface is an interface for querying namespace details.
type NamespaceInterface interface {
	Names() ([]string, error)
	InitialNamespace() string
	ProvidedNamespaces() []string
	HasNamespace(namespace string) bool
}

type namespaceClient struct {
	restClient         rest.Interface
	dynamicClient      dynamic.Interface
	initialNamespace   string
	providedNamespaces []string
}

var _ NamespaceInterface = (*namespaceClient)(nil)

func newNamespaceClient(dynamicClient dynamic.Interface, restClient rest.Interface, initialNamespace string, providedNamespaces []string) *namespaceClient {
	return &namespaceClient{
		restClient:         restClient,
		dynamicClient:      dynamicClient,
		initialNamespace:   initialNamespace,
		providedNamespaces: providedNamespaces,
	}
}

func (n *namespaceClient) Names() ([]string, error) {
	namespaces, err := namespaces(n.dynamicClient)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, namespace := range namespaces {
		names = append(names, namespace.GetName())
	}

	return names, nil
}

func (n *namespaceClient) HasNamespace(namespace string) bool {
	ns := &corev1.Namespace{}
	err := n.restClient.Get().Resource("namespaces").Name(namespace).Do().Into(ns)
	if err != nil {
		return false
	}
	return true
}

// Namespaces returns available namespaces.
func namespaces(dc dynamic.Interface) ([]corev1.Namespace, error) {
	res := schema.GroupVersionResource{
		Version:  "v1",
		Resource: "namespaces",
	}

	nri := dc.Resource(res)

	list, err := nri.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "list namespaces")
	}

	var nsList corev1.NamespaceList
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(list.UnstructuredContent(), &nsList)
	if err != nil {
		return nil, errors.Wrap(err, "convert object to namespace list")
	}

	return nsList.Items, nil
}

func (n *namespaceClient) InitialNamespace() string {
	return n.initialNamespace
}

func (n *namespaceClient) ProvidedNamespaces() []string {
	return n.providedNamespaces
}
