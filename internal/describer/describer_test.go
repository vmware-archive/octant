package describer

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	configFake "github.com/heptio/developer-dash/internal/config/fake"
	printerfake "github.com/heptio/developer-dash/internal/modules/overview/printer/fake"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/pkg/objectstoreutil"
	"github.com/heptio/developer-dash/pkg/plugin"
	"github.com/heptio/developer-dash/pkg/view/component"
)

func TestListDescriber(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	thePath := "/"

	pod := testutil.CreatePod("pod")
	pod.CreationTimestamp = metav1.Time{
		Time: time.Unix(1547472896, 0),
	}

	key, err := objectstoreutil.KeyFromObject(pod)
	require.NoError(t, err)

	ctx := context.Background()
	namespace := "default"

	dashConfig := configFake.NewMockDash(controller)
	pluginManager := plugin.NewManager(nil)
	dashConfig.EXPECT().PluginManager().Return(pluginManager)

	podListTable := createPodTable(*pod)

	objectPrinter := printerfake.NewMockPrinter(controller)
	podList := &corev1.PodList{Items: []corev1.Pod{*pod}}
	objectPrinter.EXPECT().Print(gomock.Any(), podList, pluginManager).Return(podListTable, nil)

	options := Options{
		Dash:    dashConfig,
		Printer: objectPrinter,
		LoadObjects: func(ctx context.Context, namespace string, fields map[string]string, objectStoreKeys []objectstoreutil.Key) ([]*unstructured.Unstructured, error) {
			return testutil.ToUnstructuredList(t, pod), nil
		},
	}

	d := NewList(thePath, "list", key, podListType, podObjectType, false)
	cResponse, err := d.Describe(ctx, "/path", namespace, options)
	require.NoError(t, err)

	list := component.NewList("list", nil)
	list.Add(podListTable)
	expected := component.ContentResponse{
		Components: []component.Component{list},
	}

	assert.Equal(t, expected, cResponse)
}

func TestObjectDescriber(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ctx := context.Background()
	thePath := "/"

	pod := testutil.CreatePod("pod")
	pod.CreationTimestamp = metav1.Time{
		Time: time.Unix(1547472896, 0),
	}

	key, err := objectstoreutil.KeyFromObject(pod)
	require.NoError(t, err)

	dashConfig := configFake.NewMockDash(controller)
	pluginManager := plugin.NewManager(nil)
	dashConfig.EXPECT().PluginManager().Return(pluginManager).AnyTimes()

	objectPrinter := printerfake.NewMockPrinter(controller)

	podSummary := component.NewText("summary")
	objectPrinter.EXPECT().Print(gomock.Any(),  pod, pluginManager).Return(podSummary, nil)

	options := Options{
		Dash:    dashConfig,
		Printer: objectPrinter,
		LoadObject: func(ctx context.Context, namespace string, fields map[string]string, objectStoreKey objectstoreutil.Key) (*unstructured.Unstructured, error) {
			return testutil.ToUnstructured(t, pod), nil
		},
	}

	d := NewObjectDescriber(thePath, "object", key, podObjectType, true)

	d.tabFuncDescriptors = []tabFuncDescriptor{
		{name: "summary", tabFunc: d.addSummaryTab},
	}

	cResponse, err := d.Describe(ctx, "/path", pod.Namespace, options)
	require.NoError(t, err)

	summary := component.NewText("summary")
	summary.SetAccessor("summary")

	expected := component.ContentResponse{
		Title: component.Title(component.NewText("object"), component.NewText("pod")),
		Components: []component.Component{
			summary,
		},
	}
	assert.Equal(t, expected, cResponse)

}

func TestSectionDescriber(t *testing.T) {
	namespace := "default"

	controller := gomock.NewController(t)
	defer controller.Finish()

	dashConfig := configFake.NewMockDash(controller)

	options := Options{
		Dash: dashConfig,
	}

	ctx := context.Background()

	tests := []struct {
		name     string
		d        *Section
		expected component.ContentResponse
	}{
		{
			name: "general",
			d: NewSectionDescriber(
				"/section",
				"section",
				NewStubDescriber("/foo"),
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
				NewEmptyDescriber("/foo"),
				NewEmptyDescriber("/bar"),
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
				NewStubDescriber("/foo", &emptyComponent{}),
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

			got, err := tc.d.Describe(ctx, "/prefix", namespace, options)
			require.NoError(t, err)

			assert.Equal(t, tc.expected, got)
		})
	}

}

type emptyComponent struct{}

var _ component.Component = (*emptyComponent)(nil)

func (c *emptyComponent) GetMetadata() component.Metadata {
	return component.Metadata{
		Type: "empty",
	}
}

func (c *emptyComponent) SetAccessor(string) {
	// no-op
}

func (c *emptyComponent) IsEmpty() bool {
	return true
}

func (c *emptyComponent) String() string {
	return ""
}

func (c *emptyComponent) LessThan(interface{}) bool {
	return false
}

func (c emptyComponent) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})

	return json.Marshal(m)
}

func createPodTable(pods ...corev1.Pod) *component.Table {
	tableCols := component.NewTableCols("Name", "Labels", "Age")
	table := component.NewTable("/v1, Kind=PodList", tableCols)
	for _, pod := range pods {
		table.Add(component.TableRow{
			"Age":    component.NewTimestamp(pod.CreationTimestamp.Time),
			"Labels": component.NewLabels(pod.Labels),
			"Name":   component.NewText(pod.Name),
		})
	}

	return table
}

func podListType() interface{} {
	return &corev1.PodList{}
}

func podObjectType() interface{} {
	return &corev1.Pod{}
}
