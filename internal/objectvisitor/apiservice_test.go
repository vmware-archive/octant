/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectvisitor_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"

	"github.com/vmware-tanzu/octant/internal/objectvisitor"
	"github.com/vmware-tanzu/octant/internal/objectvisitor/fake"
	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/store"
	objectStoreFake "github.com/vmware-tanzu/octant/pkg/store/fake"
)

func TestAPIService_Visit(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	service := testutil.CreateService("service")

	object := testutil.CreateAPIService("v1", "example.com")
	object.Spec.Service = &apiregistrationv1.ServiceReference{
		Namespace: service.Namespace,
		Name:      service.Name,
	}
	u := testutil.ToUnstructured(t, object)

	handler := fake.NewMockObjectHandler(controller)
	handler.EXPECT().
		AddEdge(gomock.Any(), u, testutil.ToUnstructured(t, service)).
		Return(nil)

	var visited []unstructured.Unstructured
	visitor := fake.NewMockVisitor(controller)
	visitor.EXPECT().
		Visit(gomock.Any(), gomock.Any(), handler, true).
		DoAndReturn(func(ctx context.Context, object *unstructured.Unstructured, handler objectvisitor.ObjectHandler, _ bool) error {
			visited = append(visited, *object)
			return nil
		})

	objectStore := objectStoreFake.NewMockStore(controller)

	key := store.Key{
		APIVersion: "v1",
		Kind:       "Service",
		Namespace:  service.Namespace,
		Name:       service.Name,
	}
	objectStore.EXPECT().
		Get(gomock.Any(), key).
		Return(testutil.ToUnstructured(t, service), nil)

	apiService := objectvisitor.NewAPIService(objectStore)

	ctx := context.Background()
	err := apiService.Visit(ctx, u, handler, visitor, true)

	sortObjectsByName(t, visited)

	expected := testutil.ToUnstructuredList(t, service)
	assert.Equal(t, expected.Items, visited)
	assert.NoError(t, err)
}

func TestAPIService_Visit_notfound(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	service := testutil.CreateService("service")

	object := testutil.CreateAPIService("v1", "example.com")
	object.Spec.Service = &apiregistrationv1.ServiceReference{
		Namespace: service.Namespace,
		Name:      service.Name,
	}
	u := testutil.ToUnstructured(t, object)

	handler := fake.NewMockObjectHandler(controller)

	var visited []unstructured.Unstructured
	visitor := fake.NewMockVisitor(controller)
	visitor.EXPECT().
		Visit(gomock.Any(), gomock.Any(), handler, true).
		DoAndReturn(func(ctx context.Context, object *unstructured.Unstructured, handler objectvisitor.ObjectHandler, _ bool) error {
			visited = append(visited, *object)
			return nil
		}).AnyTimes()

	objectStore := objectStoreFake.NewMockStore(controller)

	key := store.Key{
		APIVersion: "v1",
		Kind:       "Service",
		Namespace:  service.Namespace,
		Name:       service.Name,
	}
	objectStore.EXPECT().
		Get(gomock.Any(), key).
		Return(nil, kerrors.NewNotFound(schema.GroupResource{Resource: "services"}, service.Name))

	apiService := objectvisitor.NewAPIService(objectStore)

	ctx := context.Background()
	err := apiService.Visit(ctx, u, handler, visitor, true)

	sortObjectsByName(t, visited)

	expected := testutil.ToUnstructuredList(t)
	assert.Equal(t, expected.Items, visited)
	assert.NoError(t, err)
}

func TestAPIService_Visit_local(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	object := testutil.CreateAPIService("v1", "example.com")
	u := testutil.ToUnstructured(t, object)

	handler := fake.NewMockObjectHandler(controller)

	var visited []unstructured.Unstructured
	visitor := fake.NewMockVisitor(controller)
	visitor.EXPECT().
		Visit(gomock.Any(), gomock.Any(), handler, true).
		DoAndReturn(func(ctx context.Context, object *unstructured.Unstructured, handler objectvisitor.ObjectHandler, _ bool) error {
			visited = append(visited, *object)
			return nil
		}).AnyTimes()

	objectStore := objectStoreFake.NewMockStore(controller)

	apiService := objectvisitor.NewAPIService(objectStore)

	ctx := context.Background()
	err := apiService.Visit(ctx, u, handler, visitor, true)

	sortObjectsByName(t, visited)

	expected := testutil.ToUnstructuredList(t)
	assert.Equal(t, expected.Items, visited)
	assert.NoError(t, err)
}

func init() {
	apiregistrationv1.AddToScheme(scheme.Scheme)
}
