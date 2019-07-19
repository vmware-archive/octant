package service

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware/octant/internal/gvk"
	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/action"
	"github.com/vmware/octant/pkg/navigation"
	"github.com/vmware/octant/pkg/plugin"
	"github.com/vmware/octant/pkg/plugin/service/fake"
	"github.com/vmware/octant/pkg/view/component"
)

func TestHandler_Register(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashboard := fake.NewMockDashboard(controller)
	factory := func(string) (Dashboard, error) {
		return dashboard, nil
	}

	capabilities := &plugin.Capabilities{
		SupportsPrinterConfig: []schema.GroupVersionKind{gvk.PodGVK},
	}

	h := Handler{
		name:             "name",
		description:      "description",
		capabilities:     capabilities,
		dashboardFactory: factory,
	}

	got, err := h.Register("address")
	require.NoError(t, err)

	expected := plugin.Metadata{
		Name:         "name",
		Description:  "description",
		Capabilities: *capabilities,
	}

	require.Equal(t, expected, got)
}

func TestHandler_Register_with_dashboard_factory_failure(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	factory := func(string) (Dashboard, error) {
		return nil, errors.New("failure")
	}

	capabilities := &plugin.Capabilities{
		SupportsPrinterConfig: []schema.GroupVersionKind{gvk.PodGVK},
	}

	h := Handler{
		name:             "name",
		description:      "description",
		capabilities:     capabilities,
		dashboardFactory: factory,
	}

	_, err := h.Register("address")
	require.Error(t, err)
}

func TestHandler_Print_default(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashboardClient := fake.NewMockDashboard(controller)

	h := Handler{
		dashboardClient: dashboardClient,
	}

	pod := testutil.CreatePod("pod")

	got, err := h.Print(pod)
	require.NoError(t, err)

	expected := plugin.PrintResponse{}

	require.Equal(t, expected, got)
}

func TestHandler_Print_using_supplied_function(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	pod := testutil.CreatePod("pod")

	dashboardClient := fake.NewMockDashboard(controller)

	ran := false
	h := Handler{
		HandlerFuncs: HandlerFuncs{
			Print: func(gotClient Dashboard, gotObject runtime.Object) (response plugin.PrintResponse, e error) {
				ran = true
				assert.Equal(t, dashboardClient, gotClient)
				assert.Equal(t, pod, gotObject)
				return plugin.PrintResponse{}, nil
			},
		},
		dashboardClient: dashboardClient,
	}

	got, err := h.Print(pod)
	require.NoError(t, err)

	expected := plugin.PrintResponse{}

	assert.Equal(t, expected, got)
	assert.True(t, ran)
}

func TestHandler_PrintTab_default(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashboardClient := fake.NewMockDashboard(controller)

	h := Handler{
		dashboardClient: dashboardClient,
	}

	pod := testutil.CreatePod("pod")

	got, err := h.PrintTab(pod)
	require.NoError(t, err)

	expected := &component.Tab{}

	require.Equal(t, expected, got)
}

func TestHandler_PrintTab_using_supplied_function(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	pod := testutil.CreatePod("pod")

	dashboardClient := fake.NewMockDashboard(controller)

	ran := false

	h := Handler{
		dashboardClient: dashboardClient,
		HandlerFuncs: HandlerFuncs{
			PrintTab: func(gotClient Dashboard, gotObject runtime.Object) (tab *component.Tab, e error) {
				ran = true
				assert.Equal(t, dashboardClient, gotClient)
				assert.Equal(t, pod, gotObject)
				return &component.Tab{}, nil
			},
		},
	}

	got, err := h.PrintTab(pod)
	require.NoError(t, err)

	expected := &component.Tab{}
	assert.Equal(t, expected, got)
	assert.True(t, ran)
}

func TestHandler_ObjectStatus_default(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashboardClient := fake.NewMockDashboard(controller)

	h := Handler{
		dashboardClient: dashboardClient,
	}

	pod := testutil.CreatePod("pod")

	got, err := h.ObjectStatus(pod)
	require.NoError(t, err)

	expected := plugin.ObjectStatusResponse{}

	require.Equal(t, expected, got)
}

