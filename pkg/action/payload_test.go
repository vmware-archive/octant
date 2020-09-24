/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package action

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware-tanzu/octant/internal/testutil"
)

func TestCreatePayload(t *testing.T) {
	action := "action"

	fields := map[string]interface{}{"foo": "bar"}
	got := CreatePayload(action, fields)

	expected := Payload{
		"action": action,
		"foo":    "bar",
	}

	assert.Equal(t, expected, got)
}

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

func TestPayload_Uint16(t *testing.T) {
	tests := []struct {
		name     string
		payload  Payload
		key      string
		isErr    bool
		expected uint16
	}{
		{
			name:     "source is int",
			payload:  Payload{"uint16": float64(7)},
			key:      "uint16",
			expected: uint16(7),
		},
		{
			name:    "source overflows",
			payload: Payload{"uint16": 2 ^ 17},
			key:     "uint16",
			isErr:   true,
		},
		{
			name:    "source overflows",
			payload: Payload{"uint16": -1},
			key:     "uint16",
			isErr:   true,
		},
		{
			name:    "value is not int",
			payload: Payload{"uint16": true},
			key:     "uint16",
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
			got, err := test.payload.Uint16(test.key)
			if test.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, test.expected, got)
		})
	}
}

func TestPayload_Int64(t *testing.T) {
	tests := []struct {
		name     string
		payload  Payload
		key      string
		isErr    bool
		expected int64
	}{
		{
			name:     "source is int",
			payload:  Payload{"int64": float64(7)},
			key:      "int64",
			expected: int64(7),
		},
		{
			name:    "source overflows",
			payload: Payload{"int64": float64(1 << 64)},
			key:     "int64",
			isErr:   true,
		},
		{
			name:    "source overflows",
			payload: Payload{"int64": float64(-1 << 64)},
			key:     "int64",
			isErr:   true,
		},
		{
			name:    "value is not int",
			payload: Payload{"int64": true},
			key:     "int64",
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
			got, err := test.payload.Int64(test.key)
			if test.isErr {
				require.Error(t, err)
				fmt.Println(got, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, test.expected, got)
		})
	}
}

func TestPayload_Raw(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		payload Payload
		wantErr bool
		want    []byte
	}{
		{
			name: "key exists",
			args: args{
				key: "key",
			},
			payload: Payload{
				"key": "value",
			},
			want: []byte(`"value"`),
		},
		{
			name: "key does not exist",
			args: args{
				key: "key",
			},
			payload: Payload{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.payload.Raw(tt.args.key)
			testutil.RequireErrorOrNot(t, tt.wantErr, err, func() {
				require.Equal(t, tt.want, got)
			})
		})
	}
}
