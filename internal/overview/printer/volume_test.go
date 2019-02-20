package printer

import (
	"testing"

	"github.com/heptio/developer-dash/internal/view/component"

	corev1 "k8s.io/api/core/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_printVolumes(t *testing.T) {
	cases := []struct {
		name     string
		volume   corev1.Volume
		expected component.TableRow
	}{
		{
			name: "hostpath",
			volume: corev1.Volume{
				Name: "hostpath-volume",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{Path: "/"},
				},
			},
			expected: component.TableRow{
				"Name":        component.NewText("hostpath-volume"),
				"Kind":        component.NewText(volumeKindHostPath),
				"Description": component.NewText("{\"path\":\"/\"}"),
			},
		},
		{
			name: "emptydir",
			volume: corev1.Volume{
				Name: "emptydir-volume",
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				},
			},
			expected: component.TableRow{
				"Name":        component.NewText("emptydir-volume"),
				"Kind":        component.NewText(volumeKindEmptyDir),
				"Description": component.NewText("{}"),
			},
		},
		{
			name: "gce persistent disk",
			volume: corev1.Volume{
				Name: "gcePersistentDisk-volume",
				VolumeSource: corev1.VolumeSource{
					GCEPersistentDisk: &corev1.GCEPersistentDiskVolumeSource{
						PDName: "pd",
					},
				},
			},
			expected: component.TableRow{
				"Name":        component.NewText("gcePersistentDisk-volume"),
				"Kind":        component.NewText(volumeKindGCEPersistentDisk),
				"Description": component.NewText("{\"pdName\":\"pd\"}"),
			},
		},
		{
			name: "aws ebs",
			volume: corev1.Volume{
				Name: "ebs-volume",
				VolumeSource: corev1.VolumeSource{
					AWSElasticBlockStore: &corev1.AWSElasticBlockStoreVolumeSource{
						VolumeID: "vol-314159",
						FSType:   "ext4",
					},
				},
			},
			expected: component.TableRow{
				"Name":        component.NewText("ebs-volume"),
				"Kind":        component.NewText(volumeKindAWSElasticBlockStore),
				"Description": component.NewText("{\"volumeID\":\"vol-314159\",\"fsType\":\"ext4\"}"),
			},
		},
		{
			name: "git repo",
			volume: corev1.Volume{
				Name: "gitRepo-volume",
				VolumeSource: corev1.VolumeSource{
					GitRepo: &corev1.GitRepoVolumeSource{},
				},
			},
			expected: component.TableRow{
				"Name":        component.NewText("gitRepo-volume"),
				"Kind":        component.NewText(volumeKindGitRepo),
				"Description": component.NewText("{\"repository\":\"\"}"),
			},
		},
		{
			name: "secret",
			volume: corev1.Volume{
				Name: "secret-volume",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{},
				},
			},
			expected: component.TableRow{
				"Name":        component.NewText("secret-volume"),
				"Kind":        component.NewText(volumeKindSecret),
				"Description": component.NewText("{}"),
			},
		},
		{
			name: "config map",
			volume: corev1.Volume{
				Name: "configMap-volume",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{},
				},
			},
			expected: component.TableRow{
				"Name":        component.NewText("configMap-volume"),
				"Kind":        component.NewText(volumeKindConfigMap),
				"Description": component.NewText("{}"),
			},
		},
		{
			name: "nfs",
			volume: corev1.Volume{
				Name: "nfs-volume",
				VolumeSource: corev1.VolumeSource{
					NFS: &corev1.NFSVolumeSource{},
				},
			},
			expected: component.TableRow{
				"Name":        component.NewText("nfs-volume"),
				"Kind":        component.NewText(volumeKindNFS),
				"Description": component.NewText("{\"server\":\"\",\"path\":\"\"}"),
			},
		},
		{
			name: "iscsi",
			volume: corev1.Volume{
				Name: "iscsi-volume",
				VolumeSource: corev1.VolumeSource{
					ISCSI: &corev1.ISCSIVolumeSource{},
				},
			},
			expected: component.TableRow{
				"Name":        component.NewText("iscsi-volume"),
				"Kind":        component.NewText(volumeKindISCSI),
				"Description": component.NewText("{\"targetPortal\":\"\",\"iqn\":\"\",\"lun\":0}"),
			},
		},
		{
			name: "gluster",
			volume: corev1.Volume{
				Name: "gluster-volume",
				VolumeSource: corev1.VolumeSource{
					Glusterfs: &corev1.GlusterfsVolumeSource{},
				},
			},
			expected: component.TableRow{
				"Name":        component.NewText("gluster-volume"),
				"Kind":        component.NewText(volumeKindGlusterfs),
				"Description": component.NewText("{\"endpoints\":\"\",\"path\":\"\"}"),
			},
		},
		{
			name: "pvc",
			volume: corev1.Volume{
				Name: "pvc-volume",
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{},
				},
			},
			expected: component.TableRow{
				"Name":        component.NewText("pvc-volume"),
				"Kind":        component.NewText(volumeKindPersistentVolumeClaim),
				"Description": component.NewText("{\"claimName\":\"\"}"),
			},
		},
		{
			name: "rbd",
			volume: corev1.Volume{
				Name: "rbd-volume",
				VolumeSource: corev1.VolumeSource{
					RBD: &corev1.RBDVolumeSource{},
				},
			},
			expected: component.TableRow{
				"Name":        component.NewText("rbd-volume"),
				"Kind":        component.NewText(volumeKindRBD),
				"Description": component.NewText("{\"monitors\":null,\"image\":\"\"}"),
			},
		},
		{
			name: "quobyte",
			volume: corev1.Volume{
				Name: "quobyte-volume",
				VolumeSource: corev1.VolumeSource{
					Quobyte: &corev1.QuobyteVolumeSource{},
				},
			},
			expected: component.TableRow{
				"Name":        component.NewText("quobyte-volume"),
				"Kind":        component.NewText(volumeKindQuobyte),
				"Description": component.NewText("{\"registry\":\"\",\"volume\":\"\"}"),
			},
		},
		{
			name: "downward",
			volume: corev1.Volume{
				Name: "downward-volume",
				VolumeSource: corev1.VolumeSource{
					DownwardAPI: &corev1.DownwardAPIVolumeSource{},
				},
			},
			expected: component.TableRow{
				"Name":        component.NewText("downward-volume"),
				"Kind":        component.NewText(volumeKindDownwardAPI),
				"Description": component.NewText("{}"),
			},
		},
		{
			name: "azure disk",
			volume: corev1.Volume{
				Name: "azureDisk-volume",
				VolumeSource: corev1.VolumeSource{
					AzureDisk: &corev1.AzureDiskVolumeSource{},
				},
			},
			expected: component.TableRow{
				"Name":        component.NewText("azureDisk-volume"),
				"Kind":        component.NewText(volumeKindAzureDisk),
				"Description": component.NewText("{\"diskName\":\"\",\"diskURI\":\"\"}"),
			},
		},
		{
			name: "vsphere",
			volume: corev1.Volume{
				Name: "vsphere-volume",
				VolumeSource: corev1.VolumeSource{
					VsphereVolume: &corev1.VsphereVirtualDiskVolumeSource{},
				},
			},

			expected: component.TableRow{
				"Name":        component.NewText("vsphere-volume"),
				"Kind":        component.NewText(volumeKindSphereVolume),
				"Description": component.NewText("{\"volumePath\":\"\"}"),
			},
		},
		{
			name: "cinder",
			volume: corev1.Volume{
				Name: "cinder-volume",
				VolumeSource: corev1.VolumeSource{
					Cinder: &corev1.CinderVolumeSource{},
				},
			},

			expected: component.TableRow{
				"Name":        component.NewText("cinder-volume"),
				"Kind":        component.NewText(volumeKindCinder),
				"Description": component.NewText("{\"volumeID\":\"\"}"),
			},
		},
		{
			name: "photon",
			volume: corev1.Volume{
				Name: "photon-volume",
				VolumeSource: corev1.VolumeSource{
					PhotonPersistentDisk: &corev1.PhotonPersistentDiskVolumeSource{},
				},
			},

			expected: component.TableRow{
				"Name":        component.NewText("photon-volume"),
				"Kind":        component.NewText(volumeKindPhotonPersistentDisk),
				"Description": component.NewText("{\"pdID\":\"\"}"),
			},
		},
		{
			name: "portworx",
			volume: corev1.Volume{
				Name: "portworx-volume",
				VolumeSource: corev1.VolumeSource{
					PortworxVolume: &corev1.PortworxVolumeSource{},
				},
			},

			expected: component.TableRow{
				"Name":        component.NewText("portworx-volume"),
				"Kind":        component.NewText(volumeKindPortworxVolume),
				"Description": component.NewText("{\"volumeID\":\"\"}"),
			},
		},
		{
			name: "scaleIO",
			volume: corev1.Volume{
				Name: "scaleIO-volume",
				VolumeSource: corev1.VolumeSource{
					ScaleIO: &corev1.ScaleIOVolumeSource{},
				},
			},

			expected: component.TableRow{
				"Name":        component.NewText("scaleIO-volume"),
				"Kind":        component.NewText(volumeKindScaleIO),
				"Description": component.NewText("{\"gateway\":\"\",\"system\":\"\",\"secretRef\":null}"),
			},
		},
		{
			name: "fc",
			volume: corev1.Volume{
				Name: "fc-volume",
				VolumeSource: corev1.VolumeSource{
					FC: &corev1.FCVolumeSource{},
				},
			},

			expected: component.TableRow{
				"Name":        component.NewText("fc-volume"),
				"Kind":        component.NewText(volumeKindFC),
				"Description": component.NewText("{}"),
			},
		},
		{
			name: "azure file",
			volume: corev1.Volume{
				Name: "azureFile-volume",
				VolumeSource: corev1.VolumeSource{
					AzureFile: &corev1.AzureFileVolumeSource{},
				},
			},

			expected: component.TableRow{
				"Name":        component.NewText("azureFile-volume"),
				"Kind":        component.NewText(volumeKindAzureFile),
				"Description": component.NewText("{\"secretName\":\"\",\"shareName\":\"\"}"),
			},
		},
		{
			name: "flex",
			volume: corev1.Volume{
				Name: "flex-volume",
				VolumeSource: corev1.VolumeSource{
					FlexVolume: &corev1.FlexVolumeSource{},
				},
			},

			expected: component.TableRow{
				"Name":        component.NewText("flex-volume"),
				"Kind":        component.NewText(volumeKindFlexVolume),
				"Description": component.NewText("{\"driver\":\"\"}"),
			},
		},
		{
			name: "flocker",
			volume: corev1.Volume{
				Name: "flocker-volume",
				VolumeSource: corev1.VolumeSource{
					Flocker: &corev1.FlockerVolumeSource{},
				},
			},

			expected: component.TableRow{
				"Name":        component.NewText("flocker-volume"),
				"Kind":        component.NewText(volumeKindFlocker),
				"Description": component.NewText("{}"),
			},
		},
		{
			name: "ceph",
			volume: corev1.Volume{
				Name: "ceph-volume",
				VolumeSource: corev1.VolumeSource{
					CephFS: &corev1.CephFSVolumeSource{},
				},
			},

			expected: component.TableRow{
				"Name":        component.NewText("ceph-volume"),
				"Kind":        component.NewText(volumeKindCephFS),
				"Description": component.NewText("{\"monitors\":null}"),
			},
		},
		{
			name: "storage OS",
			volume: corev1.Volume{
				Name: "storageOS-volume",
				VolumeSource: corev1.VolumeSource{
					StorageOS: &corev1.StorageOSVolumeSource{},
				},
			},

			expected: component.TableRow{
				"Name":        component.NewText("storageOS-volume"),
				"Kind":        component.NewText(volumeKindStorageOS),
				"Description": component.NewText("{}"),
			},
		},
		{
			name: "unknown",
			volume: corev1.Volume{
				Name:         "unknown",
				VolumeSource: corev1.VolumeSource{},
			},

			expected: component.TableRow{
				"Name":        component.NewText("unknown"),
				"Kind":        component.NewText(volumeKindUnknown),
				"Description": component.NewText(""),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := printVolumes([]corev1.Volume{tc.volume})
			require.NoError(t, err)

			expected := component.NewTableWithRows("Volumes",
				component.NewTableCols("Name", "Kind", "Description"),
				[]component.TableRow{tc.expected})
			assert.Equal(t, expected, got)
		})
	}
}
