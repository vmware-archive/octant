/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package octant

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	sigyaml "sigs.k8s.io/yaml"

	clusterFake "github.com/vmware-tanzu/octant/internal/cluster/fake"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/pkg/action"
	actionFake "github.com/vmware-tanzu/octant/pkg/action/fake"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/store/fake"
)

func TestNewApplyYaml_NamespacedCreate(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	logger := log.NopLogger()
	objectStore := fake.NewMockStore(controller)
	clusterClient := clusterFake.NewMockClientInterface(controller)
	alerter := actionFake.NewMockAlerter(controller)

	gvr := schema.GroupVersionResource{
		Version:  "v1",
		Resource: "configmaps",
	}
	clusterClient.EXPECT().
		Resource(schema.GroupKind{Kind: "ConfigMap"}).
		Return(gvr, true, nil)

	key := store.Key{
		Namespace:  "default",
		APIVersion: "v1",
		Kind:       "ConfigMap",
		Name:       "greeting",
	}
	objectStore.EXPECT().
		Get(gomock.Any(), key).
		Return(nil, kerrors.NewNotFound(gvr.GroupResource(), "greeting"))

	uGreeting := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name":      "greeting",
				"namespace": "default",
			},
			"data": map[string]interface{}{
				"hello": "world",
			},
		},
	}
	objectStore.EXPECT().
		Create(gomock.Any(), uGreeting).
		Return(nil)

	alerter.EXPECT().
		SendAlert(gomock.Any()).
		DoAndReturn(func(alert action.Alert) {
			assert.Equal(t, action.AlertTypeInfo, alert.Type)
			assert.Equal(t, "Created ConfigMap (v1) greeting in default", alert.Message)
			assert.NotNil(t, alert.Expiration)
		})

	applyYaml := NewApplyYaml(logger, objectStore, clusterClient)

	ctx := context.Background()

	payload := action.CreatePayload(ActionApplyYaml, map[string]interface{}{
		"update": `
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: greeting
data:
  hello: world
`,
		"namespace": "default",
	})

	require.NoError(t, applyYaml.Handle(ctx, alerter, payload))
}

func TestNewApplyYaml_NamespacedUpdate(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	logger := log.NopLogger()
	objectStore := fake.NewMockStore(controller)
	clusterClient := clusterFake.NewMockClientInterface(controller)
	dynamicClient := clusterFake.NewMockDynamicInterface(controller)
	dynamicNamespaceableResourceClient := clusterFake.NewMockNamespaceableResourceInterface(controller)
	dynamicResourceClient := clusterFake.NewMockResourceInterface(controller)
	alerter := actionFake.NewMockAlerter(controller)

	gvr := schema.GroupVersionResource{
		Version:  "v1",
		Resource: "configmaps",
	}
	clusterClient.EXPECT().
		Resource(schema.GroupKind{Kind: "ConfigMap"}).
		Return(gvr, true, nil)

	key := store.Key{
		Namespace:  "default",
		APIVersion: "v1",
		Kind:       "ConfigMap",
		Name:       "greeting",
	}
	objectStore.EXPECT().
		Get(gomock.Any(), key).
		// should return the ConfigMap, but we don't actually care
		Return(nil, nil)

	clusterClient.EXPECT().
		DynamicClient().
		Return(dynamicClient, nil)
	dynamicClient.EXPECT().
		Resource(gvr).
		Return(dynamicNamespaceableResourceClient)
	dynamicNamespaceableResourceClient.EXPECT().
		Namespace("default").
		Return(dynamicResourceClient)

	unstructuredYaml, _ := sigyaml.Marshal(map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "ConfigMap",
		"metadata": map[string]interface{}{
			"name":      "greeting",
			"namespace": "default",
		},
		"data": map[string]interface{}{
			"hello": "world",
		},
	})
	withForce := true
	dynamicResourceClient.EXPECT().
		Patch(
			gomock.Any(),
			"greeting",
			types.ApplyPatchType,
			unstructuredYaml,
			metav1.PatchOptions{FieldManager: "octant", Force: &withForce},
		).
		// should return the ConfigMap, but we don't actually care
		Return(nil, nil)

	alerter.EXPECT().
		SendAlert(gomock.Any()).
		DoAndReturn(func(alert action.Alert) {
			assert.Equal(t, action.AlertTypeInfo, alert.Type)
			assert.Equal(t, "Updated ConfigMap (v1) greeting in default", alert.Message)
			assert.NotNil(t, alert.Expiration)
		})

	applyYaml := NewApplyYaml(logger, objectStore, clusterClient)

	ctx := context.Background()

	payload := action.CreatePayload(ActionApplyYaml, map[string]interface{}{
		"update": `
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: greeting
data:
  hello: world
`,
		"namespace": "default",
	})

	require.NoError(t, applyYaml.Handle(ctx, alerter, payload))
}

