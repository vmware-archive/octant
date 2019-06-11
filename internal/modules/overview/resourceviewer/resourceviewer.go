package resourceviewer

import "C"
import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"k8s.io/apimachinery/pkg/api/meta"

	"github.com/heptio/developer-dash/internal/componentcache"
	"github.com/heptio/developer-dash/internal/config"
	"github.com/heptio/developer-dash/internal/modules/overview/objectvisitor"
	"github.com/heptio/developer-dash/internal/queryer"
	"github.com/heptio/developer-dash/pkg/objectstoreutil"
	"github.com/heptio/developer-dash/pkg/view/component"
)

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
func New(dashConfig config.Dash, opts ...ViewerOpt) (*ResourceViewer, error) {
	collector, err := NewCollector(dashConfig)
	if err != nil {
		return nil, errors.Wrap(err, "create collector")
	}
	rv := &ResourceViewer{
		collector: collector,
	}

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

func (rv *ResourceViewer) factoryFunc() objectvisitor.ObjectHandlerFactory {
	return func(object objectvisitor.ClusterObject) (objectvisitor.ObjectHandler, error) {
		return rv.collector, nil
	}
}

// CachedResourceViewer returns a RV component from the componentcache and starts a new visit.
func CachedResourceViewer(ctx context.Context, object objectvisitor.ClusterObject, dashConfig config.Dash, q queryer.Queryer) componentcache.UpdateFn {
	return func(ctx context.Context, cacheChan chan componentcache.Event) (string, error) {
		var event componentcache.Event
		event.Name = "Resource Viewer"

		copyObject := object.DeepCopyObject()

		key, err := objectstoreutil.KeyFromObject(copyObject)
		if err != nil {
			return "", err
		}
		sKey := fmt.Sprintf("%s-%s", "resourceviewer", key.String())
		event.Key = sKey

		componentCache := dashConfig.ComponentCache()
		if _, ok := componentCache.Get(sKey); !ok {
			title := component.Title(component.NewText("Resource Viewer"))
			loading := component.NewLoading(title, "Resource Viewer")
			componentCache.Add(sKey, loading)
		}

		rv, err := New(dashConfig, WithDefaultQueryer(q))
		if err != nil {
			return sKey, err
		}

		go func() {
			c, err := rv.Visit(ctx, copyObject)
			event.Err = err
			event.CComponent = c
			cacheChan <- event
		}()

		return sKey, nil
	}
}
