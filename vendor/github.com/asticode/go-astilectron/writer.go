package astilectron

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/asticode/go-astikit"
)

// writer represents an object capable of writing in the TCP server
type writer struct {
	l astikit.SeverityLogger
	w io.WriteCloser
}

// newWriter creates a new writer
func newWriter(w io.WriteCloser, l astikit.SeverityLogger) *writer {
	return &writer{
		l: l,
		w: w,
	}
}

// close closes the writer properly
func (w *writer) close() error {
	return w.w.Close()
}

// write writes to the stdin
func (w *writer) write(e Event) (err error) {
	// Marshal
	var b []byte
	if b, err = json.Marshal(e); err != nil {
		return fmt.Errorf("marshaling %+v failed: %w", e, err)
	}

	// Write
	w.l.Debugf("Sending to Astilectron: %s", b)
	if _, err = w.w.Write(append(b, '\n')); err != nil {
		return fmt.Errorf("writing %s failed: %w", b, err)
	}
	return
}
