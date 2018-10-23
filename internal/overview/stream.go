package overview

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/heptio/developer-dash/internal/log"
)

const (
	defaultEventTimeout = 5 * time.Second
)

type streamFn func(ctx context.Context, w http.ResponseWriter, ch chan []byte)

type contentStreamer struct {
	generator    generator
	w            http.ResponseWriter
	path         string
	prefix       string
	namespace    string
	streamFn     streamFn
	eventTimeout time.Duration
	logger       log.Logger
}

func (cs contentStreamer) content(ctx context.Context) {
	ch := make(chan []byte, 1)

	timer := time.NewTimer(0)

	go func() {
		isRunning := true
		for isRunning {
			select {
			case <-ctx.Done():
				isRunning = false
			case <-timer.C:
				title, contents, err := cs.generator.Generate(cs.path, cs.prefix, cs.namespace)
				if err != nil {
					cs.logger.Errorf("generate error: %v", err)
				}

				cr := &contentResponse{
					Title:    title,
					Contents: contents,
				}

				data, err := json.Marshal(cr)
				if err != nil {
					cs.logger.Errorf("marshal err: %v", err)
				}

				ch <- data

				timer.Reset(cs.eventTimeout)
			}
		}
	}()

	cs.streamFn(ctx, cs.w, ch)
}

func stream(ctx context.Context, w http.ResponseWriter, ch chan []byte) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "server sent events are unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	isStreaming := true

	for isStreaming {
		select {
		case <-ctx.Done():
			isStreaming = false
		case msg := <-ch:
			fmt.Fprintf(w, "data: %s\n\n", string(msg))
			flusher.Flush()
		}
	}
}
