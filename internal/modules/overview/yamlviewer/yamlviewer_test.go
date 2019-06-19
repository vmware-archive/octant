/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package yamlviewer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/heptio/developer-dash/pkg/view/component"

	corev1 "k8s.io/api/core/v1"
)

func Test_ToComponent(t *testing.T) {
	object := &corev1.Pod{}

	got, err := ToComponent(object)
	require.NoError(t, err)

	data := "---\nmetadata:\n  creationTimestamp: null\nspec:\n  containers: null\nstatus: {}\n"
	expected := component.NewYAML(component.TitleFromString("YAML"), data)

	assert.Equal(t, expected, got)
}
