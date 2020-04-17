/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package container

import "context"

type LogEntry interface {
	Line() string
	Container() string
}

type LogStreamer interface {
	// Names returns a list of all of the containers for this log stream.
	Names() []string
	// Stream wraps the client-go GetLogs().Stream call for the configured
	// pods. Stream is responsible for aggregating the logs for all the
	// containers.
	Stream(context.Context, chan<- LogEntry)
	// Close closes all of the streams.
	Close(chan<- LogEntry)
}
