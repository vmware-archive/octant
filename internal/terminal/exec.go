/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package terminal

import (
	"k8s.io/client-go/rest"
)

type Options struct {
	Config     *rest.Config
	RESTClient rest.Interface
	Executor   executor
}

type executor interface {
	Execute(opts Options) error
}

type DefaultExecutor struct{}

func (e *DefaultExecutor) Execute(opts Options) error {
	return nil
}
