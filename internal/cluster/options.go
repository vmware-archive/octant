/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package cluster

import "k8s.io/apimachinery/pkg/runtime/schema"

// options is an internal set of options for configuring this package.
type options struct {
	groupVersionParser GroupVersionParserFunc
}

// buildOptions builds an options struct from a list of functional options.
func buildOptions(list ...Option) options {
	opts := options{
		groupVersionParser: schema.ParseGroupVersion,
	}

	for _, o := range list {
		o(&opts)
	}

	return opts
}

// Option is a functional option for configuring this package.
type Option func(o *options)

// GroupVersionParser sets the group version parser.
func GroupVersionParser(gvp GroupVersionParserFunc) Option {
	return func(o *options) {
		o.groupVersionParser = gvp
	}
}
