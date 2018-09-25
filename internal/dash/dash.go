package dash

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/heptio/developer-dash/web"
	"github.com/skratchdot/open-golang/open"
)

// Run runs the dashboard.
func Run(ctx context.Context, namespace, uiURL string) error {
	log.Printf("Initial namespace for dashboard is %s", namespace)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return err
	}

	d := newDash(listener, namespace, uiURL)

	return d.Run(ctx)
}

type dash struct {
	listener        net.Listener
	uiURL           string
	namespace       string
	defaultHandler  func() (http.Handler, error)
	willOpenBrowser bool
}

func newDash(listener net.Listener, namespace, uiURL string) *dash {
	return &dash{
		listener:        listener,
		namespace:       namespace,
		uiURL:           uiURL,
		defaultHandler:  web.Handler,
		willOpenBrowser: true,
	}
}

func (d *dash) Run(ctx context.Context) error {
	handler, err := d.uiHandler()
	if err != nil {
		return err
	}

	router := mux.NewRouter()

	router.PathPrefix("/").Handler(handler)

	server := http.Server{Handler: router}

	go func() {
		if err = server.Serve(d.listener); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	dashboardURL := fmt.Sprintf("http://%s", d.listener.Addr())
	log.Printf("Dashboard is available at %s", dashboardURL)

	if d.willOpenBrowser {
		if err = open.Run(dashboardURL); err != nil {
			log.Printf("Warning: unable to open browser: %v", err)
		}
	}

	<-ctx.Done()

	return server.Shutdown(ctx)

}

func (d *dash) uiHandler() (http.Handler, error) {
	if d.uiURL == "" {
		return d.defaultHandler()
	}

	return d.uiProxy()
}

func (d *dash) uiProxy() (*httputil.ReverseProxy, error) {
	uiURL := d.uiURL

	if !strings.HasPrefix(uiURL, "http") && !strings.HasPrefix(uiURL, "https") {
		uiURL = fmt.Sprintf("http://%s", uiURL)
	}
	u, err := url.Parse(uiURL)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "" {
		u.Scheme = "http"
	}

	log.Printf("Proxying dashboard UI to %s", u.String())

	proxy := httputil.NewSingleHostReverseProxy(u)
	return proxy, nil
}
