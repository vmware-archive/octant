package plugin

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"github.com/heptio/developer-dash/pkg/plugin/proto"
	"google.golang.org/grpc"
)

//go:generate protoc -I$GOPATH/src/github.com/heptio/developer-dash/vendor -I$GOPATH/src/github.com/heptio/developer-dash -I. --go_out=plugins=grpc:. proto/dashboard.proto
//go:generate mockgen -destination=./fake/mock_runners.go -package=fake github.com/heptio/developer-dash/pkg/plugin Runners
//go:generate mockgen -destination=./fake/mock_manager_store.go -package=fake github.com/heptio/developer-dash/pkg/plugin ManagerStore
//go:generate mockgen -destination=./fake/mock_client_factory.go -package=fake github.com/heptio/developer-dash/pkg/plugin ClientFactory
//go:generate mockgen -destination=./fake/mock_service.go -package=fake github.com/heptio/developer-dash/pkg/plugin Service
//go:generate mockgen -destination=./fake/mock_broker.go -package=fake github.com/heptio/developer-dash/pkg/plugin Broker
//go:generate mockgen -destination=./fake/mock_plugin_client.go -package=fake github.com/heptio/developer-dash/pkg/plugin/proto PluginClient
//go:generate mockgen -destination=./fake/mock_client_protocol.go -package=fake github.com/hashicorp/go-plugin ClientProtocol

var (
	pluginMap = map[string]plugin.Plugin{
		PluginName: &ServicePlugin{},
	}
)

// ServicePlugin is the GRPC plugin for Service.
type ServicePlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl Service
}

var _ plugin.GRPCPlugin = (*ServicePlugin)(nil)

// GRPCServer is the plugin's GRPC server.
func (p *ServicePlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterPluginServer(s, &GRPCServer{
		Impl:   p.Impl,
		broker: broker,
	})

	return nil
}

// GRPCClient is the plugin's GRPC client.
func (p *ServicePlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{
		client: proto.NewPluginClient(c),
		broker: broker,
	}, nil
}
