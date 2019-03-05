package overview

import (
	"context"
	"sync"

	"github.com/heptio/developer-dash/internal/log"
	"github.com/pkg/errors"
)

var (
	errPathNotFound = errors.New("path not found")
)

type pathMatcher struct {
	filters map[string]pathFilter

	sync.Mutex
}

func newPathMatcher() *pathMatcher {
	return &pathMatcher{
		filters: make(map[string]pathFilter),
	}
}

func (pm *pathMatcher) Register(ctx context.Context, pf pathFilter) {
	logger := log.From(ctx)

	pm.Lock()
	defer pm.Unlock()

	logger.Debugf("register path %s", pf.path)
	pm.filters[pf.path] = pf
}

func (pm *pathMatcher) Deregister(ctx context.Context, paths ...string) {
	logger := log.From(ctx)

	pm.Lock()
	defer pm.Unlock()

	for _, p := range paths {
		logger.Debugf("deregister path %s", p)
		delete(pm.filters, p)
	}

}

func (pm *pathMatcher) Find(path string) (pathFilter, error) {
	pm.Lock()
	defer pm.Unlock()

	for _, pf := range pm.filters {
		if pf.Match(path) {
			return pf, nil
		}
	}

	return pathFilter{}, errPathNotFound
}
