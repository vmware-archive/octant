/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package errors

import "time"

// InternalError represents an internal Octant error.
type InternalError interface {
	ID() string
	Error() string
	Timestamp() time.Time
}
