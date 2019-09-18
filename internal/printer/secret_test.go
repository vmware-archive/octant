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

func Test_SecretConfig(t *testing.T) {
	secret := testutil.CreateSecret("secret")
	secret.Type = corev1.SecretTypeOpaque

	cases := []struct {
		name     string
		secret   *corev1.Secret
		isErr    bool
		expected *component.Summary
	}{
		{
			name:   "general",
			secret: secret,
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Type",
					Content: component.NewText("Opaque"),
				},
			}...)},
		{
			name:   "secret is nil",
			secret: nil,
			isErr:  true,
		},
	}

	for _, tc := range cases {
		controller := gomock.NewController(t)
		defer controller.Finish()

		tpo := newTestPrinterOptions(controller)
		printOptions := tpo.ToOptions()

		sc := NewSecretConfiguration(tc.secret)

		summary, err := sc.Create(printOptions)
		if tc.isErr {
			require.Error(t, err)
			return
		}
		require.NoError(t, err)

		component.AssertEqual(t, tc.expected, summary)
	}
}

func Test_describeSecretData(t *testing.T) {
	secret := testutil.CreateSecret("secret")
	secret.Data = map[string][]byte{
		"foo": []byte{0, 1, 2, 3},
	}

	got, err := describeSecretData(*secret)
	require.NoError(t, err)

	cols := component.NewTableCols("Key")
	expected := component.NewTable("Data", "This secret has no data!", cols)
	expected.Add([]component.TableRow{
		{
			"Key": component.NewText("foo"),
		},
	}...)

	component.AssertEqual(t, expected, got)
}
