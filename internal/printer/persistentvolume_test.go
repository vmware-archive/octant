/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_PersistentVolumeListHandler(t *testing.T) {
	cols := component.NewTableCols("Name", "Capacity", "Access Modes", "Reclaim Policy", "Status", "Claim", "Storage Class", "Reason", "Age")
	now := testutil.Time()

	pvcObject := testutil.CreatePersistentVolumeClaim("pvc")

	object := testutil.CreatePersistentVolume("persistentVolume")
	object.CreationTimestamp = metav1.Time{Time: now}
	object.Spec.ClaimRef = &corev1.ObjectReference{
		APIVersion: pvcObject.APIVersion,
		Kind:       pvcObject.Kind,
		Name:       pvcObject.Name,
		Namespace:  pvcObject.Namespace,
	}

	list := &corev1.PersistentVolumeList{
		Items: []corev1.PersistentVolume{*object},
	}

	unbound := testutil.CreatePersistentVolume("unboundPersistentVolume")
	unbound.CreationTimestamp = metav1.Time{Time: now}
	unboundList := &corev1.PersistentVolumeList{
		Items: []corev1.PersistentVolume{*unbound},
	}

	cases := []struct {
		name     string
		list     *corev1.PersistentVolumeList
		expected *component.Table
		isErr    bool
	}{
		{
			name: "in general",
			list: list,
			expected: component.NewTableWithRows("Persistent Volumes", "We couldn't find any persistent volumes!", cols,
				[]component.TableRow{
					{
						"Name": component.NewLink("", "persistentVolume", "/persistentVolume",
							genObjectStatus(component.TextStatusOK, []string{
								"v1 PersistentVolume is OK",
							})),
						"Capacity":       component.NewText("0"),
						"Access Modes":   component.NewText(""),
						"Reclaim Policy": component.NewText(""),
						"Status":         component.NewText("Bound"),
						"Claim":          component.NewLink("", "namespace/pvc", "/pvc"),
						"Storage Class":  component.NewText(""),
						"Reason":         component.NewText(""),
						"Age":            component.NewTimestamp(now),
						component.GridActionKey: gridActionsFactory([]component.GridAction{
							buildObjectDeleteAction(t, object),
						}),
					},
				}),
		},
		{
			name: "unclaimed",
			list: unboundList,
			expected: component.NewTableWithRows("Persistent Volumes", "We couldn't find any persistent volumes!", cols,
				[]component.TableRow{
					{
						"Name": component.NewLink("", "unboundPersistentVolume", "/unboundPersistentVolume",
							genObjectStatus(component.TextStatusOK, []string{
								"v1 PersistentVolume is OK",
							})),
						"Capacity":       component.NewText("0"),
						"Access Modes":   component.NewText(""),
						"Reclaim Policy": component.NewText(""),
						"Status":         component.NewText("Bound"),
						"Claim":          component.NewLink("", "", ""),
						"Storage Class":  component.NewText(""),
						"Reason":         component.NewText(""),
						"Age":            component.NewTimestamp(now),
						component.GridActionKey: gridActionsFactory([]component.GridAction{
							buildObjectDeleteAction(t, unbound),
						}),
					},
				}),
		},
		{
			name:  "list is nil",
			list:  nil,
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			tpo := newTestPrinterOptions(controller)
			printOptions := tpo.ToOptions()

			ctx := context.Background()

			if tc.list != nil {
				tpo.PathForObject(&tc.list.Items[0], tc.list.Items[0].Name, "/"+tc.list.Items[0].Name)

				pvcKey, err := store.KeyFromObject(pvcObject)
				require.NoError(t, err)

				tpo.objectStore.EXPECT().Get(ctx, pvcKey).
					Return(testutil.ToUnstructured(t, pvcObject), nil).AnyTimes()

				tpo.PathForObject(pvcObject, pvcObject.Namespace+"/"+pvcObject.Name, "/pvc")
			}
			got, err := PersistentVolumeListHandler(ctx, tc.list, printOptions)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			component.AssertEqual(t, tc.expected, got)
		})
	}
}

