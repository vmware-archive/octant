package overview

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_handler_routes(t *testing.T) {
	dynamicContent := newFakeContent(false)

	cases := []struct {
		name         string
		path         string
		values       url.Values
		generator    generator
		expectedCode int
		expectedBody string
	}{
		{
			name:         "GET dynamic content",
			path:         "/api/real",
			values:       url.Values{"namespace": []string{"default"}},
			generator:    newStubbedGenerator([]content.Content{dynamicContent}, nil),
			expectedCode: http.StatusOK,
			expectedBody: `{"views":[{"contents":[{"type":"stubbed"}],"title":"main title"}],"viewComponents":null}`,
		},
		{
			name:         "error generating dynamic content",
			path:         "/api/real",
			values:       url.Values{"namespace": []string{"default"}},
			generator:    newStubbedGenerator(nil, errors.Errorf("broken")),
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"error":{"code":500,"message":"broken"}}`,
		},
		{
			name:         "GET invalid path",
			path:         "/api/missing",
			generator:    newStubbedGenerator([]content.Content{dynamicContent}, nil),
			expectedCode: http.StatusNotFound,
			expectedBody: `{"error":{"code":404,"message":"content not found"}}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := newHandler("/api", tc.generator, stubStream, log.NopLogger())

			ts := httptest.NewServer(h)
			defer ts.Close()

			u, err := url.Parse(ts.URL)
			require.NoError(t, err)

			u.Path = tc.path
			u.RawQuery = tc.values.Encode()

			resp, err := http.Get(u.String())
			require.NoError(t, err)
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedCode, resp.StatusCode)
			assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
			assert.Equal(t, tc.expectedBody, strings.TrimSpace(string(body)))
		})
	}
}

var (
	stubStream = func(ctx context.Context, w http.ResponseWriter, ch chan []byte) {}
)

type stubbedGenerator struct {
	Contents []content.Content
	genErr   error
}

func newStubbedGenerator(contents []content.Content, genErr error) *stubbedGenerator {
	return &stubbedGenerator{
		Contents: contents,
		genErr:   genErr,
	}
}

func (g *stubbedGenerator) Generate(ctx context.Context, path, prefix, namespace string) (ContentResponse, error) {
	switch {
	case strings.HasPrefix(path, "/real"):
		return ContentResponse{
			Views: []Content{
				{
					Title:    "main title",
					Contents: stubbedContent,
				},
			},
		}, g.genErr

	default:
		return emptyContentResponse, contentNotFound
	}
}

// generatorFunc allows a bare Generate function to implement Generator
type generatorFunc func(ctx context.Context, path, prefix, namespace string) (ContentResponse, error)

func (g generatorFunc) Generate(ctx context.Context, path, prefix, namespace string) (ContentResponse, error) {
	return g(ctx, path, prefix, namespace)
}

// Test that context is passed on from handler and includes logger
func TestHandlerContext(t *testing.T) {
	logger := log.NopLogger()
	g := generatorFunc(func(ctx context.Context, path, prefix, namespace string) (ContentResponse, error) {
		l := log.From(ctx)
		assert.True(t, logger == l, "unexpected logger from context")
		return emptyContentResponse, nil
	})

	h := newHandler("/api", g, stubStream, logger)
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/api", nil)
	require.NoError(t, err)

	h.ServeHTTP(w, r)
}
