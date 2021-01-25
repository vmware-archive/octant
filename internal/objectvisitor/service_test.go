package objectvisitor_test

import (
	"context"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"

	"github.com/vmware-tanzu/octant/internal/objectvisitor"
	"github.com/vmware-tanzu/octant/internal/objectvisitor/fake"
	queryerFake "github.com/vmware-tanzu/octant/internal/queryer/fake"
	"github.com/vmware-tanzu/octant/internal/testutil"
)

func TestService_Visit(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	object := testutil.CreateService("service")
	u := testutil.ToUnstructured(t, object)

	q := queryerFake.NewMockQueryer(controller)
	ingress := testutil.CreateIngress("ingress")
	q.EXPECT().
		IngressesForService(gomock.Any(), object).
		Return([]*extv1beta1.Ingress{ingress}, nil)
	pod := testutil.CreatePod("pod")
	q.EXPECT().
		PodsForService(gomock.Any(), object).
		Return([]*corev1.Pod{pod}, nil)
	apiService := testutil.CreateAPIService("v1", "apps")
	q.EXPECT().
		APIServicesForService(gomock.Any(), object).
		Return([]*apiregistrationv1.APIService{apiService}, nil)
	mutatingWebhookConfiguration := testutil.CreateMutatingWebhookConfiguration("mutatingWebhookConfiguration")
	q.EXPECT().
		MutatingWebhookConfigurationsForService(gomock.Any(), object).
		Return([]*admissionregistrationv1beta1.MutatingWebhookConfiguration{mutatingWebhookConfiguration}, nil)
	validatingWebhookConfiguration := testutil.CreateValidatingWebhookConfiguration("validatingWebhookConfiguration")
	q.EXPECT().
		ValidatingWebhookConfigurationsForService(gomock.Any(), object).
		Return([]*admissionregistrationv1beta1.ValidatingWebhookConfiguration{validatingWebhookConfiguration}, nil)

	handler := fake.NewMockObjectHandler(controller)
	handler.EXPECT().
		AddEdge(gomock.Any(), u, testutil.ToUnstructured(t, ingress)).
		Return(nil)
	handler.EXPECT().
		AddEdge(gomock.Any(), u, testutil.ToUnstructured(t, pod)).
		Return(nil)
	handler.EXPECT().
		AddEdge(gomock.Any(), u, testutil.ToUnstructured(t, apiService)).
		Return(nil)
	handler.EXPECT().
		AddEdge(gomock.Any(), u, testutil.ToUnstructured(t, mutatingWebhookConfiguration)).
		Return(nil)
	handler.EXPECT().
		AddEdge(gomock.Any(), u, testutil.ToUnstructured(t, validatingWebhookConfiguration)).
		Return(nil)

	var visited []unstructured.Unstructured
	var m sync.Mutex
	visitor := fake.NewMockVisitor(controller)
	visitor.EXPECT().
		Visit(gomock.Any(), gomock.Any(), handler, gomock.Any()).
		DoAndReturn(func(ctx context.Context, object *unstructured.Unstructured, handler objectvisitor.ObjectHandler, _ bool) error {
			m.Lock()
			defer m.Unlock()
			visited = append(visited, *object)
			return nil
		}).
		AnyTimes()

	service := objectvisitor.NewService(q)

	ctx := context.Background()

	err := service.Visit(ctx, u, handler, visitor, true)

	sortObjectsByName(t, visited)
	expected := testutil.ToUnstructuredList(t, ingress, mutatingWebhookConfiguration, pod, apiService, validatingWebhookConfiguration)
	assert.Equal(t, expected.Items, visited)
	assert.NoError(t, err)
}
