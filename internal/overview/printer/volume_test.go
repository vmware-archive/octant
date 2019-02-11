package printer_test

import (
	"testing"

	"github.com/heptio/developer-dash/internal/view/component"

	corev1 "k8s.io/api/core/v1"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/overview/printer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_VolumeListHandler(t *testing.T) {
	printOptions := printer.Options{
		Cache: cache.NewMemoryCache(),
	}

	object := &corev1.PodSpec{
		Volumes: []corev1.Volume{
			corev1.Volume{
				Name: "gluster-volume",
				VolumeSource: corev1.VolumeSource{
					Glusterfs: &corev1.GlusterfsVolumeSource{
						EndpointsName: "glusterfs-cluster",
						Path:          "test_vol",
						ReadOnly:      true,
					},
				},
			},
			corev1.Volume{
				Name: "ebs-volume",
				VolumeSource: corev1.VolumeSource{
					AWSElasticBlockStore: &corev1.AWSElasticBlockStoreVolumeSource{
						VolumeID: "vol-314159",
						FSType:   "ext4",
					},
				},
			},
		},
	}

	got, err := printer.VolumeListHandler(object, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Kind")
	expected := component.NewTable("Volumes", cols)
	expected.Add(component.TableRow{
		"Name": component.NewText("gluster-volume"),
		"Kind": component.NewText("Glusterfs (a Glusterfs mount on the host that shares a pod's lifetime)"),
	})
	expected.Add(component.TableRow{
		"Name": component.NewText("ebs-volume"),
		"Kind": component.NewText("AWSElasticBlockStore (a Persistent Disk resource in AWS)"),
	})

	assert.Equal(t, expected, got)
}
