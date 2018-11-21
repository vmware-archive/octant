package overview

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/heptio/developer-dash/internal/content"
	"k8s.io/kubernetes/pkg/apis/core"
)

func summarizeVolume(volume core.Volume) content.Section {
	section := content.NewSection()
	section.Title = volume.Name

	switch {
	case volume.VolumeSource.HostPath != nil:
		summarizeHostPathVolumeSource(&section, volume.VolumeSource.HostPath)
	case volume.VolumeSource.EmptyDir != nil:
		summarizeEmptyDirVolumeSource(&section, volume.VolumeSource.EmptyDir)
	case volume.VolumeSource.GCEPersistentDisk != nil:
		summarizeGCEPersistentDiskVolumeSource(&section, volume.VolumeSource.GCEPersistentDisk)
	case volume.VolumeSource.AWSElasticBlockStore != nil:
		summarizeAWSElasticBlockStoreVolumeSource(&section, volume.VolumeSource.AWSElasticBlockStore)
	case volume.VolumeSource.GitRepo != nil:
		summarizeGitRepoVolumeSource(&section, volume.VolumeSource.GitRepo)
	case volume.VolumeSource.Secret != nil:
		summarizeSecretVolumeSource(&section, volume.VolumeSource.Secret)
	case volume.VolumeSource.ConfigMap != nil:
		summarizeConfigMapVolumeSource(&section, volume.VolumeSource.ConfigMap)
	case volume.VolumeSource.NFS != nil:
		summarizeNFSVolumeSource(&section, volume.VolumeSource.NFS)
	case volume.VolumeSource.ISCSI != nil:
		summarizeISCSIVolumeSource(&section, volume.VolumeSource.ISCSI)
	case volume.VolumeSource.Glusterfs != nil:
		summarizeGlusterfsVolumeSource(&section, volume.VolumeSource.Glusterfs)
	case volume.VolumeSource.PersistentVolumeClaim != nil:
		summarizePersistentVolumeClaimVolumeSource(&section, volume.VolumeSource.PersistentVolumeClaim)
	case volume.VolumeSource.RBD != nil:
		summarizeRBDVolumeSource(&section, volume.VolumeSource.RBD)
	case volume.VolumeSource.Quobyte != nil:
		summarizeQuobyteVolumeSource(&section, volume.VolumeSource.Quobyte)
	case volume.VolumeSource.DownwardAPI != nil:
		summarizeDownwardAPIVolumeSource(&section, volume.VolumeSource.DownwardAPI)
	case volume.VolumeSource.AzureDisk != nil:
		summarizeAzureDiskVolumeSource(&section, volume.VolumeSource.AzureDisk)
	case volume.VolumeSource.VsphereVolume != nil:
		summarizeVsphereVolumeSource(&section, volume.VolumeSource.VsphereVolume)
	case volume.VolumeSource.Cinder != nil:
		summarizeCinderVolumeSource(&section, volume.VolumeSource.Cinder)
	case volume.VolumeSource.PhotonPersistentDisk != nil:
		summarizePhotonPersistentDiskVolumeSource(&section, volume.VolumeSource.PhotonPersistentDisk)
	case volume.VolumeSource.PortworxVolume != nil:
		summarizePortworxVolumeSource(&section, volume.VolumeSource.PortworxVolume)
	case volume.VolumeSource.ScaleIO != nil:
		summarizeScaleIOVolumeSource(&section, volume.VolumeSource.ScaleIO)
	case volume.VolumeSource.CephFS != nil:
		summarizeCephFSVolumeSource(&section, volume.VolumeSource.CephFS)
	case volume.VolumeSource.StorageOS != nil:
		summarizeStorageOSVolumeSource(&section, volume.VolumeSource.StorageOS)
	case volume.VolumeSource.FC != nil:
		summarizeFCVolumeSource(&section, volume.VolumeSource.FC)
	case volume.VolumeSource.AzureFile != nil:
		summarizeAzureFileVolumeSource(&section, volume.VolumeSource.AzureFile)
	case volume.VolumeSource.FlexVolume != nil:
		summarizeFlexVolumeSource(&section, volume.VolumeSource.FlexVolume)
	case volume.VolumeSource.Flocker != nil:
		summarizeFlockerVolumeSource(&section, volume.VolumeSource.Flocker)
	default:
		section.AddText("Type", "<unknown>")
	}

	return section
}

func summarizeHostPathVolumeSource(section *content.Section, hostPath *core.HostPathVolumeSource) {
	hostPathType := "<none>"
	if hostPath.Type != nil {
		hostPathType = string(*hostPath.Type)
	}

	section.AddText("Type", "HostPath (bare host directory volume)")
	section.AddText("Path", hostPath.Path)
	section.AddText("HostPathType", hostPathType)
}

