package astilectron

import (
	"strconv"
	"sync"
)

// identifier is in charge of delivering a unique identifier
type identifier struct {
	i int
	m *sync.Mutex
}

// newIdentifier creates a new identifier
func newIdentifier() *identifier {
	return &identifier{m: &sync.Mutex{}}
}

// new returns a new unique identifier
func (i *identifier) new() string {
	i.m.Lock()
	defer i.m.Unlock()
	i.i++
	return strconv.Itoa(i.i)
}
