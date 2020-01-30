package objectstore

import (
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/vmware-tanzu/octant/pkg/store"
)

var defaultBackoff = wait.Backoff{
	Duration: 1 * time.Second,
	Factor:   2.0,
	Jitter:   0.1,
	Steps:    50,
	Cap:      10 * time.Minute,
}

type backoffer interface {
	isWaiting() bool
	setWaiting(bool)
	wait() time.Duration
}

var _ backoffer = (*backoffEntry)(nil)

type backoffEntry struct {
	key     store.Key
	waiting bool
	backoff wait.Backoff

	mu sync.RWMutex
}

func newBackoffEntry(key store.Key, backoff wait.Backoff) *backoffEntry {
	return &backoffEntry{
		key:     key,
		waiting: false,
		backoff: backoff,
	}
}
func (b *backoffEntry) isWaiting() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.waiting
}

func (b *backoffEntry) setWaiting(waiting bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.waiting = waiting
}

func (b *backoffEntry) wait() time.Duration {
	return b.backoff.Step()
}
