package describer

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	configFake "github.com/heptio/developer-dash/internal/config/fake"
	printerfake "github.com/heptio/developer-dash/internal/modules/overview/printer/fake"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/pkg/objectstoreutil"
	"github.com/heptio/developer-dash/pkg/plugin"
	"github.com/heptio/developer-dash/pkg/view/component"
)

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
	objectPrinter.EXPECT().Print(gomock.Any(), pod, pluginManager).Return(podSummary, nil)

	options := Options{
		Dash:    dashConfig,
		Printer: objectPrinter,
		LoadObject: func(ctx context.Context, namespace string, fields map[string]string, objectStoreKey objectstoreutil.Key) (*unstructured.Unstructured, error) {
			return testutil.ToUnstructured(t, pod), nil
		},
	}

	d := NewObject(thePath, "object", key, podObjectType, true)

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
