/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package tsgen

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTSGen_ReflectTemplate(t *testing.T) {
	type args struct {
		names []string
	}
	tests := []struct {
		name     string
		args     args
		wantFile string
		wantErr  bool
	}{
		{
			name: "in general",
			args: args{
				names: []string{"Link", "Text"},
			},
			wantFile: "reflect.go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tg, err := NewTSGen()
			require.NoError(t, err)

			got, err := tg.ReflectTemplate(tt.args.names)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assertTemplate(t, tt.wantFile, got)
		})
	}
}

func TestTSGen_ComponentConfig(t *testing.T) {
	type args struct {
		c Component
	}
	tests := []struct {
		name     string
		args     args
		wantFile string
		wantErr  bool
	}{
		{
			name: "in general",
			args: args{
				c: Component{
					Name:   "Link",
					TSName: "link",
					Fields: []Field{
						{Name: "value", Type: "string", Optional: false},
						{Name: "ref", Type: "string", Optional: false},
						{Name: "status", Type: "number", Optional: true},
						{Name: "statusDetail", Type: "component", Optional: true},
					},
				},
			},
			wantFile: "link-config.ts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tg, err := NewTSGen()
			require.NoError(t, err)

			got, err := tg.ComponentConfig(tt.args.c)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assertTemplate(t, tt.wantFile, got)
		})
	}
}

func clean(b []byte) string {
	return strings.TrimSpace(string(b))
}

func assertTemplate(t *testing.T, wantFile string, got []byte) {
	want, err := ioutil.ReadFile(filepath.Join("testdata", wantFile))
	require.NoError(t, err)

	require.Equal(t, clean(want), clean(got))
}
