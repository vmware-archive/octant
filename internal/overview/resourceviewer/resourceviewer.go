package resourceviewer

import (
	"context"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"

	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/objectstore"
	"github.com/heptio/developer-dash/internal/overview/objectvisitor"
	"github.com/heptio/developer-dash/internal/queryer"
	"github.com/heptio/developer-dash/pkg/objectstoreutil"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

//go:generate mockgen -source=resourceviewer.go -destination=./fake/mock_component_cache.go -package=fake github.com/heptio/developer-dash/internal/overview/resourceviewer ComponentCache

// ViewerOpt is an option for ResourceViewer.
type ViewerOpt func(*ResourceViewer) error

// WithDefaultQueryer configures ResourceViewer with the default visitor.
func WithDefaultQueryer(q queryer.Queryer) ViewerOpt {
	return func(rv *ResourceViewer) error {
		visitor, err := objectvisitor.NewDefaultVisitor(q, rv.factoryFunc())
		if err != nil {
			return err
		}

		rv.visitor = visitor
		return nil
	}
}

// ResourceViewer visits an object and creates a view component.
type ResourceViewer struct {
	collector *Collector
	visitor   objectvisitor.Visitor
}

// New creates an instance of ResourceViewer.
func New(logger log.Logger, o objectstore.ObjectStore, opts ...ViewerOpt) (*ResourceViewer, error) {
	rv := &ResourceViewer{
		collector: NewCollector(o),
	}

	rv.collector.logger = logger

	for _, opt := range opts {
		if err := opt(rv); err != nil {
			return nil, errors.Wrap(err, "invalid resource viewer option")
		}
	}

	if rv.visitor == nil {
		return nil, errors.New("resource viewer visitor is nil")
	}

	return rv, nil
}

// Visit visits an object and creates a view component.
func (rv *ResourceViewer) Visit(ctx context.Context, object objectvisitor.ClusterObject) (component.Component, error) {
	rv.collector.Reset()

	ctx, span := trace.StartSpan(ctx, "resourceviewer")
	defer span.End()

	if err := rv.visitor.Visit(ctx, object); err != nil {
		return nil, err
	}

	accessor := meta.NewAccessor()
	uid, err := accessor.UID(object)
	if err != nil {
		return nil, err
	}

	return rv.collector.Component(string(uid))
}

// EmptyVisit returns a component that has not been visited yet.
// Use EmptyVisit when you are running Visit in a goroutine and want to return a component quickly.
func (rv *ResourceViewer) EmptyVisit(ctx context.Context, object objectvisitor.ClusterObject) (component.Component, error) {
	ctx, span := trace.StartSpan(ctx, "resourceviewer")
	defer span.End()

	accessor := meta.NewAccessor()
	name, err := accessor.Name(object)
	if err != nil {
		return nil, err
	}

	emptyNode := component.Node{
		Name:       name,
		APIVersion: "Loading",
		Kind:       "...",
		Status:     "ok",
	}

	r := component.NewResourceViewer("Resource Viewer")
	r.AddNode("emptyID", emptyNode)
	return r, nil
}

func (rv *ResourceViewer) factoryFunc() objectvisitor.ObjectHandlerFactory {
	return func(object objectvisitor.ClusterObject) (objectvisitor.ObjectHandler, error) {
		return rv.collector, nil
	}
}

// ComponentCache is cache of Components
type ComponentCache interface {
	Get(context.Context, runtime.Object) (component.Component, error)
	SetQueryer(queryer.Queryer)
}

type componentCache struct {
	components *lru.Cache
	queryer    queryer.Queryer
	store      objectstore.ObjectStore

	mu sync.Mutex
}

// NewComponentCache creates a new component cache.
func NewComponentCache(o objectstore.ObjectStore) (ComponentCache, error) {
	components, err := lru.New(100)
	if err != nil {
		return nil, err
	}

	return &componentCache{
		components: components,
		store:      o,
	}, nil
}

// SetQueryer sets the queryer for the component cache.
func (cc *componentCache) SetQueryer(q queryer.Queryer) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	cc.queryer = q
}

// Get creates a Resource Viewer and begins starts the Visit routine. After waiting a set amount of
// time, Get returns the Component that is in the component cache. This value may or may not be the value
// from the last Visit call. If the component cache is empty we return a Component from FakeVisit.
func (cc *componentCache) Get(ctx context.Context, object runtime.Object) (component.Component, error) {
	key, err := objectstoreutil.KeyFromObject(object)
	if err != nil {
		return nil, err
	}

	if cc.queryer == nil {
		return nil, errors.New("no queryer set")
	}

	rv, err := cc.newResourceViewer(ctx)
	if err != nil {
		return nil, err
	}

	done, errChan := cc.visit(ctx, key, object, rv)

	select {
	case err := <-errChan:
		if err != nil {
			return nil, err
		}
	case keyValue := <-done:
		return cc.getComponent(ctx, keyValue, object, rv)
	case <-time.After(750 * time.Millisecond):
		return cc.getComponent(ctx, key, object, rv)
	}
	return nil, errors.New("bad")
}

func (cc *componentCache) getComponent(
	ctx context.Context,
	key objectstoreutil.Key,
	object runtime.Object,
	rv *ResourceViewer,
) (component.Component, error) {
	componentValue, ok := cc.components.Get(key)
	if !ok {
		return rv.EmptyVisit(ctx, object)
	}
	return componentValue.(component.Component), nil
}

func (cc *componentCache) visit(
	ctx context.Context,
	key objectstoreutil.Key,
	object runtime.Object,
	rv *ResourceViewer,
) (chan objectstoreutil.Key, chan error) {
	done := make(chan objectstoreutil.Key, 1)
	errChan := make(chan error, 1)

	go func() {
		defer close(done)
		defer close(errChan)

		rvComponent, err := rv.Visit(ctx, object)
		if err != nil {
			errChan <- err
		} else {
			cc.components.Add(key, rvComponent)
			done <- key
		}
	}()

	return done, errChan
}

func (cc *componentCache) newResourceViewer(ctx context.Context) (*ResourceViewer, error) {
	logger := log.From(ctx)
	return New(logger, cc.store, WithDefaultQueryer(cc.queryer))
}
