/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package tsgen

type options struct {
	disableFormat bool
}

func makeDefaultOptions(list ...Option) options {
	opts := options{}

	for _, o := range list {
		o(&opts)
	}

	return opts
}

// Option is an option for configuration tsgen.
type Option func(o *options)

// DisableFormatter disables the typescript formatter.
func DisableFormatter() Option {
	return func(o *options) {
		o.disableFormat = true
	}
}
