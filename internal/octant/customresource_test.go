/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package octant_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/internal/testutil"
)

func TestNewCustomResourceDefinition(t *testing.T) {
	tests := []struct {
		name    string
		object  *unstructured.Unstructured
		wantErr bool
	}{
		{
			name:   "with an object",
			object: testutil.LoadUnstructuredFromFile(t, "crd-v1.yaml"),
		},
		{
			name:    "without object",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := octant.NewCustomResourceDefinitionTool(tt.object)
			testutil.RequireErrorOrNot(t, tt.wantErr, err)
		})
	}
}

func TestCustomResourceDefinition_GroupKind(t *testing.T) {
	type ctorArgs struct {
		object *unstructured.Unstructured
	}
	tests := []struct {
		name     string
		ctorArgs ctorArgs
		wantErr  bool
		want     schema.GroupKind
	}{
		{
			name: "in general",
			ctorArgs: ctorArgs{
				object: &unstructured.Unstructured{
					Object: map[string]interface{}{
						"spec": map[string]interface{}{
							"group": "group",
							"names": map[string]interface{}{
								"kind": "Kind",
							},
						},
					},
				},
			},
			want: schema.GroupKind{Group: "group", Kind: "Kind"},
		},
		{
			ctorArgs: ctorArgs{
				object: &unstructured.Unstructured{
					Object: map[string]interface{}{
						"spec": map[string]interface{}{
							"group": 7,
							"names": map[string]interface{}{
								"kind": "Kind",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "in general",
			ctorArgs: ctorArgs{
				object: &unstructured.Unstructured{
					Object: map[string]interface{}{
						"spec": map[string]interface{}{
							"group": "group",
							"names": map[string]interface{}{
								"kind": 7,
							},
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			crdTool, err := octant.NewCustomResourceDefinitionTool(test.ctorArgs.object)
			require.NoError(t, err)

			got, err := crdTool.GroupKind()
			testutil.RequireErrorOrNot(t, test.wantErr, err, func() {
				require.Equal(t, test.want, got)
			})
		})
	}
}

func TestCustomResourceDefinition_Versions(t *testing.T) {
	tests := []struct {
		name    string
		object  *unstructured.Unstructured
		want    []string
		wantErr bool
	}{
		{
			name:   "v1",
			object: testutil.LoadUnstructuredFromFile(t, "crd-v1.yaml"),
			want:   []string{"v1"},
		},
		{
			name:   "v1beta1",
			object: testutil.LoadUnstructuredFromFile(t, "crd-v1beta1.yaml"),
			want:   []string{"v1"},
		},
		{
			name:   "v1beta1 - versions",
			object: testutil.LoadUnstructuredFromFile(t, "crd-v1beta1-versions.yaml"),
			want:   []string{"v1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crdTool, err := octant.NewCustomResourceDefinitionTool(tt.object)
			require.NoError(t, err)

			got, err := crdTool.Versions()
			testutil.RequireErrorOrNot(t, tt.wantErr, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestCustomResourceDefinition_Version(t *testing.T) {
	tests := []struct {
		name    string
		object  *unstructured.Unstructured
		version string
		want    octant.CustomResourceDefinitionVersion
		wantErr bool
	}{
		{
			name:    "v1",
			object:  testutil.LoadUnstructuredFromFile(t, "crd-v1.yaml"),
			version: "v1",
			want: octant.CustomResourceDefinitionVersion{
				Version: "v1",
				PrinterColumns: []octant.CustomResourceDefinitionPrinterColumn{
					{
						Name:        "Spec",
						Type:        "string",
						Description: "The cron spec defining the interval a CronJob is run",
						JSONPath:    ".spec.cronSpec",
					},
					{
						Name:        "Replicas",
						Type:        "integer",
						Description: "The number of jobs launched by the CronJob",
						JSONPath:    ".spec.replicas",
					},
					{
						Name:     "Age",
						Type:     "date",
						JSONPath: ".metadata.creationTimestamp",
					},
				},
			},
		},
		{
			name:   "v1beta1",
			object: testutil.LoadUnstructuredFromFile(t, "crd-v1beta1.yaml"),
			want: octant.CustomResourceDefinitionVersion{
				Version: "v1",
				PrinterColumns: []octant.CustomResourceDefinitionPrinterColumn{
					{
						Name:        "Spec",
						Type:        "string",
						Description: "The cron spec defining the interval a CronJob is run",
						JSONPath:    ".spec.cronSpec",
					},
					{
						Name:        "Replicas",
						Type:        "integer",
						Description: "The number of jobs launched by the CronJob",
						JSONPath:    ".spec.replicas",
					},
					{
						Name:     "Age",
						Type:     "date",
						JSONPath: ".metadata.creationTimestamp",
					},
				},
			},
		},
		{
			name:   "v1beta1 no columns",
			object: testutil.LoadUnstructuredFromFile(t, "crd-v1beta1-versions.yaml"),
			want: octant.CustomResourceDefinitionVersion{
				Version: "v1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crdTool, err := octant.NewCustomResourceDefinitionTool(tt.object)
			require.NoError(t, err)

			got, err := crdTool.Version("v1")
			testutil.RequireErrorOrNot(t, tt.wantErr, err)
			require.Equal(t, tt.want, got)
		})
	}
}
