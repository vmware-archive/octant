/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package icon

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadIcon(t *testing.T) {
	tests := []struct {
		name     string
		iconName string
		isErr    bool
	}{
		{
			name:     "icon exists",
			iconName: OverviewSecret,
			isErr:    false,
		},
		{
			name:     "icon does not exist",
			iconName: "invalid",
			isErr:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := LoadIcon(test.iconName)
			if test.isErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.True(t, len(got) > 0)
		})
	}
}
