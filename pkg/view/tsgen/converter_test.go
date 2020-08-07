/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package tsgen

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConverter_Convert(t *testing.T) {
	// var c component.Component = component.NewText("text")

	componentNames := []string{"Text"}

	type args struct {
		in reflect.Type
	}
	tests := []struct {
		name      string
		args      args
		want      string
		wantNames []string
		wantError bool
	}{
		{
			name: "string",
			args: args{in: reflect.TypeOf("string")},
			want: "string",
		},
		{
			name: "bool",
			args: args{in: reflect.TypeOf(true)},
			want: "boolean",
		},
		{
			name: "simple map",
			args: args{in: reflect.TypeOf(map[string]string{"foo": "bar"})},
			want: "{[key:string]:string}",
		},
		{
			name: "array",
			args: args{in: reflect.TypeOf([]string{"foo"})},
			want: "string[]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConverter(componentNames)
			got, gotNames, err := c.Convert(tt.args.in)
			if tt.wantError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantNames, gotNames)
			assert.Equal(t, tt.want, got)
		})
	}
}
