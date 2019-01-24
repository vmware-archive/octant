package component

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_unmarshal(t *testing.T) {
	cases := []struct {
		name       string
		config     json.RawMessage
		configFile string
		objectType string
		expected   interface{}
	}{
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
				Metadata: Metadata{Type: "containers"},
			},
		},
		{
			name:       "grid",
			configFile: "config_grid.json",
			objectType: "grid",
			expected: &Grid{
				Config: GridConfig{
					Panels: []Panel{
						{
							Config: PanelConfig{
								Content: &Text{
									Config:   TextConfig{Text: "Panel contents"},
									Metadata: Metadata{Type: "text"},
								},
								Position: PanelPosition{
									X: 0, Y: 0, W: 12, H: 7,
								},
							},
							Metadata: Metadata{
								Type: "panel",
							},
						},
					},
				},
				Metadata: Metadata{Type: "grid"},
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
				Metadata: Metadata{Type: "labels"},
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
				Metadata: Metadata{Type: "link"},
			},
		},
		{
			name:       "list",
			configFile: "config_list.json",
			objectType: "list",
			expected: &List{
				Config: ListConfig{
					Items: []ViewComponent{
						&Link{
							Config: LinkConfig{
								Text: "nginx-deployment",
								Ref:  "/overview/deployments/nginx-deployment",
							},
							Metadata: Metadata{
								Type: "link",
							},
						},
						&Labels{
							Config: LabelsConfig{
								Labels: map[string]string{
									"app": "nginx",
								},
							},
							Metadata: Metadata{
								Type: "labels",
							},
						},
					},
				},
				Metadata: Metadata{Type: "list"},
			},
		},
		{
			name:       "panel",
			configFile: "config_panel.json",
			objectType: "panel",
			expected: &Panel{
				Config: PanelConfig{
					Content: &Text{
						Config:   TextConfig{Text: "Panel contents"},
						Metadata: Metadata{Type: "text"},
					},
					Position: PanelPosition{
						X: 1, Y: 2, W: 3, H: 4,
					},
				},
				Metadata: Metadata{
					Type: "panel",
				},
			},
		},
		{
			name:       "quadrant",
			configFile: "config_quadrant.json",
			objectType: "quadrant",
			expected: &Quadrant{
				Config: QuadrantConfig{
					NW: &Text{
						Config:   TextConfig{Text: "nw"},
						Metadata: Metadata{Type: "text"},
					},
					NE: &Text{
						Config:   TextConfig{Text: "ne"},
						Metadata: Metadata{Type: "text"},
					},
					SW: &Text{
						Config:   TextConfig{Text: "sw"},
						Metadata: Metadata{Type: "text"},
					},
					SE: &Text{
						Config:   TextConfig{Text: "se"},
						Metadata: Metadata{Type: "text"},
					},
				},
				Metadata: Metadata{Type: "quadrant"},
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
							Metadata: Metadata{
								Type: "labelSelector",
							},
						},
						&ExpressionSelector{
							Config: ExpressionSelectorConfig{
								Key:      "environment",
								Operator: "In",
								Values:   []string{"production", "qa"},
							},
							Metadata: Metadata{
								Type: "expressionSelector",
							},
						},
					},
				},
				Metadata: Metadata{Type: "selectors"},
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
									Items: []ViewComponent{
										&Text{
											Config: TextConfig{
												Text: "nginx:latest",
											},
											Metadata: Metadata{
												Type:  "text",
												Title: "Image",
											},
										},
										&Text{
											Config: TextConfig{
												Text: "80/TCP",
											},
											Metadata: Metadata{
												Type:  "text",
												Title: "Port",
											},
										},
									},
								},
								Metadata: Metadata{
									Type:  "list",
									Title: "nginx",
								},
							},
						},
						{
							Header: "Empty Section",
							Content: &Text{
								Config: TextConfig{
									Text: "Nothing to see here",
								},
								Metadata: Metadata{
									Type: "text",
								},
							},
						},
					},
				},
				Metadata: Metadata{Type: "summary"},
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
								Metadata: Metadata{
									Type: "text",
								},
							},
							"Name": &Text{
								Config: TextConfig{
									Text: "First",
								},
								Metadata: Metadata{
									Type: "text",
								},
							},
						},
						{
							"Description": &Text{
								Config: TextConfig{
									Text: "The last row",
								},
								Metadata: Metadata{
									Type: "text",
								},
							},
							"Name": &Text{
								Config: TextConfig{
									Text: "Last",
								},
								Metadata: Metadata{
									Type: "text",
								},
							},
						},
					},
				},
				Metadata: Metadata{Type: "table"},
			},
		},
		{
			name:       "text",
			configFile: "config_text.json",
			objectType: "text",
			expected: &Text{
				Config:   TextConfig{Text: "text"},
				Metadata: Metadata{Type: "text"},
			},
		},
		{
			name:       "timestamp",
			configFile: "config_timestamp.json",
			objectType: "timestamp",
			expected: &Timestamp{
				Config:   TimestampConfig{Timestamp: 1548198349},
				Metadata: Metadata{Type: "timestamp"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			configData, err := ioutil.ReadFile(filepath.Join("testdata", tc.configFile))
			require.NoError(t, err)

			to := typedObject{
				Config:   json.RawMessage(configData),
				Metadata: Metadata{Type: tc.objectType},
			}

			got, err := unmarshal(to)
			require.NoError(t, err)

			assert.Equal(t, tc.expected, got)
		})
	}

}
