package plugin

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"github.com/heptio/developer-dash/pkg/plugin/proto"
	"google.golang.org/grpc"
)

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
