/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package applications

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	configFake "github.com/vmware/octant/internal/config/fake"
	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/store"
	"github.com/vmware/octant/pkg/store/fake"
	"github.com/vmware/octant/pkg/view/component"
)

func Test_summarizer(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ctx := context.Background()

	objectStore := fake.NewMockStore(controller)

	key := store.Key{
		Namespace:  "default",
		APIVersion: "v1",
		Kind:       "Pod",
	}
	podList := testutil.ToUnstructuredList(t, testutil.CreatePod("pod", withPodLabels(map[string]string{
		appLabelName:     "name",
		appLabelInstance: "instance",
		appLabelVersion:  "version",
	})))
	objectStore.EXPECT().
		List(gomock.Any(), key).
		Return(podList, true, nil)

	dashConfig := configFake.NewMockDash(controller)
	dashConfig.EXPECT().ObjectStore().Return(objectStore)

	s := summarizer{}
	actual, err := s.Summarize(ctx, "default", dashConfig)
	require.NoError(t, err)

	expected := component.NewTableWithRows("Applications", "applications", applicationListColumns, []component.TableRow{
		{
			"Name":     component.NewLink("", "name", "/applications/namespace/default/name/instance/version"),
			"Instance": component.NewText("instance"),
			"Version":  component.NewText("version"),
			"State":    component.NewText("state"),
			"Pods":     component.NewText("1"),
		},
	})

	component.AssertEqual(t, expected, actual)
}
