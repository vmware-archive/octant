package asticontext

import (
	"context"
	"sync"
)

// Canceller represents a context with a mutex
type Canceller struct {
	ctx    context.Context
	cancel context.CancelFunc
	mutex  *sync.RWMutex
}

// NewCanceller returns a new canceller
func NewCanceller() (c *Canceller) {
	c = &Canceller{mutex: &sync.RWMutex{}}
	c.Reset()
	return
}

// Cancel cancels the canceller context
func (c *Canceller) Cancel() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cancel()
}

// Cancelled returns whether the canceller has cancelled the context
func (c *Canceller) Cancelled() bool {
	return c.ctx.Err() != nil
}

// Lock locks the canceller mutex
func (c *Canceller) Lock() {
	c.mutex.Lock()
}

// Lock locks the canceller mutex
func (c *Canceller) NewContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(c.ctx)
}

// Reset resets the canceller context
func (c *Canceller) Reset() {
	c.ctx, c.cancel = context.WithCancel(context.Background())
}

// Unlock unlocks the canceller mutex
func (c *Canceller) Unlock() {
	c.mutex.Unlock()
}
