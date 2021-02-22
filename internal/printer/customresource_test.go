/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware-tanzu/octant/internal/link"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func TestCustomResourceLister_List(t *testing.T) {
	now := time.Now()
	labels := map[string]string{"foo": "bar"}

	type args struct {
		crd     *unstructured.Unstructured
		list    *unstructured.UnstructuredList
		version string
		link    link.Interface
	}

	tests := []struct {
		name string
		args func(t *testing.T, ctrl *gomock.Controller) args

		wantErr bool
		want    component.Component
	}{
		{
			name: "in general",
			args: func(t *testing.T, ctrl *gomock.Controller) args {
				tpo := newTestPrinterOptions(ctrl)
				crd := testutil.LoadUnstructuredFromFile(t, "crd.yaml")
				resource := testutil.LoadUnstructuredFromFile(t, "crd-resource.yaml")
				resource.SetCreationTimestamp(metav1.Time{Time: now})
				tpo.PathForObject(resource, resource.GetName(), "/my-crontab")
				resource.SetLabels(labels)
				list := testutil.ToUnstructuredList(t, resource)
				return args{
					crd:     crd,
					list:    list,
					version: "v1",
					link:    tpo.link,
				}
			},

			wantErr: false,
			want: component.NewTableWithRows(
				"crontabs.stable.example.com/v1", "We could not find any crontabs.stable.example.com/v1!",
				component.NewTableCols("Name", "Labels", "Age"),
				[]component.TableRow{
					{
						"Name":   component.NewLink("", "my-crontab", "/my-crontab"),
						"Age":    component.NewTimestamp(now),
						"Labels": component.NewLabels(labels),
					},
				}),
		},
		{
			name: "custom columns",
			args: func(t *testing.T, ctrl *gomock.Controller) args {
				tpo := newTestPrinterOptions(ctrl)
				crd := testutil.LoadUnstructuredFromFile(t, "crd-additional-columns.yaml")
				resource := testutil.LoadUnstructuredFromFile(t, "crd-resource.yaml")
				resource.SetCreationTimestamp(metav1.Time{Time: now})
				tpo.PathForObject(resource, resource.GetName(), "/my-crontab")
				resource.SetLabels(labels)
				list := testutil.ToUnstructuredList(t, resource)
				return args{
					crd:     crd,
					list:    list,
					version: "v1",
					link:    tpo.link,
				}
			},

			wantErr: false,
			want: component.NewTableWithRows(
				"crontabs.stable.example.com/v1", "We could not find any crontabs.stable.example.com/v1!",
				component.NewTableCols("Name", "Labels", "Spec", "Replicas", "Errors", "Resource Name", "Age"),
				[]component.TableRow{
					{
						"Name":          component.NewLink("", "my-crontab", "/my-crontab"),
						"Age":           component.NewTimestamp(now),
						"Labels":        component.NewLabels(labels),
						"Replicas":      component.NewText("1"),
						"Spec":          component.NewText("* * * * */5"),
						"Errors":        component.NewText("1"),
						"Resource Name": component.NewText("my-crontab"),
					},
				}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			lister := NewCustomResourceLister()
			a := test.args(t, ctrl)
			got, err := lister.List(a.crd, a.list, a.version, a.link)

			testutil.RequireErrorOrNot(t, test.wantErr, err, func() {
				component.AssertEqual(t, test.want, got)
			})
		})
	}
}

func TestCustomResourceHandler(t *testing.T) {

}

func Test_printCustomResourceConfig(t *testing.T) {
	cases := []struct {
		name     string
		crd      string
		cr       string
		expected component.Component
		wantErr  bool
	}{
		{
			name: "with additional columns",
			crd:  "crd-additional-columns.yaml",
			cr:   "crd-resource.yaml",
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Spec",
					Content: component.NewText("* * * * */5"),
				},
				{
					Header:  "Replicas",
					Content: component.NewText("1"),
				},
			}...),
		},
		{
			name:     "in general",
			crd:      "crd.yaml",
			cr:       "crd-resource.yaml",
			expected: component.NewSummary("Configuration"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			crd := testutil.LoadUnstructuredFromFile(t, tc.crd)
			resource := testutil.LoadUnstructuredFromFile(t, tc.cr)

			now := time.Now()
			resource.SetCreationTimestamp(metav1.Time{Time: now})

			labels := map[string]string{"foo": "bar"}
			resource.SetLabels(labels)

			got, err := printCustomResourceConfig(crd, resource)
			testutil.RequireErrorOrNot(t, tc.wantErr, err)

			assert.Equal(t, tc.expected, got)
		})
	}
}

func Test_printCustomResourceStatus(t *testing.T) {
	cases := []struct {
		name     string
		crd      string
		cr       string
		expected component.Component
		wantErr  bool
	}{
		{
			name: "with additional columns",
			crd:  "crd-additional-columns.yaml",
			cr:   "crd-resource.yaml",
			expected: component.NewSummary("Status", []component.SummarySection{
				{
					Header:  "Errors",
					Content: component.NewText("1"),
				},
			}...),
		},
		{
			name:     "in general",
			crd:      "crd.yaml",
			cr:       "crd-resource.yaml",
			expected: component.NewSummary("Status"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			crd := testutil.LoadUnstructuredFromFile(t, tc.crd)
			resource := testutil.LoadUnstructuredFromFile(t, tc.cr)

			now := time.Now()
			resource.SetCreationTimestamp(metav1.Time{Time: now})

			labels := map[string]string{"foo": "bar"}
			resource.SetLabels(labels)

			got, err := printCustomResourceStatus(crd, resource)
			testutil.RequireErrorOrNot(t, tc.wantErr, err)

			assert.Equal(t, tc.expected, got)
		})
	}
}

func Test_printCustomColumn(t *testing.T) {
	cases := []struct {
		name       string
		objectPath string
		jsonPath   string
		expected   string
		wantErr    bool
	}{
		{
			name:       "simple",
			objectPath: "certificate.yaml",
			jsonPath:   ".metadata.name",
			expected:   "kubecon-panel",
		},
		{
			name:       "with a filter",
			objectPath: "certificate.yaml",
			jsonPath:   ".status.conditions[?(@.type==\"Ready\")].status",
			expected:   "True",
		},
		{
			name:       "invalid json path",
			objectPath: "certificate.yaml",
			jsonPath:   ".status.conditions[?(@.type==\"Ready\"].status",
			wantErr:    true,
		},
		{
			name:       "execute error: not found",
			objectPath: "certificate.yaml",
			jsonPath:   ".missing",
			expected:   "<not found>",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resource := testutil.LoadUnstructuredFromFile(t, tc.objectPath)

			def := octant.CustomResourceDefinitionPrinterColumn{
				Name:     "name",
				JSONPath: tc.jsonPath,
			}

			got, err := printCustomColumn(resource.Object, def)
			testutil.RequireErrorOrNot(t, tc.wantErr, err)

			assert.Equal(t, tc.expected, got)
		})
	}

}
