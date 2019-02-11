package printer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/heptio/developer-dash/internal/overview/printer"
	"github.com/heptio/developer-dash/internal/view/component"
)

func Test_AffinityDescriber_Create(t *testing.T) {
	cases := []struct {
		name     string
		affinity *corev1.Affinity
		expected *component.List
		isErr    bool
	}{
		{
			name: "preferred node label value in",
			affinity: &corev1.Affinity{
				NodeAffinity: &corev1.NodeAffinity{
					PreferredDuringSchedulingIgnoredDuringExecution: []corev1.PreferredSchedulingTerm{
						{
							Preference: corev1.NodeSelectorTerm{
								MatchExpressions: []corev1.NodeSelectorRequirement{
									{
										Key:      "foo",
										Operator: corev1.NodeSelectorOpIn,
										Values:   []string{"x", "y"},
									},
								},
							},
						},
					},
				},
			},
			expected: component.NewList("", []component.ViewComponent{
				component.NewText("Prefer to schedule on nodes with label foo with values x, y."),
			}),
		},
		{
			name: "preferred node label value not in",
			affinity: &corev1.Affinity{
				NodeAffinity: &corev1.NodeAffinity{
					PreferredDuringSchedulingIgnoredDuringExecution: []corev1.PreferredSchedulingTerm{
						{
							Preference: corev1.NodeSelectorTerm{
								MatchExpressions: []corev1.NodeSelectorRequirement{
									{
										Key:      "foo",
										Operator: corev1.NodeSelectorOpNotIn,
										Values:   []string{"x", "y"},
									},
								},
							},
						},
					},
				},
			},
			expected: component.NewList("", []component.ViewComponent{
				component.NewText("Prefer to schedule on nodes with label foo without values x, y."),
			}),
		},
		{
			name: "preferred node label exists",
			affinity: &corev1.Affinity{
				NodeAffinity: &corev1.NodeAffinity{
					PreferredDuringSchedulingIgnoredDuringExecution: []corev1.PreferredSchedulingTerm{
						{
							Preference: corev1.NodeSelectorTerm{
								MatchExpressions: []corev1.NodeSelectorRequirement{
									{
										Key:      "foo",
										Operator: corev1.NodeSelectorOpExists,
									},
								},
							},
						},
					},
				},
			},
			expected: component.NewList("", []component.ViewComponent{
				component.NewText("Prefer to schedule on nodes where label foo exists."),
			}),
		},
		{
			name: "preferred node label does not exists",
			affinity: &corev1.Affinity{
				NodeAffinity: &corev1.NodeAffinity{
					PreferredDuringSchedulingIgnoredDuringExecution: []corev1.PreferredSchedulingTerm{
						{
							Preference: corev1.NodeSelectorTerm{
								MatchExpressions: []corev1.NodeSelectorRequirement{
									{
										Key:      "foo",
										Operator: corev1.NodeSelectorOpDoesNotExist,
									},
								},
							},
						},
					},
				},
			},
			expected: component.NewList("", []component.ViewComponent{
				component.NewText("Prefer to schedule on nodes where label foo does not exist."),
			}),
		},
		{
			name: "preferred node label greater than",
			affinity: &corev1.Affinity{
				NodeAffinity: &corev1.NodeAffinity{
					PreferredDuringSchedulingIgnoredDuringExecution: []corev1.PreferredSchedulingTerm{
						{
							Preference: corev1.NodeSelectorTerm{
								MatchExpressions: []corev1.NodeSelectorRequirement{
									{
										Key:      "foo",
										Operator: corev1.NodeSelectorOpGt,
										Values:   []string{"1"},
									},
								},
							},
						},
					},
				},
			},
			expected: component.NewList("", []component.ViewComponent{
				component.NewText("Prefer to schedule on nodes where label foo is greater than 1."),
			}),
		},
		{
			name: "preferred node label less than",
			affinity: &corev1.Affinity{
				NodeAffinity: &corev1.NodeAffinity{
					PreferredDuringSchedulingIgnoredDuringExecution: []corev1.PreferredSchedulingTerm{
						{
							Preference: corev1.NodeSelectorTerm{
								MatchExpressions: []corev1.NodeSelectorRequirement{
									{
										Key:      "foo",
										Operator: corev1.NodeSelectorOpLt,
										Values:   []string{"1"},
									},
								},
							},
						},
					},
				},
			},
			expected: component.NewList("", []component.ViewComponent{
				component.NewText("Prefer to schedule on nodes where label foo is less than 1."),
			}),
		},
		{
			name: "preferred node field value in",
			affinity: &corev1.Affinity{
				NodeAffinity: &corev1.NodeAffinity{
					PreferredDuringSchedulingIgnoredDuringExecution: []corev1.PreferredSchedulingTerm{
						{
							Preference: corev1.NodeSelectorTerm{
								MatchFields: []corev1.NodeSelectorRequirement{
									{
										Key:      "foo",
										Operator: corev1.NodeSelectorOpIn,
										Values:   []string{"x", "y"},
									},
								},
							},
						},
					},
				},
			},
			expected: component.NewList("", []component.ViewComponent{
				component.NewText("Prefer to schedule on nodes with field foo with values x, y."),
			}),
		},
		{
			name: "preferred node field with weight",
			affinity: &corev1.Affinity{
				NodeAffinity: &corev1.NodeAffinity{
					PreferredDuringSchedulingIgnoredDuringExecution: []corev1.PreferredSchedulingTerm{
						{
							Weight: 10,
							Preference: corev1.NodeSelectorTerm{
								MatchFields: []corev1.NodeSelectorRequirement{
									{
										Key:      "foo",
										Operator: corev1.NodeSelectorOpIn,
										Values:   []string{"x", "y"},
									},
								},
							},
						},
					},
				},
			},
			expected: component.NewList("", []component.ViewComponent{
				component.NewText("Prefer to schedule on nodes with field foo with values x, y. Weight 10."),
			}),
		},
		{
			name: "required node field with weight",
			affinity: &corev1.Affinity{
				NodeAffinity: &corev1.NodeAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
						NodeSelectorTerms: []corev1.NodeSelectorTerm{
							{
								MatchFields: []corev1.NodeSelectorRequirement{
									{
										Key:      "foo",
										Operator: corev1.NodeSelectorOpIn,
										Values:   []string{"x", "y"},
									},
								},
							},
						},
					},
				},
			},
			expected: component.NewList("", []component.ViewComponent{
				component.NewText("Schedule on nodes with field foo with values x, y."),
			}),
		},
		{
			name: "affinity: required pod label selector with match labels",
			affinity: &corev1.Affinity{
				PodAffinity: &corev1.PodAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
						{
							LabelSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"foo": "bar",
									"bar": "foo",
								},
							},
							TopologyKey: "topology",
						},
					},
				},
			},
			expected: component.NewList("", []component.ViewComponent{
				component.NewText("Schedule with pod labeled bar:foo, foo:bar in topology topology."),
			}),
		},
		{
			name: "affinity: required pod label selector with match expressions",
			affinity: &corev1.Affinity{
				PodAffinity: &corev1.PodAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
						{
							LabelSelector: &metav1.LabelSelector{
								MatchExpressions: []metav1.LabelSelectorRequirement{
									{
										Key:      "key",
										Operator: metav1.LabelSelectorOpExists,
									},
								},
							},
							TopologyKey: "topology",
						},
					},
				},
			},
			expected: component.NewList("", []component.ViewComponent{
				component.NewText("Schedule with pod where key exists in topology topology."),
			}),
		},
		{
			name: "affinity: required pod label selector with match expressions and match labels",
			affinity: &corev1.Affinity{
				PodAffinity: &corev1.PodAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
						{
							LabelSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"foo": "bar",
									"bar": "foo",
								},
								MatchExpressions: []metav1.LabelSelectorRequirement{
									{
										Key:      "key",
										Operator: metav1.LabelSelectorOpExists,
									},
								},
							},
							TopologyKey: "topology",
						},
					},
				},
			},
			expected: component.NewList("", []component.ViewComponent{
				component.NewText("Schedule with pod labeled bar:foo, foo:bar where key exists in topology topology."),
			}),
		},
		{
			name: "affinity: preferred pod label selector with match labels",
			affinity: &corev1.Affinity{
				PodAffinity: &corev1.PodAffinity{
					PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
						{
							PodAffinityTerm: corev1.PodAffinityTerm{
								LabelSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"foo": "bar",
										"bar": "foo",
									},
								},
								TopologyKey: "topology",
							},
						},
					},
				},
			},
			expected: component.NewList("", []component.ViewComponent{
				component.NewText("Prefer to schedule with pod labeled bar:foo, foo:bar in topology topology."),
			}),
		},
		{
			name: "affinity: preferred pod label selector with match labels weighed",
			affinity: &corev1.Affinity{
				PodAffinity: &corev1.PodAffinity{
					PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
						{
							Weight: 5,
							PodAffinityTerm: corev1.PodAffinityTerm{
								LabelSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"foo": "bar",
										"bar": "foo",
									},
								},
								TopologyKey: "topology",
							},
						},
					},
				},
			},
			expected: component.NewList("", []component.ViewComponent{
				component.NewText("Prefer to schedule with pod labeled bar:foo, foo:bar in topology topology. Weight 5."),
			}),
		},
		{
			name: "anti-affinity: preferred pod label selector with match labels",
			affinity: &corev1.Affinity{
				PodAntiAffinity: &corev1.PodAntiAffinity{
					PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
						{
							PodAffinityTerm: corev1.PodAffinityTerm{
								LabelSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"foo": "bar",
										"bar": "foo",
									},
								},
								TopologyKey: "topology",
							},
						},
					},
				},
			},
			expected: component.NewList("", []component.ViewComponent{
				component.NewText("Prefer to not schedule with pod labeled bar:foo, foo:bar in topology topology."),
			}),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			podSpec := corev1.PodSpec{
				Affinity: tc.affinity,
			}

			ad := printer.NewAffinityDescriber(podSpec)

			got, err := ad.Create()
			if tc.isErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, tc.expected, got)
		})
	}
}
