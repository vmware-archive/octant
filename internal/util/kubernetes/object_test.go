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

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"

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