func Test_PersistentVolumeConfiguration(t *testing.T) {
	persistentVolume := testutil.CreatePersistentVolume("persistentVolume")
	persistentVolume.Spec.PersistentVolumeSource = corev1.PersistentVolumeSource{
		HostPath: &corev1.HostPathVolumeSource{
			Path: "/private/tmp",
		},
	}

	cases := []struct {
		name             string
		persistentVolume *corev1.PersistentVolume
		expected         component.Component
		isErr            bool
	}{
		{
			name:             "host path",
			persistentVolume: persistentVolume,
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Reclaim Policy",
					Content: component.NewText(""),
				},
				{
					Header:  "Storage Class",
					Content: component.NewText(""),
				},
				{
					Header:  "Access Modes",
					Content: component.NewText(""),
				},
				{
					Header:  "Capacity",
					Content: component.NewText("0"),
				},
				{
					Header:  "Host Path",
					Content: component.NewText("{\"path\":\"/private/tmp\"}"),
				},
			}...),
		},
		{
			name:             "nil persistentVolume",
			persistentVolume: nil,
			isErr:            true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			tpo := newTestPrinterOptions(controller)
			printOptions := tpo.ToOptions()

			pc := NewPersistentVolumeConfiguration(tc.persistentVolume)

			summary, err := pc.Create(printOptions)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			component.AssertEqual(t, tc.expected, summary)
		})
	}
}

