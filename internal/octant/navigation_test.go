/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package octant_test

import (
	"fmt"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware/octant/internal/icon"
	"github.com/vmware/octant/internal/octant"
)

func Test_NewNavigation(t *testing.T) {
	navPath := "/navPath"
	title := "title"

	nav, err := octant.NewNavigation(title, navPath)
	require.NoError(t, err)

	assert.Equal(t, navPath, nav.Path)
	assert.Equal(t, title, nav.Title)
}

func TestNavigationEntriesHelper(t *testing.T) {
	neh := octant.NavigationEntriesHelper{}

	neh.Add("title", "suffix", icon.OverviewService)

	list, err := neh.Generate("/prefix")
	require.NoError(t, err)

	expected := octant.Navigation{
		Title:    "title",
		Path:     path.Join("/prefix", "suffix"),
		IconName: fmt.Sprintf("internal:%s", icon.OverviewService),
	}

	assert.Len(t, list, 1)
	assert.Equal(t, expected.Title, list[0].Title)
	assert.Equal(t, expected.Path, list[0].Path)
	assert.Equal(t, expected.IconName, list[0].IconName)
	assert.NotEmpty(t, list[0].IconSource)
}
