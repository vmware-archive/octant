/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package log

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/testutil"
)

func TestOctantSink(t *testing.T) {
	validConvert := func(b []byte) (Message, error) {
		return Message{}, nil
	}

	invalidConvert := func(b []byte) (Message, error) {
		return Message{}, errors.New("invalid")
	}

	type ctorArgs struct {
		options []OctantSinkOption
	}
	type args struct {
		message string
	}
	tests := []struct {
		name     string
		ctorArgs ctorArgs
		args     args
		wantErr  bool
	}{
		{
			name: "conversion is success",
			ctorArgs: ctorArgs{
				options: []OctantSinkOption{
					func(o *OctantSink) {
						o.converter = validConvert
					},
				},
			},
		},
		{
			name: "conversion fails",
			ctorArgs: ctorArgs{
				options: []OctantSinkOption{
					func(o *OctantSink) {
						o.converter = invalidConvert
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewOctantSink(tt.ctorArgs.options...)

			defer func() {
				_ = s.Close()
			}()

			ch, cancel := s.Listen()
			defer cancel()

			done := make(chan bool, 1)

			n, err := s.Write([]byte(tt.args.message))
			testutil.RequireErrorOrNot(t, tt.wantErr, err, func() {
				go func() {
					<-ch
					done <- true
				}()
				<-done
				require.Len(t, tt.args.message, n)
			})
		})
	}
}

func TestConvertBytesToMessage(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Message
		wantErr bool
	}{
		{
			name: "in general with JSON",
			args: args{
				bytes: []byte(strings.Join([]string{
					"2020-09-03T14:39:51.115-0400",
					"INFO",
					"file.go:50",
					"message",
					`{"foo": "bar"}`,
				}, "\t") + "\n"),
			},
			want: Message{
				Date:     1599158391,
				LogLevel: "INFO",
				Location: "file.go:50",
				Text:     "message",
				JSON:     `{"foo": "bar"}`,
			},
		},
		{
			name: "in general without JSON",
			args: args{
				bytes: []byte(strings.Join([]string{
					"2020-09-03T14:39:51.115-0400",
					"INFO",
					"file.go:50",
					"message",
				}, "\t") + "\n"),
			},
			want: Message{
				Date:     1599158391,
				LogLevel: "INFO",
				Location: "file.go:50",
				Text:     "message",
			},
		},
		{
			name: "invalid format (too short)",
			args: args{
				bytes: []byte(strings.Join([]string{
					"message",
				}, "\t") + "\n"),
			},
			wantErr: true,
		},
		{
			name: "invalid format (too long)",
			args: args{
				bytes: []byte(strings.Join([]string{
					"message",
					"message",
					"message",
					"message",
					"message",
					"message",
				}, "\t") + "\n"),
			},
			wantErr: true,
		}, {
			name: "invalid timestamp (not ISO8601)",
			args: args{
				bytes: []byte(strings.Join([]string{
					"2020-09-03T14:39:51.115-0400invalid",
					"INFO",
					"file.go:50",
					"message",
				}, "\t") + "\n"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertBytesToMessage(tt.args.bytes)
			testutil.RequireErrorOrNot(t, tt.wantErr, err, func() {
				require.Equal(t, tt.want, got)
			})
		})
	}
}
