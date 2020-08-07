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

	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func TestField_TSType(t *testing.T) {
	type Nested struct {
		One string
		Two bool
	}

	type myStruct struct {
		StringField    string
		BoolField      bool
		Int64Field     int64
		AliasedField   component.TextStatus
		ComponentField component.Component
		Map1Field      map[string]string
		Map2Field      map[string]interface{}
		Nested         Nested
	}

	xType := reflect.TypeOf(myStruct{})

	type args struct {
		field reflect.StructField
	}
	tests := []struct {
		name      string
		args      args
		want      string
		wantNames []string
		wantErr   bool
	}{
		{
			name: "string",
			args: args{
				field: xType.Field(0),
			},
			want: "string",
		},
		{
			name: "bool",
			args: args{
				field: xType.Field(1),
			},
			want: "boolean",
		},
		{
			name: "int64",
			args: args{
				field: xType.Field(2),
			},
			want: "number",
		},
		{
			name: "aliased field",
			args: args{
				field: xType.Field(3),
			},
			want: "number",
		},
		{
			name: "component field",
			args: args{
				field: xType.Field(4),
			},
			want: "Component<any>",
		},
		{
			name: "string map field",
			args: args{
				field: xType.Field(5),
			},
			want: "{[key:string]:string}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotNames, err := tsType(tt.args.field, []string{})
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantNames, gotNames)
		})
	}

}
