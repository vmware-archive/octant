package resourceviewer

import (
	"context"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/pkg/objectstoreutil"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// VisitCache stores components returned from ResoruceViewer.Visit
type VisitCache interface {
	Prime(context.Context, objectstoreutil.Key, *ResourceViewer, runtime.Object) (component.Component, error)
	Get(context.Context, objectstoreutil.Key, runtime.Object) (component.Component, bool)
}

type visitCache struct {
	visitValue  *lru.Cache
	visitViewer *lru.Cache
}

// NewVisitCache creates a VisitCache
func NewVisitCache(size int) (VisitCache, error) {
	valueCache, err := lru.New(size)
	if err != nil {
		return nil, err
	}

	viewerCache, err := lru.New(size)
	if err != nil {
		return nil, err
	}

	return &visitCache{
		visitValue:  valueCache,
		visitViewer: viewerCache,
	}, nil
}

// Get retrieves the value for the key provided and an ok value
func (cc *visitCache) Get(ctx context.Context, key objectstoreutil.Key, object runtime.Object) (component.Component, bool) {
	logger := log.From(ctx)

	_, ok := cc.visitViewer.Get(key)
	if !ok {
		return nil, false
	}

	c, err := cc.update(ctx, key, object)
	if err != nil {
		logger.Errorf("failed to get Component: %s", err)
		return nil, false
	}
	return c, true
}

func (cc *visitCache) update(ctx context.Context, key objectstoreutil.Key, object runtime.Object) (component.Component, error) {
	v, ok := cc.visitViewer.Get(key)
	if !ok {
		return nil, errors.New("no resourceviewer to update")
	}
	rv := v.(*ResourceViewer)

	done := make(chan objectstoreutil.Key, 1)
	errChan := make(chan error, 1)

	go func() {
		defer close(done)
		defer close(errChan)

		component, err := rv.Visit(ctx, object)
		if err != nil {
			errChan <- err
		} else {
			cc.visitValue.Add(key, component)
			done <- key
		}
	}()

	getComponent := func(key objectstoreutil.Key) (component.Component, error) {
		componentValue, ok := cc.visitValue.Get(key)
		if !ok {
			return nil, errors.New("error loading value from cache")
		}
		return componentValue.(component.Component), nil
	}

	select {
	case err := <-errChan:
		if err != nil {
			return nil, err
		}
	case keyValue := <-done:
		return getComponent(keyValue)
	case <-time.After(750 * time.Millisecond):
		return getComponent(key)
	}

	return nil, errors.New("select failed")
}

// Prime attempts to retreive an item from the cache, but unlike Get, it will launch ResourceViwer.Visit if no value is found.
// and when Visit completes set the value in the cache.
func (cc *visitCache) Prime(ctx context.Context, key objectstoreutil.Key, rv *ResourceViewer, object runtime.Object) (component.Component, error) {
	_, ok := cc.visitViewer.Get(key)
	if !ok {
		cc.visitViewer.Add(key, rv)
		fakeComponent, err := rv.FakeVisit(ctx, object)
		if err != nil {
			return nil, err
		}
		cc.visitValue.Add(key, fakeComponent)
		return cc.update(ctx, key, object)
	}
	c, _ := cc.visitValue.Get(key)
	return c.(component.Component), nil
}
