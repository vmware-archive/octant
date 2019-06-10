package describer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	configFake "github.com/heptio/developer-dash/internal/config/fake"
	"github.com/heptio/developer-dash/internal/link"
	storefake "github.com/heptio/developer-dash/internal/objectstore/fake"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/pkg/objectstoreutil"
	"github.com/heptio/developer-dash/pkg/view/component"
)

func Test_crdListDescriber(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	o := storefake.NewMockObjectStore(controller)

	crd := testutil.CreateCRD("crd1")
	crd.Spec.Group = "foo.example.com"
	crd.Spec.Version = "v1"
	crd.Spec.Names.Kind = "Name"

	crdKey := objectstoreutil.Key{
		APIVersion: "apiextensions.k8s.io/v1beta1",
		Kind:       "CustomResourceDefinition",
		Name:       crd.Name,
	}

	o.EXPECT().HasAccess(gomock.Any(), "list").Return(nil)
	o.EXPECT().Get(gomock.Any(), gomock.Eq(crdKey)).Return(testutil.ToUnstructured(t, crd), nil)

	crKey := objectstoreutil.Key{
		Namespace:  "default",
		APIVersion: "foo.example.com/v1",
		Kind:       "Name",
	}

	var objects []*unstructured.Unstructured
	o.EXPECT().List(gomock.Any(), gomock.Eq(crKey)).Return(objects, nil)

	listPrinter := func(cld *crdList) {
		cld.printer = func(name string, crd *apiextv1beta1.CustomResourceDefinition, objects []*unstructured.Unstructured, linkGenerator link.Interface) (component.Component, error) {
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

	got, err := cld.Describe(ctx, "prefix", "default", options)
	require.NoError(t, err)

	expected := *component.NewContentResponse(nil)
	list := component.NewList("Custom Resources / crd1", []component.Component{
		component.NewText("crd list"),
	})
	expected.Add(list)

	assertJSONEqual(t, expected, got)
}
