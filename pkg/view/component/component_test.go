/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package component

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetadata_UnmarshalJSON(t *testing.T) {
	data, err := ioutil.ReadFile(filepath.Join("testdata", "metadata.json"))
	require.NoError(t, err)

	got := Metadata{}
	require.NoError(t, got.UnmarshalJSON(data))

	expected := Metadata{
		Type: "type",
		Title: []TitleComponent{
			NewText("title"),
		},
		Accessor: "accessor",
	}
	require.Equal(t, expected, got)
}
