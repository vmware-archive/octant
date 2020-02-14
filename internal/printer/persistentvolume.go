/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// PersistentVolumeListHandler is a printFunc that creates a component to display multiple Persistent Volumes
func PersistentVolumeListHandler(ctx context.Context, list *v1.PersistentVolumeList, options Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("nil list")
	}

	cols := component.NewTableCols("Name", "Capacity", "Access Modes", "Reclaim Policy", "Status", "Claim", "Storage Class", "Reason", "Age")
	table := component.NewTable("Persistent Volumes", "We couldn't find any persistent volumes!", cols)

	for _, pv := range list.Items {
		row := component.TableRow{}
		nameLink, err := options.Link.ForObject(&pv, pv.Name)
		if err != nil {
			return nil, err
		}

		storage := pv.Spec.Capacity[v1.ResourceStorage]
		capacity := storage.String()

		accessModes := getAccessModesAsString(pv.Spec.AccessModes)

		claimLink, err := createBoundPersistentVolumeClaimLink(ctx, &pv, options)
		if err != nil {
			return nil, err
		}

		row["Name"] = nameLink
		row["Capacity"] = component.NewText(capacity)
		row["Access Modes"] = component.NewText(accessModes)
		row["Reclaim Policy"] = component.NewText(string(pv.Spec.PersistentVolumeReclaimPolicy))
		row["Status"] = component.NewText(string(pv.Status.Phase))
		row["Claim"] = claimLink
		row["Storage Class"] = component.NewText(pv.Spec.StorageClassName)
		row["Reason"] = component.NewText(pv.Status.Reason)
		row["Age"] = component.NewTimestamp(pv.CreationTimestamp.Time)

		table.Add(row)
	}

	return table, nil
}

// PersistentVolumeHandler is a printFunc that creates a component to display a single Persistent Volume
func PersistentVolumeHandler(ctx context.Context, pv *v1.PersistentVolume, options Options) (component.Component, error) {
	obj := NewObject(pv)
	obj.EnableEvents()

	pvh, err := newPersistentVolumeHandler(pv, obj)
	if err != nil {
		return nil, err
	}

	if err := pvh.Config(options); err != nil {
		return nil, errors.Wrap(err, "print persistent volume configuration")
	}

	if err := pvh.Status(ctx, options); err != nil {
		return nil, errors.Wrap(err, "print persisten volume claims")
	}

	return obj.ToComponent(ctx, options)
}

type persistentVolumeHandler struct {
	configFunc       func(*v1.PersistentVolume, Options) (*component.Summary, error)
	statusFunc       func(context.Context, *v1.PersistentVolume, Options) (*component.Summary, error)
	persistentVolume *v1.PersistentVolume
	object           *Object
}

func newPersistentVolumeHandler(pv *v1.PersistentVolume, object *Object) (*persistentVolumeHandler, error) {
	if pv == nil {
		return nil, errors.New("cannot print a nil persistentVolume")
	}
	if object == nil {
		return nil, errors.New("cannot print persistentVolume using a nil object printer")
	}

	pvh := &persistentVolumeHandler{
		configFunc:       defaultPersistentVolumeConfig,
		statusFunc:       defaultPersistentVolumeStatus,
		persistentVolume: pv,
		object:           object,
	}

	return pvh, nil
}

func (pvh *persistentVolumeHandler) Config(options Options) error {
	out, err := pvh.configFunc(pvh.persistentVolume, options)
	if err != nil {
		return err
	}
	pvh.object.RegisterConfig(out)
	return nil
}

func (pvh *persistentVolumeHandler) Status(ctx context.Context, options Options) error {
	out, err := pvh.statusFunc(ctx, pvh.persistentVolume, options)
	if err != nil {
		return err
	}
	pvh.object.RegisterSummary(out)
	return nil
}

func defaultPersistentVolumeConfig(pv *v1.PersistentVolume, options Options) (*component.Summary, error) {
	return NewPersistentVolumeConfiguration(pv).Create(options)
}

func defaultPersistentVolumeStatus(ctx context.Context, pv *v1.PersistentVolume, options Options) (*component.Summary, error) {
	return NewPersistentVolumeStatus(pv).Create(ctx, options)
}

