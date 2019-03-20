package printer

import (
	"encoding/json"

	corev1 "k8s.io/api/core/v1"

	"github.com/heptio/developer-dash/pkg/view/component"
)

const (
	volumeKindHostPath              = "HostPath (bare host directory volume)"
	volumeKindEmptyDir              = "EmptyDir (a temporary directory that shares a pod's lifetime)"
	volumeKindGCEPersistentDisk     = "GCEPersistentDisk (a Persistent Disk resource in Google Compute Engine)"
	volumeKindAWSElasticBlockStore  = "AWSElasticBlockStore (a Persistent Disk resource in AWS)"
	volumeKindGitRepo               = "GitRepo (a volume that is pulled from git when the pod is created)"
	volumeKindSecret                = "Secret (a volume populated by a Secret)"
	volumeKindConfigMap             = "ConfigMap (a volume populated by a ConfigMap)"
	volumeKindNFS                   = "NFS (an NFS mount that lasts the lifetime of a pod)"
	volumeKindISCSI                 = "ISCSI (an ISCSI Disk resource that is attached to a kubelet's host machine and then exposed to the pod)"
	volumeKindGlusterfs             = "Glusterfs (a Glusterfs mount on the host that shares a pod's lifetime)"
	volumeKindPersistentVolumeClaim = "PersistentVolumeClaim"
	volumeKindRBD                   = "RBD (a Rados Block Device mount on the host that shares a pod's lifetime)"
	volumeKindQuobyte               = "Quobyte (a Quobyte mount on the host that shares a pod's lifetime)"
	volumeKindDownwardAPI           = "DownwardAPI (a volume populated by information about the pod)"
	volumeKindAzureDisk             = "AzureDisk (an Azure Data Disk mount on the host and bind mount to the pod)"
	volumeKindSphereVolume          = "SphereVolume (a Persistent Disk resource in vSphere)"
	volumeKindCinder                = "Cinder (a Persistent Disk resource in OpenStack)"
	volumeKindPhotonPersistentDisk  = "PhotonPersistentDisk (a Persistent Disk resource in photon platform)"
	volumeKindPortworxVolume        = "PortworxVolume (a Portworx Volume resource)"
	volumeKindScaleIO               = "ScaleIO (a persistent volume backed by a block device in ScaleIO)"
	volumeKindCephFS                = "CephFS (a CephFS mount on the host that shares a pod's lifetime)"
	volumeKindStorageOS             = "StorageOS (a StorageOS Persistent Disk resource)"
	volumeKindFC                    = "FC (a Fibre Channel disk)"
	volumeKindAzureFile             = "AzureFile (an Azure File Service mount on the host and bind mount to the pod)"
	volumeKindFlexVolume            = "FlexVolume (a generic volume resource that is provisioned/attached using an exec based plugin)"
	volumeKindFlocker               = "Flocker (a Flocker volume mounted by the Flocker agent)"
	volumeKindUnknown               = "Unknown"
)

