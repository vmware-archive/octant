package objectvisitor

import (
	"context"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
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
	return gvk.ServiceGVK
}

// Visit visits a service. It looks for associated pods and ingresses.
func (s *Service) Visit(ctx context.Context, object runtime.Object, handler ObjectHandler, visitor Visitor) error {
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
				return errors.Wrapf(visitor.Visit(ctx, pod, handler), "services %s visit pod %s",
					kubernetes.PrintObject(service), kubernetes.PrintObject(pod))
			})

			if err := handler.AddEdge(object, pod); err != nil {
				return err
			}
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
				if err := visitor.Visit(ctx, ingress, handler); err != nil {
					return errors.Wrapf(err, "service %s visit ingress %s",
						kubernetes.PrintObject(service), kubernetes.PrintObject(ingress))
				}

				return handler.AddEdge(object, ingress)
			})
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return handler.Process(ctx, object)
}
