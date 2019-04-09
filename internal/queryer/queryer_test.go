package queryer

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/golang/mock/gomock"
	storefake "github.com/heptio/developer-dash/internal/objectstore/fake"
	fakequeryer "github.com/heptio/developer-dash/internal/queryer/fake"
	"github.com/heptio/developer-dash/pkg/cacheutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestCacheQueryer_Children(t *testing.T) {
	deployment := &appsv1.Deployment{
		TypeMeta:   metav1.TypeMeta{APIVersion: "apps/v1", Kind: "Deployment"},
		ObjectMeta: metav1.ObjectMeta{Name: "deployment", Namespace: "default"},
	}

	rs := &appsv1.ReplicaSet{
		TypeMeta:   metav1.TypeMeta{APIVersion: "apps/v1", Kind: "ReplicaSet"},
		ObjectMeta: metav1.ObjectMeta{Name: "rs", Namespace: "default"},
	}
	rs.SetOwnerReferences([]metav1.OwnerReference{
		*metav1.NewControllerRef(deployment, deployment.GroupVersionKind()),
	})

	resourceLists := []*metav1.APIResourceList{
		nil,
		{
			GroupVersion: "apps/v1",
			APIResources: []metav1.APIResource{
				{
					Namespaced: true,
					Kind:       "Deployment",
					Verbs:      metav1.Verbs{"watch", "list"},
				},
				{
					Namespaced: true,
					Kind:       "ReplicaSet",
					Verbs:      metav1.Verbs{"watch", "list"},
				},
				{
					Namespaced: true,
					Kind:       "NotListable",
					Verbs:      metav1.Verbs{"get"},
				},
			},
		},
		{
			GroupVersion: "v1",
			APIResources: []metav1.APIResource{
				{Namespaced: false, Kind: "Namespace"},
			},
		},
	}

	cases := []struct {
		name     string
		owner    metav1.Object
		setup    func(t *testing.T, c *storefake.MockObjectStore, disco *fakequeryer.MockDiscoveryInterface)
		expected func(t *testing.T) []runtime.Object
		isErr    bool
	}{
		{
			name:  "in general",
			owner: deployment,
			setup: func(t *testing.T, o *storefake.MockObjectStore, disco *fakequeryer.MockDiscoveryInterface) {
				deploymentKey := cacheutil.Key{Namespace: "default", APIVersion: "apps/v1", Kind: "Deployment"}
				o.EXPECT().
					List(gomock.Any(), gomock.Eq(deploymentKey)).
					Return([]*unstructured.Unstructured{toUnstructured(t, deployment)}, nil).Times(1)

				rsKey := cacheutil.Key{Namespace: "default", APIVersion: "apps/v1", Kind: "ReplicaSet"}
				o.EXPECT().
					List(gomock.Any(), gomock.Eq(rsKey)).
					Return([]*unstructured.Unstructured{toUnstructured(t, rs)}, nil).Times(1)

				disco.EXPECT().
					ServerResources().
					Return(resourceLists, nil).AnyTimes()

			},
			expected: func(t *testing.T) []runtime.Object {
				return []runtime.Object{
					toUnstructured(t, rs),
				}
			},
		},
		{
			name:  "owner is nil",
			owner: nil,
			isErr: true,
		},
		{
			name:  "fetch resource lists failure",
			owner: deployment,
			setup: func(t *testing.T, o *storefake.MockObjectStore, disco *fakequeryer.MockDiscoveryInterface) {
				disco.EXPECT().
					ServerResources().
					Return(nil, errors.New("failed")).AnyTimes()
			},
			isErr: true,
		},
		{
			name:  "objectstore list fails",
			owner: deployment,
			setup: func(t *testing.T, o *storefake.MockObjectStore, disco *fakequeryer.MockDiscoveryInterface) {
				deploymentKey := cacheutil.Key{Namespace: "default", APIVersion: "apps/v1", Kind: "Deployment"}
				o.EXPECT().
					List(gomock.Any(), gomock.Eq(deploymentKey)).
					Return(nil, errors.New("failed")).Times(1)

				rsKey := cacheutil.Key{Namespace: "default", APIVersion: "apps/v1", Kind: "ReplicaSet"}
				o.EXPECT().
					List(gomock.Any(), gomock.Eq(rsKey)).
					Return([]*unstructured.Unstructured{toUnstructured(t, rs)}, nil).AnyTimes()

				disco.EXPECT().
					ServerResources().
					Return(resourceLists, nil).AnyTimes()
			},
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			o := storefake.NewMockObjectStore(controller)
			discovery := fakequeryer.NewMockDiscoveryInterface(controller)

			if tc.setup != nil {
				tc.setup(t, o, discovery)
			}

			cq := New(o, discovery)

			ctx := context.Background()
			got, err := cq.Children(ctx, tc.owner)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expected(t), got)
		})
	}
}

