package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware-tanzu/octant/internal/gvk"
	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service/fake"
)

func TestHandler_Register(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashboard := fake.NewMockDashboard(controller)
	factory := func(string) (Dashboard, error) {
		return dashboard, nil
	}

	capabilities := &plugin.Capabilities{
		SupportsPrinterConfig: []schema.GroupVersionKind{gvk.Pod},
	}

	h := Handler{
		name:             "name",
		description:      "description",
		capabilities:     capabilities,
		dashboardFactory: factory,
	}

	ctx := context.Background()
	got, err := h.Register(ctx, "address")
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
		SupportsPrinterConfig: []schema.GroupVersionKind{gvk.Pod},
	}

	h := Handler{
		name:             "name",
		description:      "description",
		capabilities:     capabilities,
		dashboardFactory: factory,
	}

	ctx := context.Background()
	_, err := h.Register(ctx, "address")
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

	ctx := context.Background()
	got, err := h.Print(ctx, pod)
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
			Print: func(r *PrintRequest) (response plugin.PrintResponse, e error) {
				ran = true
				assert.Equal(t, dashboardClient, r.DashboardClient)
				assert.Equal(t, pod, r.Object)
				return plugin.PrintResponse{}, nil
			},
		},
		dashboardClient: dashboardClient,
	}

	ctx := context.Background()
	got, err := h.Print(ctx, pod)
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

	ctx := context.Background()
	got, err := h.PrintTab(ctx, pod)
	require.NoError(t, err)

	expected := plugin.TabResponse{}

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
			PrintTab: func(r *PrintRequest) (plugin.TabResponse, error) {
				ran = true
				assert.Equal(t, dashboardClient, r.DashboardClient)
				assert.Equal(t, pod, r.Object)
				return plugin.TabResponse{}, nil
			},
		},
	}

	ctx := context.Background()
	got, err := h.PrintTab(ctx, pod)
	require.NoError(t, err)

	expected := plugin.TabResponse{}
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

	ctx := context.Background()
	got, err := h.ObjectStatus(ctx, pod)
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
			ObjectStatus: func(r *PrintRequest) (response plugin.ObjectStatusResponse, e error) {
				ran = true
				assert.Equal(t, dashboardClient, r.DashboardClient)
				assert.Equal(t, pod, r.Object)
				return plugin.ObjectStatusResponse{}, nil
			},
		},
	}

	ctx := context.Background()
	got, err := h.ObjectStatus(ctx, pod)
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

	actionName := "action.octant.dev/testDefault"
	payload := action.Payload{"foo": "bar"}

	ctx := context.Background()
	err := h.HandleAction(ctx, actionName, payload)
	require.NoError(t, err)
}

func TestHandler_HandleAction_using_supplied_function(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashboardClient := fake.NewMockDashboard(controller)

	actionName := "action.octant.dev/testAction"
	payload := action.Payload{"foo": "bar"}

	ran := false

	h := Handler{
		dashboardClient: dashboardClient,
		HandlerFuncs: HandlerFuncs{
			HandleAction: func(r *ActionRequest) error {
				ran = true
				assert.Equal(t, dashboardClient, r.DashboardClient)
				assert.Equal(t, payload, r.Payload)

				return nil
			},
		},
	}

	ctx := context.Background()
	err := h.HandleAction(ctx, actionName, payload)
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

	ctx := context.Background()
	got, err := h.Navigation(ctx)
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
			Navigation: func(r *NavigationRequest) (navigation.Navigation, error) {
				ran = true
				assert.Equal(t, dashboardClient, r.DashboardClient)
				return navigation.Navigation{}, nil
			},
		},
	}

	ctx := context.Background()
	got, err := h.Navigation(ctx)
	require.NoError(t, err)

	expected := navigation.Navigation{}
	assert.Equal(t, expected, got)
	assert.True(t, ran)
}
