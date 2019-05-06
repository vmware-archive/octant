package overview

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/heptio/developer-dash/internal/objectstore"
	storefake "github.com/heptio/developer-dash/internal/objectstore/fake"
	"github.com/heptio/developer-dash/internal/overview/printer"
	printerfake "github.com/heptio/developer-dash/internal/overview/printer/fake"
	"github.com/heptio/developer-dash/internal/queryer"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/pkg/objectstoreutil"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func Test_customResourceDefinitionNames(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	o := storefake.NewMockObjectStore(controller)

	crd1 := testutil.CreateCRD("crd1.example.com")
	crd2 := testutil.CreateCRD("crd2.example.com")

	crdList := []*unstructured.Unstructured{
		testutil.ToUnstructured(t, crd1),
		testutil.ToUnstructured(t, crd2),
	}

	crdKey := objectstoreutil.Key{
		APIVersion: "apiextensions.k8s.io/v1beta1",
		Kind:       "CustomResourceDefinition",
	}
	o.EXPECT().CheckAccess(gomock.Any()).Return(nil)
	o.EXPECT().List(gomock.Any(), gomock.Eq(crdKey)).Return(crdList, nil)

	ctx := context.Background()
	got, err := customResourceDefinitionNames(ctx, o)
	require.NoError(t, err)

	expected := []string{"crd1.example.com", "crd2.example.com"}

	assert.Equal(t, expected, got)
}

func Test_customResourceDefinition(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	o := storefake.NewMockObjectStore(controller)

	crd1 := testutil.CreateCRD("crd1.example.com")

	crdKey := objectstoreutil.Key{
		APIVersion: "apiextensions.k8s.io/v1beta1",
		Kind:       "CustomResourceDefinition",
		Name:       "crd1.example.com",
	}
	o.EXPECT().Get(gomock.Any(), gomock.Eq(crdKey)).Return(testutil.ToUnstructured(t, crd1), nil)

	name := "crd1.example.com"
	ctx := context.Background()
	got, err := customResourceDefinition(ctx, name, o)
	require.NoError(t, err)

	assert.Equal(t, crd1, got)
}

func Test_crdSectionDescriber(t *testing.T) {
	csd := newCRDSectionDescriber("/path", "title")

	d1View := component.NewText("d1")
	d1 := newStubDescriber("/d1", component.NewList("", []component.Component{d1View}))

	csd.Add("d1", d1)

	ctx := context.Background()

	view1, err := csd.Describe(ctx, "/prefix", "default", nil, DescriberOptions{})
	require.NoError(t, err)

	expect1 := component.ContentResponse{
		Title: component.TitleFromString("title"),
		Components: []component.Component{
			component.NewList("Custom Resources", []component.Component{d1View}),
		},
	}

	assert.Equal(t, expect1, view1)

	csd.Remove("d1")

	view2, err := csd.Describe(ctx, "/prefix", "default", nil, DescriberOptions{})
	require.NoError(t, err)

	expect2 := component.ContentResponse{
		Title: component.TitleFromString("title"),
		Components: []component.Component{
			component.NewList("Custom Resources", nil),
		},
	}

	assert.Equal(t, expect2, view2)
}

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

	o.EXPECT().CheckAccess(gomock.Any()).Return(nil)
	o.EXPECT().Get(gomock.Any(), gomock.Eq(crdKey)).Return(testutil.ToUnstructured(t, crd), nil)

	crKey := objectstoreutil.Key{
		Namespace:  "default",
		APIVersion: "foo.example.com/v1",
		Kind:       "Name",
	}

	objects := []*unstructured.Unstructured{}
	o.EXPECT().List(gomock.Any(), gomock.Eq(crKey)).Return(objects, nil)

	listPrinter := func(cld *crdListDescriber) {
		cld.printer = func(ctx context.Context, name, namespace string, crd *apiextv1beta1.CustomResourceDefinition, objects []*unstructured.Unstructured) (component.Component, error) {
			return component.NewText("crd list"), nil
		}
	}

	cld := newCRDListDescriber(crd.Name, "path", listPrinter)

	options := DescriberOptions{
		ObjectStore: o,
	}

	ctx := context.Background()

	got, err := cld.Describe(ctx, "prefix", "default", nil, options)
	require.NoError(t, err)

	expected := *component.NewContentResponse(nil)
	list := component.NewList("Custom Resources / crd1", []component.Component{
		component.NewText("crd list"),
	})
	expected.Add(list)

	assert.Equal(t, expected, got)
}

func Test_crdDescriber(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	o := storefake.NewMockObjectStore(controller)

	pluginPrinter := printerfake.NewMockPluginPrinter(controller)
	pluginPrinter.EXPECT().Tabs(gomock.Any()).Return([]component.Tab{}, nil)

	crd := testutil.CreateCRD("crd1")
	crd.Spec.Group = "foo.example.com"
	crd.Spec.Version = "v1"
	crd.Spec.Names.Kind = "Name"

	crdKey := objectstoreutil.Key{
		APIVersion: "apiextensions.k8s.io/v1beta1",
		Kind:       "CustomResourceDefinition",
		Name:       crd.Name,
	}

	o.EXPECT().Get(gomock.Any(), gomock.Eq(crdKey)).Return(testutil.ToUnstructured(t, crd), nil)

	crKey := objectstoreutil.Key{
		Namespace:  "default",
		APIVersion: "foo.example.com/v1",
		Kind:       "Name",
	}

	object := &unstructured.Unstructured{}
	o.EXPECT().Get(gomock.Any(), gomock.Eq(crKey)).Return(object, nil)

	crPrinter := func(cd *crdDescriber) {
		cd.summaryPrinter = func(ctx context.Context, crd *apiextv1beta1.CustomResourceDefinition, object *unstructured.Unstructured, options printer.Options) (component.Component, error) {
			return component.NewText("cr"), nil
		}

		cd.resourceViewerPrinter = func(ctx context.Context, object *unstructured.Unstructured, o objectstore.ObjectStore, q queryer.Queryer) (component.Component, error) {
			return component.NewText("rv"), nil
		}

		cd.yamlPrinter = func(runtime.Object) (*component.YAML, error) {
			return component.NewYAML(component.TitleFromString("yaml"), "data"), nil
		}
	}

	cd := newCRDDescriber(crd.Name, "path", crPrinter)

	options := DescriberOptions{
		ObjectStore:   o,
		PluginManager: pluginPrinter,
	}

	ctx := context.Background()

	got, err := cd.Describe(ctx, "prefix", "default", nil, options)
	require.NoError(t, err)

	expected := *component.NewContentResponse([]component.TitleComponent{
		component.NewLink("", "crd1", "/content/overview/namespace/default/custom-resources/crd1"),
		component.NewText(""),
	})

	crView := component.NewText("cr")
	crView.SetAccessor("summary")
	expected.Add(crView)
	rvView := component.NewText("rv")
	rvView.SetAccessor("resourceViewer")
	expected.Add(rvView)
	yView := component.NewYAML(component.TitleFromString("yaml"), "data")
	yView.SetAccessor("yaml")
	expected.Add(yView)

	assert.Equal(t, expected, got)
}
