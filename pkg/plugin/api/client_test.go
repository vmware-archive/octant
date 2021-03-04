package api_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/plugin/api"
	"github.com/vmware-tanzu/octant/pkg/plugin/api/fake"
	"github.com/vmware-tanzu/octant/pkg/plugin/api/proto"
)

func TestClient_Update(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ctx := context.Background()
	pod := testutil.ToUnstructured(t, testutil.CreatePod("pod"))

	podData, err := json.Marshal(pod)
	require.NoError(t, err)

	dashboardClient := fake.NewMockDashboardClient(controller)
	req := &proto.UpdateRequest{
		Object: podData,
	}
	dashboardClient.EXPECT().Update(gomock.Any(), req).Return(&proto.UpdateResponse{}, nil)

	conn := fake.NewMockDashboardConnection(controller)
	conn.EXPECT().Client().Return(dashboardClient)

	connOpt := MockDashboardConnection(conn)

	client, err := api.NewClient("address", connOpt)
	require.NoError(t, err)

	err = client.Update(ctx, pod)
	require.NoError(t, err)
}

func TestClient_ForceFrontendUpdate(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ctx := context.Background()

	dashboardClient := fake.NewMockDashboardClient(controller)
	dashboardClient.EXPECT().ForceFrontendUpdate(gomock.Any(), &proto.Empty{}).Return(nil, nil)

	conn := fake.NewMockDashboardConnection(controller)
	conn.EXPECT().Client().Return(dashboardClient)

	connOpt := MockDashboardConnection(conn)

	client, err := api.NewClient("address", connOpt)
	require.NoError(t, err)

	err = client.ForceFrontendUpdate(ctx)
	require.NoError(t, err)
}

func MockDashboardConnection(conn *fake.MockDashboardConnection) api.ClientOption {
	return func(client *api.Client) {
		client.DashboardConnection = conn
	}
}
