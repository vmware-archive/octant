package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/heptio/developer-dash/internal/hcli"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_navigation_handler(t *testing.T) {
	validSections := &fakeNavSections{
		sections: []*hcli.Navigation{
			{},
		},
	}

	invalidSections := &fakeNavSections{
		sectionsErr: errors.Errorf("foo"),
	}

	cases := []struct {
		name       string
		nav        *navigation
		statusCode int
		body       []byte
	}{

		{
			name:       "in general",
			nav:        newNavigation(validSections, nil),
			statusCode: http.StatusOK,
			body:       []byte("{\"sections\":[{}]}\n"),
		},
		{
			name:       "no section generator",
			nav:        newNavigation(nil, nil),
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "section generate error",
			nav:        newNavigation(invalidSections, nil),
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
	sections    []*hcli.Navigation
	sectionsErr error
}

func (ns *fakeNavSections) Sections() ([]*hcli.Navigation, error) {
	return ns.sections, ns.sectionsErr
}
