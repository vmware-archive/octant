package objectvisitor_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware-tanzu/octant/internal/objectvisitor"
	"github.com/vmware-tanzu/octant/internal/objectvisitor/fake"
	queryerFake "github.com/vmware-tanzu/octant/internal/queryer/fake"
	"github.com/vmware-tanzu/octant/internal/testutil"
)

func TestPod_Visit(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	serviceAccount := testutil.CreateServiceAccount("service-account")
	configMap := testutil.CreateConfigMap("configmap")
	secret := testutil.CreateSecret("secret")
	pvc := testutil.CreatePersistentVolumeClaim("pvc")

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
	q.EXPECT().
		ConfigMapsForPod(gomock.Any(), object).
		Return([]*corev1.ConfigMap{configMap}, nil)
	q.EXPECT().
		SecretsForPod(gomock.Any(), object).
		Return([]*corev1.Secret{secret}, nil)
	q.EXPECT().
		PersistentVolumeClaimsForPod(gomock.Any(), object).
		Return([]*corev1.PersistentVolumeClaim{pvc}, nil)

	handler := fake.NewMockObjectHandler(controller)
	handler.EXPECT().SetLevel(gomock.Any(), 1).Return(2)
	handler.EXPECT().
		AddEdge(gomock.Any(), u, testutil.ToUnstructured(t, service), gomock.Any()).
		Return(nil)
	handler.EXPECT().
		AddEdge(gomock.Any(), u, testutil.ToUnstructured(t, serviceAccount), gomock.Any()).
		Return(nil)
	handler.EXPECT().
		AddEdge(gomock.Any(), u, testutil.ToUnstructured(t, configMap), gomock.Any()).
		Return(nil)
	handler.EXPECT().AddEdge(gomock.Any(), u, testutil.ToUnstructured(t, secret), gomock.Any()).
		Return(nil)
	handler.EXPECT().AddEdge(gomock.Any(), u, testutil.ToUnstructured(t, pvc), gomock.Any()).
		Return(nil)

	var visited []unstructured.Unstructured
	visitor := fake.NewMockVisitor(controller)
	visitor.EXPECT().
		Visit(gomock.Any(), gomock.Any(), handler, true, gomock.Any()).
		DoAndReturn(func(ctx context.Context, object *unstructured.Unstructured, handler objectvisitor.ObjectHandler, _ bool, _ int) error {
			visited = append(visited, *object)
			return nil
		}).AnyTimes()

	pod := objectvisitor.NewPod(q)

	ctx := context.Background()
	err := pod.Visit(ctx, u, handler, visitor, true, 1)

	sortObjectsByName(t, visited)

	expected := testutil.ToUnstructuredList(t, service, serviceAccount)
	assert.Equal(t, expected.Items, visited)
	assert.NoError(t, err)
}
