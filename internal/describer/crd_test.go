package describer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/heptio/developer-dash/internal/config"
	configFake "github.com/heptio/developer-dash/internal/config/fake"
	linkFake "github.com/heptio/developer-dash/internal/link/fake"
	"github.com/heptio/developer-dash/internal/modules/overview/printer"
	storefake "github.com/heptio/developer-dash/internal/objectstore/fake"
	"github.com/heptio/developer-dash/internal/queryer"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/pkg/objectstoreutil"
	pluginFake "github.com/heptio/developer-dash/pkg/plugin/fake"
	"github.com/heptio/developer-dash/pkg/view/component"
)

func Test_crd(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashConfig := configFake.NewMockDash(controller)

	objectStore := storefake.NewMockObjectStore(controller)
	dashConfig.EXPECT().ObjectStore().Return(objectStore)

	crdObject := testutil.CreateCRD("crd1")
	crdObject.Spec.Group = "foo.example.com"
	crdObject.Spec.Version = "v1"
	crdObject.Spec.Names.Kind = "Name"

	crdKey := objectstoreutil.Key{
		APIVersion: "apiextensions.k8s.io/v1beta1",
		Kind:       "CustomResourceDefinition",
		Name:       crdObject.Name,
	}

	objectStore.EXPECT().Get(gomock.Any(), gomock.Eq(crdKey)).Return(testutil.ToUnstructured(t, crdObject), nil)

	crKey := objectstoreutil.Key{
		Namespace:  "default",
		APIVersion: "foo.example.com/v1",
		Kind:       "Name",
	}

	object := testutil.CreateCustomResource("cr")
	objectStore.EXPECT().Get(gomock.Any(), gomock.Eq(crKey)).Return(object, nil)

	linkGenerator := linkFake.NewMockInterface(controller)

	pluginManager := pluginFake.NewMockManagerInterface(controller)
	dashConfig.EXPECT().PluginManager().Return(pluginManager)

	var tabs []component.Tab
	pluginManager.EXPECT().Tabs(object).Return(tabs, nil)

	crPrinter := func(cd *crd) {
		cd.summaryPrinter = func(ctx context.Context, crd *apiextv1beta1.CustomResourceDefinition, object *unstructured.Unstructured, options printer.Options) (component.Component, error) {
			return component.NewText("cr"), nil
		}

		cd.resourceViewerPrinter = func(ctx context.Context, object *unstructured.Unstructured, dashConfig config.Dash, q queryer.Queryer) (component.Component, error) {
			return component.NewText("rv"), nil
		}

		cd.yamlPrinter = func(runtime.Object) (*component.YAML, error) {
			return component.NewYAML(component.TitleFromString("yaml"), "data"), nil
		}
	}

	c := newCRD(crdObject.Name, "path", crPrinter)

	options := Options{
		Dash: dashConfig,
		Link: linkGenerator,
	}

	ctx := context.Background()

	got, err := c.Describe(ctx, "prefix", "default", options)
	require.NoError(t, err)

	expected := *component.NewContentResponse([]component.TitleComponent{
		component.NewText("Custom Resources"),
		component.NewText("crd1"),
		component.NewText("cr"),
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

	assertJSONEqual(t, expected, got)
}
