package astilectron

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"

	"github.com/asticode/go-astikit"
)

// reader represents an object capable of reading in the TCP server
type reader struct {
	ctx context.Context
	d   *dispatcher
	l   astikit.SeverityLogger
	r   io.ReadCloser
}

// newReader creates a new reader
func newReader(ctx context.Context, l astikit.SeverityLogger, d *dispatcher, r io.ReadCloser) *reader {
	return &reader{
		ctx: ctx,
		d:   d,
		l:   l,
		r:   r,
	}
}

// close closes the reader properly
func (r *reader) close() error {
	return r.r.Close()
}

// isEOFErr checks whether the error is an EOF error
// wsarecv is the error sent on Windows when the client closes its connection
func (r *reader) isEOFErr(err error) bool {
	return err == io.EOF || strings.Contains(strings.ToLower(err.Error()), "wsarecv:")
}

// read reads from stdout
func (r *reader) read() {
	var reader = bufio.NewReader(r.r)
	for {
		// Check context error
		if r.ctx.Err() != nil {
			return
		}

		// Read next line
		var b []byte
		var err error
		if b, err = reader.ReadBytes('\n'); err != nil {
			if !r.isEOFErr(err) {
				r.l.Errorf("%s while reading", err)
				continue
			}
			return
		}
		b = bytes.TrimSpace(b)
		r.l.Debugf("Astilectron says: %s", b)

		// Unmarshal
		var e Event
		if err = json.Unmarshal(b, &e); err != nil {
			r.l.Errorf("%s while unmarshaling %s", err, b)
			continue
		}

		// Dispatch
		r.d.dispatch(e)
	}
}
