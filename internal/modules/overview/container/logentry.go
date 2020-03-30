/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package container

var _ LogEntry = (*logEntry)(nil)

func NewLogEntry(container, line string) logEntry {
	return logEntry{
		container: container,
		line:      line,
	}
}

type logEntry struct {
	line      string
	container string
}

func (l logEntry) Line() string {
	return l.line
}

func (l logEntry) Container() string {
	return l.container
}
