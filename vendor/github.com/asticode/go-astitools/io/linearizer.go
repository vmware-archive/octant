package astiio

import (
	"context"
	"io"
	"sync"
	"time"
)

// Linearizer represents an object capable of linearizing data coming from an io.Reader as packets
type Linearizer struct {
	bufferSize int
	bytesPool  sync.Pool
	cancel     context.CancelFunc
	ctx        context.Context
	events     []*event
	eventsSize int
	md         sync.Mutex // Lock dispatch
	me         sync.Mutex // Locks events
	mr         sync.Mutex // Locks read
	r          io.Reader
}

// event represents an event
type event struct {
	b   []byte
	err error
	n   int
	p   int
}

// NewLinearizer creates a new linearizer that will read readSize bytes at each iteration, write it in its internal
// buffer capped at bufferSize bytes and allow reading this linearized data.
func NewLinearizer(ctx context.Context, r io.Reader, readSize, bufferSize int) (l *Linearizer) {
	l = &Linearizer{
		bufferSize: bufferSize,
		bytesPool:  sync.Pool{New: func() interface{} { return make([]byte, readSize) }},
		r:          r,
	}
	l.ctx, l.cancel = context.WithCancel(ctx)
	return
}

// Close implements the io.Closer interface
func (l *Linearizer) Close() error {
	l.cancel()
	if c, ok := l.r.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

// Start reads the reader and dispatches events accordingly
func (l *Linearizer) Start() {
	for {
		// Check context error
		if l.ctx.Err() != nil {
			l.md.Lock()
			l.dispatchEvent(&event{err: l.ctx.Err()})
			return
		}

		// Get bytes from pool
		var b = l.bytesPool.Get().([]byte)

		// Read
		n, err := l.r.Read(b)
		if err != nil {
			l.md.Lock()
			l.dispatchEvent(&event{err: err})
			return
		}

		// Dispatch event in a go routine so that it doesn't block the read
		l.md.Lock()
		go l.dispatchEvent(&event{b: b, n: n})
	}
}

// dispatchEvent dispatches an event if it doesn't make the buffer overflow based on the bufferSize
// Assumption is made that l.md is locked
func (l *Linearizer) dispatchEvent(e *event) {
	defer l.md.Unlock()
	l.me.Lock()
	defer l.me.Unlock()
	if e.n+l.eventsSize > l.bufferSize {
		return
	}
	l.events = append(l.events, e)
	l.eventsSize += e.n
}

// Read implements the io.Reader interface
func (l *Linearizer) Read(b []byte) (n int, err error) {
	// Only one read at a time
	l.mr.Lock()
	defer l.mr.Unlock()

	// Loop until there's enough data to read or there's an error
	for {
		// Check events size
		l.me.Lock()
		if l.eventsSize >= len(b) {
			// Loop in events
			for idx := 0; idx < len(l.events); idx++ {
				// Copy bytes
				e := l.events[idx]
				if len(b)-n < e.n-e.p {
					// Copy a part of the bytes
					copy(b[n:], e.b[e.p:len(b)-n+e.p])
					e.p += len(b) - n
					n += len(b) - n
				} else {
					// Copy the remainder of the bytes
					copy(b[n:], e.b[e.p:e.n])
					n += e.n - e.p

					// Put bytes back in pool
					l.bytesPool.Put(e.b)

					// Remove event
					l.events = append(l.events[:idx], l.events[idx+1:]...)
					idx--
				}

				// All the bytes have been read
				if n == len(b) {
					break
				}
			}
			l.eventsSize -= n
			l.me.Unlock()
			return
		}

		// Process error in last event
		if len(l.events) > 0 {
			if err = l.events[len(l.events)-1].err; err != nil {
				l.me.Unlock()
				return
			}
		}
		l.me.Unlock()

		// Wait for 1ms
		time.Sleep(time.Millisecond)
	}
	return
}