func Test_printPersistentVolumeSource(t *testing.T) {
	cases := []struct {
		name             string
		persistentVolume corev1.PersistentVolume
		expected         *component.Summary
	}{
		{
			name: "gce persistent disk",
			persistentVolume: corev1.PersistentVolume{
				Spec: corev1.PersistentVolumeSpec{
					PersistentVolumeSource: corev1.PersistentVolumeSource{
						GCEPersistentDisk: &corev1.GCEPersistentDiskVolumeSource{
							PDName: "my-data-disk",
						},
					},
				},
			},
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "GCE Persistent Disk",
					Content: component.NewText("{\"pdName\":\"my-data-disk\"}"),
				},
			}...),
		},
		{
			name: "aws ebs",
			persistentVolume: corev1.PersistentVolume{
				Spec: corev1.PersistentVolumeSpec{
					PersistentVolumeSource: corev1.PersistentVolumeSource{
						AWSElasticBlockStore: &corev1.AWSElasticBlockStoreVolumeSource{
							VolumeID: "vol-0b402cb854e070fad",
						},
					},
				},
			},
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "AWS Elastic Block Store",
					Content: component.NewText("{\"volumeID\":\"vol-0b402cb854e070fad\"}"),
				},
			}...),
		},
		{
			name: "host path",
			persistentVolume: corev1.PersistentVolume{
				Spec: corev1.PersistentVolumeSpec{
					PersistentVolumeSource: corev1.PersistentVolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/private/tmp",
						},
					},
				},
			},
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Host Path",
					Content: component.NewText("{\"path\":\"/private/tmp\"}"),
				},
			}...),
		},
		{
			name: "glusterfs",
			persistentVolume: corev1.PersistentVolume{
				Spec: corev1.PersistentVolumeSpec{
					PersistentVolumeSource: corev1.PersistentVolumeSource{
						Glusterfs: &corev1.GlusterfsPersistentVolumeSource{
							EndpointsName: "192.168.122.221:1",
							Path:          "kube_vol",
						},
					},
				},
			},
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "GlusterFS",
					Content: component.NewText("{\"endpoints\":\"192.168.122.221:1\",\"path\":\"kube_vol\"}"),
				},
			}...),
		},
		{
			name: "nfs",
			persistentVolume: corev1.PersistentVolume{
				Spec: corev1.PersistentVolumeSpec{
					PersistentVolumeSource: corev1.PersistentVolumeSource{
						NFS: &corev1.NFSVolumeSource{
							Server: "nfs.example.com",
							Path:   "/share1",
						},
					},
				},
			},
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "NFS",
					Content: component.NewText("{\"server\":\"nfs.example.com\",\"path\":\"/share1\"}"),
				},
			}...),
		},
		{
			name: "rbd",
			persistentVolume: corev1.PersistentVolume{
				Spec: corev1.PersistentVolumeSpec{
					PersistentVolumeSource: corev1.PersistentVolumeSource{
						RBD: &corev1.RBDPersistentVolumeSource{
							CephMonitors: []string{"10.16.154.78:6789", "10.16.154.82:6789"},
							RBDImage:     "foo",
						},
					},
				},
			},
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "RBD",
					Content: component.NewText("{\"monitors\":[\"10.16.154.78:6789\",\"10.16.154.82:6789\"],\"image\":\"foo\"}"),
				},
			}...),
		},
		{
			name: "isci",
			persistentVolume: corev1.PersistentVolume{
				Spec: corev1.PersistentVolumeSpec{
					PersistentVolumeSource: corev1.PersistentVolumeSource{
						ISCSI: &corev1.ISCSIPersistentVolumeSource{
							TargetPortal: "10.0.2.15:3260",
							IQN:          "iqn.2001-04.com.example:storage.kube.sys1.xyz",
							Lun:          0,
						},
					},
				},
			},
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "ISCI",
					Content: component.NewText("{\"targetPortal\":\"10.0.2.15:3260\",\"iqn\":\"iqn.2001-04.com.example:storage.kube.sys1.xyz\",\"lun\":0}"),
				},
			}...),
		},
		{
			name: "cinder",
			persistentVolume: corev1.PersistentVolume{
				Spec: corev1.PersistentVolumeSpec{
					PersistentVolumeSource: corev1.PersistentVolumeSource{
						Cinder: &corev1.CinderPersistentVolumeSource{
							VolumeID: "573e024d-5235-49ce-8332-be1576d323f8",
						},
					},
				},
			},
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Cinder",
					Content: component.NewText("{\"volumeID\":\"573e024d-5235-49ce-8332-be1576d323f8\"}"),
				},
			}...),
		},
		{
			name: "cephfs",
			persistentVolume: corev1.PersistentVolume{
				Spec: corev1.PersistentVolumeSpec{
					PersistentVolumeSource: corev1.PersistentVolumeSource{
						CephFS: &corev1.CephFSPersistentVolumeSource{
							Monitors: []string{"10.16.154.78:6789", "10.16.154.82:6789"},
						},
					},
				},
			},
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "CephFS",
					Content: component.NewText("{\"monitors\":[\"10.16.154.78:6789\",\"10.16.154.82:6789\"]}"),
				},
			}...),
		},
		{
			name: "fc",
			persistentVolume: corev1.PersistentVolume{
				Spec: corev1.PersistentVolumeSpec{
					PersistentVolumeSource: corev1.PersistentVolumeSource{
						FC: &corev1.FCVolumeSource{
							FSType: "ext4",
						},
					},
				},
			},
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "FC",
					Content: component.NewText("{\"fsType\":\"ext4\"}"),
				},
			}...),
		},
		{
			name: "flocker",
			persistentVolume: corev1.PersistentVolume{
				Spec: corev1.PersistentVolumeSpec{
					PersistentVolumeSource: corev1.PersistentVolumeSource{
						Flocker: &corev1.FlockerVolumeSource{
							DatasetName: "my-flocker-vol",
						},
					},
				},
			},
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Flocker",
					Content: component.NewText("{\"datasetName\":\"my-flocker-vol\"}"),
				},
			}...),
		},
		{
			name: "flex volume",
			persistentVolume: corev1.PersistentVolume{
				Spec: corev1.PersistentVolumeSpec{
					PersistentVolumeSource: corev1.PersistentVolumeSource{
						FlexVolume: &corev1.FlexPersistentVolumeSource{
							Driver: "kubernetes.io/lvm",
						},
					},
				},
			},
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Flex Volume",
					Content: component.NewText("{\"driver\":\"kubernetes.io/lvm\"}"),
				},
			}...),
		},
		{
			name: "azure file",
			persistentVolume: corev1.PersistentVolume{
				Spec: corev1.PersistentVolumeSpec{
					PersistentVolumeSource: corev1.PersistentVolumeSource{
						AzureFile: &corev1.AzureFilePersistentVolumeSource{
							SecretName: "azure-secret",
							ShareName:  "k8stest",
						},
					},
				},
			},
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Azure File",
					Content: component.NewText("{\"secretName\":\"azure-secret\",\"shareName\":\"k8stest\",\"secretNamespace\":null}"),
				},
			}...),
		},
		{
			name: "vsphere volume",
			persistentVolume: corev1.PersistentVolume{
				Spec: corev1.PersistentVolumeSpec{
					PersistentVolumeSource: corev1.PersistentVolumeSource{
						VsphereVolume: &corev1.VsphereVirtualDiskVolumeSource{
							VolumePath: "/private/tmp",
						},
					},
				},
			},
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Vsphere Volume",
					Content: component.NewText("{\"volumePath\":\"/private/tmp\"}"),
				},
			}...),
		},
		{
			name: "quobyte",
			persistentVolume: corev1.PersistentVolume{
				Spec: corev1.PersistentVolumeSpec{
					PersistentVolumeSource: corev1.PersistentVolumeSource{
						Quobyte: &corev1.QuobyteVolumeSource{
							Registry: "registry:7861",
							Volume:   "testVolume",
						},
					},
				},
			},
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Quobyte",
					Content: component.NewText("{\"registry\":\"registry:7861\",\"volume\":\"testVolume\"}"),
				},
			}...),
		},
		{
			name: "azure disk",
			persistentVolume: corev1.PersistentVolume{
				Spec: corev1.PersistentVolumeSpec{
					PersistentVolumeSource: corev1.PersistentVolumeSource{
						AzureDisk: &corev1.AzureDiskVolumeSource{
							DiskName:    "test.vhd",
							DataDiskURI: "https://someaccount.blob.microsoft.net/vhds/test.vhd",
						},
					},
				},
			},
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Azure Disk",
					Content: component.NewText("{\"diskName\":\"test.vhd\",\"diskURI\":\"https://someaccount.blob.microsoft.net/vhds/test.vhd\"}"),
				},
			}...),
		},
		{
			name: "photon persistent disk",
			persistentVolume: corev1.PersistentVolume{
				Spec: corev1.PersistentVolumeSpec{
					PersistentVolumeSource: corev1.PersistentVolumeSource{
						PhotonPersistentDisk: &corev1.PhotonPersistentDiskVolumeSource{
							PdID:   "e8dbc38c-9374-4099-8a09-fcf4be0a7641",
							FSType: "ext4",
						},
					},
				},
			},
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Photon Persistent Disk",
					Content: component.NewText("{\"pdID\":\"e8dbc38c-9374-4099-8a09-fcf4be0a7641\",\"fsType\":\"ext4\"}"),
				},
			}...),
		},
		{
			name: "portworx",
			persistentVolume: corev1.PersistentVolume{
				Spec: corev1.PersistentVolumeSpec{
					PersistentVolumeSource: corev1.PersistentVolumeSource{
						PortworxVolume: &corev1.PortworxVolumeSource{
							VolumeID: "vol1",
							FSType:   "ext4",
						},
					},
				},
			},
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Portworx Volume",
					Content: component.NewText("{\"volumeID\":\"vol1\",\"fsType\":\"ext4\"}"),
				},
			}...),
		},
		{
			name: "scaleio",
			persistentVolume: corev1.PersistentVolume{
				Spec: corev1.PersistentVolumeSpec{
					PersistentVolumeSource: corev1.PersistentVolumeSource{
						ScaleIO: &corev1.ScaleIOPersistentVolumeSource{
							Gateway:    "https://localhost:443/api",
							System:     "scaleio",
							VolumeName: "vol-0",
							SecretRef: &corev1.SecretReference{
								Name: "sio-secret",
							},
							FSType: "xfs",
						},
					},
				},
			},
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "ScaleIO",
					Content: component.NewText("{\"gateway\":\"https://localhost:443/api\",\"system\":\"scaleio\",\"secretRef\":{\"name\":\"sio-secret\"},\"volumeName\":\"vol-0\",\"fsType\":\"xfs\"}"),
				},
			}...),
		},
		{
			name: "local",
			persistentVolume: corev1.PersistentVolume{
				Spec: corev1.PersistentVolumeSpec{
					PersistentVolumeSource: corev1.PersistentVolumeSource{
						Local: &corev1.LocalVolumeSource{
							Path: "/private/tmp",
						},
					},
				},
			},
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Local",
					Content: component.NewText("{\"path\":\"/private/tmp\"}"),
				},
			}...),
		},
		{
			name: "storageos",
			persistentVolume: corev1.PersistentVolume{
				Spec: corev1.PersistentVolumeSpec{
					PersistentVolumeSource: corev1.PersistentVolumeSource{
						StorageOS: &corev1.StorageOSPersistentVolumeSource{
							VolumeName: "vol-0",
						},
					},
				},
			},
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "StorageOS",
					Content: component.NewText("{\"volumeName\":\"vol-0\"}"),
				},
			}...),
		},
		{
			name: "csi",
			persistentVolume: corev1.PersistentVolume{
				Spec: corev1.PersistentVolumeSpec{
					PersistentVolumeSource: corev1.PersistentVolumeSource{
						CSI: &corev1.CSIPersistentVolumeSource{
							Driver:       "csi-driver.example.com",
							VolumeHandle: "existingVolumeName",
						},
					},
				},
			},
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "CSI",
					Content: component.NewText("{\"driver\":\"csi-driver.example.com\",\"volumeHandle\":\"existingVolumeName\"}"),
				},
			}...),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var sections component.SummarySections
			sections, err := printPersistentVolumeSource(&tc.persistentVolume, sections)
			require.NoError(t, err)

			got := component.NewSummary("Configuration", sections...)

			component.AssertEqual(t, tc.expected, got)
		})
	}
}
