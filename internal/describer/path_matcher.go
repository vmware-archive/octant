package describer

import (
	"context"
	"sync"

	"github.com/heptio/developer-dash/internal/log"
	"github.com/pkg/errors"
)

var (
	ErrPathNotFound = errors.New("path not found")
)

type PathMatcher struct {
	filters map[string]PathFilter

	sync.Mutex
}

func NewPathMatcher() *PathMatcher {
	return &PathMatcher{
		filters: make(map[string]PathFilter),
	}
}

func (pm *PathMatcher) Register(ctx context.Context, pf PathFilter) {
	logger := log.From(ctx)

	pm.Lock()
	defer pm.Unlock()

	logger.With("path", pf.path).Debugf("register path")
	pm.filters[pf.path] = pf
}

func (pm *PathMatcher) Deregister(ctx context.Context, paths ...string) {
	logger := log.From(ctx)

	pm.Lock()
	defer pm.Unlock()

	for _, p := range paths {
		logger.With("path", p).Debugf("deregister path")
		delete(pm.filters, p)
	}

}

func (pm *PathMatcher) Find(path string) (PathFilter, error) {
	pm.Lock()
	defer pm.Unlock()

	for _, pf := range pm.filters {
		if pf.Match(path) {
			return pf, nil
		}
	}

	return PathFilter{}, ErrPathNotFound
}
