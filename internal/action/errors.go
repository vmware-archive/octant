/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package action

import "fmt"

type NotFoundError struct {
	Path string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("action path %q not found", e.Path)
}

type PreviouslyConfiguredError struct {
	Path string
}

func (e *PreviouslyConfiguredError) Error() string {
	return fmt.Sprintf("action with path %q was previously configured", e.Path)
}
