package dash

import (
	"github.com/vmware-tanzu/octant/pkg/plugin/api"
	"github.com/vmware-tanzu/octant/pkg/plugin/api/proto"
	"net"

	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)


func serveGRPC(l net.Listener, service api.Service) {

	dashboardServer := api.NewGRPCServer(service)

	s := grpc.NewServer()
	proto.RegisterDashboardServer(s, dashboardServer)

	if err := s.Serve(l); err != cmux.ErrListenerClosed {
		panic(err)
	}
}
