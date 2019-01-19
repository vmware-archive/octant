package astiexec

import "bytes"

// Vars
var (
	bytesEOL = []byte("\n")
)

// StdWriter represents an object capable of writing what's coming out of stdout or stderr
type StdWriter struct {
	buffer *bytes.Buffer
	fn     func(i []byte)
}

// NewStdWriter creates a new StdWriter
func NewStdWriter(fn func(i []byte)) *StdWriter {
	return &StdWriter{buffer: &bytes.Buffer{}, fn: fn}
}

// Close closes the writer
func (w *StdWriter) Close() {
	if w.buffer.Len() > 0 {
		w.write(w.buffer.Bytes())
	}
}

// Write implements the io.Writer interface
func (w *StdWriter) Write(i []byte) (n int, err error) {
	// Update n to avoid broken pipe error
	defer func() {
		n = len(i)
	}()

	// No EOL in the log, write in buffer
	if bytes.Index(i, bytesEOL) == -1 {
		w.buffer.Write(i)
		return
	}

	// Loop in items split by EOL
	var items = bytes.Split(i, bytesEOL)
	for i := 0; i < len(items)-1; i++ {
		// If first item, add the buffer
		if i == 0 {
			items[i] = append(w.buffer.Bytes(), items[i]...)
			w.buffer.Reset()
		}

		// Log
		w.write(items[i])
	}

	// Add remaining to buffer
	w.buffer.Write(items[len(items)-1])
	return
}

// write writes the input
func (w *StdWriter) write(i []byte) {
	w.fn(i)
}
