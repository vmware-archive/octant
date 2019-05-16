package overview

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/go-plugin"
	clusterfake "github.com/heptio/developer-dash/internal/cluster/fake"
	storefake "github.com/heptio/developer-dash/internal/objectstore/fake"
	"github.com/heptio/developer-dash/internal/overview/printer"
	printerfake "github.com/heptio/developer-dash/internal/overview/printer/fake"
	pffake "github.com/heptio/developer-dash/internal/portforward/fake"
	"github.com/heptio/developer-dash/pkg/objectstoreutil"
	dashplugin "github.com/heptio/developer-dash/pkg/plugin"
	"github.com/heptio/developer-dash/pkg/plugin/fake"
	managerstorefake "github.com/heptio/developer-dash/pkg/plugin/fake"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestListDescriber(t *testing.T) {
	thePath := "/"
	key := objectstoreutil.Key{APIVersion: "v1", Kind: "Pod"}
	namespace := "default"
	fields := map[string]string{}

	controller := gomock.NewController(t)
	defer controller.Finish()

	o := storefake.NewMockObjectStore(controller)

	pf := pffake.NewMockPortForwarder(controller)

	ms := managerstorefake.NewMockManagerStore(controller)

	client := clusterfake.NewMockClientInterface(controller)

	retrieveKey := objectstoreutil.Key{Namespace: namespace, APIVersion: "v1", Kind: "Pod"}
	object := map[string]interface{}{
		"kind":       "Pod",
		"apiVersion": "v1",
		"metadata": map[string]interface{}{
			"name":              "name",
			"creationTimestamp": "2019-01-14T13:34:56+00:00",
		},
	}

	o.EXPECT().
		List(gomock.Any(), gomock.Eq(retrieveKey)).
		Return([]*unstructured.Unstructured{{Object: object}}, nil)

	listType := func() interface{} {
		return &corev1.PodList{}
	}

	objectType := func() interface{} {
		return &corev1.Pod{}
	}

	d := NewListDescriber(thePath, "list", key, listType, objectType, false)

	options := DescriberOptions{
		ObjectStore: o,
		Fields:      fields,
		Printer:     printer.NewResource(o, pf, ms),
	}

	ctx := context.Background()
	cResponse, err := d.Describe(ctx, "/path", namespace, client, options)
	require.NoError(t, err)

	list := component.NewList("list", nil)

	tableCols := component.NewTableCols("Name", "Labels", "Age")
	table := component.NewTable("/v1, Kind=PodList", tableCols)
	table.Add(component.TableRow{
		"Age":    component.NewTimestamp(time.Unix(1547472896, 0)),
		"Labels": component.NewLabels(nil),
		"Name":   component.NewText("name"),
	})
	list.Add(table)

	expected := component.ContentResponse{
		Components: []component.Component{list},
	}

	assert.Equal(t, expected, cResponse)
}

func TestObjectDescriber(t *testing.T) {
	thePath := "/"
	key := objectstoreutil.Key{APIVersion: "v1", Kind: "Pod"}
	namespace := "default"
	fields := map[string]string{}

	controller := gomock.NewController(t)
	defer controller.Finish()

	o := storefake.NewMockObjectStore(controller)
	clusterClient := clusterfake.NewMockClientInterface(controller)
	pf := pffake.NewMockPortForwarder(controller)
	pluginPrinter := printerfake.NewMockPluginPrinter(controller)
	ms := managerstorefake.NewMockManagerStore(controller)

	object := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "Pod",
			"apiVersion": "v1",
			"metadata": map[string]interface{}{
				"name":      "pod",
				"namespace": "default",
			},
			"spec": map[string]interface{}{
				"containers": []interface{}{
					map[string]interface{}{
						"name": "one",
					},
					map[string]interface{}{
						"name": "two",
					},
				},
			},
		},
	}

	retrieveKey := objectstoreutil.Key{Namespace: namespace, APIVersion: "v1", Kind: "Pod"}

	o.EXPECT().
		Get(gomock.Any(), gomock.Eq(retrieveKey)).
		Return(object, nil)

	pluginPrinter.EXPECT().Tabs(gomock.Any()).Return([]component.Tab{}, nil)

	objectType := func() interface{} {
		return &corev1.Pod{}
	}

	fn := DefaultLoader(key)

	d := NewObjectDescriber(thePath, "object", fn, objectType, true)

	p := printer.NewResource(o, pf, ms)
	err := p.Handler(func(context.Context, *corev1.Pod, printer.Options) (component.Component, error) {
		return component.NewText("*v1.Pod"), nil
	})
	require.NoError(t, err)

	options := DescriberOptions{
		ObjectStore:   o,
		Fields:        fields,
		Printer:       p,
		PluginPrinter: pluginPrinter,
	}

	ctx := context.Background()
	cResponse, err := d.Describe(ctx, "/path", namespace, clusterClient, options)
	require.NoError(t, err)

	summary := component.NewText("*v1.Pod")
	summary.SetAccessor("summary")

	yaml := component.NewYAML(component.TitleFromString("YAML"),
		"---\napiVersion: v1\nkind: Pod\nmetadata:\n  creationTimestamp: null\n  name: pod\n  namespace: default\nspec:\n  containers:\n  - name: one\n    resources: {}\n  - name: two\n    resources: {}\nstatus: {}\n")
	yaml.SetAccessor("yaml")

	logs := component.NewLogs("default", "pod", []string{"one", "two"})
	logs.SetAccessor("logs")

	expected := component.ContentResponse{
		Title: component.Title(component.NewText("object"), component.NewText("pod")),
		Components: []component.Component{
			summary,
			yaml,
			logs,
		},
	}
	assert.Equal(t, expected, cResponse)
}

