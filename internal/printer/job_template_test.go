/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware/octant/pkg/view/component"
)

func TestJobTemplateHeader(t *testing.T) {
	labels := map[string]string{
		"app": "myapp",
	}

	jth := NewJobTemplateHeader(labels)
	got, err := jth.Create()

	require.NoError(t, err)

	assert.Len(t, got.Config.Labels, 1)

	expected := []component.TitleComponent{
		component.NewText("Job Template"),
	}

	assert.Equal(t, expected, got.Metadata.Title)
}
