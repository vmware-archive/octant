/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package queryer

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	queryerFake "github.com/vmware-tanzu/octant/internal/queryer/fake"
	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/store"
	storeFake "github.com/vmware-tanzu/octant/pkg/store/fake"
)

func TestCacheQueryer_Children(t *testing.T) {
	deployment := testutil.ToUnstructured(t, testutil.CreateDeployment("deployment"))

	rs := testutil.ToUnstructured(t, testutil.CreateExtReplicaSet("rs"))
	rs.SetOwnerReferences(testutil.ToOwnerReferences(t, deployment))

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
					Kind:       "NotListable",
					Verbs:      metav1.Verbs{"get"},
				},
			},
		},
		{
			GroupVersion: "extensions/v1beta1",
			APIResources: []metav1.APIResource{
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

	rsKey, err := store.KeyFromObject(rs)
	require.NoError(t, err)
	rsKey.Name = ""

	deploymentKey, err := store.KeyFromObject(deployment)
	require.NoError(t, err)
	deploymentKey.Name = ""

	cases := []struct {
		name     string
		owner    *unstructured.Unstructured
		setup    func(t *testing.T, c *storeFake.MockStore, disco *queryerFake.MockDiscoveryInterface)
		expected func(t *testing.T) *unstructured.UnstructuredList
		isErr    bool
	}{
		{
			name:  "in general",
			owner: deployment,
			setup: func(t *testing.T, o *storeFake.MockStore, disco *queryerFake.MockDiscoveryInterface) {
				o.EXPECT().
					List(gomock.Any(), gomock.Eq(deploymentKey)).
					Return(testutil.ToUnstructuredList(t, deployment), false, nil)

				o.EXPECT().
					List(gomock.Any(), gomock.Eq(rsKey)).
					Return(testutil.ToUnstructuredList(t, rs), false, nil)

				disco.EXPECT().
					ServerPreferredResources().
					Return(resourceLists, nil)

			},
			expected: func(t *testing.T) *unstructured.UnstructuredList {
				return testutil.ToUnstructuredList(t, rs)
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
			setup: func(t *testing.T, o *storeFake.MockStore, disco *queryerFake.MockDiscoveryInterface) {
				disco.EXPECT().
					ServerPreferredResources().
					Return(nil, errors.New("failed")).AnyTimes()
			},
			isErr: true,
		},
		{
			name:  "object store list fails",
			owner: deployment,
			setup: func(t *testing.T, o *storeFake.MockStore, disco *queryerFake.MockDiscoveryInterface) {
				o.EXPECT().
					List(gomock.Any(), gomock.Eq(deploymentKey)).
					Return(nil, false, errors.New("failed")).Times(1)

				o.EXPECT().
					List(gomock.Any(), gomock.Eq(rsKey)).
					Return(testutil.ToUnstructuredList(t, rs), false, nil)

				disco.EXPECT().
					ServerPreferredResources().
					Return(resourceLists, nil).AnyTimes()
			},
			isErr: true,
		},
	}

	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			o := storeFake.NewMockStore(controller)
			discovery := queryerFake.NewMockDiscoveryInterface(controller)

			crdKey := store.Key{
				APIVersion: "apiextensions.k8s.io/v1beta1",
				Kind:       "CustomResourceDefinition",
			}
			o.EXPECT().List(gomock.Any(), crdKey).Return(&unstructured.UnstructuredList{}, false, nil).AnyTimes()

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
		setup    func(t *testing.T, o *storeFake.MockStore)
		isErr    bool
		expected []string
	}{
		{
			name:   "in general",
			object: deployment,
			setup: func(t *testing.T, o *storeFake.MockStore) {
				key := store.Key{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "Event",
				}
				o.EXPECT().
					List(gomock.Any(), gomock.Eq(key)).
					Return(testutil.ToUnstructuredList(t, events[0], events[1], events[2]), false, nil)

			},
			expected: []string{"event-0", "event-1", "event-2"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			o := storeFake.NewMockStore(controller)
			discovery := queryerFake.NewMockDiscoveryInterface(controller)

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
		setup    func(t *testing.T, o *storeFake.MockStore)
		expected []*extv1beta1.Ingress
		isErr    bool
	}{
		{
			name:    "in general",
			service: service,
			setup: func(t *testing.T, o *storeFake.MockStore) {
				ingressesKey := store.Key{
					Namespace:  "default",
					APIVersion: "extensions/v1beta1",
					Kind:       "Ingress",
				}
				o.EXPECT().
					List(gomock.Any(), gomock.Eq(ingressesKey)).
					Return(testutil.ToUnstructuredList(t, ingress1, ingress2, ingress3), false, nil)
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
			setup: func(t *testing.T, o *storeFake.MockStore) {
				ingressesKey := store.Key{
					Namespace:  "default",
					APIVersion: "extensions/v1beta1",
					Kind:       "Ingress",
				}
				o.EXPECT().
					List(gomock.Any(), gomock.Eq(ingressesKey)).
					Return(nil, false, errors.New("failed"))
			},
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			o := storeFake.NewMockStore(controller)
			discovery := queryerFake.NewMockDiscoveryInterface(controller)

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
	deployment1 := testutil.ToUnstructured(t, testutil.CreateDeployment("deployment1"))
	deployment2 := testutil.ToUnstructured(t, testutil.CreateDeployment("deployment2"))
	replicaSet1 := testutil.ToUnstructured(t, testutil.CreateAppReplicaSet("replica-set1"))
	replicaSet2 := testutil.ToUnstructured(t, testutil.CreateAppReplicaSet("replica-set2"))
	replicaSet1.SetOwnerReferences(testutil.ToOwnerReferences(t, deployment1))
	replicaSet2.SetOwnerReferences(testutil.ToOwnerReferences(t, deployment1, deployment2))

	type args struct {
		object *unstructured.Unstructured
	}
	cases := []struct {
		name     string
		setup    func(t *testing.T, o *storeFake.MockStore)
		args     args
		expected func(t *testing.T) []*unstructured.Unstructured
		isErr    bool
	}{
		{
			name: "single owner",
			setup: func(t *testing.T, o *storeFake.MockStore) {
				key := store.Key{
					Namespace:  deployment1.GetNamespace(),
					APIVersion: "apps/v1",
					Kind:       "Deployment",
					Name:       "deployment",
				}
				key, err := store.KeyFromObject(deployment1)
				require.NoError(t, err)
				o.EXPECT().
					Get(gomock.Any(), gomock.Eq(key)).
					Return(deployment1, nil).AnyTimes()
			},
			args: args{
				object: replicaSet1,
			},
			expected: func(t *testing.T) []*unstructured.Unstructured {
				return []*unstructured.Unstructured{
					testutil.ToUnstructured(t, deployment1),
				}
			},
		},
		{
			name: "multiple owner",
			setup: func(t *testing.T, o *storeFake.MockStore) {
				for _, object := range []*unstructured.Unstructured{deployment1, deployment2} {
					key, err := store.KeyFromObject(object)
					require.NoError(t, err)
					o.EXPECT().
						Get(gomock.Any(), gomock.Eq(key)).
						Return(object, nil)
				}
			},
			args: args{
				object: replicaSet2,
			},
			expected: func(t *testing.T) []*unstructured.Unstructured {
				return []*unstructured.Unstructured{
					testutil.ToUnstructured(t, deployment1),
					testutil.ToUnstructured(t, deployment2),
				}
			},
		},
		{
			name: "object store get failure",
			setup: func(t *testing.T, o *storeFake.MockStore) {
				key, err := store.KeyFromObject(deployment1)
				require.NoError(t, err)
				o.EXPECT().
					Get(gomock.Any(), gomock.Eq(key)).
					Return(nil, errors.New("failed"))
			},
			args: args{
				object: replicaSet1,
			},
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			o := storeFake.NewMockStore(controller)
			discovery := queryerFake.NewMockDiscoveryInterface(controller)

			discovery.EXPECT().
				ServerResourcesForGroupVersion("apps/v1").
				Return(&metav1.APIResourceList{
					APIResources: []metav1.APIResource{{Kind: "Deployment", Namespaced: true}},
				}, nil).AnyTimes()

			if tc.setup != nil {
				tc.setup(t, o)
			}

			oq := New(o, discovery)

			ctx := context.Background()
			found, got, err := oq.OwnerReference(ctx, tc.args.object)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.True(t, found)
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
		setup    func(t *testing.T, o *storeFake.MockStore)
		expected []*corev1.Pod
		isErr    bool
	}{
		{
			name:    "in general",
			service: service,
			setup: func(t *testing.T, o *storeFake.MockStore) {
				key := store.Key{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "Pod",
				}
				o.EXPECT().
					List(gomock.Any(), gomock.Eq(key)).
					Return(testutil.ToUnstructuredList(t, pod1, pod2), false, nil)
			},
			expected: []*corev1.Pod{pod1},
		},
		{
			name:    "service is nil",
			service: nil,
			isErr:   true,
		},
		{
			name:    "object store list failure",
			service: service,
			setup: func(t *testing.T, o *storeFake.MockStore) {
				key := store.Key{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "Pod",
				}
				o.EXPECT().
					List(gomock.Any(), gomock.Eq(key)).
					Return(nil, false, errors.New("failed"))
			},
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			o := storeFake.NewMockStore(controller)
			discovery := queryerFake.NewMockDiscoveryInterface(controller)

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

func TestCacheQueryer_ServicesForIngress_service_not_found(t *testing.T) {
	ingress := testutil.CreateIngress("ingress")
	ingress.Spec.Backend = &extv1beta1.IngressBackend{
		ServiceName: "not-found",
	}

	controller := gomock.NewController(t)
	defer controller.Finish()

	o := storeFake.NewMockStore(controller)
	o.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(nil, nil)

	discovery := queryerFake.NewMockDiscoveryInterface(controller)

	oq := New(o, discovery)

	ctx := context.Background()
	services, err := oq.ServicesForIngress(ctx, ingress)
	require.NoError(t, err)
	require.Empty(t, services)
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
		setup    func(t *testing.T, o *storeFake.MockStore)
		expected []string
		isErr    bool
	}{
		{
			name:    "in general: service defined as backend",
			ingress: ingress1,
			setup: func(t *testing.T, o *storeFake.MockStore) {
				key := store.Key{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "Service",
					Name:       "service1",
				}
				o.EXPECT().
					Get(gomock.Any(), gomock.Eq(key)).
					Return(testutil.ToUnstructured(t, service1), nil)
			},
			expected: []string{"service1"},
		},
		{
			name:    "in general: services defined in rules",
			ingress: ingress2,
			setup: func(t *testing.T, o *storeFake.MockStore) {
				key1 := store.Key{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "Service",
					Name:       "service1",
				}
				o.EXPECT().
					Get(gomock.Any(), gomock.Eq(key1)).
					Return(testutil.ToUnstructured(t, service1), nil)
				key2 := store.Key{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "Service",
					Name:       "service2",
				}
				o.EXPECT().
					Get(gomock.Any(), gomock.Eq(key2)).
					Return(testutil.ToUnstructured(t, service2), nil)
			},
			expected: []string{"service1", "service2"},
		},
		{
			name:    "ingress is nil",
			ingress: nil,
			isErr:   true,
		},
		{
			name:    "object store list failure",
			ingress: ingress1,
			setup: func(t *testing.T, c *storeFake.MockStore) {
				key := store.Key{
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

			o := storeFake.NewMockStore(controller)
			discovery := queryerFake.NewMockDiscoveryInterface(controller)

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
			for _, service := range services.Items {
				accessor, err := meta.Accessor(&service)
				require.NoError(t, err)
				got = append(got, accessor.GetName())
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
		setup    func(t *testing.T, o *storeFake.MockStore)
		expected []string
		isErr    bool
	}{
		{
			name: "in general",
			pod:  pod1,
			setup: func(t *testing.T, o *storeFake.MockStore) {
				key := store.Key{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "Service",
				}
				o.EXPECT().
					List(gomock.Any(), gomock.Eq(key)).
					Return(testutil.ToUnstructuredList(t, service1, service2), false, nil)
			},
			expected: []string{"service1"},
		},
		{
			name:  "service is nil",
			pod:   nil,
			isErr: true,
		},
		{
			name: "object store list failure",
			pod:  pod1,
			setup: func(t *testing.T, o *storeFake.MockStore) {
				key := store.Key{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "Service",
				}
				o.EXPECT().
					List(gomock.Any(), gomock.Eq(key)).
					Return(nil, false, errors.New("failed"))
			},
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			o := storeFake.NewMockStore(controller)
			discovery := queryerFake.NewMockDiscoveryInterface(controller)

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

func TestObjectStoreQueryer_ServiceAccountForPod(t *testing.T) {
	serviceAccount := testutil.CreateServiceAccount("service-account")

	pod := testutil.CreatePod("pod")
	pod.Spec.ServiceAccountName = serviceAccount.Name

	controller := gomock.NewController(t)
	defer controller.Finish()

	o := storeFake.NewMockStore(controller)
	key, err := store.KeyFromObject(serviceAccount)
	require.NoError(t, err)
	o.EXPECT().
		Get(gomock.Any(), key).
		Return(testutil.ToUnstructured(t, serviceAccount), nil)

	discovery := queryerFake.NewMockDiscoveryInterface(controller)

	q := New(o, discovery)

	ctx := context.Background()
	got, err := q.ServiceAccountForPod(ctx, pod)
	require.NoError(t, err)

	require.Equal(t, serviceAccount, got)
}

func TestObjectStoreQueryer_ConfigMapsForPod(t *testing.T) {
	configMapKeyRef := testutil.CreateConfigMap("configmap1")
	configMapEnv := testutil.CreateConfigMap("configmap2")

	pod := testutil.CreatePod("pod")
	pod.Spec.Containers = []corev1.Container{
		{
			EnvFrom: []corev1.EnvFromSource{
				{
					ConfigMapRef: &corev1.ConfigMapEnvSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "configmap2",
						},
					},
				},
			},
			Env: []corev1.EnvVar{
				{
					Name:  "configmap3",
					Value: "configmap3_value",
				},
				{
					ValueFrom: &corev1.EnvVarSource{
						ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "configmap1",
							},
						},
					},
				},
			},
		},
	}

	controller := gomock.NewController(t)
	defer controller.Finish()

	o := storeFake.NewMockStore(controller)
	key := store.Key{
		Namespace:  "namespace",
		APIVersion: "v1",
		Kind:       "ConfigMap",
	}

	discovery := queryerFake.NewMockDiscoveryInterface(controller)

	q := New(o, discovery)

	ctx := context.Background()

	o.EXPECT().
		List(gomock.Any(), gomock.Eq(key)).
		Return(testutil.ToUnstructuredList(t, configMapKeyRef, configMapEnv), false, nil)
	configMaps, err := q.ConfigMapsForPod(ctx, pod)
	require.NoError(t, err)

	var got []string
	for _, configmap := range configMaps {
		got = append(got, configmap.Name)
	}
	sort.Strings(got)

	assert.Equal(t, []string([]string{configMapKeyRef.Name, configMapEnv.Name}), got)
}

func TestObjectStoreQueryer_SecretsForPod(t *testing.T) {
	secretInVolume := testutil.CreateSecret("secret1")
	secretEnv := testutil.CreateSecret("secret2")
	secretEnvFrom := testutil.CreateSecret("secret3")

	pod := testutil.CreatePod("pod")
	pod.Spec.Containers = []corev1.Container{
		{
			EnvFrom: []corev1.EnvFromSource{
				{
					SecretRef: &corev1.SecretEnvSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "secret3",
						},
					},
				},
				{
					ConfigMapRef: &corev1.ConfigMapEnvSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "Not a secret",
						},
					},
				},
			},
			Env: []corev1.EnvVar{
				{
					Name:  "TEST_SECRET_FOR_POD",
					Value: "test_secret_for_pod_value",
				},
				{
					ValueFrom: &corev1.EnvVarSource{
						ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
							Key: "Not a secret",
						},
					},
				},
				{
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "secret2",
							},
						},
					},
				},
			},
		},
	}
	pod.Spec.Volumes = []corev1.Volume{
		{
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: "secret1",
				},
			},
		},
	}

	controller := gomock.NewController(t)
	defer controller.Finish()

	o := storeFake.NewMockStore(controller)
	key := store.Key{
		Namespace:  "namespace",
		APIVersion: "v1",
		Kind:       "Secret",
	}

	discovery := queryerFake.NewMockDiscoveryInterface(controller)

	q := New(o, discovery)

	ctx := context.Background()

	o.EXPECT().
		List(gomock.Any(), gomock.Eq(key)).
		Return(testutil.ToUnstructuredList(t, secretInVolume, secretEnv, secretEnvFrom), false, nil)
	secrets, err := q.SecretsForPod(ctx, pod)
	require.NoError(t, err)

	var got []string
	for _, secret := range secrets {
		got = append(got, secret.Name)
	}
	sort.Strings(got)

	assert.Equal(t, []string([]string{secretInVolume.Name, secretEnv.Name, secretEnvFrom.Name}), got)
}

func TestObjectStoreQueryer_ScaleTarget(t *testing.T) {
	deployment := testutil.CreateDeployment("deployment")

	hpa := testutil.CreateHorizontalPodAutoscaler("hpa")
	hpa.Spec.ScaleTargetRef = autoscalingv1.CrossVersionObjectReference{
		APIVersion: deployment.APIVersion,
		Kind:       deployment.Kind,
		Name:       deployment.Name,
	}

	controller := gomock.NewController(t)
	defer controller.Finish()

	o := storeFake.NewMockStore(controller)
	key, err := store.KeyFromObject(deployment)
	require.NoError(t, err)
	o.EXPECT().
		Get(gomock.Any(), key).
		Return(testutil.ToUnstructured(t, deployment), nil)

	discovery := queryerFake.NewMockDiscoveryInterface(controller)

	q := New(o, discovery)

	ctx := context.Background()
	got, err := q.ScaleTarget(ctx, hpa)
	require.NoError(t, err)

	u := testutil.ToUnstructured(t, deployment)
	require.Equal(t, u.Object, got)
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

			o := storeFake.NewMockStore(controller)
			discovery := queryerFake.NewMockDiscoveryInterface(controller)

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

func genEventFor(t *testing.T, object runtime.Object, name string) *corev1.Event {
	u := testutil.ToUnstructured(t, object)

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
