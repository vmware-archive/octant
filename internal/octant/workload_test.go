/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package octant_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"

	"github.com/vmware-tanzu/octant/internal/objectstatus"
	"github.com/vmware-tanzu/octant/internal/octant"
	octantFake "github.com/vmware-tanzu/octant/internal/octant/fake"
	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/store"
	storeFake "github.com/vmware-tanzu/octant/pkg/store/fake"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func podMetricsLoader(controller *gomock.Controller, pm *unstructured.Unstructured, supportsMetrics bool, supportFails bool) *octantFake.MockPodMetricsLoader {
	pml := octantFake.NewMockPodMetricsLoader(controller)

	var err error
	if supportFails {
		err = fmt.Errorf("failed")
	}

	pml.EXPECT().SupportsMetrics(context.TODO()).Return(supportsMetrics, err).AnyTimes()

	if pm != nil {
		pml.EXPECT().
			Load("namespace", gomock.Any()).
			Return(pm, true, nil).
			AnyTimes()
	}

	return pml
}

func TestClusterWorkloadLoader(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	replicaSet := testutil.ToUnstructured(t, testutil.CreateAppReplicaSet("rs-1"))

	pod1 := testutil.ToUnstructured(t, testutil.CreatePod("pod-1"))
	pod2 := testutil.ToUnstructured(t, testutil.CreatePod("pod-with-owner", func(pod *corev1.Pod) {
		pod.OwnerReferences = testutil.ToOwnerReferences(t, replicaSet)
	}))
	pods := testutil.ToUnstructuredList(t, pod1, pod2)

	pm := testutil.ToUnstructured(t, testutil.CreatePodMetrics("pm", podMetricWithContainer))
	rl := corev1.ResourceList{
		corev1.ResourceCPU:    resource.MustParse("150m"),
		corev1.ResourceMemory: resource.MustParse("3Mi"),
	}

	pod3 := testutil.ToUnstructured(t, testutil.CreatePod("pod-3", func(pod *corev1.Pod) {
		pod.OwnerReferences = []metav1.OwnerReference{
			{
				APIVersion: "apps/v1",
				Kind:       "ReplicaSet",
				Name:       "missing",
				UID:        "12345",
			},
		}
	}))

	objectStore := storeFake.NewMockStore(controller)
	objectStore.EXPECT().
		List(gomock.Any(), store.Key{
			Namespace:  "namespace",
			APIVersion: "v1",
			Kind:       "Pod",
		}).
		Return(pods, true, nil).
		AnyTimes()

	objectStore.EXPECT().
		Get(gomock.Any(), store.Key{
			Namespace:  "namespace",
			APIVersion: "apps/v1",
			Kind:       "ReplicaSet",
			Name:       "rs-1",
		}).
		Return(replicaSet, nil).
		AnyTimes()

	noOwnerObjectStore := storeFake.NewMockStore(controller)
	noOwnerObjectStore.EXPECT().
		List(gomock.Any(), store.Key{
			Namespace:  "namespace",
			APIVersion: "v1",
			Kind:       "Pod",
		}).
		Return(testutil.ToUnstructuredList(t, pod3), true, nil).
		AnyTimes()
	noOwnerObjectStore.EXPECT().
		Get(gomock.Any(), store.Key{
			Namespace:  "namespace",
			APIVersion: "apps/v1",
			Kind:       "ReplicaSet",
			Name:       "missing",
		}).
		Return(nil, nil).
		AnyTimes()

	cases := []struct {
		name            string
		namespace       string
		objectStore     store.Store
		podMetricLoader octant.PodMetricsLoader
		expected        []octant.Workload
		isErr           bool
	}{
		{
			name:            "supports metrics",
			objectStore:     objectStore,
			podMetricLoader: podMetricsLoader(controller, pm, true, false),
			namespace:       "namespace",
			expected: []octant.Workload{
				{
					Name:     "pod-1",
					IconName: "application",
					Owner:    pod1,
					SegmentCounter: map[component.NodeStatus][]octant.PodWithMetric{
						component.NodeStatusOK: {
							{
								Pod:          pod1,
								ResourceList: rl,
							},
						},
					},
				},
				{
					Name:     "rs-1",
					IconName: "application",
					Owner:    replicaSet,
					SegmentCounter: map[component.NodeStatus][]octant.PodWithMetric{
						component.NodeStatusOK: {
							{
								Pod:          pod2,
								ResourceList: rl,
							},
						},
					},
				},
			},
		},
		{
			name:            "does not support metrics",
			objectStore:     objectStore,
			podMetricLoader: podMetricsLoader(controller, pm, false, false),
			namespace:       "namespace",
			expected: []octant.Workload{
				{
					Name:     "pod-1",
					IconName: "application",
					Owner:    pod1,
					SegmentCounter: map[component.NodeStatus][]octant.PodWithMetric{
						component.NodeStatusOK: {
							{
								Pod:          pod1,
								ResourceList: corev1.ResourceList{},
							},
						},
					},
				},
				{
					Name:     "rs-1",
					IconName: "application",
					Owner:    replicaSet,
					SegmentCounter: map[component.NodeStatus][]octant.PodWithMetric{
						component.NodeStatusOK: {
							{
								Pod:          pod2,
								ResourceList: corev1.ResourceList{},
							},
						},
					},
				},
			},
		},
		{
			name:            "pod with missing owner",
			objectStore:     noOwnerObjectStore,
			podMetricLoader: podMetricsLoader(controller, pm, false, false),
			namespace:       "namespace",
			expected: []octant.Workload{
				{
					Name:     "pod-3",
					IconName: "application",
					Owner:    pod3,
					SegmentCounter: map[component.NodeStatus][]octant.PodWithMetric{
						component.NodeStatusOK: {
							{
								Pod:          pod3,
								ResourceList: corev1.ResourceList{},
							},
						},
					},
				},
			},
		},
		{
			name:            "namespace is blank",
			objectStore:     objectStore,
			podMetricLoader: podMetricsLoader(controller, pm, true, false),
			isErr:           true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			wl, err := octant.NewClusterWorkloadLoader(c.objectStore, c.podMetricLoader, func(wl *octant.ClusterWorkloadLoader) {
				wl.ObjectStatuser = func(context.Context, runtime.Object, store.Store) (status objectstatus.ObjectStatus, e error) {
					objectStatus := objectstatus.ObjectStatus{}
					return objectStatus, nil

				}
			})
			require.NoError(t, err)

			actual, err := wl.Load(ctx, c.namespace)
			if c.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			testutil.AssertJSONEqual(t, c.expected, actual)
		})
	}
}

