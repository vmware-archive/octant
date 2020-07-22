/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package clusteroverview

import (
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"

	"github.com/vmware-tanzu/octant/internal/describer"
	"github.com/vmware-tanzu/octant/pkg/icon"
	"github.com/vmware-tanzu/octant/pkg/store"
)

var (
	customResourcesDescriber = describer.NewCRDSection(
		"/custom-resources",
		"Custom Resources",
		describer.ResourceLink{Title: "Cluster Overview", Url: "/cluster-overview"},
	)

	crdsDescriber = describer.NewResource(describer.ResourceOptions{
		Path:           "/custom-resource-definitions",
		ObjectStoreKey: store.Key{APIVersion: "apiextensions.k8s.io/v1", Kind: "CustomResourceDefinition"},
		ListType:       &apiextv1.CustomResourceDefinitionList{},
		ObjectType:     &apiextv1.CustomResourceDefinition{},
		Titles:         describer.ResourceTitle{List: "Custom Resource Definitions", Object: "Custom Resource Definition"},
		ClusterWide:    true,
		IconName:       icon.CustomResourceDefinition,
		RootPath:       describer.ResourceLink{Title: "Cluster Overview", Url: "/cluster-overview"},
	})

	rbacClusterRoles = describer.NewResource(describer.ResourceOptions{
		Path:           "/rbac/cluster-roles",
		ObjectStoreKey: store.Key{APIVersion: "rbac.authorization.k8s.io/v1", Kind: "ClusterRole"},
		ListType:       &rbacv1.ClusterRoleList{},
		ObjectType:     &rbacv1.ClusterRole{},
		Titles:         describer.ResourceTitle{List: "Cluster Roles", Object: "Cluster Roles"},
		ClusterWide:    true,
		IconName:       icon.ClusterOverviewClusterRole,
		RootPath:       describer.ResourceLink{Title: "Cluster Overview", Url: "/cluster-overview"},
	})

	rbacClusterRoleBindings = describer.NewResource(describer.ResourceOptions{
		Path:           "/rbac/cluster-role-bindings",
		ObjectStoreKey: store.Key{APIVersion: "rbac.authorization.k8s.io/v1", Kind: "ClusterRoleBinding"},
		ListType:       &rbacv1.ClusterRoleBindingList{},
		ObjectType:     &rbacv1.ClusterRoleBinding{},
		Titles:         describer.ResourceTitle{List: "Cluster Role Bindings", Object: "Cluster Role Bindings"},
		ClusterWide:    true,
		IconName:       icon.ClusterOverviewClusterRoleBinding,
		RootPath:       describer.ResourceLink{Title: "Cluster Overview", Url: "/cluster-overview"},
	})

	rbacDescriber = describer.NewSection(
		"/rbac",
		"RBAC",
		rbacClusterRoles,
		rbacClusterRoleBindings,
	)

	nodesDescriber = describer.NewResource(describer.ResourceOptions{
		Path:                  "/nodes",
		ObjectStoreKey:        store.Key{APIVersion: "v1", Kind: "Node"},
		ListType:              &corev1.NodeList{},
		ObjectType:            &corev1.Node{},
		Titles:                describer.ResourceTitle{List: "Nodes", Object: "Nodes"},
		DisableResourceViewer: true,
		ClusterWide:           true,
		IconName:              icon.ClusterOverviewNode,
		RootPath:              describer.ResourceLink{Title: "Cluster Overview", Url: "/cluster-overview"},
	})

	storagePersistentVolumeDescriber = describer.NewResource(describer.ResourceOptions{
		Path:           "/storage/persistent-volumes",
		ObjectStoreKey: store.Key{APIVersion: "v1", Kind: "PersistentVolume"},
		ListType:       &corev1.PersistentVolumeList{},
		ObjectType:     &corev1.PersistentVolume{},
		Titles:         describer.ResourceTitle{List: "Persistent Volumes", Object: "Persistent Volumes"},
		ClusterWide:    true,
		IconName:       icon.ClusterOverviewPersistentVolume,
		RootPath:       describer.ResourceLink{Title: "Cluster Overview", Url: "/cluster-overview"},
	})

	storageDescriber = describer.NewSection(
		"/storage",
		"Storage",
		storagePersistentVolumeDescriber,
	)

	namespacesDescriber = describer.NewResource(describer.ResourceOptions{
		Path:                  "/namespaces",
		ObjectStoreKey:        store.Key{APIVersion: "v1", Kind: "Namespace"},
		ListType:              &corev1.NamespaceList{},
		ObjectType:            &corev1.Namespace{},
		Titles:                describer.ResourceTitle{List: "Namespaces", Object: "Namespaces"},
		DisableResourceViewer: true,
		ClusterWide:           true,
		IconName:              icon.ClusterOverviewNamespace,
		RootPath:              describer.ResourceLink{Title: "Cluster Overview", Url: "/cluster-overview"},
	})

	portForwardDescriber = NewPortForwardListDescriber()

	apiServerDescriber = describer.NewSection(
		"/api-server",
		"API Server",
		apiServerApiServices,
		apiServerMutatingWebhooks,
		apiServerValidatingWebhooks,
	)

	apiServerApiServices = describer.NewResource(describer.ResourceOptions{
		Path:           "/api-server/api-services",
		ObjectStoreKey: store.Key{APIVersion: "apiregistration.k8s.io/v1", Kind: "APIService"},
		ListType:       &apiregistrationv1.APIServiceList{},
		ObjectType:     &apiregistrationv1.APIService{},
		Titles:         describer.ResourceTitle{List: "API Services", Object: "API Service"},
		ClusterWide:    true,
		IconName:       icon.ApiServer,
		RootPath:       describer.ResourceLink{Title: "Cluster Overview", Url: "/cluster-overview"},
	})

	apiServerValidatingWebhooks = describer.NewResource(describer.ResourceOptions{
		Path:           "/api-server/validating-webhooks",
		ObjectStoreKey: store.Key{APIVersion: "admissionregistration.k8s.io/v1", Kind: "ValidatingWebhookConfiguration"},
		ListType:       &admissionregistrationv1.ValidatingWebhookConfigurationList{},
		ObjectType:     &admissionregistrationv1.ValidatingWebhookConfiguration{},
		Titles:         describer.ResourceTitle{List: "Validating Webhook Configurations", Object: "Validating Webhook Configuration"},
		ClusterWide:    true,
		IconName:       icon.ApiServer,
		RootPath:       describer.ResourceLink{Title: "Cluster Overview", Url: "/cluster-overview"},
	})

	apiServerMutatingWebhooks = describer.NewResource(describer.ResourceOptions{
		Path:           "/api-server/mutating-webhooks",
		ObjectStoreKey: store.Key{APIVersion: "admissionregistration.k8s.io/v1", Kind: "MutatingWebhookConfiguration"},
		ListType:       &admissionregistrationv1.MutatingWebhookConfigurationList{},
		ObjectType:     &admissionregistrationv1.MutatingWebhookConfiguration{},
		Titles:         describer.ResourceTitle{List: "Mutating Webhook Configurations", Object: "Mutating Webhook Configuration"},
		ClusterWide:    true,
		IconName:       icon.ApiServer,
		RootPath:       describer.ResourceLink{Title: "Cluster Overview", Url: "/cluster-overview"},
	})

	rootDescriber = describer.NewSection(
		"/",
		"Cluster Overview",
		namespacesDescriber,
		apiServerDescriber,
		customResourcesDescriber,
		crdsDescriber,
		rbacDescriber,
		nodesDescriber,
		storageDescriber,
		portForwardDescriber,
	)
)
