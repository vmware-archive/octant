/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
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

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func Test_rebindHandler(t *testing.T) {
	cases := []struct {
		name                       string
		host                       string
		origin                     string
		expectedCode               int
		listenerKey                string
		listenerAddr               string
		disableCrossOriginKey      string
		disableCrossOriginChecking bool
		errorMessage               string
	}{
		{
			name:         "in general",
			expectedCode: http.StatusOK,
		},
		{
			name:         "rebind",
			host:         "hacker.com",
			expectedCode: http.StatusForbidden,
			errorMessage: "forbidden host\n",
		},
		{
			name:         "invalid host",
			host:         ":::::::::",
			expectedCode: http.StatusBadRequest,
			errorMessage: "bad request\n",
		},
		{
			name:         "custom host",
			host:         "0.0.0.0",
			expectedCode: http.StatusOK,
			listenerKey:  "listener-addr",
			listenerAddr: "0.0.0.0:0000",
		},
		{
			name:                       "disable CORS",
			host:                       "example.com",
			origin:                     "hacker.com",
			expectedCode:               http.StatusOK,
			disableCrossOriginKey:      "disable-origin-check",
			disableCrossOriginChecking: true,
			listenerKey:                "listener-addr",
			listenerAddr:               "example.com:80",
			errorMessage:               "response",
		},
		{
			name:         "fails CORS and invalid host",
			host:         "example.com",
			origin:       "hacker.com",
			expectedCode: http.StatusForbidden,
			errorMessage: "forbidden host: forbidden bad origin\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.listenerKey != "" {
				viper.Set(tc.listenerKey, tc.listenerAddr)
				defer viper.Set(tc.listenerKey, "")
			}

			if tc.disableCrossOriginKey != "" {
				viper.Set(tc.disableCrossOriginKey, tc.disableCrossOriginChecking)
				defer viper.Set(tc.disableCrossOriginKey, false)
			}
			fake := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, "response")
			})

			wrapped := rebindHandler(context.TODO(), AcceptedHosts())(fake)

			ts := httptest.NewServer(wrapped)
			defer ts.Close()

			req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
			require.NoError(t, err)

			if tc.origin != "" {
				req.Header["Origin"] = []string{tc.origin}
			}

			if tc.host != "" {
				req.Host = tc.host
			}

			res, err := http.DefaultClient.Do(req)
			require.NoError(t, err)

			if tc.errorMessage != "" {
				message, err := ioutil.ReadAll(res.Body)
				require.NoError(t, err)
				require.Equal(t, tc.errorMessage, string(message))
			}

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
			require.Equal(t, tc.expected, ShouldAllowHost(tc.host, tc.acceptedHosts))
		})
	}
}

func Test_checkSameOrigin(t *testing.T) {
	cases := []struct {
		name     string
		host     string
		origin   string
		expected bool
	}{
		{
			name:     "host/origin match",
			host:     "192.168.1.1:7777",
			origin:   "http://192.168.1.1:7777",
			expected: true,
		},
		{
			name:     "host/origin do not match",
			host:     "192.168.1.1:7777",
			origin:   "http://127.0.0.1:7777",
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := &http.Request{
				Host:   tc.host,
				Header: make(http.Header, 1),
			}
			r.Header.Set("Origin", tc.origin)
			require.Equal(t, tc.expected, checkSameOrigin(r))
		})
	}
}
