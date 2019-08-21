/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"fmt"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_rebindHandler(t *testing.T) {
	cases := []struct {
		name         string
		host         string
		expectedCode int
		listenerKey  string
		listenerAddr string
	}{
		{
			name:         "in general",
			expectedCode: http.StatusOK,
		},
		{
			name:         "rebind",
			host:         "hacker.com",
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "invalid host",
			host:         ":::::::::",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "custom host",
			host:         "0.0.0.0",
			expectedCode: http.StatusOK,
			listenerKey:  "OCTANT_LISTENER_ADDR",
			listenerAddr: "0.0.0.0:0000",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.listenerKey != "" {
				os.Setenv(tc.listenerKey, tc.listenerAddr)
				defer os.Unsetenv(tc.listenerKey)
			}
			fake := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, "response")
			})

			wrapped := rebindHandler(context.TODO(), acceptedHosts())(fake)

			ts := httptest.NewServer(wrapped)
			defer ts.Close()

			req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
			require.NoError(t, err)

			if tc.host != "" {
				req.Host = tc.host
			}

			res, err := http.DefaultClient.Do(req)
			require.NoError(t, err)

			require.Equal(t, tc.expectedCode, res.StatusCode)
		})
	}
}

func Test_shouldAllowHost(t *testing.T) {
	cases := []struct {
		name          string
		host          string
		acceptedHosts []string
		expected      bool
	}{
		{
			name:          "0.0.0.0 allow all",
			host:          "192.168.1.1",
			acceptedHosts: []string{"127.0.0.1", "localhost", "0.0.0.0"},
			expected:      true,
		},
		{
			name:          "deny 192.168.1.1",
			host:          "192.168.1.1",
			acceptedHosts: []string{"127.0.0.1", "localhost"},
			expected:      false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, shouldAllowHost(tc.host, tc.acceptedHosts))
		})
	}
}
