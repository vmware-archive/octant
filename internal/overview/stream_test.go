package overview

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_contentStreamer(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, cancel := context.WithCancel(context.Background())
	g := newStubbedGenerator([]content.Content{newFakeContent(false)}, nil)

	rcv := make(chan bool, 1)

	fn := func(ctx context.Context, w http.ResponseWriter, ch chan []byte) {
		msg := <-ch

		assert.Equal(t, `{"contents":[{"type":"stubbed"}],"title":"title"}`, string(msg))
		rcv <- true
	}

	cs := contentStreamer{
		generator: g,
		w:         w,
		path:      "/real/foo",
		prefix:    "/real",
		namespace: "default",
		streamFn:  fn,
	}

	cs.content(ctx)

	<-rcv
	cancel()
}

func Test_stream(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan []byte)

	go stream(ctx, w, ch)

	ch <- []byte("output")

	resp := w.Result()
	defer resp.Body.Close()
	actualBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	expectedBody := fmt.Sprintf("data: output\n\n")
	assert.Equal(t, expectedBody, string(actualBody))

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

	cancel()
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
	ch := make(chan []byte, 1)

	go stream(ctx, w, ch)
	ch <- []byte("output")

	<-w.writeCh

	assert.Equal(t, http.StatusInternalServerError, w.statusCode)
}
