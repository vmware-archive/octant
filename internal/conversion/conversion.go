/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package conversion

import (
	"fmt"
)

// PtrInt32ToString convert *int32 to string
func PtrInt32ToString(p *int32) string {
	var i int32
	if p != nil {
		i = *p
	}

	return fmt.Sprintf("%d", i)
}

// PtrInt32 converts int32 to *int32
func PtrInt32(i int32) *int32 {
	return &i
}

// PtrInt64 converts int64 to *int64
func PtrInt64(i int64) *int64 {
	return &i
}

// PtrBool converts bool to *bool
func PtrBool(b bool) *bool {
	return &b
}
