package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/module"
	"github.com/heptio/developer-dash/internal/module/fake"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_contentEventGenerator(t *testing.T) {
	runEvery := 1 * time.Second

	fn := func(ctx context.Context, path, prefix, namespace string, opts module.ContentOptions) (component.ContentResponse, error) {
		return component.ContentResponse{}, nil
	}

	ceg := contentEventGenerator{
		generatorFn: fn,
		path:        "/path",
		prefix:      "/prefix",
		namespace:   "default",
		runEvery:    runEvery,
	}

	assert.Equal(t, runEvery, ceg.RunEvery())

	got, err := ceg.Generate(context.Background())
	require.NoError(t, err)

	expectedEvent := event{
		eventType: eventTypeContent,
		data:      []byte(`{"content":{"viewComponents":null}}`),
	}

	assert.Equal(t, expectedEvent, got)
}

func Test_navigationEventGenerator(t *testing.T) {
	m := fake.NewModule("module", log.NopLogger())

	neg := navigationEventGenerator{
		modules: []module.Module{m},
	}

	assert.Equal(t, 5*time.Second, neg.RunEvery())

	got, err := neg.Generate(context.Background())
	require.NoError(t, err)

	assert.Equal(t, eventTypeNavigation, got.eventType)
}

func Test_contentStreamer(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, cancel := context.WithCancel(context.Background())

	rcv := make(chan bool, 1)

	fn := func(ctx context.Context, w http.ResponseWriter, ch chan event) {
		e := <-ch

		assert.Equal(t, "data", string(e.data))
		assert.Equal(t, eventType("name"), e.eventType)

		rcv <- true
	}

	cs := contentStreamer{
		eventGenerators: []eventGenerator{&fakeEventGenerator{
			event:    fakeEvent,
			eventErr: nil,
		}},
		w:        w,
		streamFn: fn,
		logger:   log.NopLogger(),
	}

	err := cs.content(ctx)
	require.NoError(t, err)

	<-rcv
	cancel()
}

func Test_contentStream_object_not_found(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, cancel := context.WithCancel(context.Background())

	rcv := make(chan bool, 1)

	fn := func(ctx context.Context, w http.ResponseWriter, ch chan event) {
		e := <-ch

		assert.Equal(t, "/content/overview/namespace/default/workloads/deployments", string(e.data))
		assert.Equal(t, eventTypeObjectNotFound, e.eventType)

		rcv <- true
	}

	cs := contentStreamer{
		eventGenerators: []eventGenerator{&fakeEventGenerator{
			eventErr: NewNotFoundError("/path"),
		}},
		w:           w,
		streamFn:    fn,
		logger:      log.NopLogger(),
		requestPath: "/api/v1/content/overview/namespace/default/workloads/deployments/nginx-deployment/",
	}

	err := cs.content(ctx)
	require.NoError(t, err)

	<-rcv
	cancel()
}

type fakeEventGenerator struct {
	event    event
	eventErr error
}

func (g *fakeEventGenerator) Generate(ctx context.Context) (event, error) {
	return g.event, g.eventErr
}

func (g *fakeEventGenerator) RunEvery() time.Duration {
	return 0
}

func (g *fakeEventGenerator) Name() string {
	return "fakeEventGenerator"
}

var fakeEvent = event{data: []byte("data"), eventType: "name"}

func Test_stream(t *testing.T) {
	cases := []struct {
		name         string
		event        event
		expectedBody string
	}{
		{
			name:         "event with data",
			event:        event{data: []byte("output")},
			expectedBody: fmt.Sprintf("data: output\n\n"),
		},
		{
			name:         "event with name and data",
			event:        event{eventType: "name", data: []byte("output")},
			expectedBody: fmt.Sprintf("event: name\ndata: output\n\n"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, cancel := context.WithCancel(context.Background())
			ch := make(chan event)

			done := make(chan bool, 1)

			go func() {
				stream(ctx, w, ch)
				done <- true
			}()

			ch <- tc.event
			cancel()

			<-done

			resp := w.Result()
			defer resp.Body.Close()
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
	ch := make(chan event, 1)

	go stream(ctx, w, ch)
	ch <- event{data: []byte("output")}

	<-w.writeCh

	assert.Equal(t, http.StatusInternalServerError, w.statusCode)
}
