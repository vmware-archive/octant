package astilectron

import "sync"

// displayPool represents a display pool
type displayPool struct {
	d map[int64]*Display
	m *sync.Mutex
}

// newDisplayPool creates a new display pool
func newDisplayPool() *displayPool {
	return &displayPool{
		d: make(map[int64]*Display),
		m: &sync.Mutex{},
	}
}

// all returns all the displays
func (p *displayPool) all() (ds []*Display) {
	p.m.Lock()
	defer p.m.Unlock()
	ds = []*Display{}
	for _, d := range p.d {
		ds = append(ds, d)
	}
	return
}

// primary returns the primary display
// It defaults to the last display
func (p *displayPool) primary() (d *Display) {
	p.m.Lock()
	defer p.m.Unlock()
	for _, d = range p.d {
		if d.primary {
			return
		}
	}
	return
}

// update updates the pool based on event displays
func (p *displayPool) update(e *EventDisplays) {
	p.m.Lock()
	defer p.m.Unlock()
	var ids = make(map[int64]bool)
	for _, o := range e.All {
		ids[*o.ID] = true
		var primary bool
		if *o.ID == *e.Primary.ID {
			primary = true
		}
		if d, ok := p.d[*o.ID]; ok {
			d.primary = primary
			*d.o = *o
		} else {
			p.d[*o.ID] = newDisplay(o, primary)
		}
	}
	for id := range p.d {
		if _, ok := ids[id]; !ok {
			delete(p.d, id)
		}
	}
}
