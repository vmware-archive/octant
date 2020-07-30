/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 */

package dash

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/log"
)

func TestRunner_ValidateKubeconfig(t *testing.T) {
	fs := afero.NewMemMapFs()
	afero.WriteFile(fs, "/test1", []byte(""), 0755)
	afero.WriteFile(fs, "/test2", []byte(""), 0755)

	separator := string(filepath.ListSeparator)

	tests := []struct {
		name     string
		fileList string
		expected string
		isErr    bool
	}{
		{
			name:     "single path",
			fileList: "/test1",
			expected: "/test1",
			isErr:    false,
		},
		{
			name:     "multiple paths",
			fileList: "/test1" + separator + "/test2",
			expected: "/test1" + separator + "/test2",
			isErr:    false,
		},
		{
			name:     "single path not found",
			fileList: "/unknown",
			expected: "",
			isErr:    true,
		},
		{
			name:     "multiple paths not found",
			fileList: "/unknown" + separator + "/unknown2",
			expected: "",
			isErr:    true,
		},
		{
			name:     "multiple file path; missing a config",
			fileList: "/test1" + separator + "/unknown",
			expected: "/test1",
			isErr:    false,
		},
		{
			name:     "invalid path",
			fileList: "not a filepath",
			expected: "",
			isErr:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logger := log.NopLogger()
			path, err := ValidateKubeConfig(logger, test.fileList, fs)
			if test.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, path, test.expected)
		})
	}
}