func summarizeEmptyDirVolumeSource(section *content.Section, emptyDir *core.EmptyDirVolumeSource) {
	section.AddText("Type", "EmptyDir (a temporary directory that shares a pod's lifetime)")
	section.AddText("Medium", string(emptyDir.Medium))
}

func summarizeGCEPersistentDiskVolumeSource(section *content.Section, gce *core.GCEPersistentDiskVolumeSource) {
	section.AddText("Type", "GCEPersistentDisk (a Persistent Disk resource in Google Compute Engine)")
	section.AddText("PDName", gce.PDName)
	section.AddText("FSType", gce.FSType)
	section.AddText("Partition", strconv.Itoa(int(gce.Partition)))
	section.AddText("ReadOnly", fmt.Sprintf("%v", gce.ReadOnly))
}

func summarizeAWSElasticBlockStoreVolumeSource(section *content.Section, aws *core.AWSElasticBlockStoreVolumeSource) {
	section.AddText("Type", "AWSElasticBlockStore (a Persistent Disk resource in AWS)")
	section.AddText("VolumeID", aws.VolumeID)
	section.AddText("FSType", aws.FSType)
	section.AddText("Partition", strconv.Itoa(int(aws.Partition)))
	section.AddText("ReadOnly", fmt.Sprintf("%v", aws.ReadOnly))
}

func summarizeGitRepoVolumeSource(section *content.Section, git *core.GitRepoVolumeSource) {
	section.AddText("Type", "GitRepo (a volume that is pulled from git when the pod is created)")
	section.AddText("Repository", git.Repository)
	section.AddText("Revision", git.Revision)
}

func summarizeSecretVolumeSource(section *content.Section, secret *core.SecretVolumeSource) {
	optional := secret.Optional != nil && *secret.Optional
	section.AddText("Type", "Secret (a volume populated by a Secret)")
	section.AddText("SecretName", secret.SecretName)
	section.AddText("Optional", fmt.Sprintf("%v", optional))
}

func summarizeConfigMapVolumeSource(section *content.Section, configMap *core.ConfigMapVolumeSource) {
	optional := configMap.Optional != nil && *configMap.Optional
	section.AddText("Type", "ConfigMapSecret (a volume populated by a ConfigMap)")
	section.AddText("Name", configMap.Name)
	section.AddText("Optional", fmt.Sprintf("%v", optional))
}

func summarizeNFSVolumeSource(section *content.Section, nfs *core.NFSVolumeSource) {
	section.AddText("Type", "NFS (an NFS mount that lasts the lifetime of a pod)")
	section.AddText("Server", nfs.Server)
	section.AddText("Path", nfs.Path)
	section.AddText("ReadOnly", fmt.Sprintf("%v", nfs.ReadOnly))
}

func summarizeQuobyteVolumeSource(section *content.Section, quobyte *core.QuobyteVolumeSource) {
	section.AddText("Type", "Quobyte (a Quobyte mount on the host that shares a pod's lifetime)")
	section.AddText("Registry", quobyte.Registry)
	section.AddText("Volume", quobyte.Volume)
	section.AddText("ReadOnly", fmt.Sprintf("%v", quobyte.ReadOnly))
}

func summarizePortworxVolumeSource(section *content.Section, pwxVolume *core.PortworxVolumeSource) {
	section.AddText("Type", "PortworxVolume (a Portworx Volume resource)")
	section.AddText("VolumeID", pwxVolume.VolumeID)
}

func summarizeISCSIVolumeSource(section *content.Section, iscsi *core.ISCSIVolumeSource) {
	initiator := "<none>"
	if iscsi.InitiatorName != nil {
		initiator = *iscsi.InitiatorName
	}
	secretRef := printObjectRef(iscsi.SecretRef)
	section.AddText("Type", "ISCSI (an ISCSI Disk resource that is attached to a kubelet's host machine and then exposed to the pod)")
	section.AddText("TargetPortal", iscsi.TargetPortal)
	section.AddText("IQN", iscsi.IQN)
	section.AddText("Lun", fmt.Sprintf("%v", iscsi.Lun))
	section.AddText("ISCSIInterface", iscsi.ISCSIInterface)
	section.AddText("FSType", iscsi.FSType)
	section.AddText("ReadOnly", fmt.Sprintf("%v", iscsi.ReadOnly))
	section.AddText("Portals", strings.Join(iscsi.Portals, ","))
	section.AddText("DiscoveryCHAPAuth", fmt.Sprintf("%v", iscsi.DiscoveryCHAPAuth))
	section.AddText("SessionCHAPAuth", fmt.Sprintf("%v", iscsi.SessionCHAPAuth))
	section.AddText("SecretRef", secretRef)
	section.AddText("InitiatorName", initiator)
}

