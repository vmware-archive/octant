package objectvisitor

import (
	"fmt"
	"sort"
	"testing"

	"github.com/golang/mock/gomock"
	queryerfake "github.com/heptio/developer-dash/internal/queryer/fake"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

func Test_DefaultVisitor_Visit(t *testing.T) {
	cases := []struct {
		name            string
		init            func(t *testing.T, q *queryerfake.MockQueryer) []ClusterObject
		expectedObjects []string
		expectedEdges   map[string][]string
	}{
		{
			name: "workload with pod",
			init: func(t *testing.T, q *queryerfake.MockQueryer) []ClusterObject {
				daemonSet := createDaemonSet("daemonset")
				pod := createPod("pod")
				pod.SetOwnerReferences(toOwnerReferences(t, daemonSet))

				q.EXPECT().
					Children(gomock.Eq(daemonSet)).
					Return([]runtime.Object{pod}, nil).AnyTimes()

				q.EXPECT().
					ServicesForPod(gomock.Eq(pod)).
					Return([]*corev1.Service{}, nil).AnyTimes()

				q.EXPECT().
					OwnerReference(gomock.Eq("namespace"), gomock.Eq(pod.OwnerReferences[0])).
					Return(daemonSet, nil).AnyTimes()

				q.EXPECT().
					Children(gomock.Eq(pod)).
					Return([]runtime.Object{}, nil).AnyTimes()

				return []ClusterObject{daemonSet, pod}
			},
			expectedObjects: []string{
				"apps/v1, Kind=DaemonSet:daemonset",
				"/v1, Kind=Pod:pod",
			},
			expectedEdges: map[string][]string{
				"daemonset": []string{"pod"},
			},
		},
		{
			name: "service with pod",
			init: func(t *testing.T, q *queryerfake.MockQueryer) []ClusterObject {
				service := createService("service")
				pod := createPod("pod")

				q.EXPECT().
					PodsForService(gomock.Eq(service)).
					Return([]*corev1.Pod{pod}, nil).AnyTimes()

				q.EXPECT().
					ServicesForPod(gomock.Eq(pod)).
					Return([]*corev1.Service{service}, nil).AnyTimes()

				q.EXPECT().
					Children(gomock.Eq(pod)).
					Return([]runtime.Object{}, nil).AnyTimes()

				q.EXPECT().
					Children(gomock.Eq(service)).
					Return([]runtime.Object{}, nil).AnyTimes()

				q.EXPECT().
					IngressesForService(gomock.Eq(service)).
					Return([]*extv1beta1.Ingress{}, nil).AnyTimes()

				return []ClusterObject{service}
			},
			expectedObjects: []string{
				"/v1, Kind=Service:service",
				"/v1, Kind=Pod:pod",
			},
			expectedEdges: map[string][]string{
				"service": []string{"pod"},
			},
		},
		{
			name: "ingress with service and pod",
			init: func(t *testing.T, q *queryerfake.MockQueryer) []ClusterObject {
				ingress := createIngress("ingress")
				service := createService("service")
				pod := createPod("pod")

				q.EXPECT().
					ServicesForIngress(gomock.Eq(ingress)).
					Return([]*corev1.Service{service}, nil).AnyTimes()

				q.EXPECT().
					IngressesForService(gomock.Eq(service)).
					Return([]*extv1beta1.Ingress{ingress}, nil).AnyTimes()

				q.EXPECT().
					PodsForService(gomock.Eq(service)).
					Return([]*corev1.Pod{pod}, nil).AnyTimes()

				q.EXPECT().
					ServicesForPod(gomock.Eq(pod)).
					Return([]*corev1.Service{service}, nil).AnyTimes()

				q.EXPECT().
					Children(gomock.Eq(ingress)).
					Return([]runtime.Object{}, nil).AnyTimes()

				q.EXPECT().
					Children(gomock.Eq(service)).
					Return([]runtime.Object{}, nil).AnyTimes()

				q.EXPECT().
					Children(gomock.Eq(pod)).
					Return([]runtime.Object{}, nil).AnyTimes()

				return []ClusterObject{ingress, service, pod}
			},
			expectedObjects: []string{
				"extensions/v1beta1, Kind=Ingress:ingress",
				"/v1, Kind=Pod:pod",
				"/v1, Kind=Service:service",
			},
			expectedEdges: map[string][]string{
				"ingress": []string{"service"},
				"service": []string{"pod"},
			},
		},
		{
			name: "full workload",
			init: func(t *testing.T, q *queryerfake.MockQueryer) []ClusterObject {
				ingress := createIngress("ingress")
				service := createService("service")
				pod := createPod("pod")
				deployment := createDeployment("deployment")
				replicaSet := createReplicaSet("replicaSet")

				replicaSet.SetOwnerReferences(toOwnerReferences(t, deployment))
				pod.SetOwnerReferences(toOwnerReferences(t, replicaSet))

				q.EXPECT().
					ServicesForIngress(gomock.Eq(ingress)).
					Return([]*corev1.Service{service}, nil).AnyTimes()

				q.EXPECT().
					IngressesForService(gomock.Eq(service)).
					Return([]*extv1beta1.Ingress{ingress}, nil).AnyTimes()

				q.EXPECT().
					PodsForService(gomock.Eq(service)).
					Return([]*corev1.Pod{pod}, nil).AnyTimes()

				q.EXPECT().
					ServicesForPod(gomock.Eq(pod)).
					Return([]*corev1.Service{service}, nil).AnyTimes()

				q.EXPECT().
					Children(gomock.Eq(ingress)).
					Return([]runtime.Object{}, nil).AnyTimes()

				q.EXPECT().
					Children(gomock.Eq(service)).
					Return([]runtime.Object{}, nil).AnyTimes()

				q.EXPECT().
					Children(gomock.Eq(pod)).
					Return([]runtime.Object{}, nil).AnyTimes()

				q.EXPECT().
					Children(gomock.Eq(replicaSet)).
					Return([]runtime.Object{pod}, nil).AnyTimes()

				q.EXPECT().
					Children(gomock.Eq(deployment)).
					Return([]runtime.Object{replicaSet}, nil).AnyTimes()

				q.EXPECT().
					OwnerReference(gomock.Eq("namespace"), gomock.Eq(pod.OwnerReferences[0])).
					Return(replicaSet, nil).AnyTimes()

				q.EXPECT().
					OwnerReference(gomock.Eq("namespace"), gomock.Eq(replicaSet.OwnerReferences[0])).
					Return(deployment, nil).AnyTimes()

				return []ClusterObject{ingress, service, pod, replicaSet, deployment}
			},
			expectedObjects: []string{
				"extensions/v1beta1, Kind=Ingress:ingress",
				"/v1, Kind=Pod:pod",
				"/v1, Kind=Service:service",
				"apps/v1, Kind=ReplicaSet:replicaSet",
				"apps/v1, Kind=Deployment:deployment",
			},
			expectedEdges: map[string][]string{
				"service":    []string{"pod"},
				"replicaSet": []string{"pod"},
				"ingress":    []string{"service"},
				"deployment": []string{"replicaSet"},
			},
		},
		{
			name: "multiple workloads/services, single ingress",
			init: func(t *testing.T, q *queryerfake.MockQueryer) []ClusterObject {
				d1 := createDeployment("d1")
				d1rs1 := createReplicaSet("d1rs1")
				d1rs1.SetOwnerReferences(toOwnerReferences(t, d1))
				d1rs1p1 := createPod("d1rs1p1")
				d1rs1p1.SetOwnerReferences(toOwnerReferences(t, d1rs1))
				d1rs1p2 := createPod("d1rs1p2")
				d1rs1p2.SetOwnerReferences(toOwnerReferences(t, d1rs1))
				s1 := createService("s1")

				d2 := createDeployment("d2")
				d2rs1 := createReplicaSet("d2rs1")
				d2rs1.SetOwnerReferences(toOwnerReferences(t, d2))
				d2rs1p1 := createPod("d2rs1p1")
				d2rs1p1.SetOwnerReferences(toOwnerReferences(t, d2rs1))
				s2 := createService("s2")

				ingress := createIngress("i1")

				q.EXPECT().
					Children(gomock.Eq(d1)).
					Return([]runtime.Object{d1rs1}, nil).AnyTimes()

				q.EXPECT().
					Children(gomock.Eq(d2)).
					Return([]runtime.Object{d2rs1}, nil).AnyTimes()

				q.EXPECT().
					Children(gomock.Eq(d1rs1)).
					Return([]runtime.Object{d1rs1p1, d1rs1p2}, nil).AnyTimes()

				q.EXPECT().
					Children(gomock.Eq(d2rs1)).
					Return([]runtime.Object{d2rs1p1}, nil).AnyTimes()

				q.EXPECT().
					Children(gomock.Eq(d1rs1p1)).
					Return([]runtime.Object{}, nil).AnyTimes()

				q.EXPECT().
					Children(gomock.Eq(d1rs1p2)).
					Return([]runtime.Object{}, nil).AnyTimes()

				q.EXPECT().
					Children(gomock.Eq(d2rs1p1)).
					Return([]runtime.Object{}, nil).AnyTimes()

				q.EXPECT().
					Children(gomock.Eq(s1)).
					Return([]runtime.Object{}, nil).AnyTimes()

				q.EXPECT().
					Children(gomock.Eq(s2)).
					Return([]runtime.Object{}, nil).AnyTimes()

				q.EXPECT().
					Children(gomock.Eq(ingress)).
					Return([]runtime.Object{}, nil).AnyTimes()

				q.EXPECT().
					OwnerReference(gomock.Eq("namespace"), gomock.Eq(d1rs1.OwnerReferences[0])).
					Return(d1, nil).AnyTimes()

				q.EXPECT().
					OwnerReference(gomock.Eq("namespace"), gomock.Eq(d2rs1.OwnerReferences[0])).
					Return(d2, nil).AnyTimes()

				q.EXPECT().
					OwnerReference(gomock.Eq("namespace"), gomock.Eq(d1rs1p1.OwnerReferences[0])).
					Return(d1rs1, nil).AnyTimes()

				q.EXPECT().
					OwnerReference(gomock.Eq("namespace"), gomock.Eq(d1rs1p2.OwnerReferences[0])).
					Return(d1rs1, nil).AnyTimes()

				q.EXPECT().
					OwnerReference(gomock.Eq("namespace"), gomock.Eq(d2rs1p1.OwnerReferences[0])).
					Return(d2rs1, nil).AnyTimes()

				q.EXPECT().
					ServicesForPod(gomock.Eq(d1rs1p1)).
					Return([]*corev1.Service{s1}, nil).AnyTimes()

				q.EXPECT().
					ServicesForPod(gomock.Eq(d1rs1p2)).
					Return([]*corev1.Service{s1}, nil).AnyTimes()

				q.EXPECT().
					ServicesForPod(gomock.Eq(d2rs1p1)).
					Return([]*corev1.Service{s2}, nil).AnyTimes()

				q.EXPECT().
					PodsForService(gomock.Eq(s1)).
					Return([]*corev1.Pod{d1rs1p1, d1rs1p2}, nil).AnyTimes()

				q.EXPECT().
					PodsForService(gomock.Eq(s2)).
					Return([]*corev1.Pod{d2rs1p1}, nil).AnyTimes()

				q.EXPECT().
					IngressesForService(gomock.Eq(s1)).
					Return([]*extv1beta1.Ingress{ingress}, nil).AnyTimes()

				q.EXPECT().
					IngressesForService(gomock.Eq(s2)).
					Return([]*extv1beta1.Ingress{ingress}, nil).AnyTimes()

				q.EXPECT().
					ServicesForIngress(gomock.Eq(ingress)).
					Return([]*corev1.Service{s1, s2}, nil).AnyTimes()

				return []ClusterObject{d1, d1rs1, d1rs1p1, d1rs1p2, d2, d2rs1,
					d2rs1p1, s1, s2, ingress}
			},
			expectedObjects: []string{
				"apps/v1, Kind=Deployment:d1",
				"apps/v1, Kind=ReplicaSet:d1rs1",
				"/v1, Kind=Pod:d1rs1p1",
				"/v1, Kind=Pod:d1rs1p2",
				"apps/v1, Kind=Deployment:d2",
				"apps/v1, Kind=ReplicaSet:d2rs1",
				"/v1, Kind=Pod:d2rs1p1",
				"/v1, Kind=Service:s1",
				"/v1, Kind=Service:s2",
				"extensions/v1beta1, Kind=Ingress:i1",
			},
			expectedEdges: map[string][]string{
				"d1":    []string{"d1rs1"},
				"d1rs1": []string{"d1rs1p1", "d1rs1p2"},
				"s1":    []string{"d1rs1p1", "d1rs1p2"},
				"d2":    []string{"d2rs1"},
				"d2rs1": []string{"d2rs1p1"},
				"s2":    []string{"d2rs1p1"},
				"i1":    []string{"s1", "s2"},
			},
		},
	}

	gvks := []schema.GroupVersionKind{DaemonSetGVK, DeploymentGVK, IngressGVK, PodGVK,
		ServiceGVK, ReplicaSetGSK, ReplicationControllerGSK, StatefulSetGSK}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			q := queryerfake.NewMockQueryer(ctrl)

			require.NotNil(t, tc.init, "init func is required")
			objects := tc.init(t, q)

			for _, object := range objects {
				t.Run(fmt.Sprintf("seeded with %T", object), func(t *testing.T) {
					factoryGen := NewDefaultFactoryGenerator()

					ic := identityCollector{t: t}

					for _, gvk := range gvks {
						factoryRegister(t, factoryGen, gvk, ic.factoryFn)
					}

					dv, err := NewDefaultVisitor(q, factoryGen.FactoryFunc())
					require.NoError(t, err)

					err = dv.Visit(object)
					require.NoError(t, err)

					ic.assertMatch(tc.expectedObjects)
					ic.assertChildren(tc.expectedEdges)
				})
			}
		})
	}
}

