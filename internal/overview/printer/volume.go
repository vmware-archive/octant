package printer

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/pkg/errors"
)

// VolumeListHandler is a printFunc that lists volumes
func VolumeListHandler(spec *corev1.PodSpec, opts Options) (component.ViewComponent, error) {
	if spec == nil {
		return nil, errors.New("nil list")
	}

	cols := component.NewTableCols("Name", "Kind")
	tbl := component.NewTable("Volumes", cols)

	for _, volume := range spec.Volumes {
		row := component.TableRow{}
		row["Name"] = component.NewText("", volume.Name)

		switch {
		case volume.VolumeSource.HostPath != nil:
			row["Kind"] = component.NewText("", "HostPath (bare host directory volume)")
		case volume.VolumeSource.EmptyDir != nil:
			row["Kind"] = component.NewText("", "EmptyDir (a temporary directory that shares a pod's lifetime)")
		case volume.VolumeSource.GCEPersistentDisk != nil:
			row["Kind"] = component.NewText("", "GCEPersistentDisk (a Persistent Disk resource in Google Compute Engine)")
		case volume.VolumeSource.AWSElasticBlockStore != nil:
			row["Kind"] = component.NewText("", "AWSElasticBlockStore (a Persistent Disk resource in AWS)")
		case volume.VolumeSource.GitRepo != nil:
			row["Kind"] = component.NewText("", "GitRepo (a volume that is pulled from git when the pod is created)")
		case volume.VolumeSource.Secret != nil:
			row["Kind"] = component.NewText("", "Secret (a volume populated by a Secret)")
		case volume.VolumeSource.ConfigMap != nil:
			row["Kind"] = component.NewText("", "ConfigMapSecret (a volume populated by a ConfigMap)")
		case volume.VolumeSource.NFS != nil:
			row["Kind"] = component.NewText("", "NFS (an NFS mount that lasts the lifetime of a pod)")
		case volume.VolumeSource.ISCSI != nil:
			row["Kind"] = component.NewText("", "ISCSI (an ISCSI Disk resource that is attached to a kubelet's host machine and then exposed to the pod)")
		case volume.VolumeSource.Glusterfs != nil:
			row["Kind"] = component.NewText("", "Glusterfs (a Glusterfs mount on the host that shares a pod's lifetime)")
		case volume.VolumeSource.PersistentVolumeClaim != nil:
			row["Kind"] = component.NewText("", "PersistentVolumeClaim")
		case volume.VolumeSource.RBD != nil:
			row["Kind"] = component.NewText("", "RBD (a Rados Block Device mount on the host that shares a pod's lifetime)")
		case volume.VolumeSource.Quobyte != nil:
			row["Kind"] = component.NewText("", "Quobyte (a Quobyte mount on the host that shares a pod's lifetime)")
		case volume.VolumeSource.DownwardAPI != nil:
			row["Kind"] = component.NewText("", "DownwardAPI (a volume populated by information about the pod)")
		case volume.VolumeSource.AzureDisk != nil:
			row["Kind"] = component.NewText("", "AzureDisk (an Azure Data Disk mount on the host and bind mount to the pod)")
		case volume.VolumeSource.VsphereVolume != nil:
			row["Kind"] = component.NewText("", "SphereVolume (a Persistent Disk resource in vSphere)")
		case volume.VolumeSource.Cinder != nil:
			row["Kind"] = component.NewText("", "Cinder (a Persistent Disk resource in OpenStack)")
		case volume.VolumeSource.PhotonPersistentDisk != nil:
			row["Kind"] = component.NewText("", "PhotonPersistentDisk (a Persistent Disk resource in photon platform)")
		case volume.VolumeSource.PortworxVolume != nil:
			row["Kind"] = component.NewText("", "PortworxVolume (a Portworx Volume resource)")
		case volume.VolumeSource.ScaleIO != nil:
			row["Kind"] = component.NewText("", "ScaleIO (a persistent volume backed by a block device in ScaleIO)")
		case volume.VolumeSource.CephFS != nil:
			row["Kind"] = component.NewText("", "CephFS (a CephFS mount on the host that shares a pod's lifetime)")
		case volume.VolumeSource.StorageOS != nil:
			row["Kind"] = component.NewText("", "StorageOS (a StorageOS Persistent Disk resource)")
		case volume.VolumeSource.FC != nil:
			row["Kind"] = component.NewText("", "FC (a Fibre Channel disk)")
		case volume.VolumeSource.AzureFile != nil:
			row["Kind"] = component.NewText("", "AzureFile (an Azure File Service mount on the host and bind mount to the pod)")
		case volume.VolumeSource.FlexVolume != nil:
			row["Kind"] = component.NewText("", "FlexVolume (a generic volume resource that is provisioned/attached using an exec based plugin)")
		case volume.VolumeSource.Flocker != nil:
			row["Kind"] = component.NewText("", "Flocker (a Flocker volume mounted by the Flocker agent)")
		default:
			row["Kind"] = component.NewText("", "Unknown")
		}

		tbl.Add(row)
	}

	return tbl, nil
}

// VolumeHandler is a printFunc that prints Volumes.
// TODO: This handler is incomplete.
func VolumeHandler(volume *corev1.Volume, options Options) (component.ViewComponent, error) {
	grid := component.NewGrid("Summary")

	detailsSummary := component.NewSummary("Details")

	detailsPanel := component.NewPanel("", detailsSummary)
	grid.Add(*detailsPanel)

	list := component.NewList("", []component.ViewComponent{grid})

	return list, nil
}