func summarizeISCSIPersistentVolumeSource(section *content.Section, iscsi *core.ISCSIPersistentVolumeSource) {
	initiator := "<none>"
	if iscsi.InitiatorName != nil {
		initiator = *iscsi.InitiatorName
	}
	secretRef := printSecretRef(iscsi.SecretRef)
	section.AddText("Type", "ISCSI (an ISCSI Disk resource that is attached to a kubelet's host machine and then exposed to the pod)")
	section.AddText("TargetPortal", iscsi.TargetPortal)
	section.AddText("IQN", iscsi.IQN)
	section.AddText("Lun", fmt.Sprintf("%v", iscsi.Lun))
	section.AddText("ISCSIInterface", iscsi.ISCSIInterface)
	section.AddText("FSType", iscsi.FSType)
	section.AddText("ReadOnly", fmt.Sprintf("%v", iscsi.ReadOnly))
	section.AddText("Portals", strings.Join(iscsi.Portals, ","))
	section.AddText("DiscoveryCHAPAuth", fmt.Sprintf("%v", iscsi.DiscoveryCHAPAuth))
	section.AddText("SessionCHAPAuth", fmt.Sprintf("%v", iscsi.SessionCHAPAuth))
	section.AddText("SecretRef", secretRef)
	section.AddText("InitiatorName", initiator)
}

func summarizeGlusterfsVolumeSource(section *content.Section, glusterfs *core.GlusterfsVolumeSource) {
	section.AddText("Type", "Glusterfs (a Glusterfs mount on the host that shares a pod's lifetime)")
	section.AddText("EndpointsName", glusterfs.EndpointsName)
	section.AddText("Path", glusterfs.Path)
	section.AddText("ReadOnly", fmt.Sprintf("%v", glusterfs.ReadOnly))
}

func summarizePersistentVolumeClaimVolumeSource(section *content.Section, claim *core.PersistentVolumeClaimVolumeSource) {
	section.AddText("Type", "PersistentVolumeClaim")
	section.AddLink("ClaimName", claim.ClaimName, gvkPath("v1", "PersistentVolumeClaim", claim.ClaimName))
	section.AddText("ReadOnly", fmt.Sprintf("%t", claim.ReadOnly))
}

func summarizeRBDVolumeSource(section *content.Section, rbd *core.RBDVolumeSource) {
	secretRef := printObjectRef(rbd.SecretRef)
	section.AddText("Type", "RBD (a Rados Block Device mount on the host that shares a pod's lifetime)")
	section.AddText("CephMonitors", strings.Join(rbd.CephMonitors, ","))
	section.AddText("RBDImage", rbd.RBDImage)
	section.AddText("FSType", rbd.FSType)
	section.AddText("RadosUser", rbd.RadosUser)
	section.AddText("Keyring", rbd.Keyring)
	section.AddText("SecretRef", secretRef)
	section.AddText("ReadOnly", fmt.Sprintf("%v", rbd.ReadOnly))
}

func summarizeRBDPersistentVolumeSource(section *content.Section, rbd *core.RBDPersistentVolumeSource) {
	secretRef := printSecretRef(rbd.SecretRef)
	section.AddText("Type", "RBD (a Rados Block Device mount on the host that shares a pod's lifetime)")
	section.AddText("CephMonitors", strings.Join(rbd.CephMonitors, ","))
	section.AddText("RBDImage", rbd.RBDImage)
	section.AddText("FSType", rbd.FSType)
	section.AddText("RadosUser", rbd.RadosUser)
	section.AddText("Keyring", rbd.Keyring)
	section.AddText("SecretRef", secretRef)
	section.AddText("ReadOnly", fmt.Sprintf("%v", rbd.ReadOnly))
}

func summarizeDownwardAPIVolumeSource(section *content.Section, d *core.DownwardAPIVolumeSource) {
	section.AddText("Type", "DownwardAPI (a volume populated by information about the pod)")
	for _, mapping := range d.Items {
		if mapping.FieldRef != nil {
			section.AddText(mapping.FieldRef.FieldPath, mapping.Path)
		}
		if mapping.ResourceFieldRef != nil {
			section.AddText(mapping.ResourceFieldRef.Resource, mapping.Path)
		}
	}
}

