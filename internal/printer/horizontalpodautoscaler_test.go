/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/apis/autoscaling"
	"k8s.io/kubernetes/pkg/apis/core"
)

func Test_HorizontalPodAutoscalerListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	objectLabels := map[string]string{
		"foo": "bar",
	}

	now := testutil.Time()
	var minReplicas int32 = 1
	var currentPercentCPU int32 = 5
	var targetPercentCPU int32 = 50

	object := &autoscalingv1.HorizontalPodAutoscaler{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "autoscaling/v1",
			Kind:       "HorizontalPodAutoscaler",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "horizontalpodautoscaler",
			Namespace: "default",
			CreationTimestamp: metav1.Time{
				Time: now,
			},
			Labels: objectLabels,
		},
		Spec: autoscalingv1.HorizontalPodAutoscalerSpec{
			MinReplicas:                    &minReplicas,
			MaxReplicas:                    10,
			TargetCPUUtilizationPercentage: &targetPercentCPU,
		},
		Status: autoscalingv1.HorizontalPodAutoscalerStatus{
			CurrentReplicas:                 2,
			CurrentCPUUtilizationPercentage: &currentPercentCPU,
		},
	}

	tpo.PathForObject(object, object.Name, "/path")

	list := &autoscalingv1.HorizontalPodAutoscalerList{
		Items: []autoscalingv1.HorizontalPodAutoscaler{*object},
	}

	ctx := context.Background()
	got, err := HorizontalPodAutoscalerListHandler(ctx, list, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Labels", "Targets", "Minimum Pods", "Maximum Pods", "Replicas", "Age")
	expected := component.NewTable("Horizontal Pod Autoscalers", "We couldn't find any horizontal pod autoscalers", cols)
	expected.Add(component.TableRow{
		"Name":         component.NewLink("", "horizontalpodautoscaler", "/path"),
		"Labels":       component.NewLabels(objectLabels),
		"Targets":      component.NewText("5/50%"),
		"Minimum Pods": component.NewText("1"),
		"Maximum Pods": component.NewText("10"),
		"Replicas":     component.NewText("2"),
		"Age":          component.NewTimestamp(now),
	})

	component.AssertEqual(t, expected, got)
}

func Test_HorizontalPodAutoscalerConfiguration(t *testing.T) {
	var replicas int32 = 1
	hpa := testutil.CreateHorizontalPodAutoscaler("hpa")
	hpa.Spec.MinReplicas = &replicas
	hpa.Spec.MaxReplicas = 10

	deployment := testutil.CreateDeployment("deployment")
	hpa.Spec.ScaleTargetRef = autoscalingv1.CrossVersionObjectReference{
		Kind:       deployment.Kind,
		APIVersion: deployment.APIVersion,
		Name:       deployment.Name,
	}

	cases := []struct {
		name                    string
		horizontalPodAutoscaler *autoscalingv1.HorizontalPodAutoscaler
		expected                component.Component
		isErr                   bool
	}{
		{
			name:                    "general",
			horizontalPodAutoscaler: hpa,
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Scale target",
					Content: component.NewLink("", "deployment", "/deployment"),
				},
				{
					Header:  "Min Replicas",
					Content: component.NewText("1"),
				},
				{
					Header:  "Max Replicas",
					Content: component.NewText("10"),
				},
			}...),
		},
		{
			name:                    "nil horizontalpodautoscaler",
			horizontalPodAutoscaler: nil,
			isErr:                   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			ctx := context.Background()

			tpo := newTestPrinterOptions(controller)
			printOptions := tpo.ToOptions()

			hc := NewHorizontalPodAutoscalerConfiguration(tc.horizontalPodAutoscaler)

			scaleTarget := component.NewLink("", "deployment", "/deployment")
			tpo.link.EXPECT().
				ForGVK("namespace", "apps/v1", "Deployment", "deployment", "deployment").
				Return(scaleTarget, nil).
				AnyTimes()

			if tc.horizontalPodAutoscaler != nil {
				key := store.Key{
					APIVersion: deployment.APIVersion,
					Kind:       deployment.Kind,
					Name:       deployment.Name,
					Namespace:  deployment.Namespace,
				}
				tpo.objectStore.EXPECT().Get(ctx, gomock.Eq(key)).Return(testutil.ToUnstructured(t, deployment), nil)
			}

			summary, err := hc.Create(ctx, printOptions)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			component.AssertEqual(t, tc.expected, summary)
		})
	}
}

