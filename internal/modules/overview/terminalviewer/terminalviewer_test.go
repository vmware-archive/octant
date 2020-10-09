/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package terminalviewer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	clusterFake "github.com/vmware-tanzu/octant/internal/cluster/fake"
	configFake "github.com/vmware-tanzu/octant/internal/config/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"

	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_ToComponent(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashConfig := configFake.NewMockDash(controller)

	clusterClient := clusterFake.NewMockClientInterface(controller)
	dashConfig.EXPECT().ClusterClient().Return(clusterClient).AnyTimes()

	discoveryInterface := clusterFake.NewMockDiscoveryInterface(controller)
	clusterClient.EXPECT().DiscoveryClient().Return(discoveryInterface, nil).AnyTimes()

	discoveryInterface.EXPECT().ServerGroupsAndResources().AnyTimes()

	object := &corev1.Pod{}

	got, err := ToComponent(context.Background(), object, log.NopLogger(), dashConfig)
	require.NoError(t, err)

	details := component.TerminalDetails{
		Container: "",
		Command:   "/bin/sh",
		Active:    true,
	}
	expected := component.NewTerminal("", "Terminal", "", nil, details)

	assert.Equal(t, expected, got)
}
