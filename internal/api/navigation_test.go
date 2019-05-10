package api

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/heptio/developer-dash/internal/clustereye"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_navigation_handler(t *testing.T) {
	validSections := &fakeNavSections{
		sections: []*clustereye.Navigation{
			{},
		},
	}

	invalidSections := &fakeNavSections{
		sectionsErr: errors.Errorf("foo"),
	}

	logger := log.TestLogger(t)

	cases := []struct {
		name       string
		nav        *navigation
		statusCode int
		body       []byte
	}{

		{
			name:       "in general",
			nav:        newNavigation(validSections, logger),
			statusCode: http.StatusOK,
			body:       []byte("{\"sections\":[{}]}\n"),
		},
		{
			name:       "no section generator",
			nav:        newNavigation(nil, logger),
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "section generate error",
			nav:        newNavigation(invalidSections, logger),
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
	sections    []*clustereye.Navigation
	sectionsErr error
}

func (ns *fakeNavSections) Sections(ctx context.Context, namespace string) ([]*clustereye.Navigation, error) {
	return ns.sections, ns.sectionsErr
}