func TestWorkloadSummaryChart(t *testing.T) {
	pwm := octant.PodWithMetric{}

	w := octant.NewWorkload("workload", "icon")

	w.SegmentCounter = map[component.NodeStatus][]octant.PodWithMetric{
		component.NodeStatusOK:      {pwm},
		component.NodeStatusWarning: {pwm, pwm},
		component.NodeStatusError:   {pwm, pwm, pwm},
	}

	cases := []struct {
		name     string
		workload *octant.Workload
		expected *component.DonutChart
		isErr    bool
	}{
		{
			name:     "in general",
			workload: w,
			expected: &component.DonutChart{
				Config: component.DonutChartConfig{
					Size: component.DonutChartSizeMedium,
					Segments: []component.DonutSegment{
						{
							Count:  3,
							Status: component.NodeStatusError,
						},
						{
							Count:  1,
							Status: component.NodeStatusOK,
						},
						{
							Count:  2,
							Status: component.NodeStatusWarning,
						},
					},
					Labels: component.DonutChartLabels{
						Plural:   "Pods",
						Singular: "Pod",
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual, err := c.workload.DonutChart(component.DonutChartSizeMedium)
			if c.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			component.AssertEqual(t, c.expected, actual)
		})
	}
}

// nolint:dupl
func TestPodCPUStat(t *testing.T) {
	expectedTitle := "CPU(cores)"

	cases := []struct {
		name     string
		workload *octant.Workload
		expected *component.SingleStat
		isErr    bool
	}{
		{
			name: "pod with limit and request; measure below request",
			workload: createWorkload(t, "10m", "1Mi",
				podOptionWithLimitAndRequestValues(
					"500m", "14Mi", "250m", "7Mi")),
			expected: component.NewSingleStat(expectedTitle, "30m", octant.WorkloadStatusColorOK),
		},
		{
			name: "pod with limit and request; measure below limit",
			workload: createWorkload(t, "400m", "1Mi",
				podOptionWithLimitAndRequestValues(
					"500m", "14Mi", "250m", "7Mi")),
			expected: component.NewSingleStat(expectedTitle, "1200m", octant.WorkloadStatusColorWarning),
		},
		{
			name: "pod with limit and request; measure above limit",
			workload: createWorkload(t, "501m", "1Mi",
				podOptionWithLimitAndRequestValues(
					"500m", "14Mi", "250m", "7Mi")),
			expected: component.NewSingleStat(expectedTitle, "1503m", octant.WorkloadStatusColorError),
		},
		{
			name: "pod with limit; measure below limit",
			workload: createWorkload(t, "10m", "1Mi",
				podOptionWithLimitValues("500m", "14Mi")),
			expected: component.NewSingleStat(expectedTitle, "30m", octant.WorkloadStatusColorOK),
		},
		{
			name: "pod with request; measure below above request",
			workload: createWorkload(t, "501m", "1Mi",
				podOptionWithRequestValues("500m", "14Mi")),
			expected: component.NewSingleStat(expectedTitle, "1503m", octant.WorkloadStatusColorWarning),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual, err := octant.PodCPUStat(c.workload)
			if c.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			component.AssertEqual(t, c.expected, actual)
		})
	}
}

// nolint:dupl
func TestPodMemoryStat(t *testing.T) {
	expectedTitle := "Memory(bytes)"

	cases := []struct {
		name     string
		workload *octant.Workload
		expected *component.SingleStat
		isErr    bool
	}{
		{
			name: "pod with limit and request; measure below request",
			workload: createWorkload(t, "10m", "1Mi",
				podOptionWithLimitAndRequestValues(
					"500m", "14Mi", "250m", "7Mi")),
			expected: component.NewSingleStat(expectedTitle, "3Mi", octant.WorkloadStatusColorOK),
		},
		{
			name: "pod with limit and request; measure below limit, but above request",
			workload: createWorkload(t, "400m", "8Mi",
				podOptionWithLimitAndRequestValues(
					"500m", "14Mi", "250m", "7Mi")),
			expected: component.NewSingleStat(expectedTitle, "24Mi", octant.WorkloadStatusColorWarning),
		},
		{
			name: "pod with limit and request; measure above limit",
			workload: createWorkload(t, "501m", "15Mi",
				podOptionWithLimitAndRequestValues(
					"500m", "14Mi", "250m", "7Mi")),
			expected: component.NewSingleStat(expectedTitle, "45Mi", octant.WorkloadStatusColorError),
		},
		{
			name: "pod with limit; measure below limit",
			workload: createWorkload(t, "10m", "1Mi",
				podOptionWithLimitValues("500m", "14Mi")),
			expected: component.NewSingleStat(expectedTitle, "3Mi", octant.WorkloadStatusColorOK),
		},
		{
			name: "pod with request; measure below above request",
			workload: createWorkload(t, "501m", "15Mi",
				podOptionWithRequestValues("500m", "14Mi")),
			expected: component.NewSingleStat(expectedTitle, "45Mi", octant.WorkloadStatusColorWarning),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual, err := octant.PodMemoryStat(c.workload)
			if c.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			component.AssertEqual(t, c.expected, actual)
		})
	}
}

