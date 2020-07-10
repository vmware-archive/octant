/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectstatus

import (
	"context"

	"github.com/pkg/errors"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

func cronJob(_ context.Context, object runtime.Object, _ store.Store) (ObjectStatus, error) {
	if object == nil {
		return ObjectStatus{}, errors.Errorf("cronjob is nil")
	}

	cronjob := &batchv1beta1.CronJob{}

	if err := scheme.Scheme.Convert(object, cronjob, 0); err != nil {
		return ObjectStatus{}, errors.Wrap(err, "convert object to batch/v1beta1 cronjob")
	}

	if cronjob.Spec.Suspend != nil && *cronjob.Spec.Suspend {
		return ObjectStatus{
			nodeStatus: component.NodeStatusWarning,
			Details:    []component.Component{component.NewText("Cronjob is suspended")},
		}, nil
	}
	return ObjectStatus{
		nodeStatus: component.NodeStatusOK,
		Details:    []component.Component{component.NewText("batch/v1beta1 CronJob is OK")},
	}, nil
}
