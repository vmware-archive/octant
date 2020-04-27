/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package yamlviewer

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/view/component"

	corev1 "k8s.io/api/core/v1"
)

func Test_ToComponent(t *testing.T) {
	object := &corev1.Pod{}

	got, err := ToComponent(object)
	require.NoError(t, err)

	data := "---\nmetadata:\n  creationTimestamp: null\nspec:\n  containers: null\nstatus: {}\n"
	expected := component.NewEditor(component.TitleFromString("YAML"), data, false)
	require.NoError(t, expected.SetValueFromObject(object))

	testutil.AssertJSONEqual(t, expected, got)
}
