package printer

import (
	"fmt"
	"strings"

	"github.com/heptio/developer-dash/internal/overview/link"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
)

// PersistentVolumeClaimListHandler is a printFunc that prints persistentvolumeclaims
func PersistentVolumeClaimListHandler(list *corev1.PersistentVolumeClaimList, options Options) (component.ViewComponent, error) {
	if list == nil {
		return nil, errors.New("nil list")
	}

	cols := component.NewTableCols("Name", "Status", "Volume", "Capacity", "Access Modes", "Storage Class", "Age")
	tbl := component.NewTable("Persistent Volume Claims", cols)

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

		row["Name"] = link.ForObject(&persistentVolumeClaim, persistentVolumeClaim.Name)
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
func PersistentVolumeClaimHandler(persistentVolumeClaim *corev1.PersistentVolumeClaim, options Options) (component.ViewComponent, error) {
	o := NewObject(persistentVolumeClaim)

	o.RegisterConfig(func() (component.ViewComponent, error) {
		return printPersistentVolumeClaimConfig(persistentVolumeClaim)
	}, 12)

	o.RegisterSummary(func() (component.ViewComponent, error) {
		return printPersistentVolumeClaimStatus(persistentVolumeClaim)
	}, 12)

	o.RegisterItems(ItemDescriptor{
		Func: func() (component.ViewComponent, error) {
			return createMountedPodListView(persistentVolumeClaim.Namespace, persistentVolumeClaim.Name, options)
		},
		Width: 24,
	})

	o.EnableEvents()

	return o.ToComponent(options)
}

func printPersistentVolumeClaimConfig(persistentVolumeClaim *corev1.PersistentVolumeClaim) (component.ViewComponent, error) {
	if persistentVolumeClaim == nil {
		return nil, errors.New("persistentvolumeclaim is nil")
	}

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

func printPersistentVolumeClaimStatus(persistentVolumeClaim *corev1.PersistentVolumeClaim) (component.ViewComponent, error) {
	if persistentVolumeClaim == nil {
		return nil, errors.New("persistentvolumeclaim is nil")
	}

	var sections component.SummarySections

	if storageStatus := persistentVolumeClaim.Status.Phase; &storageStatus != nil {
		sections.AddText("Claim Status", string(storageStatus))
	}

	if requestedStorage := persistentVolumeClaim.Spec.Resources.Requests[corev1.ResourceStorage]; &requestedStorage != nil {
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

		if availableStorage := persistentVolumeClaim.Status.Capacity[corev1.ResourceStorage]; &availableStorage != nil {
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
