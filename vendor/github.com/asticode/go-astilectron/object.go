package astilectron

import (
	"context"
)

// object represents a base object
type object struct {
	cancel context.CancelFunc
	ctx    context.Context
	d      *dispatcher
	i      *identifier
	id     string
	w      *writer
}

// newObject returns a new base object
func newObject(ctx context.Context, d *dispatcher, i *identifier, w *writer, id string) (o *object) {
	o = &object{
		d:  d,
		i:  i,
		id: id,
		w:  w,
	}
	o.ctx, o.cancel = context.WithCancel(ctx)
	return
}

// On implements the Listenable interface
func (o *object) On(eventName string, l Listener) {
	o.d.addListener(o.id, eventName, l)
}
