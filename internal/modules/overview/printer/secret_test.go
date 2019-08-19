/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
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

	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/view/component"
)

func Test_SecretListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	labels := map[string]string{
		"foo": "bar",
	}

	now := testutil.Time()

	object := &corev1.SecretList{
		Items: []corev1.Secret{
			{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "default",
					CreationTimestamp: metav1.Time{
						Time: now,
					},
					Labels: labels,
				},
				Data: map[string][]byte{
					"key": []byte("value"),
				},
				Type: corev1.SecretTypeOpaque,
			},
		},
	}

	tpo.PathForObject(&object.Items[0], object.Items[0].Name, "/secret")

	ctx := context.Background()
	got, err := SecretListHandler(ctx, object, printOptions)
	require.NoError(t, err)

	expected := component.NewTable("Secrets", "We couldn't find any secrets!", secretTableCols)
	expected.Add(component.TableRow{
		"Name":   component.NewLink("", "secret", "/secret"),
		"Labels": component.NewLabels(labels),
		"Type":   component.NewText("Opaque"),
		"Data":   component.NewText("1"),
		"Age":    component.NewTimestamp(now),
	})

	component.AssertEqual(t, expected, got)
}
