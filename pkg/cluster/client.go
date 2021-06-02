/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package cluster

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// ClientInterface is a client for cluster operations.
type ClientInterface interface {
	DefaultNamespace() string
	ResourceExists(schema.GroupVersionResource) bool
	Resource(schema.GroupKind) (schema.GroupVersionResource, bool, error)
	ResetMapper()
	KubernetesClient() (kubernetes.Interface, error)
	DynamicClient() (dynamic.Interface, error)
	DiscoveryClient() (discovery.DiscoveryInterface, error)
	NamespaceClient() (NamespaceInterface, error)
	InfoClient() (InfoInterface, error)
	Close()
	RESTInterface
}

type RESTInterface interface {
	RESTClient() (rest.Interface, error)
	RESTConfig() *rest.Config
}

type RESTConfigOptions struct {
	QPS       float32
	Burst     int
	UserAgent string
}
