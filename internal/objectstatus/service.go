/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectstatus

import (
	"context"

	"github.com/vmware-tanzu/octant/internal/link"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func service(ctx context.Context, object runtime.Object, o store.Store, _ link.Interface) (ObjectStatus, error) {
	if object == nil {
		return ObjectStatus{}, errors.Errorf("service is nil")
	}

	service := &corev1.Service{}

	if err := scheme.Scheme.Convert(object, service, 0); err != nil {
		return ObjectStatus{}, errors.Wrap(err, "convert object to service")
	}

	if service.Spec.ExternalName == "" {
		key := store.Key{
			Namespace:  service.Namespace,
			APIVersion: "v1",
			Kind:       "Endpoints",
			Name:       service.Name,
		}

		endpoints := &corev1.Endpoints{}

		found, err := store.GetAs(ctx, o, key, endpoints)
		if err != nil {
			return ObjectStatus{}, errors.Wrapf(err, "get endpoints for service %s", service.Name)
		}

		if !found {
			return ObjectStatus{
				nodeStatus: component.NodeStatusWarning,
				Details:    []component.Component{component.NewText("Service has no endpoints")},
			}, nil
		}

		addressCount := 0

		for _, subset := range endpoints.Subsets {
			addressCount += len(subset.Addresses)
		}

		if addressCount == 0 {
			return ObjectStatus{
				nodeStatus: component.NodeStatusWarning,
				Details:    []component.Component{component.NewText("Service has no endpoint addresses")},
			}, nil
		}
	}
	properties := []component.Property{{Label: "Type", Value: component.NewText(string(service.Spec.Type))},
		{Label: "Session Affinity", Value: component.NewText(string(service.Spec.SessionAffinity))}}

	return ObjectStatus{
		nodeStatus: component.NodeStatusOK,
		Details:    []component.Component{component.NewText("Service is OK")},
		Properties: properties,
	}, nil
}
