package api

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/vmware-tanzu/octant/internal/config"
)

type fakeWebsocketClientManager struct {
	WebsocketClientManager
}

func (c *fakeWebsocketClientManager) ClientFromRequest(dashConfig config.Dash, w http.ResponseWriter, r *http.Request) (*WebsocketClient, error) {
	return nil, fmt.Errorf("test: error")
}

func TestWebsocketService_serveWebsocket(t *testing.T) {
	f := &fakeWebsocketClientManager{}
	serveWebsocket(f, nil, nil, nil)
}
