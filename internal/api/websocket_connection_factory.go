package api

import (
	"context"
	"net"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/vmware-tanzu/octant/internal/config"
)

type WebsocketConnectionFactory struct {
	upgrader websocket.Upgrader
}

func NewWebsocketConnectionFactory() *WebsocketConnectionFactory {
	return &WebsocketConnectionFactory{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				host, _, err := net.SplitHostPort(r.RemoteAddr)
				if err != nil {
					return false
				}

				return shouldAllowHost(host, acceptedHosts())
			},
		},
	}
}

var _ StreamingClientFactory = (*WebsocketConnectionFactory)(nil)

func (wcf *WebsocketConnectionFactory) NewConnection(
	clientID uuid.UUID, w http.ResponseWriter, r *http.Request, m *StreamingConnectionManager, dashConfig config.Dash,
) (StreamingClient, context.CancelFunc, error) {
	conn, err := wcf.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithCancel(m.ctx)
	client := NewWebsocketClient(ctx, conn, m, dashConfig, m.actionDispatcher, clientID)

	go client.readPump()
	go client.writePump()

	return client, cancel, nil
}

func (wcf *WebsocketConnectionFactory) NewTemporaryConnection(
	clientID uuid.UUID, w http.ResponseWriter, r *http.Request, m *StreamingConnectionManager,
) (StreamingClient, context.CancelFunc, error) {
	conn, err := wcf.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithCancel(m.ctx)
	client := NewTemporaryWebsocketClient(ctx, conn, m, m.actionDispatcher, clientID)

	go client.readPump()
	go client.writePump()

	return client, cancel, nil
}
