/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectstatus

import (
	"context"

	"github.com/vmware-tanzu/octant/internal/link"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"

	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// apiService creates status for an apiregistration.k8s.io/v1 apiservice.
// This is not the final implementation. It is included to generate output.
func apiService(_ context.Context, object runtime.Object, _ store.Store, _ link.Interface) (ObjectStatus, error) {
	if object == nil {
		return ObjectStatus{}, errors.Errorf("apiservice is nil")
	}

	apiService := &apiregistrationv1.APIService{}

	if err := scheme.Scheme.Convert(object, apiService, 0); err != nil {
		return ObjectStatus{}, errors.Wrap(err, "convert object to apiregistration.k8s.io/v1 apiservice")
	}

	var availableCondition *apiregistrationv1.APIServiceCondition
	for _, cond := range apiService.Status.Conditions {
		if cond.Type == apiregistrationv1.Available {
			availableCondition = &cond
			break
		}
	}

	switch {
	case availableCondition == nil:
		return ObjectStatus{
			NodeStatus: component.NodeStatusWarning,
			Details:    []component.Component{component.NewText("No available condition for this apiservice")},
		}, nil
	case availableCondition.Status == apiregistrationv1.ConditionFalse:
		return ObjectStatus{
			NodeStatus: component.NodeStatusError,
			Details:    []component.Component{component.NewTextf("Not available: (%s) %s", availableCondition.Reason, availableCondition.Message)},
		}, nil
	case availableCondition.Status == apiregistrationv1.ConditionTrue:
		return ObjectStatus{
			NodeStatus: component.NodeStatusOK,
			Details:    []component.Component{component.NewText("API Service is OK")},
		}, nil
	default:
		return ObjectStatus{
			NodeStatus: component.NodeStatusWarning,
			Details: []component.Component{
				component.NewTextf("Unknown availability for apiservice")},
		}, nil
	}
}
