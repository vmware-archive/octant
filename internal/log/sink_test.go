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
	validConvert := func(b []byte, options ...Option) (Message, error) {
		return Message{}, nil
	}

	invalidConvert := func(b []byte, options ...Option) (Message, error) {
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

func TestOctantSink_Close(t *testing.T) {
	s := NewOctantSink()
	ch, _ := s.Listen()

	done := make(chan struct{}, 1)

	closed := false
	go func() {
		<-ch
		closed = true
		done <- struct{}{}
	}()

	require.NoError(t, s.Close())

	<-done
	require.True(t, closed)
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
			name: "info message with JSON",
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
				ID:       "12345",
				Date:     1599158391115000000,
				LogLevel: "INFO",
				Location: "file.go:50",
				Text:     "message",
				JSON:     `{"foo": "bar"}`,
			},
		},
		{
			name: "info message without JSON",
			args: args{
				bytes: []byte(strings.Join([]string{
					"2020-09-03T14:39:51.115-0400",
					"INFO",
					"file.go:50",
					"message",
				}, "\t") + "\n"),
			},
			want: Message{
				ID:       "12345",
				Date:     1599158391115000000,
				LogLevel: "INFO",
				Location: "file.go:50",
				Text:     "message",
			},
		},
		{
			name: "error message with JSON",
			args: args{
				bytes: []byte("2020-09-23T09:10:52.181-0400\tERROR\tapi/content_manager.go:154\tgenerate content\t{\"client-id\": \"2c3d1670-fd9e-11ea-b34c-e450ebbc2e8c\", \"err\": \"generate content: calling contentHandler: check access to list CacheKey[Namespace='default', APIVersion='cluster.x-k8s.io/v1alpha3', Kind='Cluster']: unable to get resource for group kind Cluster.cluster.x-k8s.io: no matches for kind \\\"Cluster\\\" in group \\\"cluster.x-k8s.io\\\"\", \"content-path\": \"octant-clusterapi-plugin\"}\ngithub.com/vmware-tanzu/octant/internal/api.(*ContentManager).runUpdate.func1\n\t/Users/bryan/Development/projects/octant/internal/api/content_manager.go:154\ngithub.com/vmware-tanzu/octant/internal/api.(*InterruptiblePoller).Run.func1\n\t/Users/bryan/Development/projects/octant/internal/api/poller.go:86\ngithub.com/vmware-tanzu/octant/internal/api.(*InterruptiblePoller).Run\n\t/Users/bryan/Development/projects/octant/internal/api/poller.go:95\ngithub.com/vmware-tanzu/octant/internal/api.(*ContentManager).Start\n\t/Users/bryan/Development/projects/octant/internal/api/content_manager.go:128\n"),
			},
			want: Message{
				ID:       "12345",
				Date:     1600866652181000000,
				LogLevel: "ERROR",
				Location: "api/content_manager.go:154",
				Text:     "generate content",
				JSON:     "{\"client-id\": \"2c3d1670-fd9e-11ea-b34c-e450ebbc2e8c\", \"err\": \"generate content: calling contentHandler: check access to list CacheKey[Namespace='default', APIVersion='cluster.x-k8s.io/v1alpha3', Kind='Cluster']: unable to get resource for group kind Cluster.cluster.x-k8s.io: no matches for kind \\\"Cluster\\\" in group \\\"cluster.x-k8s.io\\\"\", \"content-path\": \"octant-clusterapi-plugin\"}",
				StackTrace: "github.com/vmware-tanzu/octant/internal/api.(*ContentManager).runUpdate.func1\n" +
					"\t/Users/bryan/Development/projects/octant/internal/api/content_manager.go:154\n" +
					"github.com/vmware-tanzu/octant/internal/api.(*InterruptiblePoller).Run.func1\n" +
					"\t/Users/bryan/Development/projects/octant/internal/api/poller.go:86\n" +
					"github.com/vmware-tanzu/octant/internal/api.(*InterruptiblePoller).Run\n" +
					"\t/Users/bryan/Development/projects/octant/internal/api/poller.go:95\n" +
					"github.com/vmware-tanzu/octant/internal/api.(*ContentManager).Start\n" +
					"\t/Users/bryan/Development/projects/octant/internal/api/content_manager.go:128",
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
			got, err := ConvertBytesToMessage(tt.args.bytes, WithIDGenerator(newTestIDGen("12345")))
			testutil.RequireErrorOrNot(t, tt.wantErr, err, func() {
				require.Equal(t, tt.want, got)
			})
		})
	}
}

type testIDGen struct {
	id string
}

var _ MessageIDGenerator = &testIDGen{}

func newTestIDGen(id string) *testIDGen {
	return &testIDGen{id: id}
}

func (i testIDGen) Generate() string {
	return i.id
}