func summarizeAzureDiskVolumeSource(section *content.Section, d *core.AzureDiskVolumeSource) {
	kind := ""
	if d.Kind != nil {
		kind = string(*d.Kind)
	}
	fsType := ""
	if d.FSType != nil {
		fsType = *d.FSType
	}
	cachingMode := ""
	if d.CachingMode != nil {
		cachingMode = string(*d.CachingMode)
	}

	section.AddText("Type", "AzureDisk (an Azure Data Disk mount on the host and bind mount to the pod)")
	section.AddText("DiskName", d.DiskName)
	section.AddText("DiskURI", d.DataDiskURI)
	section.AddText("Kind", kind)
	section.AddText("FSType", fsType)
	section.AddText("CachingMode", cachingMode)
	section.AddText("ReadOnly", fmt.Sprintf("%v", *d.ReadOnly))
}

func summarizeVsphereVolumeSource(section *content.Section, vsphere *core.VsphereVirtualDiskVolumeSource) {
	section.AddText("Type", "SphereVolume (a Persistent Disk resource in vSphere)")
	section.AddText("VolumePath", vsphere.VolumePath)
	section.AddText("FSType", vsphere.FSType)
	section.AddText("StoragePolicyName", vsphere.StoragePolicyName)
}

func summarizePhotonPersistentDiskVolumeSource(section *content.Section, photon *core.PhotonPersistentDiskVolumeSource) {
	section.AddText("Type", "PhotonPersistentDisk (a Persistent Disk resource in photon platform)")
	section.AddText("PdID", photon.PdID)
	section.AddText("FSType", photon.FSType)
}

func summarizeCinderVolumeSource(section *content.Section, cinder *core.CinderVolumeSource) {
	secretRef := printObjectRef(cinder.SecretRef)
	section.AddText("Type", "Cinder (a Persistent Disk resource in OpenStack)")
	section.AddText("VolumeID", cinder.VolumeID)
	section.AddText("FSType", cinder.FSType)
	section.AddText("ReadOnly", fmt.Sprintf("%v", cinder.ReadOnly))
	section.AddText("SecretRef", secretRef)
}

func summarizeCinderPersistentVolumeSource(section *content.Section, cinder *core.CinderPersistentVolumeSource) {
	secretRef := printSecretRef(cinder.SecretRef)
	section.AddText("Type", "Cinder (a Persistent Disk resource in OpenStack)")
	section.AddText("VolumeID", cinder.VolumeID)
	section.AddText("FSType", cinder.FSType)
	section.AddText("ReadOnly", fmt.Sprintf("%v", cinder.ReadOnly))
	section.AddText("SecretRef", secretRef)
}

func summarizeScaleIOVolumeSource(section *content.Section, sio *core.ScaleIOVolumeSource) {
	section.AddText("Type", "ScaleIO (a persistent volume backed by a block device in ScaleIO)")
	section.AddText("Gateway", sio.Gateway)
	section.AddText("System", sio.System)
	section.AddText("Protection Domain", sio.ProtectionDomain)
	section.AddText("Storage Pool", sio.StoragePool)
	section.AddText("Storage Mode", sio.StorageMode)
	section.AddText("VolumeName", sio.VolumeName)
	section.AddText("FSType", sio.FSType)
	section.AddText("ReadOnly", fmt.Sprintf("%v", sio.ReadOnly))
}

func summarizeScaleIOPersistentVolumeSource(section *content.Section, sio *core.ScaleIOPersistentVolumeSource) {
	section.AddText("Type", "ScaleIO (a persistent volume backed by a block device in ScaleIO)")
	section.AddText("Gateway", sio.Gateway)
	section.AddText("System", sio.System)
	section.AddText("Protection Domain", sio.ProtectionDomain)
	section.AddText("Storage Pool", sio.StoragePool)
	section.AddText("Storage Mode", sio.StorageMode)
	section.AddText("VolumeName", sio.VolumeName)
	section.AddText("FSType", sio.FSType)
	section.AddText("ReadOnly", fmt.Sprintf("%v", sio.ReadOnly))
}

func summarizeCephFSVolumeSource(section *content.Section, cephfs *core.CephFSVolumeSource) {
	secretRef := printObjectRef(cephfs.SecretRef)

	section.AddText("Type", "CephFS (a CephFS mount on the host that shares a pod's lifetime)")
	section.AddText("Monitors", strings.Join(cephfs.Monitors, ""))
	section.AddText("Path", cephfs.Path)
	section.AddText("User", cephfs.User)
	section.AddText("SecretFile", cephfs.SecretFile)
	section.AddText("SecretRef", secretRef)
	section.AddText("ReadOnly", fmt.Sprintf("%v", cephfs.ReadOnly))
}

