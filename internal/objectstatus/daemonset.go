/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectstatus

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/vmware-tanzu/octant/internal/link"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func daemonSet(_ context.Context, object runtime.Object, _ store.Store, _ link.Interface) (ObjectStatus, error) {
	if object == nil {
		return ObjectStatus{}, errors.Errorf("daemon set is nil")
	}

	ds := &appsv1.DaemonSet{}

	if err := scheme.Scheme.Convert(object, ds, 0); err != nil {
		return ObjectStatus{}, errors.Wrap(err, "convert object to daemon set")
	}

	status := ds.Status

	var properties []component.Property

	if selector := ds.Spec.Selector; selector != nil {
		properties = append(properties, component.Property{Label: "Selectors", Value: getSelectors(selector)})
	}

	if nodeSelector := ds.Spec.Template.Spec.NodeSelector; nodeSelector != nil {
		properties = append(properties, component.Property{Label: "Node Selectors", Value: getSelectorMap(nodeSelector)})
	}

	switch {
	case status.NumberMisscheduled > 0:
		return ObjectStatus{
			nodeStatus: component.NodeStatusWarning,
			Details:    []component.Component{component.NewText("Daemon Set pods are running on nodes that aren't supposed to run Daemon Set pods")},
			Properties: properties,
		}, nil
	case status.DesiredNumberScheduled != status.NumberReady:
		return ObjectStatus{
			nodeStatus: component.NodeStatusWarning,
			Details:    []component.Component{component.NewText("Daemon Set pods are not ready")},
			Properties: properties,
		}, nil
	default:
		return ObjectStatus{
			nodeStatus: component.NodeStatusOK,
			Details:    []component.Component{component.NewText("Daemon Set is OK")},
			Properties: properties,
		}, nil
	}
}

func getSelectors(selector *metav1.LabelSelector) component.Component {
	selectorComponent := component.NewSelectors(nil)
	if selector == nil {
		return selectorComponent
	}

	for k, v := range selector.MatchLabels {
		selectorComponent.Add(component.NewLabelSelector(k, v))
	}

	for _, e := range selector.MatchExpressions {
		es := component.NewExpressionSelector(e.Key, component.Operator(e.Operator), e.Values)
		selectorComponent.Add(es)
	}
	return selectorComponent
}

func getSelectorMap(selector map[string]string) *component.Selectors {
	s := component.NewSelectors(nil)
	if len(selector) == 0 {
		return s
	}

	for k, v := range selector {
		s.Add(component.NewLabelSelector(k, v))
	}

	return s
}
