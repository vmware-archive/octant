/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vmware/octant/internal/log"
	navigation2 "github.com/vmware/octant/pkg/navigation"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_navigation_handler(t *testing.T) {
	validSections := &fakeNavSections{
		sections: []navigation2.Navigation{
			{},
		},
	}

	invalidSections := &fakeNavSections{
		sectionsErr: errors.Errorf("foo"),
	}

	logger := log.TestLogger(t)

	cases := []struct {
		name       string
		nav        *navigationHandler
		statusCode int
		body       []byte
	}{

		{
			name:       "in general",
			nav:        newNavigationHandler(validSections, logger),
			statusCode: http.StatusOK,
			body:       []byte("{\"sections\":[{}]}\n"),
		},
		{
			name:       "no section generator",
			nav:        newNavigationHandler(nil, logger),
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "section generate error",
			nav:        newNavigationHandler(invalidSections, logger),
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ts := httptest.NewServer(tc.nav)
			defer ts.Close()

			res, err := http.Get(ts.URL)
			require.NoError(t, err)

			assert.Equal(t, tc.statusCode, res.StatusCode)
			defer res.Body.Close()

			if tc.body != nil {
				got, err := ioutil.ReadAll(res.Body)
				if assert.NoError(t, err) {
					assert.Equal(t, string(tc.body), string(got))
				}
			}
		})
	}
}

type fakeNavSections struct {
	sections    []navigation2.Navigation
	sectionsErr error
}

func (ns *fakeNavSections) Sections(ctx context.Context, namespace string) ([]navigation2.Navigation, error) {
	return ns.sections, ns.sectionsErr
}