func summarizeStorageOSVolumeSource(section *content.Section, storageos *core.StorageOSVolumeSource) {
	section.AddText("Type", "StorageOS (a StorageOS Persistent Disk resource)")
	section.AddText("VolumeName", storageos.VolumeName)
	section.AddText("VolumeNamespace", storageos.VolumeNamespace)
	section.AddText("FSType", storageos.FSType)
	section.AddText("ReadOnly", fmt.Sprintf("%v", storageos.ReadOnly))
}

func summarizeFCVolumeSource(section *content.Section, fc *core.FCVolumeSource) {
	lun := "<none>"
	if fc.Lun != nil {
		lun = strconv.Itoa(int(*fc.Lun))
	}

	section.AddText("Type", "FC (a Fibre Channel disk)")
	section.AddText("TargetWWMNs", strings.Join(fc.TargetWWNs, ","))
	section.AddText("LUN", lun)
	section.AddText("FSType", fc.FSType)
	section.AddText("ReadOnly", fmt.Sprintf("%v", fc.ReadOnly))
}

func summarizeAzureFileVolumeSource(section *content.Section, azureFile *core.AzureFileVolumeSource) {
	section.AddText("Type", "AzureFile (an Azure File Service mount on the host and bind mount to the pod)")
	section.AddText("SecretName", azureFile.SecretName)
	section.AddText("ShareName", azureFile.ShareName)
	section.AddText("ReadOnly", fmt.Sprintf("%v", azureFile.ReadOnly))
}

func summarizeAzureFilePersistentVolumeSource(section *content.Section, azureFile *core.AzureFilePersistentVolumeSource) {
	ns := ""
	if azureFile.SecretNamespace != nil {
		ns = *azureFile.SecretNamespace
	}
	section.AddText("Type", "AzureFile (an Azure File Service mount on the host and bind mount to the pod)")
	section.AddText("SecretName", azureFile.SecretName)
	section.AddText("SecretNamespace", ns)
	section.AddText("ShareName", azureFile.ShareName)
	section.AddText("ReadOnly", fmt.Sprintf("%v", azureFile.ReadOnly))
}

func summarizeFlexPersistentVolumeSource(section *content.Section, flex *core.FlexPersistentVolumeSource) {
	secretRef := printSecretRef(flex.SecretRef)

	section.AddText("Type", "FlexVolume (a generic volume resource that is provisioned/attached using an exec based plugin)")
	section.AddText("Driver", flex.Driver)
	section.AddText("FSType", flex.FSType)
	section.AddText("SecretRef", secretRef)
	section.AddText("ReadOnly", fmt.Sprintf("%v", flex.ReadOnly))
	section.AddText("Options", fmt.Sprintf("%v", flex.Options)) // TODO
}

func summarizeFlexVolumeSource(section *content.Section, flex *core.FlexVolumeSource) {
	secretRef := printObjectRef(flex.SecretRef)

	section.AddText("Type", "FlexVolume (a generic volume resource that is provisioned/attached using an exec based plugin)")
	section.AddText("Driver", flex.Driver)
	section.AddText("FSType", flex.FSType)
	section.AddText("SecretRef", secretRef)
	section.AddText("ReadOnly", fmt.Sprintf("%v", flex.ReadOnly))
	section.AddText("Options", fmt.Sprintf("%v", flex.Options)) // TODO
}

func summarizeFlockerVolumeSource(section *content.Section, flocker *core.FlockerVolumeSource) {
	section.AddText("Type", "Flocker (a Flocker volume mounted by the Flocker agent)")
	section.AddText("DatasetName", flocker.DatasetName)
	section.AddText("DatasetUUID", flocker.DatasetUUID)
}

func summarizeCSIPersistentVolumeSource(section *content.Section, csi *core.CSIPersistentVolumeSource) {
	section.AddText("Type", "CSI (a Container Storage Interface (CSI) volume source)")
	section.AddText("Driver", csi.Driver)
	section.AddText("Driver", csi.Driver)
	section.AddText("VolumeHandle", csi.VolumeHandle)
	section.AddText("ReadOnly", fmt.Sprintf("%v", csi.ReadOnly))
	// TODO TBD
	// summarizeCSIPersistentVolumeAttributesMultiline(w, "VolumeAttributes", csi.VolumeAttributes)
}

func printSecretRef(secretRef *core.SecretReference) string {
	if secretRef == nil {
		return "<none>"
	}
	return secretRef.Name
}

func printObjectRef(objectRef *core.LocalObjectReference) string {
	if objectRef == nil {
		return "<none>"
	}
	return objectRef.Name
}
