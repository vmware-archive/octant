package printer_test

import (
	"testing"
	"time"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/overview/printer"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_ConfigMapListHandler(t *testing.T) {
	printOptions := printer.Options{
		Cache: cache.NewMemoryCache(),
	}

	labels := map[string]string{
		"foo": "bar",
	}

	object := &corev1.ConfigMapList{
		Items: []corev1.ConfigMap{
			{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "ConfigMap",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "configmap",
					CreationTimestamp: metav1.Time{
						Time: now,
					},
					Labels: labels,
				},
				Data: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
		},
	}

	got, err := printer.ConfigMapListHandler(object, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Labels", "Data", "Age")
	expected := component.NewTable("ConfigMaps", cols)
	expected.Add(component.TableRow{
		"Name":   component.NewLink("", "configmap", "/content/overview/config-and-storage/config-maps/configmap"),
		"Labels": component.NewLabels(labels),
		"Data":   component.NewText("", "2"),
		"Age":    component.NewTimestamp(now),
	})

	assert.Equal(t, expected, got)
}

func TestConfigMapConfiguration(t *testing.T) {
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
			cc := printer.NewConfigMapConfiguration(tc.configmap)

			summary, err := cc.Create()
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expected, summary)
		})
	}
}
