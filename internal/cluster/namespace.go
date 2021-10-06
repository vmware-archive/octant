/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package cluster

import (
	"context"
	"fmt"

	authv1 "k8s.io/api/authorization/v1"
	"k8s.io/client-go/kubernetes"
	authClient "k8s.io/client-go/kubernetes/typed/authorization/v1"

	internalLog "github.com/vmware-tanzu/octant/internal/log"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"

	clusterTypes "github.com/vmware-tanzu/octant/pkg/cluster"
)

type namespaceClient struct {
	restClient         rest.Interface
	dynamicClient      dynamic.Interface
	authClient         authClient.AuthorizationV1Interface
	initialNamespace   string
	providedNamespaces []string
}

var _ clusterTypes.NamespaceInterface = (*namespaceClient)(nil)

func newNamespaceClient(dynamicClient dynamic.Interface, restClient rest.Interface, kubernetesClient kubernetes.Interface, initialNamespace string, providedNamespaces []string) *namespaceClient {
	var authClient authClient.AuthorizationV1Interface
	if kubernetesClient != nil {
		authClient = kubernetesClient.AuthorizationV1()
	}
	return &namespaceClient{
		restClient:         restClient,
		dynamicClient:      dynamicClient,
		authClient:         authClient,
		initialNamespace:   initialNamespace,
		providedNamespaces: providedNamespaces,
	}
}

func (n *namespaceClient) Names(ctx context.Context) ([]string, error) {
	namespaces, err := n.namespaces(ctx, n.dynamicClient)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, namespace := range namespaces {
		names = append(names, namespace.GetName())
	}

	return names, nil
}

func (n *namespaceClient) HasNamespace(ctx context.Context, namespace string) bool {
	ns := &corev1.Namespace{}
	err := n.restClient.Get().Resource("namespaces").Name(namespace).Do(ctx).Into(ns)
	if err != nil {
		return false
	}
	return true
}

// Namespaces returns available namespaces.
func (n *namespaceClient) namespaces(ctx context.Context, dc dynamic.Interface) ([]corev1.Namespace, error) {
	logger := internalLog.From(ctx)

	res := schema.GroupVersionResource{
		Version:  "v1",
		Resource: "namespaces",
	}

	var nsList corev1.NamespaceList

	ssar := &authv1.SelfSubjectAccessReview{
		Spec: authv1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &authv1.ResourceAttributes{
				Verb:     "list",
				Version:  "v1",
				Resource: "namespaces",
			},
		},
	}

	var resp *authv1.SelfSubjectAccessReview
	var err error
	if n.authClient != nil {
		resp, err = n.authClient.SelfSubjectAccessReviews().Create(ctx, ssar, metav1.CreateOptions{})
		if err != nil {
			return nil, fmt.Errorf("auth client: %w", err)
		}
	}

	if resp == nil || resp.Status.Allowed {
		nri := dc.Resource(res)

		list, err := nri.List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, errors.Wrap(err, "list namespaces")
		}

		err = runtime.DefaultUnstructuredConverter.FromUnstructured(list.UnstructuredContent(), &nsList)
		if err != nil {
			return nil, errors.Wrap(err, "convert object to namespace list")
		}
	}

	if len(nsList.Items) == 0 {
		logger.Debugf("no namespaces found")
	}

	return nsList.Items, nil
}

// InitialNamespace returns the initial namespace for Octant
func (n *namespaceClient) InitialNamespace() string {
	return n.initialNamespace
}

// ProvidedNamespaces returns the list of namespaces provided.
// If no namespaces are provided, it will default to returning the InitialNamespace
func (n *namespaceClient) ProvidedNamespaces(ctx context.Context) []string {
	if len(n.providedNamespaces) == 0 {
		n.providedNamespaces = []string{n.initialNamespace}
	}
	return n.providedNamespaces
}