func TestCacheQueryer_Events(t *testing.T) {
	deployment := &appsv1.Deployment{
		TypeMeta:   metav1.TypeMeta{APIVersion: "apps/v1", Kind: "Deployment"},
		ObjectMeta: metav1.ObjectMeta{Name: "deployment", Namespace: "default"},
	}

	rs := &appsv1.ReplicaSet{
		TypeMeta:   metav1.TypeMeta{APIVersion: "apps/v1", Kind: "ReplicaSet"},
		ObjectMeta: metav1.ObjectMeta{Name: "rs", Namespace: "default"},
	}

	var events []*corev1.Event
	for i := 0; i < 3; i++ {
		events = append(events, genEventFor(t, deployment, fmt.Sprintf("event-%d", i)))
	}

	events = append(events, genEventFor(t, rs, fmt.Sprintf("event-rs")))

	cases := []struct {
		name     string
		object   metav1.Object
		setup    func(t *testing.T, o *storefake.MockObjectStore)
		isErr    bool
		expected []string
	}{
		{
			name:   "in general",
			object: deployment,
			setup: func(t *testing.T, o *storefake.MockObjectStore) {
				key := cacheutil.Key{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "Event",
				}
				o.EXPECT().
					List(gomock.Any(), gomock.Eq(key)).
					Return([]*unstructured.Unstructured{
						toUnstructured(t, events[0]),
						toUnstructured(t, events[1]),
						toUnstructured(t, events[2]),
					}, nil).AnyTimes()

			},
			expected: []string{"event-0", "event-1", "event-2"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			o := storefake.NewMockObjectStore(controller)
			discovery := fakequeryer.NewMockDiscoveryInterface(controller)

			if tc.setup != nil {
				tc.setup(t, o)
			}

			oq := New(o, discovery)

			ctx := context.Background()
			events, err := oq.Events(ctx, tc.object)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			var got []string
			for _, event := range events {
				got = append(got, event.GetName())
			}

			sort.Strings(tc.expected)
			sort.Strings(got)

			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestCacheQueryer_IngressesForService(t *testing.T) {
	service := &corev1.Service{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Service"},
		ObjectMeta: metav1.ObjectMeta{Name: "service", Namespace: "default"},
	}

	ingress1 := &extv1beta1.Ingress{
		TypeMeta:   metav1.TypeMeta{APIVersion: "extensions/v1beta1", Kind: "Ingress"},
		ObjectMeta: metav1.ObjectMeta{Name: "ingress1", Namespace: "default"},
		Spec: extv1beta1.IngressSpec{
			Backend: &extv1beta1.IngressBackend{
				ServiceName: "service",
			},
		},
	}

	ingress2 := &extv1beta1.Ingress{
		TypeMeta:   metav1.TypeMeta{APIVersion: "extensions/v1beta1", Kind: "Ingress"},
		ObjectMeta: metav1.ObjectMeta{Name: "ingress2", Namespace: "default"},
		Spec: extv1beta1.IngressSpec{
			Rules: []extv1beta1.IngressRule{
				{
					IngressRuleValue: extv1beta1.IngressRuleValue{
						HTTP: &extv1beta1.HTTPIngressRuleValue{
							Paths: []extv1beta1.HTTPIngressPath{
								{
									Backend: extv1beta1.IngressBackend{
										ServiceName: "service",
									},
								},
								{
									Backend: extv1beta1.IngressBackend{
										ServiceName: "",
									},
								},
							},
						},
					},
				},
				{
					IngressRuleValue: extv1beta1.IngressRuleValue{},
				},
			},
		},
	}

	ingress3 := &extv1beta1.Ingress{
		TypeMeta:   metav1.TypeMeta{APIVersion: "extensions/v1beta1", Kind: "Ingress"},
		ObjectMeta: metav1.ObjectMeta{Name: "ingress2", Namespace: "default"},
	}

	cases := []struct {
		name     string
		service  *corev1.Service
		setup    func(t *testing.T, o *storefake.MockObjectStore)
		expected []*extv1beta1.Ingress
		isErr    bool
	}{
		{
			name:    "in general",
			service: service,
			setup: func(t *testing.T, o *storefake.MockObjectStore) {
				ingressesKey := cacheutil.Key{
					Namespace:  "default",
					APIVersion: "extensions/v1beta1",
					Kind:       "Ingress",
				}
				o.EXPECT().
					List(gomock.Any(), gomock.Eq(ingressesKey)).
					Return([]*unstructured.Unstructured{
						toUnstructured(t, ingress1),
						toUnstructured(t, ingress2),
						toUnstructured(t, ingress3),
					}, nil)
			},
			expected: []*extv1beta1.Ingress{
				ingress1, ingress2,
			},
		},
		{
			name:    "service is nil",
			service: nil,
			isErr:   true,
		},
		{
			name:    "ingress list failure",
			service: service,
			setup: func(t *testing.T, o *storefake.MockObjectStore) {
				ingressesKey := cacheutil.Key{
					Namespace:  "default",
					APIVersion: "extensions/v1beta1",
					Kind:       "Ingress",
				}
				o.EXPECT().
					List(gomock.Any(), gomock.Eq(ingressesKey)).
					Return(nil, errors.New("failed"))
			},
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			o := storefake.NewMockObjectStore(controller)
			discovery := fakequeryer.NewMockDiscoveryInterface(controller)

			if tc.setup != nil {
				tc.setup(t, o)
			}

			oq := New(o, discovery)

			ctx := context.Background()
			got, err := oq.IngressesForService(ctx, tc.service)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestCacheQueryer_OwnerReference(t *testing.T) {
	ownerReference := metav1.OwnerReference{
		APIVersion: "apps/v1",
		Kind:       "Deployment",
		Name:       "deployment",
	}

	deployment := &appsv1.Deployment{
		TypeMeta:   metav1.TypeMeta{APIVersion: "apps/v1", Kind: "Deployment"},
		ObjectMeta: metav1.ObjectMeta{Name: "deployment", Namespace: "default"},
	}

	cases := []struct {
		name     string
		setup    func(t *testing.T, o *storefake.MockObjectStore)
		expected func(t *testing.T) runtime.Object
		isErr    bool
	}{
		{
			name: "in general",
			setup: func(t *testing.T, o *storefake.MockObjectStore) {
				key := cacheutil.Key{
					Namespace:  "default",
					APIVersion: "apps/v1",
					Kind:       "Deployment",
					Name:       "deployment",
				}
				o.EXPECT().
					Get(gomock.Any(), gomock.Eq(key)).
					Return(toUnstructured(t, deployment), nil)
			},
			expected: func(t *testing.T) runtime.Object {
				return toUnstructured(t, deployment)
			},
		},
		{
			name: "objectstore get failure",
			setup: func(t *testing.T, o *storefake.MockObjectStore) {
				key := cacheutil.Key{
					Namespace:  "default",
					APIVersion: "apps/v1",
					Kind:       "Deployment",
					Name:       "deployment",
				}
				o.EXPECT().
					Get(gomock.Any(), gomock.Eq(key)).
					Return(nil, errors.New("failed"))
			},
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			o := storefake.NewMockObjectStore(controller)
			discovery := fakequeryer.NewMockDiscoveryInterface(controller)

			if tc.setup != nil {
				tc.setup(t, o)
			}

			oq := New(o, discovery)

			ctx := context.Background()
			got, err := oq.OwnerReference(ctx, "default", ownerReference)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expected(t), got)
		})
	}
}

func TestCacheQueryer_PodsForService(t *testing.T) {
	service := &corev1.Service{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Service"},
		ObjectMeta: metav1.ObjectMeta{Name: "service", Namespace: "default"},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "one",
			},
		},
	}

	pod1 := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod1",
			Namespace: "default",
			Labels: map[string]string{
				"app": "one",
			},
		},
	}

	pod2 := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod2",
			Namespace: "default",
			Labels: map[string]string{
				"app": "two",
			},
		},
	}

	cases := []struct {
		name     string
		service  *corev1.Service
		setup    func(t *testing.T, o *storefake.MockObjectStore)
		expected []*corev1.Pod
		isErr    bool
	}{
		{
			name:    "in general",
			service: service,
			setup: func(t *testing.T, o *storefake.MockObjectStore) {
				key := cacheutil.Key{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "Pod",
				}
				o.EXPECT().
					List(gomock.Any(), gomock.Eq(key)).
					Return([]*unstructured.Unstructured{
						toUnstructured(t, pod1),
						toUnstructured(t, pod2),
					}, nil)
			},
			expected: []*corev1.Pod{pod1},
		},
		{
			name:    "service is nil",
			service: nil,
			isErr:   true,
		},
		{
			name:    "objectstore list failure",
			service: service,
			setup: func(t *testing.T, o *storefake.MockObjectStore) {
				key := cacheutil.Key{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "Pod",
				}
				o.EXPECT().
					List(gomock.Any(), gomock.Eq(key)).
					Return(nil, errors.New("failed"))
			},
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			o := storefake.NewMockObjectStore(controller)
			discovery := fakequeryer.NewMockDiscoveryInterface(controller)

			if tc.setup != nil {
				tc.setup(t, o)
			}

			oq := New(o, discovery)

			ctx := context.Background()
			got, err := oq.PodsForService(ctx, tc.service)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestCacheQueryer_ServicesForIngress(t *testing.T) {
	ingress1 := &extv1beta1.Ingress{
		TypeMeta:   metav1.TypeMeta{APIVersion: "extensions/v1beta1", Kind: "Ingress"},
		ObjectMeta: metav1.ObjectMeta{Name: "ingress1", Namespace: "default"},
		Spec: extv1beta1.IngressSpec{
			Backend: &extv1beta1.IngressBackend{
				ServiceName: "service1",
			},
		},
	}

	ingress2 := &extv1beta1.Ingress{
		TypeMeta:   metav1.TypeMeta{APIVersion: "extensions/v1beta1", Kind: "Ingress"},
		ObjectMeta: metav1.ObjectMeta{Name: "ingress2", Namespace: "default"},
		Spec: extv1beta1.IngressSpec{
			Rules: []extv1beta1.IngressRule{
				{
					IngressRuleValue: extv1beta1.IngressRuleValue{
						HTTP: &extv1beta1.HTTPIngressRuleValue{
							Paths: []extv1beta1.HTTPIngressPath{
								{
									Backend: extv1beta1.IngressBackend{
										ServiceName: "service2",
									},
								},
								{
									Backend: extv1beta1.IngressBackend{
										ServiceName: "service1",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	service1 := &corev1.Service{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Service"},
		ObjectMeta: metav1.ObjectMeta{Name: "service1", Namespace: "default"},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "one",
			},
		},
	}

	service2 := &corev1.Service{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Service"},
		ObjectMeta: metav1.ObjectMeta{Name: "service2", Namespace: "default"},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "two",
			},
		},
	}

	cases := []struct {
		name     string
		ingress  *extv1beta1.Ingress
		setup    func(t *testing.T, o *storefake.MockObjectStore)
		expected []string
		isErr    bool
	}{
		{
			name:    "in general: service defined as backend",
			ingress: ingress1,
			setup: func(t *testing.T, o *storefake.MockObjectStore) {
				key := cacheutil.Key{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "Service",
					Name:       "service1",
				}
				o.EXPECT().
					Get(gomock.Any(), gomock.Eq(key)).
					Return(toUnstructured(t, service1), nil)
			},
			expected: []string{"service1"},
		},
		{
			name:    "in general: services defined in rules",
			ingress: ingress2,
			setup: func(t *testing.T, o *storefake.MockObjectStore) {
				key1 := cacheutil.Key{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "Service",
					Name:       "service1",
				}
				o.EXPECT().
					Get(gomock.Any(), gomock.Eq(key1)).
					Return(toUnstructured(t, service1), nil)
				key2 := cacheutil.Key{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "Service",
					Name:       "service2",
				}
				o.EXPECT().
					Get(gomock.Any(), gomock.Eq(key2)).
					Return(toUnstructured(t, service2), nil)
			},
			expected: []string{"service1", "service2"},
		},
		{
			name:    "ingress is nil",
			ingress: nil,
			isErr:   true,
		},
		{
			name:    "objectstore list failure",
			ingress: ingress1,
			setup: func(t *testing.T, c *storefake.MockObjectStore) {
				key := cacheutil.Key{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "Service",
					Name:       "service1",
				}
				c.EXPECT().
					Get(gomock.Any(), gomock.Eq(key)).
					Return(nil, errors.New("failed"))
			},
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			o := storefake.NewMockObjectStore(controller)
			discovery := fakequeryer.NewMockDiscoveryInterface(controller)

			if tc.setup != nil {
				tc.setup(t, o)
			}

			oq := New(o, discovery)

			ctx := context.Background()
			services, err := oq.ServicesForIngress(ctx, tc.ingress)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			var got []string
			for _, service := range services {
				got = append(got, service.Name)
			}
			sort.Strings(got)
			sort.Strings(tc.expected)

			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestCacheQueryer_ServicesForPods(t *testing.T) {
	service1 := &corev1.Service{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Service"},
		ObjectMeta: metav1.ObjectMeta{Name: "service1", Namespace: "default"},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "one",
			},
		},
	}

	service2 := &corev1.Service{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Service"},
		ObjectMeta: metav1.ObjectMeta{Name: "service2", Namespace: "default"},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "two",
			},
		},
	}

	pod1 := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod1",
			Namespace: "default",
			Labels: map[string]string{
				"app": "one",
			},
		},
	}

	cases := []struct {
		name     string
		pod      *corev1.Pod
		setup    func(t *testing.T, o *storefake.MockObjectStore)
		expected []string
		isErr    bool
	}{
		{
			name: "in general",
			pod:  pod1,
			setup: func(t *testing.T, o *storefake.MockObjectStore) {
				key := cacheutil.Key{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "Service",
				}
				o.EXPECT().
					List(gomock.Any(), gomock.Eq(key)).
					Return([]*unstructured.Unstructured{
						toUnstructured(t, service1),
						toUnstructured(t, service2),
					}, nil)
			},
			expected: []string{"service1"},
		},
		{
			name:  "service is nil",
			pod:   nil,
			isErr: true,
		},
		{
			name: "objectstore list failure",
			pod:  pod1,
			setup: func(t *testing.T, o *storefake.MockObjectStore) {
				key := cacheutil.Key{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "Service",
				}
				o.EXPECT().
					List(gomock.Any(), gomock.Eq(key)).
					Return(nil, errors.New("failed"))
			},
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			o := storefake.NewMockObjectStore(controller)
			discovery := fakequeryer.NewMockDiscoveryInterface(controller)

			if tc.setup != nil {
				tc.setup(t, o)
			}

			oq := New(o, discovery)

			ctx := context.Background()
			services, err := oq.ServicesForPod(ctx, tc.pod)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			var got []string
			for _, service := range services {
				got = append(got, service.Name)
			}
			sort.Strings(got)
			sort.Strings(tc.expected)

			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestCacheQueryer_getSelector(t *testing.T) {
	selector := &metav1.LabelSelector{
		MatchLabels: map[string]string{"foo": "bar"},
	}

	cases := []struct {
		name     string
		object   runtime.Object
		expected *metav1.LabelSelector
		isErr    bool
	}{
		{
			name:     "cron job",
			object:   &batchv1beta1.CronJob{},
			expected: nil,
		},
		{
			name: "daemon set",
			object: &appsv1.DaemonSet{
				Spec: appsv1.DaemonSetSpec{
					Selector: selector,
				},
			},
			expected: selector,
		},
		{
			name: "deployment",
			object: &appsv1.Deployment{
				Spec: appsv1.DeploymentSpec{
					Selector: selector,
				},
			},
			expected: selector,
		},
		{
			name: "replication controller",
			object: &corev1.ReplicationController{
				Spec: corev1.ReplicationControllerSpec{
					Selector: selector.MatchLabels,
				},
			},
			expected: selector,
		},
		{
			name: "replica set",
			object: &appsv1.ReplicaSet{
				Spec: appsv1.ReplicaSetSpec{
					Selector: selector,
				},
			},
			expected: selector,
		},
		{
			name: "service",
			object: &corev1.Service{
				Spec: corev1.ServiceSpec{
					Selector: selector.MatchLabels,
				},
			},
			expected: selector,
		},
		{
			name: "stateful set",
			object: &appsv1.StatefulSet{
				Spec: appsv1.StatefulSetSpec{
					Selector: selector,
				},
			},
			expected: selector,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			o := storefake.NewMockObjectStore(controller)
			discovery := fakequeryer.NewMockDiscoveryInterface(controller)

			oq := New(o, discovery)

			got, err := oq.getSelector(tc.object)
			if tc.isErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, tc.expected, got)
		})
	}
}

func toUnstructured(t *testing.T, object metav1.Object) *unstructured.Unstructured {
	m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(object)
	require.NoError(t, err)
	return &unstructured.Unstructured{Object: m}
}

func genEventFor(t *testing.T, object metav1.Object, name string) *corev1.Event {
	u := toUnstructured(t, object)

	return &corev1.Event{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Event"},
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "default"},
		InvolvedObject: corev1.ObjectReference{
			Namespace:  u.GetNamespace(),
			APIVersion: u.GetAPIVersion(),
			Kind:       u.GetKind(),
			Name:       u.GetName(),
		},
	}
}
