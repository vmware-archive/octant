/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package octant

import (
	"context"
	"fmt"
	"testing"

	"k8s.io/apimachinery/pkg/types"

	"github.com/vmware-tanzu/octant/internal/objectstore"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

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
	clusterClient := clusterFake.NewMockClientInterface(controller)
	dynamicCache := fake.NewMockStore(controller)
	alerter := actionFake.NewMockAlerter(controller)

	gvr := schema.GroupVersionResource{
		Version:  "v1",
		Resource: "configmaps",
	}

	key := store.Key{
		Namespace:  "default",
		APIVersion: "v1",
		Kind:       "ConfigMap",
		Name:       "greeting",
	}
	update := `
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: greeting
data:
  hello: world
`
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

	notFound := kerrors.NewNotFound(schema.GroupResource{}, "greeting-system")
	clusterClient.EXPECT().Resource(schema.GroupKind{Kind: "ConfigMap"}).Return(gvr, true, nil)
	dynamicCache.EXPECT().Get(context.TODO(), key).Return(nil, notFound)
	dynamicCache.EXPECT().Create(context.TODO(), uGreeting)

	output, err := objectstore.CreateOrUpdateFromHandler(
		context.TODO(),
		key.Namespace,
		update,
		dynamicCache.Get,
		dynamicCache.Create,
		clusterClient)
	dynamicCache.EXPECT().CreateOrUpdateFromYAML(gomock.Any(), key.Namespace, update).Return(output, err)

	alerter.EXPECT().
		SendAlert(gomock.Any()).
		DoAndReturn(func(alert action.Alert) {
			assert.Equal(t, action.AlertTypeInfo, alert.Type)
			assert.Equal(t, "Created ConfigMap (v1) greeting in default", alert.Message)
			assert.NotNil(t, alert.Expiration)
		})

	applyYaml := NewApplyYaml(logger, dynamicCache)

	ctx := context.Background()

	payload := action.CreatePayload(action.ActionApplyYaml, map[string]interface{}{
		"update":    update,
		"namespace": "default",
	})

	require.NoError(t, applyYaml.Handle(ctx, alerter, payload))
}

func TestNewApplyYaml_NamespacedUpdate(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	logger := log.NopLogger()
	clusterClient := clusterFake.NewMockClientInterface(controller)
	dynamicClient := clusterFake.NewMockDynamicInterface(controller)
	nsResourceClient := clusterFake.NewMockNamespaceableResourceInterface(controller)
	dynamicCache := fake.NewMockStore(controller)
	alerter := actionFake.NewMockAlerter(controller)

	gvr := schema.GroupVersionResource{
		Version:  "v1",
		Resource: "configmaps",
	}
	key := store.Key{
		Kind:       "ConfigMap",
		APIVersion: "v1",
		Name:       "greeting",
		Namespace:  "default",
	}
	update := `
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: greeting
data:
  hello: world
`
	clusterClient.EXPECT().Resource(gomock.Any()).Return(gvr, true, nil).AnyTimes()
	clusterClient.EXPECT().DynamicClient().Return(dynamicClient, nil)
	dynamicClient.EXPECT().Resource(gvr).Return(nsResourceClient)
	nsResourceClient.EXPECT().Namespace("default").Return(nsResourceClient)
	nsResourceClient.EXPECT().Patch(context.TODO(), key.Name, types.ApplyPatchType, gomock.Any(), gomock.Any(), gomock.Any())
	dynamicCache.EXPECT().Get(context.TODO(), key)

	output, err := objectstore.CreateOrUpdateFromHandler(
		context.TODO(),
		"default",
		update,
		dynamicCache.Get,
		dynamicCache.Create,
		clusterClient)
	dynamicCache.EXPECT().CreateOrUpdateFromYAML(gomock.Any(), "default", update).Return(output, err)

	alerter.EXPECT().
		SendAlert(gomock.Any()).
		DoAndReturn(func(alert action.Alert) {
			assert.Equal(t, action.AlertTypeInfo, alert.Type)
			assert.Equal(t, "Updated ConfigMap (v1) greeting in default", alert.Message)
			assert.NotNil(t, alert.Expiration)
		})

	applyYaml := NewApplyYaml(logger, dynamicCache)

	ctx := context.Background()

	payload := action.CreatePayload(action.ActionApplyYaml, map[string]interface{}{
		"update":    update,
		"namespace": "default",
	})

	require.NoError(t, applyYaml.Handle(ctx, alerter, payload))
}

