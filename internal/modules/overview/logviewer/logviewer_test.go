/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package logviewer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_ToComponent(t *testing.T) {
	cases := []struct {
		name     string
		object   runtime.Object
		expected component.Component
		isErr    bool
	}{
		{
			name: "with containers",
			object: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pod",
					Namespace: "default",
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{Name: "one"},
						{Name: "two"},
					},
				},
			},
			expected: component.NewLogs("default", "pod", []string{"one", "two"}),
		},
		{
			name: "with init containers",
			object: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pod",
					Namespace: "default",
				},
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{Name: "init"},
					},
					Containers: []corev1.Container{
						{Name: "one"},
						{Name: "two"},
					},
				},
			},
			expected: component.NewLogs("default", "pod", []string{"init", "one", "two"}),
		},
		{
			name:   "nil",
			object: nil,
			isErr:  true,
		},
		{
			name:   "not a v1 Pod",
			object: &corev1.Service{},
			isErr:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ToComponent(tc.object)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expected, got)
		})
	}

}
