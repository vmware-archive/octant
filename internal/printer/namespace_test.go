/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_NamespaceListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	namespace := testutil.CreateNamespace("ns-test-1")
	namespace.CreationTimestamp = *testutil.CreateTimestamp()
	namespace.Status = corev1.NamespaceStatus{Phase: corev1.NamespaceActive}

	list := &corev1.NamespaceList{
		Items: []corev1.Namespace{
			*namespace,
		},
	}

	ctx := context.Background()
	got, err := NamespaceListHandler(ctx, list, printOptions)
	require.NoError(t, err)

	expected := component.NewTableWithRows("Namespaces", "We couldn't find any namespaces!", namespaceListCols, []component.TableRow{
		{
			"Name":   component.NewLink("", "ns-test-1", "/cluster-overview/namespaces/ns-test-1", genObjectStatus(component.TextStatusOK, []string{"v1 Namespace is OK"})),
			"Labels": component.NewLabels(make(map[string]string)),
			"Status": component.NewText("Active"),
			"Age":    component.NewTimestamp(namespace.CreationTimestamp.Time),
			component.GridActionKey: gridActionsFactory([]component.GridAction{
				buildObjectDeleteAction(t, namespace),
			}),
		},
	})

	testutil.AssertJSONEqual(t, expected, got)
}

func Test_printNamespaceResourceQuotas(t *testing.T) {
	max, err := resource.ParseQuantity("10")
	require.NoError(t, err)

	used, err := resource.ParseQuantity("0")
	require.NoError(t, err)

	quotas := []corev1.ResourceQuota{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-2",
				Namespace: "ns-test-2",
				UID:       types.UID("test-2"),
			},
			Spec: corev1.ResourceQuotaSpec{
				Hard: corev1.ResourceList{corev1.ResourceStorage: max},
			},
			Status: corev1.ResourceQuotaStatus{
				Hard: corev1.ResourceList{corev1.ResourceStorage: max},
				Used: corev1.ResourceList{corev1.ResourceStorage: used},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-3",
				Namespace: "ns-test-2",
				UID:       types.UID("test-3"),
			},
			Spec: corev1.ResourceQuotaSpec{
				Hard: corev1.ResourceList{corev1.ResourcePods: max},
			},
			Status: corev1.ResourceQuotaStatus{
				Hard: corev1.ResourceList{corev1.ResourcePods: max},
				Used: corev1.ResourceList{corev1.ResourcePods: used},
			},
		},
	}

	table1 := component.NewTableWithRows("test-2", "There are no resource quotas", namespaceResourceQuotasCols, []component.TableRow{
		{
			"Resource": component.NewText("storage"),
			"Used":     component.NewText("0"),
			"Limit":    component.NewText("10"),
		},
	})
	table2 := component.NewTableWithRows("test-3", "There are no resource quotas", namespaceResourceQuotasCols, []component.TableRow{
		{
			"Resource": component.NewText("pods"),
			"Used":     component.NewText("0"),
			"Limit":    component.NewText("10"),
		},
	})

	expected := map[string]component.FlexLayoutItem{
		"test-2": component.FlexLayoutItem{Width: component.WidthHalf, View: table1},
		"test-3": component.FlexLayoutItem{Width: component.WidthHalf, View: table2},
	}

	got := printNamespaceResourceQuotas(quotas)

	for k := range got {
		g := got[k]
		e := expected[k]
		component.AssertEqual(t, e.View, g.View)
		assert.Equal(t, e.Width, g.Width)
	}
}

func Test_printNamespaceResourceLimits(t *testing.T) {
	min, err := resource.ParseQuantity("200")
	require.NoError(t, err)
	max, err := resource.ParseQuantity("400")
	require.NoError(t, err)

	limits := corev1.LimitRangeList{
		Items: []corev1.LimitRange{
			corev1.LimitRange{
				Spec: corev1.LimitRangeSpec{
					Limits: []corev1.LimitRangeItem{
						{
							Type: corev1.LimitTypeContainer,
							Min:  corev1.ResourceList{corev1.ResourceCPU: min},
							Max:  corev1.ResourceList{corev1.ResourceCPU: max},
						},
					},
				},
			},
			corev1.LimitRange{
				Spec: corev1.LimitRangeSpec{
					Limits: []corev1.LimitRangeItem{
						{
							Type: corev1.LimitTypeContainer,
							Min:  corev1.ResourceList{corev1.ResourceMemory: min},
							Max:  corev1.ResourceList{corev1.ResourceMemory: max},
						},
					},
				},
			},
		},
	}

	got, err := printNamespaceResourceLimits(&limits)
	require.NoError(t, err)

	// ("Type", "Resource", "Min", "Max", "Default Request", "Default Limit", "Limit/Request Ratio")
	expected := component.NewTableWithRows("Resource Limits", "There are no resource limits", namespaceResourceLimitsCols, []component.TableRow{
		{
			"Type":                component.NewText("Container"),
			"Resource":            component.NewText("cpu"),
			"Min":                 component.NewText("200"),
			"Max":                 component.NewText("400"),
			"Default Request":     component.NewText("0"),
			"Default Limit":       component.NewText("0"),
			"Limit/Request Ratio": component.NewText("0"),
		},
		{
			"Type":                component.NewText("Container"),
			"Resource":            component.NewText("memory"),
			"Min":                 component.NewText("200"),
			"Max":                 component.NewText("400"),
			"Default Request":     component.NewText("0"),
			"Default Limit":       component.NewText("0"),
			"Limit/Request Ratio": component.NewText("0"),
		},
	})
	component.AssertEqual(t, expected, got)
}

func Test_createResourceLimitCPURow(t *testing.T) {
	min, err := resource.ParseQuantity("150")
	require.NoError(t, err)
	max, err := resource.ParseQuantity("300")
	require.NoError(t, err)

	cpu := corev1.LimitRangeItem{
		Type: corev1.LimitTypeContainer,
		Min:  corev1.ResourceList{corev1.ResourceCPU: min},
		Max:  corev1.ResourceList{corev1.ResourceCPU: max},
	}
	_, _, created := createResourceLimitCPURow(cpu)
	require.True(t, created)

	mem := corev1.LimitRangeItem{
		Type: corev1.LimitTypeContainer,
		Min:  corev1.ResourceList{corev1.ResourceMemory: min},
		Max:  corev1.ResourceList{corev1.ResourceMemory: max},
	}

	_, _, created = createResourceLimitCPURow(mem)
	require.False(t, created)
}

func Test_createResourceLimitMemoryRow(t *testing.T) {
	min, err := resource.ParseQuantity("150")
	require.NoError(t, err)
	max, err := resource.ParseQuantity("300")
	require.NoError(t, err)

	mem := corev1.LimitRangeItem{
		Type: corev1.LimitTypeContainer,
		Min:  corev1.ResourceList{corev1.ResourceMemory: min},
		Max:  corev1.ResourceList{corev1.ResourceMemory: max},
	}

	_, _, created := createResourceLimitMemoryRow(mem)
	require.True(t, created)

	cpu := corev1.LimitRangeItem{
		Type: corev1.LimitTypeContainer,
		Min:  corev1.ResourceList{corev1.ResourceCPU: min},
		Max:  corev1.ResourceList{corev1.ResourceCPU: max},
	}
	_, _, created = createResourceLimitMemoryRow(cpu)
	require.False(t, created)
}
