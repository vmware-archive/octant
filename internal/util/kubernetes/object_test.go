/*
 *  Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 *  SPDX-License-Identifier: Apache-2.0
 *
 */

package kubernetes

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/install"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/vmware-tanzu/octant/internal/testutil"
)

func TestReadObject(t *testing.T) {
	data, err := ioutil.ReadFile(filepath.Join("testdata", "deployment.yaml"))
	require.NoError(t, err)

	r := bytes.NewReader(data)

	o, err := ReadObject(r)
	require.NoError(t, err)

	require.Equal(t, "nginx-deployment", o.GetName())
}

func TestSerializeToString(t *testing.T) {
	pod := testutil.CreatePod("pod")

	tests := []struct {
		name    string
		object  runtime.Object
		wanted  string
		wantErr bool
	}{
		{
			name:   "in general",
			object: pod,
			wanted: string(testutil.LoadTestData(t, "pod.yaml")),
		},
		{
			name:    "nil object",
			object:  nil,
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := SerializeToString(test.object)
			if test.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, test.wanted, actual)
		})
	}
}

func TestFromUnstructured(t *testing.T) {
	install.Install(scheme.Scheme)

	type args struct {
		as   runtime.Object
		path string
	}
	tests := []struct {
		name    string
		args    args
		check   func(object runtime.Object)
		wantErr bool
	}{
		{
			name: "crd",
			args: args{
				as:   &apiextv1.CustomResourceDefinition{},
				path: "crd.yaml",
			},
			check: func(object runtime.Object) {
				crd, ok := object.(*apiextv1.CustomResourceDefinition)
				require.True(t, ok)

				assert.Equal(t, "minioinstances.operator.min.io", crd.Name)
			},
		},
		{
			name: "deployment",
			args: args{
				as:   &appsv1.Deployment{},
				path: "deployment.yaml",
			},
			check: func(object runtime.Object) {
				d, ok := object.(*appsv1.Deployment)
				require.True(t, ok)

				assert.Equal(t, "nginx-deployment", d.Name)
				labels := map[string]string{
					"app.kubernetes.io/name":     "nginx",
					"app.kubernetes.io/instance": "sample",
					"app.kubernetes.io/version":  "v1",
				}
				assert.Equal(t, labels, d.Labels)
			},
		},
		{
			name: "nil as",
			args: args{
				as:   nil,
				path: "deployment.yaml",
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			u := testutil.LoadUnstructuredFromFile(t, test.args.path)

			err := FromUnstructured(u, test.args.as)
			if test.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			test.check(test.args.as)
		})
	}
}
