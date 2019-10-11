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

	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/view/component"
)

func TestNamespaceListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	namespace := testutil.CreateNamespace("ns-test-1")
	namespace.CreationTimestamp = *testutil.CreateTimestamp()
	namespace.Status = corev1.NamespaceStatus{Phase: corev1.NamespaceActive}

	list := &corev1.NamespaceList{
		Items: []corev1.Namespace{
			*namespace,
		},
	}

	ctx := context.Background()
	got, err := NamespaceListHandler(ctx, list, printOptions)
	require.NoError(t, err)

	expected := component.NewTableWithRows("Namespaces", "We couldn't find any namespaces!", namespaceListCols, []component.TableRow{
		{
			"Name":   component.NewLink("", "ns-test-1", "/cluster-overview/namespaces/ns-test-1"),
			"Labels": component.NewLabels(make(map[string]string)),
			"Status": component.NewText("Active"),
			"Age":    component.NewTimestamp(namespace.CreationTimestamp.Time),
		},
	})

	component.AssertEqual(t, expected, got)
}
