/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/vmware-tanzu/octant/pkg/action"

	jsoniter "github.com/json-iterator/go"
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
			name:       "accordion",
			configFile: "config_accordion.json",
			objectType: TypeAccordion,
			expected: &Accordion{
				Base: newBase(TypeAccordion, nil),
				Config: AccordionConfig{
					Rows: []AccordionRow{
						{
							Title:   "Accordion title",
							Content: NewText("accordion content"),
						},
					},
					AllowMultipleExpanded: true,
				},
			},
		},
		{
			name:       "annotations",
			configFile: "config_annotations.json",
			objectType: TypeAnnotations,
			expected: &Annotations{
				Base: newBase(TypeAnnotations, nil),
				Config: AnnotationsConfig{
					map[string]string{
						"foo": "bar",
					},
				},
			},
		},
		{
			name:       "cardList",
			configFile: "config_card_list.json",
			objectType: TypeCardList,
			expected: &CardList{
				Config: CardListConfig{
					Cards: []Card{
						{
							Base: newBase(TypeCard, TitleFromString("card title")),
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
									Status:  AlertStatusWarning,
									Type:    AlertTypeBanner,
									Message: "warning",
								},
							},
						},
					},
				},
			},
		},
		{
			name:       "code",
			configFile: "config_code.json",
			objectType: "codeBlock",
			expected: &Code{
				Config: CodeConfig{
					Code: "echo HELLO_WORLD",
				},
				Base: newBase(TypeCode, nil),
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
				Base: newBase(TypeContainers, nil),
			},
		},
		{
			name:       "donutchart",
			configFile: "config_donutchart.json",
			objectType: "donutChart",
			expected: &DonutChart{
				Config: DonutChartConfig{
					Segments: []DonutSegment{
						{
							Count:  1,
							Status: "ok",
						},
					},
					Labels: DonutChartLabels{
						Plural:   "tests",
						Singular: "test",
					},
					Size: DonutChartSizeSmall,
				},
				Base: newBase(TypeDonutChart, nil),
			},
		},
		{
			name:       "dropdown",
			configFile: "config_dropdown.json",
			objectType: "dropdown",
			expected: &Dropdown{
				Config: DropdownConfig{
					DropdownType:     DropdownButton,
					DropdownPosition: BottomLeft,
					Action:           "action.octant.dev/dropdownTest",
					UseSelection:     true,
					Items: []DropdownItemConfig{{
						Name:  "first",
						Type:  "text",
						Label: "First Item",
					},
						{
							Name:  "second",
							Type:  "link",
							Label: "Second Item",
							Url:   "/items/second",
						},
					},
				},
				Base: newBase(TypeDropdown, nil),
			},
		},
		{
			name:       "editor",
			configFile: "config_editor.json",
			objectType: "editor",
			expected: &Editor{
				Config: EditorConfig{
					Value:    "code",
					ReadOnly: true,
				},
				Base: newBase(TypeEditor, nil),
			},
		},
		{
			name:       "error",
			configFile: "config_error.json",
			objectType: "error",
			expected: &Error{
				Config: ErrorConfig{
					Data: "error test",
				},
				Base: newBase(TypeError, nil),
			},
		},
		{
			name:       "expandable row detail",
			configFile: "config_expandable_row_detail.json",
			objectType: "expandableRowDetail",
			expected: &ExpandableRowDetail{
				Config: ExpandableRowDetailConfig{
					Body:    []Component{NewText("test")},
					Replace: true,
				},
				Base: newBase(TypeExpandableRowDetail, nil),
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
					ButtonGroup: &ButtonGroup{
						Base: Base{},
						Config: ButtonGroupConfig{
							Buttons: []Button{{
								Base: newBase(TypeButton, nil),
								Config: ButtonConfig{
									Name:    "test",
									Payload: action.Payload{"foo": "bar"},
								},
							}},
						},
					},
				},
				Base: newBase(TypeFlexLayout, nil),
			},
		},
		{
			name:       "form field text",
			configFile: "config_form_text.json",
			objectType: TypeFormField,
			expected: &FormField{
				Base: newBase(TypeFormField, nil),
				Config: FormFieldConfig{
					Type:        FieldTypeText,
					Label:       "label",
					Name:        "name",
					Value:       "value",
					Placeholder: "placeholder",
					Error:       "error message",
					Validators: map[FormValidator]interface{}{
						FormValidatorMaxLength: 100,
					},
				},
			},
		},
		{
			name:       "form field checkbox",
			configFile: "config_form_checkbox.json",
			objectType: TypeFormField,
			expected: &FormField{
				Base: newBase(TypeFormField, nil),
				Config: FormFieldConfig{
					Type:  FieldTypeCheckBox,
					Label: "label",
					Name:  "name",
					Configuration: &FormFieldOptions{
						Choices: []InputChoice{
							{
								Label:   "dog",
								Value:   "dog",
								Checked: true,
							},
							{
								Label:   "cat",
								Value:   "cat",
								Checked: false,
							},
						},
					},
				},
			},
		},
		{
			name:       "grid actions",
			configFile: "config_grid_actions.json",
			objectType: TypeGridActions,
			expected: &GridActions{
				Config: GridActionsConfig{
					Actions: []GridAction{
						{
							Name:       "name",
							ActionPath: "/path",
							Payload:    action.Payload{"foo": "bar"},
						},
					},
				},
			},
		},
		{
			name:       "json editor",
			configFile: "config_json_editor.json",
			objectType: TypeJSONEditor,
			expected: &JSONEditor{
				Config: JSONEditorConfig{
					Mode:    ViewMode,
					Content: "{ \"hello\": 123 }",
				},
				Base: newBase(TypeJSONEditor, nil),
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
				Base: newBase(TypeLabels, nil),
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
				Base: newBase(TypeLink, nil),
			},
		},
		{
			name:       "link with status",
			configFile: "config_link_status.json",
			objectType: "link",
			expected: &Link{
				Config: LinkConfig{
					Text:   "text",
					Ref:    "ref",
					Status: TextStatusOK,
					StatusDetail: &Text{
						Config: TextConfig{
							Text: "Ready",
						},
						Base: newBase(TypeText, nil),
					},
				},
				Base: newBase(TypeLink, nil),
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
							Base: newBase(TypeLink, nil),
						},
						&Labels{
							Config: LabelsConfig{
								Labels: map[string]string{
									"app": "nginx",
								},
							},
							Base: newBase(TypeLabels, nil),
						},
					},
				},
				Base: newBase(TypeList, nil),
			},
		},
		{
			name:       "logs",
			configFile: "config_logs.json",
			objectType: TypeLogs,
			expected: &Logs{
				Config: LogsConfig{
					Namespace:  "test",
					Name:       "nginx-deployment-7cb4fc6c56-29pbw",
					Containers: []string{"nginx"},
					Durations:  []Since{{Label: "5 minutes", Seconds: 300}},
				},
			},
		},
		{
			name:       "mfcomponent",
			configFile: "config_mfcomponent.json",
			objectType: TypeMFComponent,
			expected: &MfComponent{
				Base: newBase(TypeMFComponent, nil),
				Config: MfComponentConfig{
					Name:          "template-mf-angular",
					RemoteEntry:   "http://localhost:3000/remoteEntry.js",
					RemoteName:    "templatemfangular",
					ExposedModule: "./web-components",
					ElementName:   "template-mf-angular",
				},
			},
		},
		{
			name:       "modal",
			configFile: "config_modal.json",
			objectType: TypeModal,
			expected: &Modal{
				Config: ModalConfig{
					Body:      NewText("test"),
					Opened:    true,
					ModalSize: ModalSizeSmall,
				},
				Base: newBase(TypeModal, nil),
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
				Base: newBase(TypeQuadrant, nil),
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
							Details: []Component{
								NewText("my-nginx"),
							},
							Path: NewLink("", "my-nginx", "/overview/namespace/default/discovery-and-load-balancing/services/my-nginx"),
						},
						"71c2b4eb-2949-11e9-b356-42010a8000e5": Node{
							Name:       "nginx-deployment",
							APIVersion: "apps/v1",
							Kind:       "Deployment",
							Status:     "ok",
							Path:       NewLink("", "nginx-deployment", "/overview/namespace/default/workloads/deployments/nginx-deployment"),
						},
						"8682460a-29b5-11e9-b356-42010a8000e5": Node{
							Name:       "nginx-deployment-56c74bb7cd",
							APIVersion: "extensions/v1beta1",
							Kind:       "ReplicaSet",
							Status:     "ok",
							Path:       NewLink("", "nginx-deployment-56c74bb7cd", "/overview/namespace/default/workloads/replica-sets/nginx-deployment-56c74bb7cd"),
						},
						"bf4800b5b6602c4c78ba3b654af02b3b": Node{
							Name:       "nginx-deployment-56c74bb7cd pods",
							APIVersion: "v1",
							Kind:       "Pod",
							Status:     "ok",
							Path:       NewLink("", "nginx-deployment-56c74bb7cd pods", "/overview/namespace/default/workloads/pods/nginx-deployment-56c74bb7cd pods"),
						},
					},
				},
				Base: newBase(TypeResourceViewer, nil),
			},
		},
		{
			name:       "selectFile",
			configFile: "config_select_file.json",
			objectType: "selectFile",
			expected: &SelectFile{
				Config: SelectFileConfig{
					Label:         "Open File",
					Multiple:      false,
					Status:        "success",
					StatusMessage: "Success message",
					Layout:        "compact",
					Action:        "action.octant.dev/SelectFileAction",
				},
				Base: newBase(TypeSelectFile, nil),
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
							Base: newBase(TypeLabelSelector, nil),
						},
						&ExpressionSelector{
							Config: ExpressionSelectorConfig{
								Key:      "environment",
								Operator: "In",
								Values:   []string{"production", "qa"},
							},
							Base: newBase(TypeExpressionSelector, nil),
						},
					},
				},
				Base: newBase(TypeSelectors, nil),
			},
		},
		{
			name:       "singleStat",
			configFile: "config_single_stat.json",
			objectType: "singleStat",
			expected: &SingleStat{
				Config: SingleStatConfig{
					Title: "testing",
					Value: SingleStatValue{
						Text:  "30m",
						Color: "#60b515",
					},
				},
				Base: newBase(TypeSingleStat, nil),
			},
		},
		{
			name:       "stepper",
			configFile: "config_stepper.json",
			objectType: "stepper",
			expected: &Stepper{
				Config: StepperConfig{
					Action: "action.octant.dev/stepperTest",
					Steps: []StepConfig{{
						Name:        "Step 1",
						Title:       "First Step",
						Description: "Setup step",
						Form:        Form{Fields: []FormField{NewFormFieldText("test", "test", "test")}},
					}, {
						Name:        "Step 2",
						Title:       "Second Step",
						Description: "Confirmation step",
						Form:        Form{},
					}},
				},
				Base: newBase(TypeStepper, nil),
			},
		},
		{
			name:       "tabsView",
			configFile: "config_tabs.json",
			objectType: "tabsView",
			expected: &TabsView{
				Config: TabsViewConfig{
					Tabs: []SingleTab{{
						Name: "Tab 1",
						Contents: FlexLayout{
							Base: newBase(TypeFlexLayout, nil),
							Config: FlexLayoutConfig{
								Sections: []FlexLayoutSection{{{
									Width: WidthFull,
									View:  NewButton("test", action.Payload{"foo": "bar"}),
								}}},
								ButtonGroup: nil,
							},
						},
					}},
				},
				Base: newBase(TypeTabsView, nil),
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
											Base: newBase(TypeText, TitleFromString("Image")),
										},
										&Text{
											Config: TextConfig{
												Text: "80/TCP",
											},
											Base: newBase(TypeText, TitleFromString("Port")),
										},
									},
								},
								Base: newBase(TypeList, TitleFromString("nginx")),
							},
						},
						{
							Header: "Empty Section",
							Content: &Text{
								Config: TextConfig{
									Text: "Nothing to see here",
								},
								Base: newBase(TypeText, nil),
							},
						},
					},
				},
				Base: newBase(TypeSummary, nil),
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
								Base: newBase(TypeText, nil),
							},
							"Name": &Text{
								Config: TextConfig{
									Text: "First",
								},
								Base: newBase(TypeText, nil),
							},
						},
						{
							"Description": &Text{
								Config: TextConfig{
									Text: "The last row",
								},
								Base: newBase(TypeText, nil),
							},
							"Name": &Text{
								Config: TextConfig{
									Text: "Last",
								},
								Base: newBase(TypeText, nil),
							},
						},
					},
				},
				Base: newBase(TypeTable, nil),
			},
		},
		{
			name:       "table with button group",
			configFile: "config_table_buttongroup.json",
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
								Base: newBase(TypeText, nil),
							},
							"Name": &Text{
								Config: TextConfig{
									Text: "First",
								},
								Base: newBase(TypeText, nil),
							},
						},
						{
							"Description": &Text{
								Config: TextConfig{
									Text: "The last row",
								},
								Base: newBase(TypeText, nil),
							},
							"Name": &Text{
								Config: TextConfig{
									Text: "Last",
								},
								Base: newBase(TypeText, nil),
							},
						},
					},
					ButtonGroup: &ButtonGroup{
						Config: ButtonGroupConfig{
							Buttons: []Button{{
								Base: Base{},
								Config: ButtonConfig{
									Name: "Create",
									Payload: action.Payload{
										"action": "action.local/create",
										"prop":   "value",
									},
								},
							}},
						},
					},
				},
				Base: newBase(TypeTable, nil),
			},
		},
		{
			name:       "text",
			configFile: "config_text.json",
			objectType: "text",
			expected: &Text{
				Config: TextConfig{Text: "text"},
				Base:   newBase(TypeText, nil),
			},
		},
		{
			name:       "timestamp",
			configFile: "config_timestamp.json",
			objectType: "timestamp",
			expected: &Timestamp{
				Config: TimestampConfig{Timestamp: 1548198349},
				Base:   newBase(TypeTimestamp, nil),
			},
		},
		{
			name:       "timeline",
			configFile: "config_timeline.json",
			objectType: "timeline",
			expected: &Timeline{
				Config: TimelineConfig{
					Steps: []TimelineStep{
						{
							State:       TimelineStepCurrent,
							Header:      "Header",
							Title:       "Title",
							Description: "Description",
						},
					},
					Vertical: true,
				},
				Base: newBase(TypeTimeline, nil),
			},
		},
		{
			name:       "button",
			configFile: "config_button.json",
			objectType: "button",
			expected: &Button{
				Config: ButtonConfig{
					Name: "test",
					Payload: action.Payload{
						"foo": "bar",
					}},
				Base: newBase(TypeButton, nil),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			configData, err := ioutil.ReadFile(filepath.Join("testdata", tc.configFile))
			require.NoError(t, err)

			to := TypedObject{
				Config:   jsoniter.RawMessage(configData),
				Metadata: Metadata{Type: tc.objectType},
			}

			got, err := unmarshal(to)
			require.NoError(t, err)
			AssertEqual(t, tc.expected, got)
		})
	}
}