func TestCombineResourceRequirements(t *testing.T) {
	a := corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("25m"),
			corev1.ResourceMemory: resource.MustParse("1Mi"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("50m"),
			corev1.ResourceMemory: resource.MustParse("2Mi"),
		},
	}

	actual := octant.CombineResourceRequirements(a, a)

	assert.Equal(t, "100m", actual.Limits.Cpu().String())
	assert.Equal(t, "4Mi", actual.Limits.Memory().String())
	assert.Equal(t, "50m", actual.Requests.Cpu().String())
	assert.Equal(t, "2Mi", actual.Requests.Memory().String())
}

func createWorkload(t *testing.T, cpu, memory string, options ...testutil.PodOption) *octant.Workload {
	resourceList := createResourceList(cpu, memory)
	pod := testutil.ToUnstructured(t, testutil.CreatePod("pod", options...))

	workload := octant.NewWorkload("workload", "application")

	workload.AddPodStatus(component.NodeStatusOK, pod, resourceList.DeepCopy())
	workload.AddPodStatus(component.NodeStatusWarning, pod, resourceList.DeepCopy())
	workload.AddPodStatus(component.NodeStatusError, pod, resourceList.DeepCopy())

	return workload
}

func podOptionWithLimitAndRequestValues(limitCPU, limitMemory, requestCPU, requestMemory string) testutil.PodOption {
	return func(pod *corev1.Pod) {
		podOptionWithLimitValues(limitCPU, limitMemory)(pod)
		podOptionWithRequestValues(requestCPU, requestMemory)(pod)
	}
}

func podOptionWithLimitValues(cpu, memory string) testutil.PodOption {
	return func(pod *corev1.Pod) {
		if len(pod.Spec.Containers) != 1 {
			pod.Spec.Containers = []corev1.Container{{Name: "container"}}
		}

		pod.Spec.Containers[0].Resources.Limits = createResourceList(cpu, memory)
	}
}
func podOptionWithRequestValues(cpu, memory string) testutil.PodOption {
	return func(pod *corev1.Pod) {
		if len(pod.Spec.Containers) != 1 {
			pod.Spec.Containers = []corev1.Container{{Name: "container"}}
		}

		pod.Spec.Containers[0].Resources.Requests = createResourceList(cpu, memory)
	}
}

func createResourceList(cpu, memory string) corev1.ResourceList {
	return corev1.ResourceList{
		corev1.ResourceCPU:    resource.MustParse(cpu),
		corev1.ResourceMemory: resource.MustParse(memory),
	}
}

func podMetricWithContainer(podMetric *metricsv1beta1.PodMetrics) {
	podMetric.Containers = []metricsv1beta1.ContainerMetrics{
		{
			Name: "container",
			Usage: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("150m"),
				corev1.ResourceMemory: resource.MustParse("3Mi"),
			},
		},
	}
}
