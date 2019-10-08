/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_printMatchLabels(t *testing.T) {
	cases := []struct {
		name        string
		matchLabels map[string]string
		expected    string
	}{
		{
			name:        "single label",
			matchLabels: map[string]string{"foo": "bar"},
			expected:    "foo:bar",
		},
		{
			name: "multiple labels",
			matchLabels: map[string]string{
				"foo": "bar",
				"bar": "foo",
			},
			expected: "bar:foo, foo:bar",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := printMatchLabels(tc.matchLabels)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func Test_printLabelSelectorRequirement(t *testing.T) {
	cases := []struct {
		name         string
		requirements []metav1.LabelSelectorRequirement
		expected     string
	}{
		{
			name:         "empty",
			requirements: nil,
			expected:     "",
		},
		{
			name: "in",
			requirements: []metav1.LabelSelectorRequirement{
				{
					Key:      "key",
					Operator: metav1.LabelSelectorOpIn,
					Values:   []string{"x", "y"},
				},
			},
			expected: "key in [x, y]",
		},
		{
			name: "not in",
			requirements: []metav1.LabelSelectorRequirement{
				{
					Key:      "key",
					Operator: metav1.LabelSelectorOpNotIn,
					Values:   []string{"x", "y"},
				},
			},
			expected: "key not in [x, y]",
		},
		{
			name: "exists",
			requirements: []metav1.LabelSelectorRequirement{
				{
					Key:      "key",
					Operator: metav1.LabelSelectorOpExists,
				},
			},
			expected: "key exists",
		},
		{
			name: "does not exist",
			requirements: []metav1.LabelSelectorRequirement{
				{
					Key:      "key",
					Operator: metav1.LabelSelectorOpDoesNotExist,
				},
			},
			expected: "key does not exist",
		},
		{
			name: "multiple (2)",
			requirements: []metav1.LabelSelectorRequirement{
				{
					Key:      "key1",
					Operator: metav1.LabelSelectorOpIn,
					Values:   []string{"x", "y"},
				},
				{
					Key:      "key2",
					Operator: metav1.LabelSelectorOpDoesNotExist,
				},
			},
			expected: "key1 in [x, y], key2 does not exist",
		},
		{
			name: "multiple (3+)",
			requirements: []metav1.LabelSelectorRequirement{
				{
					Key:      "key1",
					Operator: metav1.LabelSelectorOpIn,
					Values:   []string{"x", "y"},
				},
				{
					Key:      "key2",
					Operator: metav1.LabelSelectorOpDoesNotExist,
				},
				{
					Key:      "key3",
					Operator: metav1.LabelSelectorOpExists,
				},
			},
			expected: "key1 in [x, y], key2 does not exist, key3 exists",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := printLabelSelectorRequirement(tc.requirements)
			assert.Equal(t, tc.expected, got)
		})
	}
}
