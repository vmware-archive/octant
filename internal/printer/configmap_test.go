/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_ConfigMapListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	labels := map[string]string{
		"foo": "bar",
	}

	now := testutil.Time()

	configMap := testutil.CreateConfigMap("configMap")
	configMap.Data = map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	configMap.CreationTimestamp = *testutil.CreateTimestamp()
	configMap.Labels = labels

	tpo.PathForObject(configMap, configMap.Name, "/configMap")

	object := &corev1.ConfigMapList{
		Items: []corev1.ConfigMap{*configMap},
	}

	ctx := context.Background()
	tpo.pluginManager.EXPECT().ObjectStatus(ctx, configMap)
	got, err := ConfigMapListHandler(ctx, object, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Labels", "Data", "Age")
	expected := component.NewTable("ConfigMaps", "We couldn't find any config maps!", cols)
	expected.Add(component.TableRow{
		"Name": component.NewLink("", configMap.Name, "/configMap",
			genObjectStatus(component.TextStatusOK, []string{"v1 ConfigMap is OK"})),
		"Labels": component.NewLabels(labels),
		"Data":   component.NewText("2"),
		"Age":    component.NewTimestamp(now),
		component.GridActionKey: gridActionsFactory([]component.GridAction{
			buildObjectDeleteAction(t, configMap),
		}),
	})

	component.AssertEqual(t, expected, got)
}

func Test_ConfigMapConfiguration(t *testing.T) {
	var validConfigMap = &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "env-config",
			CreationTimestamp: metav1.Time{
				Time: testutil.Time(),
			},
		},
		Data: map[string]string{
			"log_level": "INFO",
		},
	}

	cases := []struct {
		name      string
		configMap *corev1.ConfigMap
		isErr     bool
		expected  *component.Summary
	}{
		{
			name:      "configmap",
			configMap: validConfigMap,
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Age",
					Content: component.NewTimestamp(testutil.Time()),
				},
			}...),
		},
		{
			name:      "configmap is nil",
			configMap: nil,
			isErr:     true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			tpo := newTestPrinterOptions(controller)
			printOptions := tpo.ToOptions()

			cc := NewConfigMapConfiguration(tc.configMap)

			summary, err := cc.Create(printOptions)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			component.AssertEqual(t, tc.expected, summary)
		})
	}
}

func Test_describeConfigMapData(t *testing.T) {
	configMap := testutil.CreateConfigMap("configmap")
	configMap.Data = map[string]string{
		"foo": "bar",
		"bar": "foo",
	}

	got, err := describeConfigMapData(configMap)
	require.NoError(t, err)

	cols := component.NewTableCols("Key", "Value")
	expected := component.NewTable("Data", "No data has been configured for this config map!", cols)
	expected.Add([]component.TableRow{
		{"Key": component.NewText("bar"), "Value": component.NewText("foo")},
		{"Key": component.NewText("foo"), "Value": component.NewText("bar")},
	}...)

	component.AssertEqual(t, expected, got)
}