// printVolumes prints volumes as a table.
func printVolumes(volumes []corev1.Volume) (component.Component, error) {
	cols := component.NewTableCols("Name", "Kind", "Description")
	table := component.NewTable("Volumes", cols)

	for _, volume := range volumes {
		row := component.TableRow{}
		row["Name"] = component.NewText(volume.Name)

		switch {
		case volume.VolumeSource.HostPath != nil:
			row["Kind"] = component.NewText(volumeKindHostPath)
			row["Description"] = component.NewText(describeVolumeSource(volume.VolumeSource.HostPath))
		case volume.VolumeSource.EmptyDir != nil:
			row["Kind"] = component.NewText(volumeKindEmptyDir)
			row["Description"] = component.NewText(describeVolumeSource(volume.VolumeSource.EmptyDir))
		case volume.VolumeSource.GCEPersistentDisk != nil:
			row["Kind"] = component.NewText(volumeKindGCEPersistentDisk)
			row["Description"] = component.NewText(describeVolumeSource(volume.VolumeSource.GCEPersistentDisk))
		case volume.VolumeSource.AWSElasticBlockStore != nil:
			row["Kind"] = component.NewText(volumeKindAWSElasticBlockStore)
			row["Description"] = component.NewText(describeVolumeSource(volume.VolumeSource.AWSElasticBlockStore))
		case volume.VolumeSource.GitRepo != nil:
			row["Kind"] = component.NewText(volumeKindGitRepo)
			row["Description"] = component.NewText(describeVolumeSource(volume.VolumeSource.GitRepo))
		case volume.VolumeSource.Secret != nil:
			row["Kind"] = component.NewText(volumeKindSecret)
			row["Description"] = component.NewText(describeVolumeSource(volume.VolumeSource.Secret))
		case volume.VolumeSource.ConfigMap != nil:
			row["Kind"] = component.NewText(volumeKindConfigMap)
			row["Description"] = component.NewText(describeVolumeSource(volume.VolumeSource.ConfigMap))
		case volume.VolumeSource.NFS != nil:
			row["Kind"] = component.NewText(volumeKindNFS)
			row["Description"] = component.NewText(describeVolumeSource(volume.VolumeSource.NFS))
		case volume.VolumeSource.ISCSI != nil:
			row["Kind"] = component.NewText(volumeKindISCSI)
			row["Description"] = component.NewText(describeVolumeSource(volume.VolumeSource.ISCSI))
		case volume.VolumeSource.Glusterfs != nil:
			row["Kind"] = component.NewText(volumeKindGlusterfs)
			row["Description"] = component.NewText(describeVolumeSource(volume.VolumeSource.Glusterfs))
		case volume.VolumeSource.PersistentVolumeClaim != nil:
			row["Kind"] = component.NewText(volumeKindPersistentVolumeClaim)
			row["Description"] = component.NewText(describeVolumeSource(volume.VolumeSource.PersistentVolumeClaim))
		case volume.VolumeSource.RBD != nil:
			row["Kind"] = component.NewText(volumeKindRBD)
			row["Description"] = component.NewText(describeVolumeSource(volume.VolumeSource.RBD))
		case volume.VolumeSource.Quobyte != nil:
			row["Kind"] = component.NewText(volumeKindQuobyte)
			row["Description"] = component.NewText(describeVolumeSource(volume.VolumeSource.Quobyte))
		case volume.VolumeSource.DownwardAPI != nil:
			row["Kind"] = component.NewText(volumeKindDownwardAPI)
			row["Description"] = component.NewText(describeVolumeSource(volume.VolumeSource.DownwardAPI))
		case volume.VolumeSource.AzureDisk != nil:
			row["Kind"] = component.NewText(volumeKindAzureDisk)
			row["Description"] = component.NewText(describeVolumeSource(volume.VolumeSource.AzureDisk))
		case volume.VolumeSource.VsphereVolume != nil:
			row["Kind"] = component.NewText(volumeKindSphereVolume)
			row["Description"] = component.NewText(describeVolumeSource(volume.VolumeSource.VsphereVolume))
		case volume.VolumeSource.Cinder != nil:
			row["Kind"] = component.NewText(volumeKindCinder)
			row["Description"] = component.NewText(describeVolumeSource(volume.VolumeSource.Cinder))
		case volume.VolumeSource.PhotonPersistentDisk != nil:
			row["Kind"] = component.NewText(volumeKindPhotonPersistentDisk)
			row["Description"] = component.NewText(describeVolumeSource(volume.VolumeSource.PhotonPersistentDisk))
		case volume.VolumeSource.PortworxVolume != nil:
			row["Kind"] = component.NewText(volumeKindPortworxVolume)
			row["Description"] = component.NewText(describeVolumeSource(volume.VolumeSource.PortworxVolume))
		case volume.VolumeSource.ScaleIO != nil:
			row["Kind"] = component.NewText(volumeKindScaleIO)
			row["Description"] = component.NewText(describeVolumeSource(volume.VolumeSource.ScaleIO))
		case volume.VolumeSource.CephFS != nil:
			row["Kind"] = component.NewText(volumeKindCephFS)
			row["Description"] = component.NewText(describeVolumeSource(volume.VolumeSource.CephFS))
		case volume.VolumeSource.StorageOS != nil:
			row["Kind"] = component.NewText(volumeKindStorageOS)
			row["Description"] = component.NewText(describeVolumeSource(volume.VolumeSource.StorageOS))
		case volume.VolumeSource.FC != nil:
			row["Kind"] = component.NewText(volumeKindFC)
			row["Description"] = component.NewText(describeVolumeSource(volume.VolumeSource.FC))
		case volume.VolumeSource.AzureFile != nil:
			row["Kind"] = component.NewText(volumeKindAzureFile)
			row["Description"] = component.NewText(describeVolumeSource(volume.VolumeSource.AzureFile))
		case volume.VolumeSource.FlexVolume != nil:
			row["Kind"] = component.NewText(volumeKindFlexVolume)
			row["Description"] = component.NewText(describeVolumeSource(volume.VolumeSource.FlexVolume))
		case volume.VolumeSource.Flocker != nil:
			row["Kind"] = component.NewText(volumeKindFlocker)
			row["Description"] = component.NewText(describeVolumeSource(volume.VolumeSource.Flocker))
		default:
			row["Kind"] = component.NewText(volumeKindUnknown)
			row["Description"] = component.NewText("")
		}

		table.Add(row)
	}

	return table, nil
}

func describeVolumeSource(source interface{}) string {
	data, _ := json.Marshal(source)
	return string(data)
}
