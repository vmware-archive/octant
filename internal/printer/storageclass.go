/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/vmware-tanzu/octant/pkg/view/component"

	storagev1 "k8s.io/api/storage/v1"
)

// StorageClassListHandler is a printFunc that creates a component to display multiple Storage Class
func StorageClassListHandler(ctx context.Context, list *storagev1.StorageClassList, options Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("storage class list is nil")
	}

	cols := component.NewTableCols("Name", "Provisioner", "Age")
	ot := NewObjectTable("Storage Class", "We couldn't find any storage class!", cols, options.DashConfig.ObjectStore())

	for _, sc := range list.Items {
		row := component.TableRow{}
		nameLink, err := options.Link.ForObject(&sc, sc.Name)
		if err != nil {
			return nil, err
		}

		row["Name"] = nameLink
		row["Provisioner"] = component.NewText(sc.Provisioner)
		row["Age"] = component.NewTimestamp(sc.CreationTimestamp.Time)

		if err := ot.AddRowForObject(ctx, &sc, row); err != nil {
			return nil, fmt.Errorf("add row for object: %w", err)
		}
	}

	return ot.ToComponent()
}

// StorageClassHandler is a printFunc that creates a component to display a single Storage Class
func StorageClassHandler(ctx context.Context, sc *storagev1.StorageClass, options Options) (component.Component, error) {
	obj := NewObject(sc)
	obj.EnableEvents()

	sch, err := newStorageClassHandler(sc, obj)
	if err != nil {
		return nil, err
	}

	if err := sch.Config(options); err != nil {
		return nil, errors.Wrap(err, "print storage class configuration")
	}

	return obj.ToComponent(ctx, options)
}

type storageClassHandler struct {
	configFunc   func(*storagev1.StorageClass, Options) (*component.Summary, error)
	dataFunc     func(*storagev1.StorageClass, Options) (*component.Summary, error)
	storageClass *storagev1.StorageClass
	object       *Object
}

func newStorageClassHandler(sc *storagev1.StorageClass, object *Object) (*storageClassHandler, error) {
	if sc == nil {
		return nil, errors.New("cannot print a nil storageclass")
	}
	if object == nil {
		return nil, errors.New("cannot print storageclass using a nil object printer")
	}

	sch := &storageClassHandler{
		configFunc:   defaultStorageClassConfig,
		storageClass: sc,
		object:       object,
	}

	return sch, nil
}

func (sch *storageClassHandler) Config(options Options) error {
	out, err := sch.configFunc(sch.storageClass, options)
	if err != nil {
		return err
	}
	sch.object.RegisterConfig(out)
	return nil
}

func (sch *storageClassHandler) Data(ctx context.Context, options Options) error {
	out, err := sch.dataFunc(sch.storageClass, options)
	if err != nil {
		return err
	}
	sch.object.RegisterSummary(out)
	return nil
}

func defaultStorageClassConfig(sc *storagev1.StorageClass, options Options) (*component.Summary, error) {
	return NewStorageClassConfiguration(sc).Create(options)
}

// StorageClassConfiguration is used to create the Storage Class' configuration component
// when displaying a single Storage Class
type StorageClassConfiguration struct {
	storageClass *storagev1.StorageClass
}

// NewStorageClassConfiguration creates a new StorageClassConfiguration using the supplied Storage Class
func NewStorageClassConfiguration(sc *storagev1.StorageClass) *StorageClassConfiguration {
	return &StorageClassConfiguration{
		storageClass: sc,
	}
}

// Create the Configuration Summary component for a Stoage Class
func (scc *StorageClassConfiguration) Create(options Options) (*component.Summary, error) {
	if scc.storageClass == nil {
		return nil, errors.New("Storage Class is nil")
	}
	sc := scc.storageClass

	provisioner := sc.Provisioner

	sections := component.SummarySections{}
	sections.AddText("Provisioner", provisioner)
	if reclaimPolicy := sc.ReclaimPolicy; reclaimPolicy != nil {
		sections.AddText("Reclaim Policy", string(*reclaimPolicy))
	}

	if volumeBindingMode := sc.VolumeBindingMode; volumeBindingMode != nil {
		sections.AddText("Volume Binding Mode", string(*volumeBindingMode))
	}

	summary := component.NewSummary("Configuration", sections...)
	return summary, nil
}
