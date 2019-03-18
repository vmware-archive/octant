package plugin

import "google.golang.org/grpc"

// Broker is a Plugin Broker.
type Broker interface {
	NextId() uint32
	AcceptAndServe(id uint32, s func([]grpc.ServerOption) *grpc.Server)
}