func TestNewApplyYaml_ClusterCreate(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	logger := log.NopLogger()
	clusterClient := clusterFake.NewMockClientInterface(controller)
	dynamicCache := fake.NewMockStore(controller)
	alerter := actionFake.NewMockAlerter(controller)

	key := store.Key{
		APIVersion: "v1",
		Kind:       "Namespace",
		Name:       "greeting-system",
	}
	update := `
---
apiVersion: v1
kind: Namespace
metadata:
  name: greeting-system
`
	uNamespace := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Namespace",
			"metadata": map[string]interface{}{
				"name": "greeting-system",
			},
		},
	}
	notFound := kerrors.NewNotFound(schema.GroupResource{}, "greeting-system")
	clusterClient.EXPECT().Resource(gomock.Any()).Return(schema.GroupVersionResource{}, false, nil).AnyTimes()
	dynamicCache.EXPECT().Get(context.TODO(), key).Return(nil, notFound)
	dynamicCache.EXPECT().Create(context.TODO(), uNamespace)

	output, err := objectstore.CreateOrUpdateFromHandler(
		context.TODO(),
		"",
		update,
		dynamicCache.Get,
		dynamicCache.Create,
		clusterClient)
	dynamicCache.EXPECT().CreateOrUpdateFromYAML(gomock.Any(), "default", update).Return(output, err)

	alerter.EXPECT().
		SendAlert(gomock.Any()).
		DoAndReturn(func(alert action.Alert) {
			assert.Equal(t, action.AlertTypeInfo, alert.Type)
			assert.Equal(t, "Created Namespace (v1) greeting-system", alert.Message)
			assert.NotNil(t, alert.Expiration)
		})

	applyYaml := NewApplyYaml(logger, dynamicCache)

	ctx := context.Background()

	payload := action.CreatePayload(action.ActionApplyYaml, map[string]interface{}{
		"update":    update,
		"namespace": "default",
	})

	require.NoError(t, applyYaml.Handle(ctx, alerter, payload))
}

func TestNewApplyYaml_ClusterUpdate(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	logger := log.NopLogger()
	clusterClient := clusterFake.NewMockClientInterface(controller)
	dynamicClient := clusterFake.NewMockDynamicInterface(controller)
	nsResourceClient := clusterFake.NewMockNamespaceableResourceInterface(controller)
	dynamicCache := fake.NewMockStore(controller)
	alerter := actionFake.NewMockAlerter(controller)

	key := store.Key{
		APIVersion: "v1",
		Kind:       "Namespace",
		Name:       "greeting-system",
	}
	update := `
---
apiVersion: v1
kind: Namespace
metadata:
  name: greeting-system
`
	clusterClient.EXPECT().Resource(gomock.Any()).Return(schema.GroupVersionResource{}, false, nil).AnyTimes()
	clusterClient.EXPECT().DynamicClient().Return(dynamicClient, nil)
	dynamicClient.EXPECT().Resource(gomock.Any()).Return(nsResourceClient)
	nsResourceClient.EXPECT().Patch(context.TODO(), key.Name, types.ApplyPatchType, gomock.Any(), gomock.Any(), gomock.Any())
	dynamicCache.EXPECT().Get(context.TODO(), key)

	output, err := objectstore.CreateOrUpdateFromHandler(
		context.TODO(),
		"",
		update,
		dynamicCache.Get,
		dynamicCache.Create,
		clusterClient)
	dynamicCache.EXPECT().CreateOrUpdateFromYAML(gomock.Any(), "default", update).Return(output, err)

	alerter.EXPECT().
		SendAlert(gomock.Any()).
		DoAndReturn(func(alert action.Alert) {
			assert.Equal(t, action.AlertTypeInfo, alert.Type)
			assert.Equal(t, "Updated Namespace (v1) greeting-system", alert.Message)
			assert.NotNil(t, alert.Expiration)
		})

	applyYaml := NewApplyYaml(logger, dynamicCache)

	ctx := context.Background()

	payload := action.CreatePayload(action.ActionApplyYaml, map[string]interface{}{
		"update":    update,
		"namespace": "default",
	})

	require.NoError(t, applyYaml.Handle(ctx, alerter, payload))
}

func TestNewApplyYaml_Error(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	logger := log.NopLogger()
	objectStore := fake.NewMockStore(controller)
	alerter := actionFake.NewMockAlerter(controller)

	update := `
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: greeting
data: {
`
	objectStore.EXPECT().
		CreateOrUpdateFromYAML(gomock.Any(), "default", update).
		Return(nil, fmt.Errorf("Unable to apply yaml:"))

	alerter.EXPECT().
		SendAlert(gomock.Any()).
		DoAndReturn(func(alert action.Alert) {
			assert.Equal(t, action.AlertTypeError, alert.Type)
			assert.Contains(t, alert.Message, "Unable to apply yaml:")
			assert.NotNil(t, alert.Expiration)
		})

	applyYaml := NewApplyYaml(logger, objectStore)

	ctx := context.Background()

	payload := action.CreatePayload(action.ActionApplyYaml, map[string]interface{}{
		"update":    update,
		"namespace": "default",
	})

	require.NoError(t, applyYaml.Handle(ctx, alerter, payload))
}
