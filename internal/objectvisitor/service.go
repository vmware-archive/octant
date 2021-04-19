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

	"github.com/vmware-tanzu/octant/internal/gvk"
	"github.com/vmware-tanzu/octant/internal/queryer"
	"github.com/vmware-tanzu/octant/internal/util/kubernetes"
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
func (s *Service) Visit(ctx context.Context, object *unstructured.Unstructured, handler ObjectHandler, visitor Visitor, visitDescendants bool, level int) error {
	ctx, span := trace.StartSpan(ctx, "visitService")
	defer span.End()

	service := &corev1.Service{}
	if err := kubernetes.FromUnstructured(object, service); err != nil {
		return err
	}
	level = handler.SetLevel(service.Kind, level)

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
				if err := visitor.Visit(ctx, u, handler, visitDescendants, level); err != nil {
					return err
				}
				return handler.AddEdge(ctx, object, u, level)
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
				if visitDescendants {
					if err := visitor.Visit(ctx, u, handler, false, level); err != nil {
						return errors.Wrapf(err, "service %s visit ingress %s",
							kubernetes.PrintObject(service), kubernetes.PrintObject(ingress))
					}
				}

				return handler.AddEdge(ctx, object, u, level)
			})
		}

		return nil
	})

	g.Go(func() error {
		apiservices, err := s.queryer.APIServicesForService(ctx, service)
		if err != nil {
			return err
		}

		for i := range apiservices {
			apiservice := apiservices[i]
			g.Go(func() error {
				m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(apiservice)
				if err != nil {
					return err
				}
				u := &unstructured.Unstructured{Object: m}
				if visitDescendants {
					if err := visitor.Visit(ctx, u, handler, false, level); err != nil {
						return err
					}
				}

				return handler.AddEdge(ctx, object, u, level)
			})

		}

		return nil
	})

	g.Go(func() error {
		mutatingwebhookconfigurations, err := s.queryer.MutatingWebhookConfigurationsForService(ctx, service)
		if err != nil {
			return err
		}

		for i := range mutatingwebhookconfigurations {
			mutatingwebhookconfiguration := mutatingwebhookconfigurations[i]
			g.Go(func() error {
				m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(mutatingwebhookconfiguration)
				if err != nil {
					return err
				}
				u := &unstructured.Unstructured{Object: m}
				if visitDescendants {
					if err := visitor.Visit(ctx, u, handler, false, level); err != nil {
						return err
					}
				}

				return handler.AddEdge(ctx, object, u, level)
			})

		}

		return nil
	})

	g.Go(func() error {
		validatingwebhookconfigurations, err := s.queryer.ValidatingWebhookConfigurationsForService(ctx, service)
		if err != nil {
			return err
		}

		for i := range validatingwebhookconfigurations {
			validatingwebhookconfiguration := validatingwebhookconfigurations[i]
			g.Go(func() error {
				m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(validatingwebhookconfiguration)
				if err != nil {
					return err
				}
				u := &unstructured.Unstructured{Object: m}
				if visitDescendants {
					if err := visitor.Visit(ctx, u, handler, false, level); err != nil {
						return err
					}
				}

				return handler.AddEdge(ctx, object, u, level)
			})

		}

		return nil
	})

	return g.Wait()
}
