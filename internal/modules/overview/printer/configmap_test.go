/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/pkg/view/component"
)

func Test_ConfigMapListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	labels := map[string]string{
		"foo": "bar",
	}

	now := time.Unix(1547211430, 0)

	configMap := testutil.CreateConfigMap("configMap")
	configMap.Data = map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	configMap.CreationTimestamp = metav1.Time{Time: now}
	configMap.Labels = labels

	tpo.PathForObject(configMap, configMap.Name, "/configMap")

	object := &corev1.ConfigMapList{
		Items: []corev1.ConfigMap{*configMap},
	}

	ctx := context.Background()
	got, err := ConfigMapListHandler(ctx, object, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Labels", "Data", "Age")
	expected := component.NewTable("ConfigMaps", cols)
	expected.Add(component.TableRow{
		"Name":   component.NewLink("", configMap.Name, "/configMap"),
		"Labels": component.NewLabels(labels),
		"Data":   component.NewText("2"),
		"Age":    component.NewTimestamp(now),
	})

	assertComponentEqual(t, expected, got)
}

func Test_describeConfigMapConfiguration(t *testing.T) {
	var validConfigMap = &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "env-config",
			CreationTimestamp: metav1.Time{
				Time: time.Unix(1548377609, 0),
			},
		},
		Data: map[string]string{
			"log_level": "INFO",
		},
	}

	cases := []struct {
		name      string
		configmap *corev1.ConfigMap
		isErr     bool
		expected  *component.Summary
	}{
		{
			name:      "configmap",
			configmap: validConfigMap,
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Age",
					Content: component.NewTimestamp(time.Unix(1548377609, 0)),
				},
			}...),
		},
		{
			name:      "configmap is nil",
			configmap: nil,
			isErr:     true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			summary, err := describeConfigMapConfig(tc.configmap)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expected, summary)
		})
	}
}
