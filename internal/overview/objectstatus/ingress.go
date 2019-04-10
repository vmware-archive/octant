package objectstatus

import (
	"context"
	"strconv"

	"github.com/heptio/developer-dash/internal/objectstore"
	"github.com/heptio/developer-dash/pkg/objectstoreutil"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/scheme"
)

const (
	ingressNoBackendsDefined = "No backends defined. All traffic will be sent to the default backend configured for the ingress controller."
)

func runIngressStatus(ctx context.Context, object runtime.Object, o objectstore.ObjectStore) (ObjectStatus, error) {
	if object == nil {
		return ObjectStatus{}, errors.Errorf("ingress is nil")
	}

	ingress := &extv1beta1.Ingress{}

	if err := scheme.Scheme.Convert(object, ingress, 0); err != nil {
		return ObjectStatus{}, errors.Wrap(err, "convert object to ingress")
	}

	is := ingressStatus{
		ingress:     *ingress,
		objectstore: o,
	}
	status, err := is.run(ctx)
	if err != nil {
		return ObjectStatus{}, errors.Wrap(err, "build status for ingress")
	}

	if len(status.Details) == 0 {
		status.AddDetail("Ingress is OK")
	}

	return status, nil
}

type ingressStatus struct {
	ingress     extv1beta1.Ingress
	objectstore objectstore.ObjectStore
}

func (is *ingressStatus) run(ctx context.Context) (ObjectStatus, error) {
	status := ObjectStatus{}

	ingress := is.ingress

	o := is.objectstore
	if o == nil {
		return status, errors.New("ingress status requires a non nil objectstore")
	}

	backends := is.backends()
	if len(backends) == 0 {
		status.SetWarning()
		status.AddDetail(ingressNoBackendsDefined)
	}

	for _, backend := range backends {
		key := objectstoreutil.Key{
			Namespace:  ingress.Namespace,
			APIVersion: "v1",
			Kind:       "Service",
			Name:       backend.ServiceName,
		}

		service := &corev1.Service{}

		if err := objectstore.GetAs(ctx, o, key, service); err != nil {
			return ObjectStatus{}, errors.Wrap(err, "get service from objectstore")
		}

		if service.Name == "" {
			status.SetError()
			status.AddDetailf("Backend refers to service %q which doesn't exist",
				backend.ServiceName)
			continue
		}

		if !matchBackendPort(backend, service.Spec.Ports) {
			status.SetError()
			status.AddDetailf("Backend for service %q specifies an invalid port",
				backend.ServiceName)
			continue
		}
	}

	tlsHosts := is.tlsHostMap()
	if len(tlsHosts) > 0 {
		for _, rule := range ingress.Spec.Rules {
			if rule.Host == "" {
				continue
			}

			if ok := tlsHosts[rule.Host]; !ok {
				status.SetError()
				status.AddDetailf("No matching TLS host for rule %q", rule.Host)
			}
		}
	}

	for _, tls := range ingress.Spec.TLS {
		if tls.SecretName == "" {
			status.SetError()
			status.AddDetail("TLS configuration did not define a secret name")
			continue
		}

		key := objectstoreutil.Key{
			Namespace:  ingress.Namespace,
			APIVersion: "v1",
			Kind:       "Secret",
			Name:       tls.SecretName,
		}

		secret, err := is.objectstore.Get(ctx, key)
		if err != nil {
			// assume an error means the secret couldn't be accessed
			break
		}

		if secret == nil {
			status.SetError()
			status.AddDetailf("Secret %q does not exist", tls.SecretName)
		}
	}

	return status, nil
}

func (is *ingressStatus) backends() []extv1beta1.IngressBackend {
	var list []extv1beta1.IngressBackend

	ingress := is.ingress

	if ingress.Spec.Backend != nil && ingress.Spec.Backend.ServiceName != "" {
		list = append(list, *ingress.Spec.Backend)
	}

	for _, rule := range ingress.Spec.Rules {
		if rule.IngressRuleValue.HTTP == nil {
			continue
		}

		for _, p := range rule.IngressRuleValue.HTTP.Paths {
			if p.Backend.ServiceName == "" {
				continue
			}

			list = append(list, p.Backend)
		}
	}

	return list
}

func (is *ingressStatus) tlsHostMap() map[string]bool {
	result := make(map[string]bool)

	for _, tls := range is.ingress.Spec.TLS {
		for _, host := range tls.Hosts {
			result[host] = true
		}
	}

	return result
}

// matchBackendPort returns true if a matching port is founded for the provided backend
// in the slice of service ports.
func matchBackendPort(b extv1beta1.IngressBackend, ports []corev1.ServicePort) bool {
	for _, p := range ports {
		switch b.ServicePort.Type {
		case intstr.String:
			if i, err := strconv.Atoi(b.ServicePort.StrVal); err == nil {
				if int32(i) == p.Port {
					return true
				}
			}
			if b.ServicePort.StrVal == p.Name {
				return true
			}
		case intstr.Int:
			if int32(b.ServicePort.IntVal) == p.Port {
				return true
			}
		default:
			continue
		}
	}

	return false
}