// PersistentVolumeConfiguration is used to create the Persistent Volume's configuration component
// when displaying a single Persistent Volume
type PersistentVolumeConfiguration struct {
	persistentVolume *v1.PersistentVolume
}

// NewPersistentVolumeConfiguration creates a new PersistentVolumeConfiguration using the supplied Persistent Volume
func NewPersistentVolumeConfiguration(pv *v1.PersistentVolume) *PersistentVolumeConfiguration {
	return &PersistentVolumeConfiguration{
		persistentVolume: pv,
	}
}

// Create the Configuration Summary component for a Persistent Volume
func (pvc *PersistentVolumeConfiguration) Create(options Options) (*component.Summary, error) {
	if pvc.persistentVolume == nil {
		return nil, errors.New("Persistent Volume is nil")
	}
	pv := pvc.persistentVolume

	accessModes := getAccessModesAsString(pv.Spec.AccessModes)
	storage := pv.Spec.Capacity[v1.ResourceStorage]
	capacity := storage.String()

	var sections component.SummarySections
	sections.AddText("Reclaim Policy", string(pv.Spec.PersistentVolumeReclaimPolicy))
	sections.AddText("Storage Class", pv.Spec.StorageClassName)
	sections.AddText("Access Modes", accessModes)
	sections.AddText("Capacity", capacity)
	sections.AddText("Host Path", pv.Spec.HostPath.Path)

	summary := component.NewSummary("Configuration", sections...)
	return summary, nil
}

// PersistentVolumeStatus is used to create the Persistent Volume's status component
// when displaying a single Persistent Volume
type PersistentVolumeStatus struct {
	persistentVolume *v1.PersistentVolume
}

// NewPersistentVolumeStatus creates a new PersistentVolumeStatus using the supplied Persistent Volume
func NewPersistentVolumeStatus(pv *v1.PersistentVolume) *PersistentVolumeStatus {
	return &PersistentVolumeStatus{
		persistentVolume: pv,
	}
}

// Create the Status Summary component for a Persistent Volume
func (pvs *PersistentVolumeStatus) Create(ctx context.Context, options Options) (*component.Summary, error) {
	if pvs.persistentVolume == nil {
		return nil, errors.New("Persistent Volume is nil")
	}
	pv := pvs.persistentVolume

	var sections component.SummarySections

	sections.AddText("Phase Status", string(pv.Status.Phase))

	claimLink, err := createBoundPersistentVolumeClaimLink(ctx, pv, options)
	if err != nil {
		return nil, err
	}

	sections = append(sections, component.SummarySection{
		Header:  "Claim",
		Content: claimLink,
	})

	summary := component.NewSummary("Status", sections...)
	return summary, nil
}

func getBoundPersistentVolumeClaim(ctx context.Context, pv *v1.PersistentVolume, options Options) (*v1.PersistentVolumeClaim, error) {
	objectStore := options.DashConfig.ObjectStore()
	pvc := &v1.PersistentVolumeClaim{}

	cr := pv.Spec.ClaimRef
	key := store.Key{
		APIVersion: cr.APIVersion,
		Kind:       cr.Kind,
		Name:       cr.Name,
		Namespace:  cr.Namespace,
	}

	o, err := objectStore.Get(ctx, key)
	if err != nil {
		return nil, errors.Wrapf(err, "get persistent volume claim for key %+v", key)
	}

	if o != nil {
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(o.Object, pvc)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Persistent Volume Claim not found")
	}

	return pvc, nil
}

func createBoundPersistentVolumeClaimLink(ctx context.Context, pv *v1.PersistentVolume, options Options) (*component.Link, error) {
	pvc, err := getBoundPersistentVolumeClaim(ctx, pv, options)
	if err != nil {
		return nil, err
	}

	cr := pv.Spec.ClaimRef
	claimText := fmt.Sprintf("%s/%s", cr.Namespace, cr.Name)
	claimLink, err := options.Link.ForObject(pvc, claimText)
	if err != nil {
		return nil, err
	}

	return claimLink, nil
}
