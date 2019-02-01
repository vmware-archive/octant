package printer_test

import (
	"testing"

	corev1 "k8s.io/api/core/v1"

	"github.com/heptio/developer-dash/internal/overview/printer"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_TolerationDescriber_Create(t *testing.T) {
	cases := []struct {
		name        string
		tolerations []corev1.Toleration
		expected    *component.List
		isErr       bool
	}{
		{
			name: "key,value",
			tolerations: []corev1.Toleration{
				{
					Key:   "key",
					Value: "value",
				},
			},
			expected: component.NewList("", []component.ViewComponent{
				component.NewText("", "Schedule on nodes with key:value taint."),
			}),
		},
		{
			name: "multiple tolerations",
			tolerations: []corev1.Toleration{
				{
					Key:   "key1",
					Value: "value1",
				},
				{
					Key:   "key2",
					Value: "value2",
				},
			},
			expected: component.NewList("", []component.ViewComponent{
				component.NewText("", "Schedule on nodes with key1:value1 taint."),
				component.NewText("", "Schedule on nodes with key2:value2 taint."),
			}),
		},
		{
			name: "key,value",
			tolerations: []corev1.Toleration{
				{
					Key:    "key",
					Value:  "value",
					Effect: corev1.TaintEffectNoSchedule,
				},
			},
			expected: component.NewList("", []component.ViewComponent{
				component.NewText("", "Schedule on nodes with key:value:NoSchedule taint."),
			}),
		},
		{
			name: "effect",
			tolerations: []corev1.Toleration{
				{
					Effect: corev1.TaintEffectNoSchedule,
				},
			},
			expected: component.NewList("", []component.ViewComponent{
				component.NewText("", "Schedule on nodes with NoSchedule taint."),
			}),
		},
		{
			name: "key,value with evict secs",
			tolerations: []corev1.Toleration{
				{
					Key:               "key",
					Value:             "value",
					TolerationSeconds: ptr64(3600),
				},
			},
			expected: component.NewList("", []component.ViewComponent{
				component.NewText("",
					"Schedule on nodes with key:value taint. Evict after 3600 seconds."),
			}),
		},
		{
			name: "key exists",
			tolerations: []corev1.Toleration{
				{
					Key:      "key",
					Operator: corev1.TolerationOpExists,
				},
			},
			expected: component.NewList("", []component.ViewComponent{
				component.NewText("",
					"Schedule on nodes with key taint."),
			}),
		},
		{
			name: "exists with no key",
			tolerations: []corev1.Toleration{
				{
					Operator: corev1.TolerationOpExists,
				},
			},
			expected: component.NewList("", []component.ViewComponent{
				component.NewText("",
					"Schedule on all nodes."),
			}),
		},
		{
			name: "unsupported toleration",
			tolerations: []corev1.Toleration{
				{
					Key: "key",
				},
			},
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			podSpec := corev1.PodSpec{
				Tolerations: tc.tolerations,
			}

			td := printer.NewTolerationDescriber(podSpec)

			got, err := td.Create()
			if tc.isErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, tc.expected, got)
		})
	}
}
