package astilectron

import (
	"encoding/json"
	"io"

	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

// writer represents an object capable of writing in the TCP server
type writer struct {
	w io.WriteCloser
}

// newWriter creates a new writer
func newWriter(w io.WriteCloser) *writer {
	return &writer{
		w: w,
	}
}

// close closes the writer properly
func (r *writer) close() error {
	return r.w.Close()
}

// write writes to the stdin
func (r *writer) write(e Event) (err error) {
	// Marshal
	var b []byte
	if b, err = json.Marshal(e); err != nil {
		return errors.Wrapf(err, "Marshaling %+v failed", e)
	}

	// Write
	astilog.Debugf("Sending to Astilectron: %s", b)
	if _, err = r.w.Write(append(b, '\n')); err != nil {
		return errors.Wrapf(err, "Writing %s failed", b)
	}
	return
}