func TestHandler_ObjectStatus_using_supplied_function(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashboardClient := fake.NewMockDashboard(controller)
	pod := testutil.CreatePod("pod")

	ran := false

	h := Handler{
		dashboardClient: dashboardClient,
		HandlerFuncs: HandlerFuncs{
			ObjectStatus: func(gotClient Dashboard, gotObject runtime.Object) (response plugin.ObjectStatusResponse, e error) {
				ran = true
				assert.Equal(t, dashboardClient, gotClient)
				assert.Equal(t, pod, gotObject)
				return plugin.ObjectStatusResponse{}, nil
			},
		},
	}

	got, err := h.ObjectStatus(pod)
	require.NoError(t, err)

	expected := plugin.ObjectStatusResponse{}
	assert.Equal(t, expected, got)
	assert.True(t, ran)
}

func TestHandler_HandleAction_default(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashboardClient := fake.NewMockDashboard(controller)

	h := Handler{
		dashboardClient: dashboardClient,
	}

	payload := action.Payload{"foo": "bar"}

	err := h.HandleAction(payload)
	require.NoError(t, err)
}

func TestHandler_HandleAction_using_supplied_function(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashboardClient := fake.NewMockDashboard(controller)

	payload := action.Payload{"foo": "bar"}

	ran := false

	h := Handler{
		dashboardClient: dashboardClient,
		HandlerFuncs: HandlerFuncs{
			HandleAction: func(gotClient Dashboard, gotPayload action.Payload) error {
				ran = true
				assert.Equal(t, dashboardClient, gotClient)
				assert.Equal(t, payload, gotPayload)

				return nil
			},
		},
	}

	err := h.HandleAction(payload)
	assert.NoError(t, err)
	assert.True(t, ran)
}

func TestHandler_Navigation_default(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashboardClient := fake.NewMockDashboard(controller)

	h := Handler{
		dashboardClient: dashboardClient,
	}

	got, err := h.Navigation()
	require.NoError(t, err)

	expected := navigation.Navigation{}

	require.Equal(t, expected, got)
}

func TestHandler_Navigation_using_supplied_function(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashboardClient := fake.NewMockDashboard(controller)

	ran := false

	h := Handler{
		dashboardClient: dashboardClient,
		HandlerFuncs: HandlerFuncs{
			Navigation: func(gotClient Dashboard) (navigation.Navigation, error) {
				ran = true
				assert.Equal(t, dashboardClient, gotClient)
				return navigation.Navigation{}, nil
			},
		},
	}

	got, err := h.Navigation()
	require.NoError(t, err)

	expected := navigation.Navigation{}
	assert.Equal(t, expected, got)
	assert.True(t, ran)
}

func TestHandler_Content_default(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashboardClient := fake.NewMockDashboard(controller)

	h := Handler{
		dashboardClient: dashboardClient,
	}

	contentPath := "/test-path-default"

	got, err := h.Content(contentPath)
	require.NoError(t, err)

	expected := component.ContentResponse{}

	require.Equal(t, expected, got)
}

func TestHandler_Content_using_supplied_function(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashboardClient := fake.NewMockDashboard(controller)

	ran := false
	contentPath := "/test-path"

	h := Handler{
		dashboardClient: dashboardClient,
		HandlerFuncs: HandlerFuncs{
			Content: func(gotDashboard Dashboard, gotPath string) (component.ContentResponse, error) {
				ran = true
				assert.Equal(t, contentPath, gotPath)
				assert.Equal(t, dashboardClient, gotDashboard)
				return component.ContentResponse{}, nil
			},
		},
	}

	got, err := h.Content(contentPath)
	require.NoError(t, err)

	expected := component.ContentResponse{}
	assert.Equal(t, expected, got)
	assert.True(t, ran)
}
