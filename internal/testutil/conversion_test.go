/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package testutil

import "testing"

func TestToUnstructured(t *testing.T) {
	tests := []struct{
		name string
	} {
		{
			name: "in general",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T){

		})
	}
}
