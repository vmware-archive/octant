/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package conversion

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/utils/pointer"
)

func Test_PtrInt32ToString(t *testing.T) {
	cases := []struct {
		name     string
		in       *int32
		expected string
	}{
		{
			name:     "*int32",
			in:       pointer.Int32Ptr(1),
			expected: "1",
		},
		{
			name:     "nil",
			in:       nil,
			expected: "0",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := PtrInt32ToString(tc.in)
			assert.Equal(t, tc.expected, got)
		})
	}
}