func Test_createHorizontalPodAutoscalerSummaryStatus(t *testing.T) {
	var observedGeneration int64 = 1
	var currentCPU int32 = 3
	now := testutil.Time()

	hpa := testutil.CreateHorizontalPodAutoscaler("hpa")
	hpa.Status.ObservedGeneration = &observedGeneration
	hpa.Status.LastScaleTime = &metav1.Time{Time: now}
	hpa.Status.CurrentReplicas = 2
	hpa.Status.DesiredReplicas = 7
	hpa.Status.CurrentCPUUtilizationPercentage = &currentCPU

	cases := []struct {
		name                    string
		horizontalPodAutoscaler *autoscalingv1.HorizontalPodAutoscaler
		expected                *component.Summary
		isErr                   bool
	}{
		{
			name:                    "in general",
			horizontalPodAutoscaler: hpa,
			expected: component.NewSummary("Status", []component.SummarySection{
				{
					Header:  "Targets",
					Content: component.NewText("3/80%"),
				},
				{
					Header:  "Observed Generation",
					Content: component.NewText("1"),
				},
				{
					Header:  "Last Scale Time",
					Content: component.NewTimestamp(now),
				},
				{
					Header:  "Current Replicas",
					Content: component.NewText("2"),
				},
				{
					Header:  "Desired Replicas",
					Content: component.NewText("7"),
				},
				{
					Header:  "Current CPU Utilization Percentage",
					Content: component.NewText("3"),
				},
			}...),
		},
		{
			name:                    "nil horizontalpodautoscaler",
			horizontalPodAutoscaler: nil,
			isErr:                   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			summary, err := createHorizontalPodAutoscalerSummaryStatus(tc.horizontalPodAutoscaler)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			component.AssertEqual(t, tc.expected, summary)
		})
	}
}

func Test_HorizontalPodAutoscalerMetrics(t *testing.T) {
	var averageUtilization int32 = 20
	metricResource := autoscaling.MetricStatus{
		Type: autoscaling.ResourceMetricSourceType,
		Resource: &autoscaling.ResourceMetricStatus{
			Name: core.ResourceCPU,
			Current: autoscaling.MetricValueStatus{
				AverageUtilization: &averageUtilization,
			},
		},
	}

	metricPods := autoscaling.MetricStatus{
		Type: autoscaling.PodsMetricSourceType,
		Pods: &autoscaling.PodsMetricStatus{
			Metric: autoscaling.MetricIdentifier{
				Name: "packets-per-second",
			},
			Current: autoscaling.MetricValueStatus{
				Value: resource.NewMilliQuantity(1000, resource.DecimalSI),
			},
		},
	}

	cases := []struct {
		name         string
		metricStatus *autoscaling.MetricStatus
		expected     *component.Summary
		isErr        bool
	}{
		{
			name:         "resource type",
			metricStatus: &metricResource,
			expected: component.NewSummary("Metric", []component.SummarySection{
				{
					Header:  "Type",
					Content: component.NewText("Resource"),
				},
				{
					Header:  "Name",
					Content: component.NewText("cpu"),
				},
				{
					Header:  "Average Utilization",
					Content: component.NewText("20%"),
				},
			}...),
		},
		{
			name:         "pods type",
			metricStatus: &metricPods,
			expected: component.NewSummary("Metric", []component.SummarySection{
				{
					Header:  "Type",
					Content: component.NewText("Pods"),
				},
				{
					Header:  "Name",
					Content: component.NewText("packets-per-second"),
				},
				{
					Header:  "Value",
					Content: component.NewText("1"),
				},
			}...),
		},
		{
			name:         "nil metric",
			metricStatus: nil,
			isErr:        true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			tpo := newTestPrinterOptions(controller)
			printOptions := tpo.ToOptions()

			summary, err := createHorizontalPodAutoscalerMetricsStatusView(tc.metricStatus, printOptions)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			component.AssertEqual(t, tc.expected, summary)
		})
	}
}

func Test_createHorizontalPodAutoscalerConditionsView(t *testing.T) {
	now := metav1.Time{Time: time.Now()}

	horizontalPodAutoscaler := testutil.CreateHorizontalPodAutoscaler("hpa")
	converted, err := convertToAutoscaling(horizontalPodAutoscaler)
	require.NoError(t, err)

	converted.Status.Conditions = []autoscaling.HorizontalPodAutoscalerCondition{
		{
			Type:               autoscaling.AbleToScale,
			Reason:             "reason",
			Status:             autoscaling.ConditionTrue,
			Message:            "message",
			LastTransitionTime: now,
		},
	}

	v1, err := convertToV1(converted)
	require.NoError(t, err)

	got, err := createHorizontalPodAutoscalerConditionsView(v1)
	require.NoError(t, err)

	cols := component.NewTableCols("Type", "Reason", "Status", "Message", "Last Transition")
	expected := component.NewTable("Conditions", "There are no horizontalpodautoscaler conditions!", cols)
	expected.Add([]component.TableRow{
		{
			"Type":            component.NewText("AbleToScale"),
			"Reason":          component.NewText("reason"),
			"Status":          component.NewText("True"),
			"Message":         component.NewText("message"),
			"Last Transition": component.NewTimestamp(now.Time),
		},
	}...)

	component.AssertEqual(t, expected, got)
}
