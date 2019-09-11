/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	configFake "github.com/vmware/octant/internal/config/fake"
	"github.com/vmware/octant/internal/link"
	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/icon"
	"github.com/vmware/octant/pkg/store"
	storefake "github.com/vmware/octant/pkg/store/fake"
	"github.com/vmware/octant/pkg/view/component"
)

func Test_crdListDescriber(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	o := storefake.NewMockStore(controller)

	crd := testutil.CreateCRD("crd1")
	crd.Spec.Group = "foo.example.com"
	crd.Spec.Version = "v1"
	crd.Spec.Names.Kind = "Name"

	crdKey := store.Key{
		APIVersion: "apiextensions.k8s.io/v1beta1",
		Kind:       "CustomResourceDefinition",
		Name:       crd.Name,
	}

	o.EXPECT().Get(gomock.Any(), gomock.Eq(crdKey)).Return(testutil.ToUnstructured(t, crd), true, nil)

	crKey := store.Key{
		Namespace:  "default",
		APIVersion: "foo.example.com/v1",
		Kind:       "Name",
	}

	objects := &unstructured.UnstructuredList{}
	o.EXPECT().List(gomock.Any(), gomock.Eq(crKey)).Return(objects, false, nil)

	listPrinter := func(cld *crdList) {
		cld.printer = func(name string, crd *apiextv1beta1.CustomResourceDefinition, objects *unstructured.UnstructuredList, linkGenerator link.Interface, loading bool) (component.Component, error) {
			return component.NewText("crd list"), nil
		}
	}

	dashConfig := configFake.NewMockDash(controller)
	dashConfig.EXPECT().ObjectStore().Return(o).AnyTimes()

	options := Options{
		Dash: dashConfig,
	}
	cld := newCRDList(crd.Name, "path", listPrinter)

	ctx := context.Background()

	got, err := cld.Describe(ctx, "default", options)
	require.NoError(t, err)

	expected := *component.NewContentResponse(nil)
	list := component.NewList("Custom Resources / crd1", []component.Component{
		component.NewText("crd list"),
	})
	iconName, iconSource := loadIcon(icon.CustomResourceDefinition)
	list.SetIcon(iconName, iconSource)
	expected.Add(list)

	testutil.AssertJSONEqual(t, expected, got)
}
