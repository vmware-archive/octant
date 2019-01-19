package astiio

import (
	"context"
	"io"
)

// Reader represents a reader with a context
type Reader struct {
	ctx    context.Context
	reader io.Reader
}

// NewReader creates a new Reader
func NewReader(ctx context.Context, r io.Reader) *Reader {
	return &Reader{
		ctx:    ctx,
		reader: r,
	}
}

// Read allows Reader to implement the io.Reader interface
func (r *Reader) Read(p []byte) (n int, err error) {
	// Check context
	if err = r.ctx.Err(); err != nil {
		return
	}

	// Read
	return r.reader.Read(p)
}