func TestNewApplyYaml_ClusterCreate(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	logger := log.NopLogger()
	objectStore := fake.NewMockStore(controller)
	clusterClient := clusterFake.NewMockClientInterface(controller)
	alerter := actionFake.NewMockAlerter(controller)

	gvr := schema.GroupVersionResource{
		Version:  "v1",
		Resource: "namespaces",
	}
	clusterClient.EXPECT().
		Resource(schema.GroupKind{Kind: "Namespace"}).
		Return(gvr, false, nil)

	key := store.Key{
		APIVersion: "v1",
		Kind:       "Namespace",
		Name:       "greeting-system",
	}
	objectStore.EXPECT().
		Get(gomock.Any(), key).
		Return(nil, kerrors.NewNotFound(gvr.GroupResource(), "greeting-system"))

	uGreeting := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Namespace",
			"metadata": map[string]interface{}{
				"name": "greeting-system",
			},
		},
	}
	objectStore.EXPECT().
		Create(gomock.Any(), uGreeting).
		Return(nil)

	alerter.EXPECT().
		SendAlert(gomock.Any()).
		DoAndReturn(func(alert action.Alert) {
			assert.Equal(t, action.AlertTypeInfo, alert.Type)
			assert.Equal(t, "Created Namespace (v1) greeting-system", alert.Message)
			assert.NotNil(t, alert.Expiration)
		})

	applyYaml := NewApplyYaml(logger, objectStore, clusterClient)

	ctx := context.Background()

	payload := action.CreatePayload(ActionApplyYaml, map[string]interface{}{
		"update": `
---
apiVersion: v1
kind: Namespace
metadata:
  name: greeting-system
`,
		"namespace": "default",
	})

	require.NoError(t, applyYaml.Handle(ctx, alerter, payload))
}

func TestNewApplyYaml_ClusterUpdate(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	logger := log.NopLogger()
	objectStore := fake.NewMockStore(controller)
	clusterClient := clusterFake.NewMockClientInterface(controller)
	dynamicClient := clusterFake.NewMockDynamicInterface(controller)
	dynamicResourceClient := clusterFake.NewMockNamespaceableResourceInterface(controller)
	alerter := actionFake.NewMockAlerter(controller)

	gvr := schema.GroupVersionResource{
		Version:  "v1",
		Resource: "namespaces",
	}
	clusterClient.EXPECT().
		Resource(schema.GroupKind{Kind: "Namespace"}).
		Return(gvr, false, nil)

	key := store.Key{
		APIVersion: "v1",
		Kind:       "Namespace",
		Name:       "greeting-system",
	}
	objectStore.EXPECT().
		Get(gomock.Any(), key).
		// should return the ConfigMap, but we don't actually care
		Return(nil, nil)

	clusterClient.EXPECT().
		DynamicClient().
		Return(dynamicClient, nil)
	dynamicClient.EXPECT().
		Resource(gvr).
		Return(dynamicResourceClient)

	unstructuredYaml, _ := sigyaml.Marshal(map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Namespace",
		"metadata": map[string]interface{}{
			"name": "greeting-system",
		},
	})
	withForce := true
	dynamicResourceClient.EXPECT().
		Patch(
			gomock.Any(),
			"greeting-system",
			types.ApplyPatchType,
			unstructuredYaml,
			metav1.PatchOptions{FieldManager: "octant", Force: &withForce},
		).
		// should return the ConfigMap, but we don't actually care
		Return(nil, nil)

	alerter.EXPECT().
		SendAlert(gomock.Any()).
		DoAndReturn(func(alert action.Alert) {
			assert.Equal(t, action.AlertTypeInfo, alert.Type)
			assert.Equal(t, "Updated Namespace (v1) greeting-system", alert.Message)
			assert.NotNil(t, alert.Expiration)
		})

	applyYaml := NewApplyYaml(logger, objectStore, clusterClient)

	ctx := context.Background()

	payload := action.CreatePayload(ActionApplyYaml, map[string]interface{}{
		"update": `
---
apiVersion: v1
kind: Namespace
metadata:
  name: greeting-system
`,
		"namespace": "default",
	})

	require.NoError(t, applyYaml.Handle(ctx, alerter, payload))
}

func TestNewApplyYaml_Error(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	logger := log.NopLogger()
	objectStore := fake.NewMockStore(controller)
	clusterClient := clusterFake.NewMockClientInterface(controller)
	alerter := actionFake.NewMockAlerter(controller)

	alerter.EXPECT().
		SendAlert(gomock.Any()).
		DoAndReturn(func(alert action.Alert) {
			assert.Equal(t, action.AlertTypeError, alert.Type)
			assert.Contains(t, alert.Message, "Unable to apply yaml:")
			assert.NotNil(t, alert.Expiration)
		})

	applyYaml := NewApplyYaml(logger, objectStore, clusterClient)

	ctx := context.Background()

	payload := action.CreatePayload(ActionApplyYaml, map[string]interface{}{
		"update": `
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: greeting
data: {
`,
		"namespace": "default",
	})

	require.NoError(t, applyYaml.Handle(ctx, alerter, payload))
}
