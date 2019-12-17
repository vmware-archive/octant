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
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_CustomResourceListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)

	crd := testutil.LoadUnstructuredFromFile(t, "crd.yaml")
	resource := testutil.LoadUnstructuredFromFile(t, "crd-resource.yaml")

	now := time.Now()
	resource.SetCreationTimestamp(metav1.Time{Time: now})

	tpo.PathForObject(resource, resource.GetName(), "/my-crontab")

	labels := map[string]string{"foo": "bar"}
	resource.SetLabels(labels)

	list := testutil.ToUnstructuredList(t, resource)
	got, err := CustomResourceListHandler(crd, list, "v1", tpo.link)
	require.NoError(t, err)

	expected := component.NewTableWithRows(
		"crontabs.stable.example.com/v1", "We couldn't find any custom resources!",
		component.NewTableCols("Name", "Labels", "Age"),
		[]component.TableRow{
			{
				"Name":   component.NewLink("", resource.GetName(), "/my-crontab"),
				"Age":    component.NewTimestamp(now),
				"Labels": component.NewLabels(labels),
			},
		})

	component.AssertEqual(t, expected, got)
}

func Test_CustomResourceListHandler_custom_columns(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)

	crd := testutil.LoadUnstructuredFromFile(t, "crd-additional-columns.yaml")
	resource := testutil.LoadUnstructuredFromFile(t, "crd-resource.yaml")

	now := time.Now()
	resource.SetCreationTimestamp(metav1.Time{Time: now})

	tpo.PathForObject(resource, resource.GetName(), "/my-crontab")

	labels := map[string]string{"foo": "bar"}
	resource.SetLabels(labels)

	list := testutil.ToUnstructuredList(t, resource)

	got, err := CustomResourceListHandler(crd, list, "v1", tpo.link)
	require.NoError(t, err)

	expected := component.NewTableWithRows(
		"crontabs.stable.example.com/v1", "We couldn't find any custom resources!",
		component.NewTableCols("Name", "Labels", "Spec", "Replicas", "Errors", "Resource Age", "Age"),
		[]component.TableRow{
			{
				"Name":         component.NewLink("", resource.GetName(), "/my-crontab"),
				"Age":          component.NewTimestamp(now),
				"Labels":       component.NewLabels(labels),
				"Replicas":     component.NewText("1"),
				"Spec":         component.NewText("* * * * */5"),
				"Errors":       component.NewText("1"),
				"Resource Age": component.NewText(resource.GetCreationTimestamp().UTC().Format(time.RFC3339)),
			},
		})

	component.AssertEqual(t, expected, got)
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
