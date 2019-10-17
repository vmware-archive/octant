/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"

	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// PersistentVolumeClaimListHandler is a printFunc that prints persistentvolumeclaims
func PersistentVolumeClaimListHandler(ctx context.Context, list *corev1.PersistentVolumeClaimList, options Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("nil list")
	}

	cols := component.NewTableCols("Name", "Status", "Volume", "Capacity", "Access Modes", "Storage Class", "Age")
	tbl := component.NewTable("Persistent Volume Claims",
		"We couldn't find any persistent volume claims!", cols)

	for _, persistentVolumeClaim := range list.Items {
		row := component.TableRow{}

		accessModes := ""
		storage := persistentVolumeClaim.Spec.Resources.Requests[corev1.ResourceStorage]
		capacity := ""
		// Check if pvc is bound
		if persistentVolumeClaim.Spec.VolumeName != "" {
			accessModes = getAccessModesAsString(persistentVolumeClaim.Spec.AccessModes)
			storage = persistentVolumeClaim.Status.Capacity[corev1.ResourceStorage]
			capacity = storage.String()
		}

		nameLink, err := options.Link.ForObject(&persistentVolumeClaim, persistentVolumeClaim.Name)
		if err != nil {
			return nil, err
		}

		row["Name"] = nameLink

		row["Status"] = component.NewText(string(persistentVolumeClaim.Status.Phase))
		// TODO: Link to volume
		row["Volume"] = component.NewText(persistentVolumeClaim.Spec.VolumeName)
		row["Capacity"] = component.NewText(capacity)
		row["Access Modes"] = component.NewText(accessModes)
		row["Storage Class"] = component.NewText(printPersistentVolumeClaimClass(&persistentVolumeClaim))
		ts := persistentVolumeClaim.CreationTimestamp.Time
		row["Age"] = component.NewTimestamp(ts)

		tbl.Add(row)
	}
	return tbl, nil
}

// PersistentVolumeClaimHandler is a printFunc that prints a PersistentVolumeClaim
func PersistentVolumeClaimHandler(ctx context.Context, persistentVolumeClaim *corev1.PersistentVolumeClaim, options Options) (component.Component, error) {
	o := NewObject(persistentVolumeClaim)
	o.EnableEvents()

	ph, err := newPersistentVolumeClaimHandler(persistentVolumeClaim, o)
	if err != nil {
		return nil, err
	}

	if err := ph.Config(options); err != nil {
		return nil, errors.Wrap(err, "print persistentvolumeclaim configuration")
	}

	if err := ph.Status(options); err != nil {
		return nil, errors.Wrap(err, "print persistentvolumeclaim status")
	}

	if err := ph.MountedPodList(ctx, options); err != nil {
		return nil, errors.Wrap(err, "print peristentvolumeclaim mounted pod list")
	}
	return o.ToComponent(ctx, options)
}

// PersistentVolumeClaimConfiguration generates a persistenvolumeclaim configuration
type PersistentVolumeClaimConfiguration struct {
	persistentVolumeClaim *corev1.PersistentVolumeClaim
}

// NewPersistentVolumeClaimConfiguration creates an instance of PersistentVolumeClaimConfiguration
func NewPersistentVolumeClaimConfiguration(pvc *corev1.PersistentVolumeClaim) *PersistentVolumeClaimConfiguration {
	return &PersistentVolumeClaimConfiguration{
		persistentVolumeClaim: pvc,
	}
}

// Create creates a persistentvolumeclaim configuration summary
func (p *PersistentVolumeClaimConfiguration) Create(options Options) (*component.Summary, error) {
	if p == nil || p.persistentVolumeClaim == nil {
		return nil, errors.New("persistentvolumeclaim is nil")
	}

	persistentVolumeClaim := p.persistentVolumeClaim

	var sections component.SummarySections

	if volumeMode := persistentVolumeClaim.Spec.VolumeMode; volumeMode != nil {
		sections.AddText("Volume Mode", string(*volumeMode))
	}

	if accessMode := persistentVolumeClaim.Spec.AccessModes; accessMode != nil {
		sections.AddText("Access Modes", getAccessModesAsString(accessMode))
	}

	if finalizers := persistentVolumeClaim.ObjectMeta.Finalizers; finalizers != nil {
		sections.AddText("Finalizers", fmt.Sprint(finalizers))
	}

	if storageClassName := persistentVolumeClaim.Spec.StorageClassName; storageClassName != nil {
		sections.AddText("Storage Class Name", string(*storageClassName))
	}

	if labels := persistentVolumeClaim.Labels; labels != nil {
		sections.Add("Labels", component.NewLabels(labels))
	}

	if selector := persistentVolumeClaim.Spec.Selector; selector != nil {
		sections.Add("Selectors", printSelector(selector))
	}

	summary := component.NewSummary("Configuration", sections...)

	return summary, nil
}

