/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package printer_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/vmware-tanzu/octant/internal/printer"
	"github.com/vmware-tanzu/octant/internal/printer/fake"
	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func TestCustomResourceDefinitionSummary_BuildConfig(t *testing.T) {
	baseCRD := func(options ...testutil.CRDOption) *apiextv1.CustomResourceDefinition {
		c := testutil.CreateCRD("crd", append([]testutil.CRDOption{
			func(crd *apiextv1.CustomResourceDefinition) {
				crd.Spec.Group = "group"
				crd.Spec.Names = apiextv1.CustomResourceDefinitionNames{
					Plural:     "plural",
					Singular:   "singular",
					ShortNames: []string{"short1", "short2"},
					Kind:       "kind",
					ListKind:   "list kind",
				}

			},
		}, options...)...)

		return c
	}

	type args struct {
		object     func(ctrl *gomock.Controller) printer.ObjectInterface
		crdOptions []testutil.CRDOption
		options    []printer.CustomResourceDefinitionSummaryOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "in general",
			args: args{
				object: func(ctrl *gomock.Controller) printer.ObjectInterface {
					o := fake.NewMockObjectInterface(ctrl)

					sections := component.SummarySections{}
					sections.AddText("Conversion Strategy", string(apiextv1.NoneConverter))
					sections.AddText("Group", "group")
					sections.AddText("Kind", "kind")
					sections.AddText("List Kind", "list kind")
					sections.AddText("Plural", "plural")
					sections.AddText("Singular", "singular")
					sections.AddText("Short Names", "short1, short2")
					sections.AddText("Categories", "cat1, cat2")
					summary := component.NewSummary("Configuration", sections...)

					o.EXPECT().
						RegisterConfig(summary)

					return o
				},
				crdOptions: []testutil.CRDOption{
					func(crd *apiextv1.CustomResourceDefinition) {
						crd.Spec.Conversion = &apiextv1.CustomResourceConversion{
							Strategy: apiextv1.NoneConverter,
						}
						crd.Spec.Names.Categories = []string{"cat1", "cat2"}
					},
				},
			},
		},
		{
			name: "no conversion strategy",
			args: args{
				object: func(ctrl *gomock.Controller) printer.ObjectInterface {
					o := fake.NewMockObjectInterface(ctrl)

					sections := component.SummarySections{}
					sections.AddText("Group", "group")
					sections.AddText("Kind", "kind")
					sections.AddText("List Kind", "list kind")
					sections.AddText("Plural", "plural")
					sections.AddText("Singular", "singular")
					sections.AddText("Short Names", "short1, short2")
					sections.AddText("Categories", "cat1, cat2")
					summary := component.NewSummary("Configuration", sections...)

					o.EXPECT().
						RegisterConfig(summary)

					return o
				},
				crdOptions: []testutil.CRDOption{
					func(crd *apiextv1.CustomResourceDefinition) {
						crd.Spec.Names.Categories = []string{"cat1", "cat2"}
					},
				},
			},
		},
		{
			name: "no categories",
			args: args{
				object: func(ctrl *gomock.Controller) printer.ObjectInterface {
					o := fake.NewMockObjectInterface(ctrl)

					sections := component.SummarySections{}
					sections.AddText("Conversion Strategy", string(apiextv1.NoneConverter))
					sections.AddText("Group", "group")
					sections.AddText("Kind", "kind")
					sections.AddText("List Kind", "list kind")
					sections.AddText("Plural", "plural")
					sections.AddText("Singular", "singular")
					sections.AddText("Short Names", "short1, short2")
					summary := component.NewSummary("Configuration", sections...)

					o.EXPECT().
						RegisterConfig(summary)

					return o
				},
				crdOptions: []testutil.CRDOption{
					func(crd *apiextv1.CustomResourceDefinition) {
						crd.Spec.Conversion = &apiextv1.CustomResourceConversion{
							Strategy: apiextv1.NoneConverter,
						}
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			crd := baseCRD(test.args.crdOptions...)

			h := printer.NewCustomResourceDefinitionSummary(crd, test.args.object(ctrl), test.args.options...)

			err := h.BuildConfig()
			if test.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestCustomResourceDefinitionSummary_BuildItems(t *testing.T) {
	type args struct {
		object         func(ctrl *gomock.Controller) printer.ObjectInterface
		printerOptions printer.Options
		options        []printer.CustomResourceDefinitionSummaryOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "add supplied items",
			args: args{
				object: func(ctrl *gomock.Controller) printer.ObjectInterface {
					o := fake.NewMockObjectInterface(ctrl)

					o.EXPECT().
						RegisterItems(printer.ItemDescriptor{
							Component: component.NewText("output"),
							Width:     component.WidthFull,
						})

					return o
				},
				options: []printer.CustomResourceDefinitionSummaryOption{
					printer.CustomResourceDefinitionSummaryItems(
						func(_ *apiextv1.CustomResourceDefinition, _ printer.Options) (component.Component, error) {
							return component.NewText("output"), nil
						},
					),
				},
			},
		},
		{
			name: "item component generation failed",
			args: args{
				object: func(ctrl *gomock.Controller) printer.ObjectInterface {
					o := fake.NewMockObjectInterface(ctrl)
					return o
				},
				options: []printer.CustomResourceDefinitionSummaryOption{
					printer.CustomResourceDefinitionSummaryItems(
						func(_ *apiextv1.CustomResourceDefinition, _ printer.Options) (component.Component, error) {
							return nil, fmt.Errorf("failed")
						},
					),
				},
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			crd := testutil.CreateCRD("crd")
			h := printer.NewCustomResourceDefinitionSummary(crd, test.args.object(ctrl), test.args.options...)

			err := h.BuildItems(test.args.printerOptions)
			if test.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

		})
	}
}

func TestCreateCRDConditionsTable(t *testing.T) {
	now := time.Now()

	type args struct {
		crd *apiextv1.CustomResourceDefinition
	}
	tests := []struct {
		name    string
		args    args
		want    *component.Table
		wantErr bool
	}{
		{
			name: "in general",
			args: args{
				crd: testutil.CreateCRD("crd", func(crd *apiextv1.CustomResourceDefinition) {
					crd.Status.Conditions = []apiextv1.CustomResourceDefinitionCondition{
						{
							Type:   apiextv1.Established,
							Status: apiextv1.ConditionTrue,
							LastTransitionTime: metav1.Time{
								Time: now,
							},
							Reason:  "reason",
							Message: "message",
						},
					}
				}),
			},
			want: component.NewTableWithRows("Conditions", "", printer.CRDConditionsColumns,
				[]component.TableRow{
					{
						"Type":                 component.NewText("Established"),
						"Status":               component.NewText("True"),
						"Last Transition Time": component.NewTimestamp(now),
						"Message":              component.NewText("message"),
						"Reason":               component.NewText("reason"),
					},
				}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := printer.CreateCRDConditionsTable(tt.args.crd)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			testutil.AssertJSONEqual(t, tt.want, got)
		})
	}
}
