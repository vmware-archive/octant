/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_unmarshal(t *testing.T) {
	cases := []struct {
		name       string
		configFile string
		objectType string
		expected   Component
	}{
		{
			name:       "cardList",
			configFile: "config_card_list.json",
			objectType: typeCardList,
			expected: &CardList{
				Config: CardListConfig{
					Cards: []Card{
						{
							base: newBase(typeCard, TitleFromString("card title")),
							Config: CardConfig{
								Body: NewText("text value"),
								Actions: []Action{
									{
										Name:  "Edit",
										Title: "Edit",
										Form: Form{
											Fields: []FormField{
												NewFormFieldText("Revision", "revision", "12345"),
											},
										},
									},
								},
								Alert: &Alert{
									Type:    AlertTypeWarning,
									Message: "warning",
								},
							},
						},
					},
				},
			},
		},
		{
			name:       "containers",
			configFile: "config_containers.json",
			objectType: "containers",
			expected: &Containers{
				Config: ContainersConfig{
					Containers: []ContainerDef{
						{Name: "nginx", Image: "nginx:1.15"},
						{Name: "kuard", Image: "gcr.io/kuar-demo/kuard-amd64:1"},
					},
				},
				base: newBase(typeContainers, nil),
			},
		},
		{
			name:       "flexlayout",
			configFile: "config_flexlayout.json",
			objectType: "flexlayout",
			expected: &FlexLayout{
				Config: FlexLayoutConfig{
					Sections: []FlexLayoutSection{
						{
							{
								Width: WidthFull,
								View:  NewText("text"),
							},
						},
					},
				},
				base: newBase(typeFlexLayout, nil),
			},
		},
		{
			name:       "labels",
			configFile: "config_labels.json",
			objectType: "labels",
			expected: &Labels{
				Config: LabelsConfig{Labels: map[string]string{
					"foo": "bar",
				}},
				base: newBase(typeLabels, nil),
			},
		},
		{
			name:       "link",
			configFile: "config_link.json",
			objectType: "link",
			expected: &Link{
				Config: LinkConfig{
					Text: "text",
					Ref:  "ref",
				},
				base: newBase(typeLink, nil),
			},
		},
		{
			name:       "list",
			configFile: "config_list.json",
			objectType: "list",
			expected: &List{
				Config: ListConfig{
					Items: []Component{
						&Link{
							Config: LinkConfig{
								Text: "nginx-deployment",
								Ref:  "/overview/deployments/nginx-deployment",
							},
							base: newBase(typeLink, nil),
						},
						&Labels{
							Config: LabelsConfig{
								Labels: map[string]string{
									"app": "nginx",
								},
							},
							base: newBase(typeLabels, nil),
						},
					},
				},
				base: newBase(typeList, nil),
			},
		},
		{
			name:       "quadrant",
			configFile: "config_quadrant.json",
			objectType: "quadrant",
			expected: &Quadrant{
				Config: QuadrantConfig{
					NW: QuadrantValue{Label: "nw", Value: "1"},
					NE: QuadrantValue{Label: "ne", Value: "1"},
					SW: QuadrantValue{Label: "sw", Value: "1"},
					SE: QuadrantValue{Label: "se", Value: "1"},
				},
				base: newBase(typeQuadrant, nil),
			},
		},
		{
			name:       "resourceViewer",
			configFile: "config_resource_viewer.json",
			objectType: "resourceViewer",
			expected: &ResourceViewer{
				Config: ResourceViewerConfig{
					Edges: AdjList{
						"69e4ea11-2985-11e9-b356-42010a8000e5": []Edge{
							{
								Node: "bf4800b5b6602c4c78ba3b654af02b3b",
								Type: "explicit",
							},
						},
						"71c2b4eb-2949-11e9-b356-42010a8000e5": []Edge{
							{
								Node: "8682460a-29b5-11e9-b356-42010a8000e5",
								Type: "explicit",
							},
						},
						"8682460a-29b5-11e9-b356-42010a8000e5": []Edge{
							{
								Node: "bf4800b5b6602c4c78ba3b654af02b3b",
								Type: "explicit",
							},
						},
					},
					Nodes: Nodes{
						"69e4ea11-2985-11e9-b356-42010a8000e5": Node{
							Name:       "my-nginx",
							APIVersion: "v1",
							Kind:       "Service",
							Status:     "ok",
						},
						"71c2b4eb-2949-11e9-b356-42010a8000e5": Node{
							Name:       "nginx-deployment",
							APIVersion: "apps/v1",
							Kind:       "Deployment",
							Status:     "ok",
						},
						"8682460a-29b5-11e9-b356-42010a8000e5": Node{
							Name:       "nginx-deployment-56c74bb7cd",
							APIVersion: "extensions/v1beta1",
							Kind:       "ReplicaSet",
							Status:     "ok",
						},
						"bf4800b5b6602c4c78ba3b654af02b3b": Node{
							Name:       "nginx-deployment-56c74bb7cd pods",
							APIVersion: "v1",
							Kind:       "Pod",
							Status:     "ok",
						},
					},
				},
				base: newBase(typeResourceViewer, nil),
			},
		},
		{
			name:       "selectors",
			configFile: "config_selectors.json",
			objectType: "selectors",
			expected: &Selectors{
				Config: SelectorsConfig{
					Selectors: []Selector{
						&LabelSelector{
							Config: LabelSelectorConfig{
								Key:   "app",
								Value: "nginx",
							},
							base: newBase(typeLabelSelector, nil),
						},
						&ExpressionSelector{
							Config: ExpressionSelectorConfig{
								Key:      "environment",
								Operator: "In",
								Values:   []string{"production", "qa"},
							},
							base: newBase(typeExpressionSelector, nil),
						},
					},
				},
				base: newBase(typeSelectors, nil),
			},
		},
		{
			name:       "summary",
			configFile: "config_summary.json",
			objectType: "summary",
			expected: &Summary{
				Config: SummaryConfig{
					Sections: []SummarySection{
						{
							Header: "Containers",
							Content: &List{
								Config: ListConfig{
									Items: []Component{
										&Text{
											Config: TextConfig{
												Text: "nginx:latest",
											},
											base: newBase(typeText, TitleFromString("Image")),
										},
										&Text{
											Config: TextConfig{
												Text: "80/TCP",
											},
											base: newBase(typeText, TitleFromString("Port")),
										},
									},
								},
								base: newBase(typeList, TitleFromString("nginx")),
							},
						},
						{
							Header: "Empty Section",
							Content: &Text{
								Config: TextConfig{
									Text: "Nothing to see here",
								},
								base: newBase(typeText, nil),
							},
						},
					},
				},
				base: newBase(typeSummary, nil),
			},
		},
		{
			name:       "table",
			configFile: "config_table.json",
			objectType: "table",
			expected: &Table{
				Config: TableConfig{
					Columns: NewTableCols("Name", "Description"),
					Rows: []TableRow{
						{
							"Description": &Text{
								Config: TextConfig{
									Text: "The first row",
								},
								base: newBase(typeText, nil),
							},
							"Name": &Text{
								Config: TextConfig{
									Text: "First",
								},
								base: newBase(typeText, nil),
							},
						},
						{
							"Description": &Text{
								Config: TextConfig{
									Text: "The last row",
								},
								base: newBase(typeText, nil),
							},
							"Name": &Text{
								Config: TextConfig{
									Text: "Last",
								},
								base: newBase(typeText, nil),
							},
						},
					},
				},
				base: newBase(typeTable, nil),
			},
		},
		{
			name:       "text",
			configFile: "config_text.json",
			objectType: "text",
			expected: &Text{
				Config: TextConfig{Text: "text"},
				base:   newBase(typeText, nil),
			},
		},
		{
			name:       "timestamp",
			configFile: "config_timestamp.json",
			objectType: "timestamp",
			expected: &Timestamp{
				Config: TimestampConfig{Timestamp: 1548198349},
				base:   newBase(typeTimestamp, nil),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			configData, err := ioutil.ReadFile(filepath.Join("testdata", tc.configFile))
			require.NoError(t, err)

			to := TypedObject{
				Config:   json.RawMessage(configData),
				Metadata: Metadata{Type: tc.objectType},
			}

			got, err := unmarshal(to)
			require.NoError(t, err)

			AssertEqual(t, tc.expected, got)
		})
	}
}
