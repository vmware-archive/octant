/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package terminalviewer

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vmware-tanzu/octant/internal/log"
	terminalFake "github.com/vmware-tanzu/octant/internal/terminal/fake"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"testing"

	"github.com/stretchr/testify/require"

	corev1 "k8s.io/api/core/v1"
)

func Test_ToComponent(t *testing.T) {
	controller := gomock.NewController(t)
	object := &corev1.Pod{}
	tm := terminalFake.NewMockManager(controller)
	tm.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),gomock.Any(), gomock.Any()).AnyTimes()

	got, err := ToComponent(context.Background(), object, tm, log.NopLogger())
	require.NoError(t, err)

	details := component.TerminalDetails{
		Container: "",
		Command:   "/bin/sh",
		UUID:      "",
		Active:    true,
	}
	expected := component.NewTerminal("", "Terminal", "", details)

	assert.Equal(t, expected, got)
}
