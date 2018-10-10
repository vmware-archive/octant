package overview

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
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
				contents, err := cs.generator.Generate(cs.path, cs.prefix, cs.namespace)
				if err != nil {
					log.Printf("generate error: %v", err)
				}

				cr := &contentResponse{
					Contents: contents,
				}

				data, err := json.Marshal(cr)
				if err != nil {
					log.Printf("marshal err: %v", err)
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