func createPersistentVolumeClaimStatusView(persistentVolumeClaim *corev1.PersistentVolumeClaim) (*component.Summary, error) {
	if persistentVolumeClaim == nil {
		return nil, errors.New("persistentvolumeclaim is nil")
	}

	sections := component.SummarySections{}

	if persistentVolumeClaim.Status.Phase != "" {
		sections.AddText("Claim Status", string(persistentVolumeClaim.Status.Phase))
	}

	if requestedStorage, ok := persistentVolumeClaim.Spec.Resources.Requests[corev1.ResourceStorage]; ok {
		sections.AddText("Storage Requested", requestedStorage.String())
	}

	if persistentVolumeClaim.Spec.VolumeName != "" {
		if boundVolume := persistentVolumeClaim.Spec.VolumeName; boundVolume != "" {
			sections = append(sections, component.SummarySection{
				Header: "Bound Volume",
				// TODO: Link to volume
				Content: component.NewText(boundVolume),
			})
		}

		if availableStorage, ok := persistentVolumeClaim.Status.Capacity[corev1.ResourceStorage]; ok {
			sections.AddText("Total Volume Capacity", availableStorage.String())
		}
	}

	summary := component.NewSummary("Status", sections...)

	return summary, nil
}

func printPersistentVolumeClaimClass(persistentVolumeClaim *corev1.PersistentVolumeClaim) string {
	if class, found := persistentVolumeClaim.Annotations[corev1.BetaStorageClassAnnotation]; found {
		return class
	}

	if persistentVolumeClaim.Spec.StorageClassName != nil {
		return *persistentVolumeClaim.Spec.StorageClassName
	}

	return ""
}

func getAccessModesAsString(modes []corev1.PersistentVolumeAccessMode) string {
	modes = removeDuplicateAccessModes(modes)
	modesStr := []string{}

	if containsAccessMode(modes, corev1.ReadWriteOnce) {
		modesStr = append(modesStr, "RWO")
	}
	if containsAccessMode(modes, corev1.ReadOnlyMany) {
		modesStr = append(modesStr, "ROX")
	}
	if containsAccessMode(modes, corev1.ReadWriteMany) {
		modesStr = append(modesStr, "RWX")
	}
	return strings.Join(modesStr, ",")
}

func containsAccessMode(modes []corev1.PersistentVolumeAccessMode, mode corev1.PersistentVolumeAccessMode) bool {
	for _, m := range modes {
		if m == mode {
			return true
		}
	}
	return false
}

// removeDuplicateAccessModes returns an array of access modes without any duplicates
func removeDuplicateAccessModes(modes []corev1.PersistentVolumeAccessMode) []corev1.PersistentVolumeAccessMode {
	accessModes := []corev1.PersistentVolumeAccessMode{}
	for _, m := range modes {
		if !containsAccessMode(accessModes, m) {
			accessModes = append(accessModes, m)
		}
	}
	return accessModes
}

type persistentVolumeClaimObject interface {
	Config(options Options) error
	Status(options Options) error
	MountedPodList(ctx context.Context, options Options) error
}

type persistentVolumeClaimHandler struct {
	persistentVolumeClaim *corev1.PersistentVolumeClaim
	configFunc            func(*corev1.PersistentVolumeClaim, Options) (*component.Summary, error)
	statusFunc            func(*corev1.PersistentVolumeClaim, Options) (*component.Summary, error)
	mountedPodListFunc    func(context.Context, string, string, Options) (component.Component, error)
	object                *Object
}

var _ persistentVolumeClaimObject = (*persistentVolumeClaimHandler)(nil)

func newPersistentVolumeClaimHandler(pvc *corev1.PersistentVolumeClaim, object *Object) (*persistentVolumeClaimHandler, error) {
	if pvc == nil {
		return nil, errors.New("can't print a nil persistentvolumeclaim")
	}

	if object == nil {
		return nil, errors.New("can't print persistentvolume claim using a nil object printer")
	}

	ph := &persistentVolumeClaimHandler{
		persistentVolumeClaim: pvc,
		configFunc:            defaultPersistentVolumeClaimConfig,
		statusFunc:            defaultPersistentVolumClaimStatus,
		mountedPodListFunc:    defaultPersistentVolumeMountedPodList,
		object:                object,
	}
	return ph, nil
}

func (p *persistentVolumeClaimHandler) Config(options Options) error {
	out, err := p.configFunc(p.persistentVolumeClaim, options)
	if err != nil {
		return err
	}
	p.object.RegisterConfig(out)
	return nil
}

func defaultPersistentVolumeClaimConfig(pvc *corev1.PersistentVolumeClaim, options Options) (*component.Summary, error) {
	return NewPersistentVolumeClaimConfiguration(pvc).Create(options)
}

func (p *persistentVolumeClaimHandler) Status(options Options) error {
	out, err := p.statusFunc(p.persistentVolumeClaim, options)
	if err != nil {
		return err
	}
	p.object.RegisterSummary(out)
	return nil
}

func defaultPersistentVolumClaimStatus(pvc *corev1.PersistentVolumeClaim, options Options) (*component.Summary, error) {
	return createPersistentVolumeClaimStatusView(pvc)
}

func (p *persistentVolumeClaimHandler) MountedPodList(ctx context.Context, options Options) error {
	if p.persistentVolumeClaim == nil {
		return errors.New("can't display mounted pod list for nil persistentvolumeclaim")
	}

	p.object.RegisterItems(ItemDescriptor{
		Width: component.WidthFull,
		Func: func() (component.Component, error) {
			return p.mountedPodListFunc(ctx, p.persistentVolumeClaim.Namespace, p.persistentVolumeClaim.Name, options)
		},
	})
	return nil
}

func defaultPersistentVolumeMountedPodList(ctx context.Context, namespace string, name string, options Options) (component.Component, error) {
	return createMountedPodListView(ctx, namespace, name, options)
}
