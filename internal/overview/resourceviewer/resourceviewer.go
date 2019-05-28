package resourceviewer

import (
	"context"

	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/objectstore"
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

// Component calls Component for the ResourveViewer collector
func (rv *ResourceViewer) Component(ctx context.Context, object objectvisitor.ClusterObject) (component.Component, error) {
	accessor := meta.NewAccessor()
	uid, err := accessor.UID(object)
	if err != nil {
		return nil, err
	}
	return rv.collector.Component(string(uid))
}

// Visit visits an object and creates a view component.
func (rv *ResourceViewer) Visit(ctx context.Context, object objectvisitor.ClusterObject) (component.Component, error) {
	rv.collector.Reset()
	rv.visitor.Reset()

	ctx, span := trace.StartSpan(ctx, "resourceviewer")
	defer span.End()

	if err := rv.visitor.Visit(ctx, object); err != nil {
		return nil, err
	}

	return rv.Component(ctx, object)
}

// FakeVisit returns a component that has not been visited yet.
// Use FakeVisit when you are running Visit in a goroutine and want to return a component quickly.
func (rv *ResourceViewer) FakeVisit(ctx context.Context, object objectvisitor.ClusterObject) (component.Component, error) {
	ctx, span := trace.StartSpan(ctx, "resourceviewer")
	defer span.End()

	accessor := meta.NewAccessor()
	name, err := accessor.Name(object)
	if err != nil {
		return nil, err
	}

	fakeNode := component.Node{
		Name:       name,
		APIVersion: "Loading",
		Kind:       "...",
		Status:     "ok",
	}

	r := component.NewResourceViewer("Resource Viewer")
	r.AddNode("fakeID", fakeNode)
	return r, nil
}

func (rv *ResourceViewer) factoryFunc() objectvisitor.ObjectHandlerFactory {
	return func(object objectvisitor.ClusterObject) (objectvisitor.ObjectHandler, error) {
		return rv.collector, nil
	}
}
