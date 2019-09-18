package objectvisitor

import (
	"context"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware/octant/internal/gvk"
	"github.com/vmware/octant/internal/queryer"
	"github.com/vmware/octant/internal/util/kubernetes"
)

// Service is a typed visitor for services.
type Service struct {
	queryer queryer.Queryer
}

var _ TypedVisitor = (*Service)(nil)

// NewService creates an instance of Service.
func NewService(q queryer.Queryer) *Service {
	return &Service{queryer: q}
}

// Supports returns the gvk this typed visitor supports.
func (Service) Supports() schema.GroupVersionKind {
	return gvk.Service
}

// Visit visits a service. It looks for associated pods and ingresses.
func (s *Service) Visit(ctx context.Context, object *unstructured.Unstructured, handler ObjectHandler, visitor Visitor, visitDescendants bool) error {
	ctx, span := trace.StartSpan(ctx, "visitService")
	defer span.End()

	service := &corev1.Service{}
	if err := convertToType(object, service); err != nil {
		return err
	}

	var g errgroup.Group

	g.Go(func() error {
		pods, err := s.queryer.PodsForService(ctx, service)
		if err != nil {
			return err
		}

		for i := range pods {
			pod := pods[i]
			g.Go(func() error {
				m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(pod)
				if err != nil {
					return err
				}
				u := &unstructured.Unstructured{Object: m}
				if err := visitor.Visit(ctx, u, handler, true); err != nil {
					return err
				}

				return handler.AddEdge(ctx, object, u)
			})

		}

		return nil
	})

	g.Go(func() error {
		ingresses, err := s.queryer.IngressesForService(ctx, service)
		if err != nil {
			return err
		}

		for i := range ingresses {
			ingress := ingresses[i]
			g.Go(func() error {
				m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(ingress)
				if err != nil {
					return err
				}
				u := &unstructured.Unstructured{Object: m}
				if err := visitor.Visit(ctx, u, handler, true); err != nil {
					return errors.Wrapf(err, "service %s visit ingress %s",
						kubernetes.PrintObject(service), kubernetes.PrintObject(ingress))
				}

				return handler.AddEdge(ctx, object, u)
			})
		}

		return nil
	})

	return g.Wait()
}
