package objectvisitor_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware/octant/internal/modules/overview/objectvisitor"
	"github.com/vmware/octant/internal/modules/overview/objectvisitor/fake"
	queryerFake "github.com/vmware/octant/internal/queryer/fake"
	"github.com/vmware/octant/internal/testutil"
)

func TestPod_Visit(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	serviceAccount := testutil.CreateServiceAccount("service-account")

	object := testutil.CreatePod("pod")
	object.Spec.ServiceAccountName = serviceAccount.Name
	u := testutil.ToUnstructured(t, object)

	q := queryerFake.NewMockQueryer(controller)
	service := testutil.CreateService("service")
	q.EXPECT().
		ServicesForPod(gomock.Any(), object).
		Return([]*corev1.Service{service}, nil)
	q.EXPECT().
		ServiceAccountForPod(gomock.Any(), object).
		Return(serviceAccount, nil)

	handler := fake.NewMockObjectHandler(controller)
	handler.EXPECT().
		AddEdge(gomock.Any(), u, service).
		Return(nil)
	handler.EXPECT().
		AddEdge(gomock.Any(), u, serviceAccount).
		Return(nil)
	handler.EXPECT().
		Process(gomock.Any(), u).Return(nil)

	var visited []runtime.Object
	visitor := fake.NewMockVisitor(controller)
	visitor.EXPECT().
		Visit(gomock.Any(), gomock.Any(), handler).
		DoAndReturn(func(ctx context.Context, object runtime.Object, handler objectvisitor.ObjectHandler) error {
			visited = append(visited, object)
			return nil
		}).AnyTimes()

	pod := objectvisitor.NewPod(q)

	ctx := context.Background()
	err := pod.Visit(ctx, u, handler, visitor)

	sortObjectsByName(t, visited)

	expected := []runtime.Object{service, serviceAccount}
	assert.Equal(t, expected, visited)
	assert.NoError(t, err)
}
