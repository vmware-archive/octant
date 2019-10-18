package resourceviewer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	linkFake "github.com/vmware-tanzu/octant/internal/link/fake"
	"github.com/vmware-tanzu/octant/internal/objectstatus"
	"github.com/vmware-tanzu/octant/internal/resourceviewer/fake"
	"github.com/vmware-tanzu/octant/internal/testutil"
	pluginFake "github.com/vmware-tanzu/octant/pkg/plugin/fake"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_objectNode(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	deployment := testutil.ToUnstructured(t, testutil.CreateDeployment("deployment"))
	deploymentLink := component.NewLink("", deployment.GetName(), "/deployment")

	l := linkFake.NewMockInterface(controller)
	l.EXPECT().
		ForObjectWithQuery(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(deploymentLink, nil)

	pluginPrinter := pluginFake.NewMockManagerInterface(controller)
	objectStatus := fake.NewMockObjectStatus(controller)
	objectStatus.EXPECT().
		Status(gomock.Any(), gomock.Any()).
		Return(&objectstatus.ObjectStatus{}, nil)

	on := objectNode{
		link:          l,
		pluginPrinter: pluginPrinter,
		objectStatus:  objectStatus,
	}

	ctx := context.Background()

	got, err := on.Create(ctx, deployment)
	require.NoError(t, err)

	expected := &component.Node{
		Name:       deployment.GetName(),
		APIVersion: deployment.GetAPIVersion(),
		Kind:       deployment.GetKind(),
		Status:     component.NodeStatusOK,
		Path:       deploymentLink,
	}

	testutil.AssertJSONEqual(t, expected, got)
}
