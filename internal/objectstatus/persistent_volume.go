/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectstatus

import (
	"context"
	"fmt"

	"github.com/vmware-tanzu/octant/internal/link"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func persistentVolume(ctx context.Context, object runtime.Object, o store.Store, _ link.Interface) (ObjectStatus, error) {
	if object == nil {
		return ObjectStatus{}, errors.Errorf("cronjob is nil")
	}

	pv := &corev1.PersistentVolume{}

	if err := scheme.Scheme.Convert(object, pv, 0); err != nil {
		return ObjectStatus{}, errors.Wrap(err, "convert object to v1 PersistentVolume")
	}

	if pv.Spec.ClaimRef != nil {
		claim := pv.Spec.ClaimRef
		pvc, err := o.Get(ctx, store.Key{
			Kind:       claim.Kind,
			Namespace:  claim.Namespace,
			Name:       claim.Name,
			APIVersion: claim.APIVersion,
		})
		if pvc == nil || err != nil {
			claimWarning := fmt.Sprintf("PVC %s/%s cannot be found", claim.Namespace, claim.Name)
			return ObjectStatus{
				NodeStatus: component.NodeStatusWarning,
				Details:    []component.Component{component.NewText(claimWarning)}}, nil
		}
	}

	return ObjectStatus{
		NodeStatus: component.NodeStatusOK,
		Details:    []component.Component{component.NewText("v1 PersistentVolume is OK")},
	}, nil
}
