/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package plugin

import "google.golang.org/grpc"

// Broker is a Plugin Broker.
type Broker interface {
	NextId() uint32
	AcceptAndServe(id uint32, s func([]grpc.ServerOption) *grpc.Server)
}
