package api

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/vmware-tanzu/octant/pkg/api"
	"github.com/vmware-tanzu/octant/pkg/config"
)

type fakeWebsocketClientManager struct {
	api.StreamingConnectionManager
}

func (c *fakeWebsocketClientManager) ClientFromRequest(dashConfig config.Dash, w http.ResponseWriter, r *http.Request) (api.StreamingClient, error) {
	return nil, fmt.Errorf("test: error")
}

func TestWebsocketService_serveWebsocket(t *testing.T) {
	f := &fakeWebsocketClientManager{}
	serveStreamingApi(f, nil, nil, nil)
}
