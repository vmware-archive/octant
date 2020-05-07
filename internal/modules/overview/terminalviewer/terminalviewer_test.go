/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package terminalviewer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/pkg/view/component"

	"github.com/stretchr/testify/require"

	corev1 "k8s.io/api/core/v1"
)

func Test_ToComponent(t *testing.T) {
	object := &corev1.Pod{}

	got, err := ToComponent(context.Background(), object, log.NopLogger())
	require.NoError(t, err)

	details := component.TerminalDetails{
		Container: "",
		Command:   "/bin/sh",
		Active:    true,
	}
	expected := component.NewTerminal("", "Terminal", "", []string{}, details)

	assert.Equal(t, expected, got)
}
