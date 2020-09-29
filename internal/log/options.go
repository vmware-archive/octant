/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package log

type options struct {
	messageIDGenerator MessageIDGenerator
}

func makeDefaultOptions(list ...Option) options {
	opts := options{
		messageIDGenerator: &UUIDMessageIDGenerator{},
	}

	for _, o := range list {
		o(&opts)
	}

	return opts
}

// Option is a log option.
type Option func(o *options)

// WithIDGenerator sets the message id generator option.
func WithIDGenerator(g MessageIDGenerator) Option {
	return func(o *options) {
		o.messageIDGenerator = g
	}
}
