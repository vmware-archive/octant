/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectvisitor_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	configFake "github.com/vmware/octant/internal/config/fake"
	"github.com/vmware/octant/internal/gvk"
	"github.com/vmware/octant/internal/modules/overview/objectvisitor"
	ovFake "github.com/vmware/octant/internal/modules/overview/objectvisitor/fake"
	queryerFake "github.com/vmware/octant/internal/queryer/fake"
	"github.com/vmware/octant/internal/testutil"
)

func TestDefaultVisitor_Visit_use_typed_visitor(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashConfig := configFake.NewMockDash(controller)

	pod := testutil.CreatePod("pod")
	unstructuredPod := testutil.ToUnstructured(t, pod)

	q := queryerFake.NewMockQueryer(controller)

	handler := ovFake.NewMockObjectHandler(controller)

	defaultHandler := ovFake.NewMockDefaultTypedVisitor(controller)
	defaultHandler.EXPECT().
		Visit(gomock.Any(), unstructuredPod, handler, gomock.Any()).Return(nil)

	tv := ovFake.NewMockTypedVisitor(controller)
	tv.EXPECT().Supports().Return(gvk.PodGVK).AnyTimes()
	tv.EXPECT().
		Visit(gomock.Any(), unstructuredPod, handler, gomock.Any())
	tvList := []objectvisitor.TypedVisitor{tv}

	dv, err := objectvisitor.NewDefaultVisitor(dashConfig, q,
		objectvisitor.SetDefaultHandler(defaultHandler),
		objectvisitor.SetTypedVisitors(tvList))
	require.NoError(t, err)

	ctx := context.Background()
	err = dv.Visit(ctx, pod, handler)
	require.NoError(t, err)
}