func createDaemonSet(name string) *appsv1.DaemonSet {
	return &appsv1.DaemonSet{
		TypeMeta:   genTypeMeta(DaemonSetGVK),
		ObjectMeta: genObjectMeta(name),
	}
}

func createDeployment(name string) *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta:   genTypeMeta(DeploymentGVK),
		ObjectMeta: genObjectMeta(name),
	}
}

func createIngress(name string) *extv1beta1.Ingress {
	return &extv1beta1.Ingress{
		TypeMeta:   genTypeMeta(IngressGVK),
		ObjectMeta: genObjectMeta(name),
	}
}

func createPod(name string) *corev1.Pod {
	return &corev1.Pod{
		TypeMeta:   genTypeMeta(PodGVK),
		ObjectMeta: genObjectMeta(name),
	}
}

func toOwnerReferences(t *testing.T, object ClusterObject) []metav1.OwnerReference {
	apiVersion, kind := object.GroupVersionKind().ToAPIVersionAndKind()

	return []metav1.OwnerReference{
		{
			APIVersion: apiVersion,
			Kind:       kind,
			Name:       object.GetName(),
			UID:        object.GetUID(),
		},
	}
}

func createReplicaSet(name string) *appsv1.ReplicaSet {
	return &appsv1.ReplicaSet{
		TypeMeta:   genTypeMeta(ReplicaSetGSK),
		ObjectMeta: genObjectMeta(name),
	}
}

