/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestPayload_GroupVersionKind(t *testing.T) {
	payload := Payload{
		"group":   "group",
		"version": "version",
		"kind":    "kind",
	}

	got, err := payload.GroupVersionKind()
	require.NoError(t, err)

	expected := schema.GroupVersionKind{
		Group:   "group",
		Version: "version",
		Kind:    "kind",
	}

	assert.Equal(t, expected, got)
}

func TestPayload_String(t *testing.T) {
	tests := []struct {
		name     string
		payload  Payload
		key      string
		expected string
		isErr    bool
	}{
		{
			name:     "valid",
			payload:  Payload{"string": "string"},
			key:      "string",
			expected: "string",
		},
		{
			name:    "not string",
			payload: Payload{"string": 7},
			key:     "string",
			isErr:   true,
		},
		{
			name:    "key does not exist",
			payload: Payload{"string": "string"},
			key:     "invalid",
			isErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.payload.String(test.key)
			if test.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, got)
		})
	}
}

func TestPayload_StringSlice(t *testing.T) {
	tests := []struct {
		name     string
		payload  Payload
		key      string
		expected []string
		isErr    bool
	}{
		{
			name:     "valid",
			payload:  Payload{"slice": []interface{}{"string"}},
			key:      "slice",
			expected: []string{"string"},
		},
		{
			name:    "not slice",
			payload: Payload{"slice": 7},
			key:     "slice",
			isErr:   true,
		},
		{
			name:    "not string slice",
			payload: Payload{"slice": []int{7}},
			key:     "slice",
			isErr:   true,
		},
		{
			name:    "key does not exist",
			payload: Payload{},
			key:     "invalid",
			isErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.payload.StringSlice(test.key)
			if test.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, got)
		})
	}
}

func TestPayload_Float64(t *testing.T) {
	tests := []struct {
		name     string
		payload  Payload
		key      string
		isErr    bool
		expected float64
	}{
		{
			name:     "source is string",
			payload:  Payload{"float64": "7"},
			key:      "float64",
			expected: float64(7),
		},
		{
			name:     "source is float64",
			payload:  Payload{"float64": float64(7)},
			key:      "float64",
			expected: float64(7),
		},
		{
			name:    "value is not string or float64",
			payload: Payload{"float64": true},
			key:     "float64",
			isErr:   true,
		},
		{
			name:    "key does not exist",
			payload: Payload{},
			key:     "invalid",
			isErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.payload.Float64(test.key)
			if test.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, test.expected, got)
		})
	}

}
