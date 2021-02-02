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
	)

	crdsDescriber = describer.NewResource(describer.ResourceOptions{
		Path:           "/custom-resource-definitions",
		ObjectStoreKey: store.Key{APIVersion: "apiextensions.k8s.io/v1", Kind: "CustomResourceDefinition"},
		ListType:       &apiextv1.CustomResourceDefinitionList{},
		ObjectType:     &apiextv1.CustomResourceDefinition{},
		Titles:         describer.ResourceTitle{List: "Custom Resource Definitions", Object: "Custom Resource Definitions"},
		ClusterWide:    true,
		IconName:       icon.CustomResourceDefinition,
	})

	rbacClusterRoles = describer.NewResource(describer.ResourceOptions{
		Path:           "/rbac/cluster-roles",
		ObjectStoreKey: store.Key{APIVersion: "rbac.authorization.k8s.io/v1", Kind: "ClusterRole"},
		ListType:       &rbacv1.ClusterRoleList{},
		ObjectType:     &rbacv1.ClusterRole{},
		Titles:         describer.ResourceTitle{List: "Cluster Roles", Object: "Cluster Roles"},
		ClusterWide:    true,
		IconName:       icon.ClusterOverviewClusterRole,
	})

	rbacClusterRoleBindings = describer.NewResource(describer.ResourceOptions{
		Path:           "/rbac/cluster-role-bindings",
		ObjectStoreKey: store.Key{APIVersion: "rbac.authorization.k8s.io/v1", Kind: "ClusterRoleBinding"},
		ListType:       &rbacv1.ClusterRoleBindingList{},
		ObjectType:     &rbacv1.ClusterRoleBinding{},
		Titles:         describer.ResourceTitle{List: "Cluster Role Bindings", Object: "Cluster Role Bindings"},
		ClusterWide:    true,
		IconName:       icon.ClusterOverviewClusterRoleBinding,
	})

	rbacDescriber = describer.NewSection(
		"/rbac",
		"RBAC",
		rbacClusterRoles,
		rbacClusterRoleBindings,
	)

	webhooksDescriber = describer.NewSection(
		"/webhooks",
		"Webhooks",
		webhooksMutatingWebhooks,
		webhooksValidatingWebhooks,
	)

	webhooksValidatingWebhooks = describer.NewResource(describer.ResourceOptions{
		Path:           "/webhooks/validating-webhooks",
		ObjectStoreKey: store.Key{APIVersion: "admissionregistration.k8s.io/v1", Kind: "ValidatingWebhookConfiguration"},
		ListType:       &admissionregistrationv1.ValidatingWebhookConfigurationList{},
		ObjectType:     &admissionregistrationv1.ValidatingWebhookConfiguration{},
		Titles:         describer.ResourceTitle{List: "Validating Webhooks", Object: "Validating Webhook Configuration"},
		ClusterWide:    true,
		IconName:       icon.Webhooks,
	})

	webhooksMutatingWebhooks = describer.NewResource(describer.ResourceOptions{
		Path:           "/webhooks/mutating-webhooks",
		ObjectStoreKey: store.Key{APIVersion: "admissionregistration.k8s.io/v1", Kind: "MutatingWebhookConfiguration"},
		ListType:       &admissionregistrationv1.MutatingWebhookConfigurationList{},
		ObjectType:     &admissionregistrationv1.MutatingWebhookConfiguration{},
		Titles:         describer.ResourceTitle{List: "Mutating Webhooks", Object: "Mutating Webhook Configuration"},
		ClusterWide:    true,
		IconName:       icon.Webhooks,
	})

	nodesDescriber = describer.NewResource(describer.ResourceOptions{
		Path:                  "/nodes",
		ObjectStoreKey:        store.Key{APIVersion: "v1", Kind: "Node"},
		ListType:              &corev1.NodeList{},
		ObjectType:            &corev1.Node{},
		Titles:                describer.ResourceTitle{List: "Nodes", Object: "Nodes"},
		DisableResourceViewer: true,
		ClusterWide:           true,
		IconName:              icon.ClusterOverviewNode,
	})

	storagePersistentVolumeDescriber = describer.NewResource(describer.ResourceOptions{
		Path:           "/storage/persistent-volumes",
		ObjectStoreKey: store.Key{APIVersion: "v1", Kind: "PersistentVolume"},
		ListType:       &corev1.PersistentVolumeList{},
		ObjectType:     &corev1.PersistentVolume{},
		Titles:         describer.ResourceTitle{List: "Persistent Volumes", Object: "Persistent Volumes"},
		ClusterWide:    true,
		IconName:       icon.ClusterOverviewPersistentVolume,
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
	})

	portForwardDescriber = NewPortForwardListDescriber()

	apiServerDescriber = describer.NewSection(
		"/api-server",
		"API Server",
		apiServerApiServices,
	)

	apiServerApiServices = describer.NewResource(describer.ResourceOptions{
		Path:           "/api-server/api-services",
		ObjectStoreKey: store.Key{APIVersion: "apiregistration.k8s.io/v1", Kind: "APIService"},
		ListType:       &apiregistrationv1.APIServiceList{},
		ObjectType:     &apiregistrationv1.APIService{},
		Titles:         describer.ResourceTitle{List: "API Services", Object: "API Services"},
		ClusterWide:    true,
		IconName:       icon.ApiServer,
	})

	rootDescriber = describer.NewSection(
		"/",
		"Cluster Overview",
		namespacesDescriber,
		customResourcesDescriber,
		crdsDescriber,
		rbacDescriber,
		webhooksDescriber,
		nodesDescriber,
		storageDescriber,
		portForwardDescriber,
		apiServerDescriber,
	)
)