func TestSectionDescriber(t *testing.T) {
	namespace := "default"

	controller := gomock.NewController(t)
	defer controller.Finish()

	clusterClient := clusterfake.NewMockClientInterface(controller)
	o := storefake.NewMockObjectStore(controller)

	options := DescriberOptions{
		ObjectStore: o,
	}

	ctx := context.Background()

	tests := []struct {
		name     string
		d        *SectionDescriber
		expected component.ContentResponse
	}{
		{
			name: "general",
			d: NewSectionDescriber(
				"/section",
				"section",
				newStubDescriber("/foo"),
			),
			expected: component.ContentResponse{
				Title: component.Title(component.NewText("section")),
				Components: []component.Component{
					component.NewList("section", nil),
				},
			},
		},
		{
			name: "empty",
			d: NewSectionDescriber(
				"/section",
				"section",
				newEmptyDescriber("/foo"),
				newEmptyDescriber("/bar"),
			),
			expected: component.ContentResponse{
				Title: component.Title(component.NewText("section")),
				Components: []component.Component{
					component.NewList("section", nil),
				},
			},
		},
		{
			name: "empty component",
			d: NewSectionDescriber(
				"/section",
				"section",
				newStubDescriber("/foo", &emptyComponent{}),
			),
			expected: component.ContentResponse{
				Title: component.Title(component.NewText("section")),
				Components: []component.Component{
					component.NewList("section", nil),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			got, err := tc.d.Describe(ctx, "/prefix", namespace, clusterClient, options)
			require.NoError(t, err)

			assert.Equal(t, tc.expected, got)
		})
	}

}

func TestPluginDescriber(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	name := "plugin-test"
	namespace := "default"
	metadata := dashplugin.Metadata{
		Name:         name,
		Description:  "this is a test",
		Capabilities: dashplugin.Capabilities{},
	}

	store := dashplugin.NewDefaultStore()
	clusterClient := clusterfake.NewMockClientInterface(controller)
	client := newFakePluginClient(name, controller)
	store.Store(name, client, metadata)

	p := NewPluginListDescriber()

	options := DescriberOptions{
		PluginManagerStore: store,
	}

	ctx := context.Background()
	cResponse, err := p.Describe(ctx, "/plugins", namespace, clusterClient, options)
	require.NoError(t, err)

	list := component.NewList("Plugins", nil)
	tableCols := component.NewTableCols("Name", "Description", "Capabilities")
	table := component.NewTable("Plugins", tableCols)
	table.Add(component.TableRow{
		"Name":        component.NewText(name),
		"Description": component.NewText("this is a test"),
		"Capability":  component.NewText("{\"SupportsPrinterConfig\":null,\"SupportsPrinterStatus\":null,\"SupportsPrinterItems\":null,\"SupportsObjectStatus\":null,\"SupportsTab\":null}"),
	})

	list.Add(table)

	expected := component.ContentResponse{
		Components: []component.Component{list},
	}

	assert.Equal(t, expected, cResponse)
}

type fakePluginClient struct {
	clientProtocol *fake.MockClientProtocol
	service        *fake.MockService
	name           string
}

var _ dashplugin.Client = (*fakePluginClient)(nil)

func newFakePluginClient(name string, controller *gomock.Controller) *fakePluginClient {
	service := fake.NewMockService(controller)
	metadata := dashplugin.Metadata{
		Name: name,
	}
	service.EXPECT().Register(gomock.Eq("localhost:54321")).Return(metadata, nil).AnyTimes()

	clientProtocol := fake.NewMockClientProtocol(controller)
	clientProtocol.EXPECT().Dispense("plugin").Return(service, nil).AnyTimes()

	return &fakePluginClient{
		service:        service,
		clientProtocol: clientProtocol,
		name:           name,
	}
}

func (c *fakePluginClient) Client() (plugin.ClientProtocol, error) {
	return c.clientProtocol, nil
}

func (c *fakePluginClient) Kill() {}