func createService(name string) *corev1.Service {
	return &corev1.Service{
		TypeMeta:   genTypeMeta(ServiceGVK),
		ObjectMeta: genObjectMeta(name),
	}
}

func genTypeMeta(gvk schema.GroupVersionKind) metav1.TypeMeta {
	apiVersion, kind := gvk.ToAPIVersionAndKind()
	return metav1.TypeMeta{
		APIVersion: apiVersion,
		Kind:       kind,
	}
}

func genObjectMeta(name string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: "namespace",
		UID:       types.UID(name),
	}
}

func factoryRegister(
	t *testing.T,
	gen *DefaultFactoryGenerator,
	gvk schema.GroupVersionKind,
	factory ObjectHandlerFactory) {
	err := gen.Register(gvk, factory)
	require.NoError(t, err)
}

type testObject struct {
	processFn  func(object ClusterObject)
	addChildFn func(parent ClusterObject, children ...ClusterObject) error
}

func (o *testObject) Process(object ClusterObject) {
	o.processFn(object)
}

func (o *testObject) AddChild(parent ClusterObject, children ...ClusterObject) error {
	return o.addChildFn(parent, children...)
}

type identityCollector struct {
	t           *testing.T
	gotVisits   []string
	gotChildren map[string][]string

	o *testObject
}

