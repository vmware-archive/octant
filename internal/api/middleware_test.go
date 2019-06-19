/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_rebindHandler(t *testing.T) {
	cases := []struct {
		name         string
		host         string
		expectedCode int
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
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			fake := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, "response")
			})

			wrapped := rebindHandler(acceptedHosts)(fake)

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
