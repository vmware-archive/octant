/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectstatus

import (
	"context"

	"strconv"

	"github.com/pkg/errors"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/vmware-tanzu/octant/internal/link"

	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func cronJob(_ context.Context, object runtime.Object, _ store.Store, _ link.Interface) (ObjectStatus, error) {
	if object == nil {
		return ObjectStatus{}, errors.Errorf("cronjob is nil")
	}

	cronjob := &batchv1beta1.CronJob{}

	if err := scheme.Scheme.Convert(object, cronjob, 0); err != nil {
		return ObjectStatus{}, errors.Wrap(err, "convert object to batch/v1beta1 cronjob")
	}

	properties := []component.Property{{Label: "Schedule", Value: component.NewText(cronjob.Spec.Schedule)}}

	if cronjob.Spec.Suspend != nil && *cronjob.Spec.Suspend {
		properties = append(properties, component.Property{Label: "Suspend", Value: component.NewText(strconv.FormatBool(*cronjob.Spec.Suspend))})
		return ObjectStatus{
			NodeStatus: component.NodeStatusWarning,
			Details:    []component.Component{component.NewText("Cronjob is suspended")},
			Properties: properties,
		}, nil
	}
	return ObjectStatus{
		NodeStatus: component.NodeStatusOK,
		Details:    []component.Component{component.NewText("batch/v1beta1 CronJob is OK")},
		Properties: properties,
	}, nil
}
