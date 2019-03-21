package resourceviewer

import (
	"context"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/overview/objectvisitor"
	"github.com/heptio/developer-dash/internal/queryer"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"k8s.io/apimachinery/pkg/api/meta"
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
func New(logger log.Logger, c cache.Cache, opts ...ViewerOpt) (*ResourceViewer, error) {
	rv := &ResourceViewer{
		collector: NewCollector(c),
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

func (rv *ResourceViewer) factoryFunc() objectvisitor.ObjectHandlerFactory {
	return func(object objectvisitor.ClusterObject) (objectvisitor.ObjectHandler, error) {
		return rv.collector, nil
	}
}
