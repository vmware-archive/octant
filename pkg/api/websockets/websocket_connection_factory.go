package websockets

import (
	"context"
	"net"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	internalAPI "github.com/vmware-tanzu/octant/internal/api"
	"github.com/vmware-tanzu/octant/pkg/api"
	"github.com/vmware-tanzu/octant/pkg/config"
)

type WebsocketConnectionFactory struct {
	upgrader websocket.Upgrader
}

var DefaultWebsocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return false
		}

		return internalAPI.ShouldAllowHost(host, internalAPI.AcceptedHosts())
	},
}

func NewWebsocketConnectionFactory() *WebsocketConnectionFactory {
	return &WebsocketConnectionFactory{
		upgrader: DefaultWebsocketUpgrader,
	}
}

var _ api.StreamingClientFactory = (*WebsocketConnectionFactory)(nil)

func (wcf *WebsocketConnectionFactory) NewConnection(
	w http.ResponseWriter, r *http.Request, m api.ClientManager, dashConfig config.Dash,
) (api.StreamingClient, context.CancelFunc, error) {
	clientID, err := uuid.NewUUID()
	if err != nil {
		return nil, nil, err
	}

	conn, err := wcf.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithCancel(m.Context())
	client := NewWebsocketClient(ctx, conn, m, dashConfig, m.ActionDispatcher(), clientID)

	return client, cancel, nil
}

func (wcf *WebsocketConnectionFactory) NewTemporaryConnection(
	w http.ResponseWriter, r *http.Request, m api.ClientManager,
) (api.StreamingClient, context.CancelFunc, error) {
	clientID, err := uuid.NewUUID()
	if err != nil {
		return nil, nil, err
	}

	conn, err := wcf.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithCancel(m.Context())
	client := NewTemporaryWebsocketClient(ctx, conn, m, m.ActionDispatcher(), clientID)

	return client, cancel, nil
}