func (ic *identityCollector) factoryFn(object ClusterObject) (ObjectHandler, error) {
	if ic.o == nil {
		ic.gotChildren = make(map[string][]string)

		objectKind := object.GetObjectKind()
		if objectKind == nil {
			return nil, errors.Errorf("object kind is nil")
		}

		ic.o = &testObject{
			processFn: func(object ClusterObject) {
				ic.gotVisits = append(ic.gotVisits, fmt.Sprintf("%s:%s", object.GroupVersionKind(), object.GetName()))
			},
			addChildFn: func(parent ClusterObject, children ...ClusterObject) error {
				pUID := string(parent.GetUID())

				for _, child := range children {
					cUID := string(child.GetUID())
					ic.gotChildren[pUID] = append(ic.gotChildren[pUID], cUID)
				}
				return nil
			},
		}
	}

	return ic.o, nil
}

func (ic *identityCollector) assertMatch(expected []string) {
	got := ic.gotVisits

	sort.Strings(expected)
	sort.Strings(got)

	assert.Equal(ic.t, expected, got)
}

func (ic *identityCollector) assertChildren(expected map[string][]string) {
	got := ic.gotChildren
	for k := range expected {
		sort.Strings(expected[k])
	}
	for k := range got {
		sort.Strings(got[k])
	}

	assert.Equal(ic.t, expected, got)
}
