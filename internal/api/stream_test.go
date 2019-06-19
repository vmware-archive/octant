/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware/octant/internal/octant"
)

func Test_stream(t *testing.T) {
	cases := []struct {
		name         string
		event        octant.Event
		expectedBody string
	}{
		{
			name:         "event with data",
			event:        octant.Event{Data: []byte("output")},
			expectedBody: fmt.Sprintf("data: output\n\n"),
		},
		{
			name:         "event with name and data",
			event:        octant.Event{Type: "name", Data: []byte("output")},
			expectedBody: fmt.Sprintf("event: name\ndata: output\n\n"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, cancel := context.WithCancel(context.Background())
			ch := make(chan octant.Event)

			done := make(chan bool, 1)

			s := &eventSourceStreamer{
				w: w,
			}

			go func() {
				s.Stream(ctx, ch)
				done <- true
			}()

			ch <- tc.event
			cancel()

			<-done

			resp := w.Result()
			defer func() {
				require.NoError(t, resp.Body.Close())
			}()
			actualBody, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedBody, string(actualBody))

			actualHeaders := w.Header()
			expectedHeaders := http.Header{
				"Content-Type":                []string{"text/event-stream"},
				"Cache-Control":               []string{"no-cache"},
				"Connection":                  []string{"keep-alive"},
				"Access-Control-Allow-Origin": []string{"*"},
			}

			for k := range expectedHeaders {
				expected := expectedHeaders.Get(k)
				actual := actualHeaders.Get(k)
				assert.Equalf(t, expected, actual, "expected header %s to be %s; actual %s",
					k, expected, actual)
			}

		})
	}
}

func Test_notFoundRedirectPath(t *testing.T) {
	cases := []struct {
		name     string
		expected string
	}{
		{
			name:     "/api/v1/content/overview/namespace/default/workloads/deployments/nginx-deployment/",
			expected: "/content/overview/namespace/default/workloads/deployments",
		},
		{
			name:     "/api/v1/content/overview/namespace/default/workloads/invalid/",
			expected: "/content/overview/namespace/default/workloads",
		},
		{
			name:     "/api/v1/content/overview/namespace/default/invalid/",
			expected: "/content/overview/namespace/default",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := notFoundRedirectPath(tc.name)
			assert.Equal(t, tc.expected, got)
		})
	}
}

type simpleResponseWriter struct {
	data       []byte
	statusCode int

	writeCh chan bool
}

func newSimpleResponseWriter() *simpleResponseWriter {
	return &simpleResponseWriter{
		writeCh: make(chan bool, 1),
	}
}

func (w *simpleResponseWriter) Header() http.Header {
	return http.Header{}
}
func (w *simpleResponseWriter) Write(data []byte) (int, error) {
	w.data = data
	w.writeCh <- true
	return 0, nil
}
func (w *simpleResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func Test_stream_errors_without_flusher(t *testing.T) {
	w := newSimpleResponseWriter()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := make(chan octant.Event, 1)

	s := eventSourceStreamer{
		w: w,
	}

	go s.Stream(ctx, ch)
	ch <- octant.Event{Data: []byte("output")}

	<-w.writeCh

	assert.Equal(t, http.StatusInternalServerError, w.statusCode)
}
