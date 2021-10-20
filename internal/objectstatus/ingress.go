/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectstatus

import (
	"context"
	"fmt"

	"github.com/vmware-tanzu/octant/internal/link"
	"github.com/vmware-tanzu/octant/pkg/view/component"

	"github.com/gobwas/glob"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/vmware-tanzu/octant/pkg/store"
)

const (
	ingressNoBackendsDefined   = "No backends defined. All traffic will be sent to the default backend configured for the ingress controller."
	ingressAlbActionAnnotation = "alb.ingress.kubernetes.io/actions."
)

func runIngressStatus(ctx context.Context, object runtime.Object, o store.Store, _ link.Interface) (ObjectStatus, error) {
	if object == nil {
		return ObjectStatus{}, errors.Errorf("ingress is nil")
	}

	ingress := &networkingv1.Ingress{}
	if err := scheme.Scheme.Convert(object, ingress, 0); err != nil {
		return ObjectStatus{}, errors.Wrap(err, "convert object to ingress")
	}

	is := ingressStatus{
		ingress:     *ingress,
		objectStore: o,
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
	ingress     networkingv1.Ingress
	objectStore store.Store
}

func (is *ingressStatus) run(ctx context.Context) (ObjectStatus, error) {
	status := ObjectStatus{}

	ingress := is.ingress

	o := is.objectStore
	if o == nil {
		return status, errors.New("ingress status requires a non nil objectStore")
	}

	backends := is.backends()
	if len(backends) == 0 {
		status.SetWarning()
		status.AddDetail(ingressNoBackendsDefined)
	}

	for _, backend := range backends {
		if backend.Service.Port.String() == "use-annotation" {
			albAction := ingressAlbActionAnnotation + backend.Service.Name
			if _, ok := ingress.Annotations[albAction]; !ok {
				status.SetError()
				status.AddDetailf("Backend refers to annotations %q which doesn't exist", albAction)
			}
			continue
		}

		key := store.Key{
			Namespace:  ingress.Namespace,
			APIVersion: "v1",
			Kind:       "Service",
			Name:       backend.Service.Name,
		}

		service := &corev1.Service{}

		found, err := store.GetAs(ctx, o, key, service)
		if err != nil {
			return ObjectStatus{}, errors.Wrap(err, "get service from object store")
		}

		if !found {
			status.SetError()
			status.AddDetailf("Backend refers to service %q which doesn't exist", key.Name)
			continue
		}

		if service.Name == "" {
			status.SetError()
			status.AddDetailf("Backend refers to service %q which doesn't exist",
				backend.Service.Name)
			continue
		}

		if !matchBackendPort(backend, service.Spec.Ports) {
			status.SetError()
			status.AddDetailf("Backend for service %q specifies an invalid port",
				backend.Service.Name)
			continue
		}
	}

	hm, err := is.createHostMatcher()
	if err != nil {
		status.SetError()
		status.AddDetailf("TLS Hosts: %v", err)
	} else {
		if len(hm.globs) > 0 {
			for _, rule := range ingress.Spec.Rules {
				if rule.Host == "" {
					continue
				}

				if !hm.Match(rule.Host) {
					status.SetError()
					status.AddDetailf("No matching TLS host for rule %q", rule.Host)
				}
			}

		}
	}

	for _, tls := range ingress.Spec.TLS {
		if tls.SecretName == "" {
			status.SetError()
			status.AddDetail("TLS configuration did not define a secret name")
			continue
		}

		key := store.Key{
			Namespace:  ingress.Namespace,
			APIVersion: "v1",
			Kind:       "Secret",
			Name:       tls.SecretName,
		}

		u, err := is.objectStore.Get(ctx, key)
		if err != nil {
			status.SetError()
			status.AddDetailf("Unable to load Secret %q: %s", tls.SecretName, err)
			continue
		}

		if u == nil {
			status.SetError()
			status.AddDetailf("Secret %q does not exist", tls.SecretName)
		}
	}

	var backendText = "Not configured"
	if backend := ingress.Spec.DefaultBackend; backend != nil {
		backendText = backend.Service.Name
	}

	status.Properties = []component.Property{{Label: "Default Backend", Value: component.NewText(backendText)}}
	return status, nil
}

func (is *ingressStatus) backends() []networkingv1.IngressBackend {
	var list []networkingv1.IngressBackend

	ingress := is.ingress
	if ingress.Spec.DefaultBackend != nil && ingress.Spec.DefaultBackend.Service.Name != "" {
		list = append(list, *ingress.Spec.DefaultBackend)
	}

	for _, rule := range ingress.Spec.Rules {
		if rule.IngressRuleValue.HTTP == nil {
			continue
		}

		for _, p := range rule.IngressRuleValue.HTTP.Paths {
			if p.Backend.Service.Name == "" {
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

func (is *ingressStatus) createHostMatcher() (*hostMatcher, error) {
	hm := hostMatcher{}

	for _, tls := range is.ingress.Spec.TLS {
		for _, host := range tls.Hosts {
			if err := hm.AddHost(host); err != nil {
				return nil, fmt.Errorf("parsing TLS host %s: %w", host, err)
			}
		}
	}

	return &hm, nil
}

type hostMatcher struct {
	globs []glob.Glob
}

func (hm *hostMatcher) AddHost(host string) error {
	g, err := glob.Compile(host)
	if err != nil {
		return fmt.Errorf("unable to compile host glob: %w", err)
	}

	hm.globs = append(hm.globs, g)

	return nil
}

func (hm hostMatcher) Match(s string) bool {
	for _, g := range hm.globs {
		if g.Match(s) {
			return true
		}
	}

	return false
}

// matchBackendPort returns true if a matching port is found for the provided backend
// in the slice of service ports.
func matchBackendPort(b networkingv1.IngressBackend, ports []corev1.ServicePort) bool {
	if len(ports) == 0 && b.Service.Port.Number > 0 {
		return false
	}

	var isMatch bool
	for _, p := range ports {
		if b.Service.Port.Number != 0 && b.Service.Port.Number == p.Port {
			isMatch = true
		}

		if b.Service.Port.Name != "" && b.Service.Port.Name == p.Name {
			isMatch = true
		}
	}

	return isMatch
}
