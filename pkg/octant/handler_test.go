/*
 *  Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 *  SPDX-License-Identifier: Apache-2.0
 *
 */

package octant

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandlerFactory_Handler(t *testing.T) {
	frontend := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "frontend")
	})

	backend := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "backend")
	})

	proxied := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "proxied")
	})

	proxiedTs := httptest.NewServer(proxied)
	defer proxiedTs.Close()

	tests := []struct {
		name         string
		options      []Option
		wantFrontend string
		wantBackend  string
		wantErr      bool
	}{
		{
			name: "wire backend and frontend by default",
			options: []Option{
				func(o *options) {
					o.frontendHandler = func(ctx context.Context) (handler http.Handler, err error) {
						return frontend, nil
					}
					o.backendHandler = func(ctx context.Context) (handler http.Handler, err error) {
						return backend, nil
					}
				},
			},
			wantFrontend: "frontend",
			wantBackend:  "backend",
		},
		{
			name: "proxy frontend",
			options: []Option{
				func(o *options) {
					o.backendHandler = func(ctx context.Context) (handler http.Handler, err error) {
						return backend, nil
					}
				},
				FrontendURL(proxiedTs.URL),
			},
			wantFrontend: "proxied",
			wantBackend:  "backend",
		},
		{
			name: "frontend factory returns an error",
			options: []Option{
				func(o *options) {
					o.frontendHandler = func(ctx context.Context) (handler http.Handler, err error) {
						return nil, fmt.Errorf("error")
					}
					o.backendHandler = func(ctx context.Context) (handler http.Handler, err error) {
						return backend, nil
					}
				},
			},
			wantErr: true,
		},
		{
			name: "backend factory returns an error",
			options: []Option{
				func(o *options) {
					o.frontendHandler = func(ctx context.Context) (handler http.Handler, err error) {
						return frontend, nil
					}
					o.backendHandler = func(ctx context.Context) (handler http.Handler, err error) {
						return nil, fmt.Errorf("error")
					}
				},
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hf := NewHandlerFactory(test.options...)

			ctx := context.Background()

			h, err := hf.Handler(ctx)
			if test.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			ts := httptest.NewServer(h)
			defer ts.Close()

			frontendPath := genTestURL(t, ts.URL, "/")
			backendPath := genTestURL(t, ts.URL, "api", "v1", "foo")

			resFrontend, err := http.Get(frontendPath)
			noCache := resFrontend.Header.Get("Cache-Control")
			require.Equal(t, "no-cache, no-store", noCache)
			require.NoError(t, err)

			resBackend, err := http.Get(backendPath)
			require.NoError(t, err)

			require.Equal(t, test.wantFrontend, string(readFromCloser(t, resFrontend.Body)))
			require.Equal(t, test.wantBackend, string(readFromCloser(t, resBackend.Body)))
		})
	}
}

func TestNewProxiedFrontend(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "content")
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	tests := []struct {
		name      string
		targetURL string
		wantErr   bool
	}{
		{
			name:      "in general",
			targetURL: ts.URL,
		},
		{
			name:      "URL without scheme",
			targetURL: "example.com",
		},
		{
			name:      "invalid URL",
			targetURL: "%",
			wantErr:   true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()

			_, err := NewProxiedFrontend(ctx, test.targetURL)
			if test.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func readFromCloser(t *testing.T, rc io.ReadCloser) []byte {
	data, err := ioutil.ReadAll(rc)
	require.NoError(t, err)
	return data
}

func genTestURL(t *testing.T, base string, parts ...string) string {
	u, err := url.Parse(base)
	require.NoError(t, err)

	u.Path = path.Join(append([]string{u.Path}, parts...)...)
	return u.String()
}
