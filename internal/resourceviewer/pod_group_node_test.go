package resourceviewer

import (
	"context"
	"testing"

	linkFake "github.com/vmware-tanzu/octant/internal/link/fake"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/objectstatus"
	"github.com/vmware-tanzu/octant/internal/resourceviewer/fake"
	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_podGroupNode(t *testing.T) {
	controller := gomock.NewController(t)
	linkInterface := linkFake.NewMockInterface(controller)
	defer controller.Finish()

	pod := testutil.ToUnstructured(t, testutil.CreatePod("pod"))
	objectStatus := fake.NewMockObjectStatus(controller)
	objectStatus.EXPECT().
		Status(gomock.Any(), pod, gomock.Any()).
		Return(&objectstatus.ObjectStatus{}, nil)

	pgn := podGroupNode{objectStatus: objectStatus}

	objects := testutil.ToUnstructuredList(t, pod)
	name := "foo pods"

	ctx := context.Background()

	got, err := pgn.Create(ctx, name, objects.Items, linkInterface)
	require.NoError(t, err)

	podStatus := component.NewPodStatus()
	podStatus.AddSummary(pod.GetName(), nil, nil, component.NodeStatusOK)

	expected := &component.Node{
		Name:       "foo pods",
		APIVersion: "v1",
		Kind:       "Pod",
		Status:     component.NodeStatusOK,
		Details:    []component.Component{podStatus},
	}

	testutil.AssertJSONEqual(t, expected, got)
}
