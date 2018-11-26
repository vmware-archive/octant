package dash

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"runtime"
	"strings"

	"github.com/pkg/errors"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/heptio/developer-dash/internal/api"
	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/module"
	"github.com/heptio/developer-dash/web"
	"github.com/heptio/go-telemetry/pkg/telemetry"
	"github.com/skratchdot/open-golang/open"
)

const (
	apiPathPrefix       = "/api/v1"
	defaultListenerAddr = "127.0.0.1:0"
)

// Run runs the dashboard.
func Run(ctx context.Context, namespace, uiURL, kubeconfig string, logger log.Logger, telemetryClient telemetry.Interface) error {
	logger.Debugf("Loading configuration: %v", kubeconfig)
	clusterClient, err := cluster.FromKubeconfig(kubeconfig)
	if err != nil {
		return errors.Wrap(err, "failed to init cluster client")
	}

	nsClient, err := clusterClient.NamespaceClient()
	if err != nil {
		return errors.Wrap(err, "failed to create namespace client")
	}

	// If not overridden, use initial namespace from current context in KUBECONFIG
	if namespace == "" {
		namespace = nsClient.InitialNamespace()
	}

	logger.Debugf("initial namespace for dashboard is %s", namespace)

	moduleManager, err := module.NewManager(clusterClient, namespace, logger)
	if err != nil {
		return errors.Wrap(err, "create module manager")
	}

	listener, err := buildListener()
	if err != nil {
		return errors.Wrap(err, "failed to create net listener")
	}

	version, err := clusterClient.Version()
	if err != nil {
		// Fail-open
		logger.Errorf("failed to get kubernetes version from cluster")
	}

	telemetryClient = telemetryClient.With(telemetry.Labels{
		"os":                 runtime.GOOS,
		"kubernetes.version": version,
	})

	d, err := newDash(listener, namespace, uiURL, nsClient, moduleManager, logger, telemetryClient)
	if err != nil {
		return errors.Wrap(err, "failed to create dash instance")
	}

	if os.Getenv("DASH_DISABLE_OPEN_BROWSER") != "" {
		d.willOpenBrowser = false
	}

	go func() {
		if err := d.Run(ctx); err != nil {
			logger.Debugf("running dashboard service: %v", err)
		}
	}()

	<-ctx.Done()
	moduleManager.Unload()

	return nil
}

func buildListener() (net.Listener, error) {
	listenerAddr := defaultListenerAddr
	if customListenerAddr := os.Getenv("DASH_LISTENER_ADDR"); customListenerAddr != "" {
		listenerAddr = customListenerAddr
	}

	return net.Listen("tcp", listenerAddr)
}

type dash struct {
	listener        net.Listener
	uiURL           string
	namespace       string
	defaultHandler  func() (http.Handler, error)
	apiHandler      api.Service
	willOpenBrowser bool
	logger          log.Logger
	telemetryClient telemetry.Interface
}

func newDash(listener net.Listener, namespace, uiURL string, nsClient cluster.NamespaceInterface, moduleManager module.ManagerInterface, logger log.Logger, telemetryClient telemetry.Interface) (*dash, error) {
	ah := api.New(apiPathPrefix, nsClient, moduleManager, logger, telemetryClient)

	for _, m := range moduleManager.Modules() {
		if err := ah.RegisterModule(m); err != nil {
			return nil, err
		}
	}

	return &dash{
		listener:        listener,
		namespace:       namespace,
		uiURL:           uiURL,
		defaultHandler:  web.Handler,
		willOpenBrowser: true,
		apiHandler:      ah,
		logger:          logger,
		telemetryClient: telemetryClient,
	}, nil
}

func (d *dash) Run(ctx context.Context) error {
	d.telemetryClient.SendEvent("dash.startup", telemetry.Measurements{"count": 1})

	handler, err := d.handler()
	if err != nil {
		return err
	}

	server := http.Server{Handler: handler}

	go func() {
		if err = server.Serve(d.listener); err != nil && err != http.ErrServerClosed {
			d.logger.Errorf("http server: %v", err)
			os.Exit(1) // TODO graceful shutdown for other goroutines
		}
	}()

	dashboardURL := fmt.Sprintf("http://%s", d.listener.Addr())
	d.logger.Infof("Dashboard is available at %s\n", dashboardURL)

	if d.willOpenBrowser {
		d.telemetryClient.SendEvent("dash.browser.open", telemetry.Measurements{"count": 1})
		if err = open.Run(dashboardURL); err != nil {
			d.telemetryClient.SendEvent("dash.browser.failure", telemetry.Measurements{"count": 1})
			d.logger.Warnf("unable to open browser: %v", err)
		}
	}

	<-ctx.Done()

	// TODO context is already done - pass a different one too allow time for graceful shutdown
	return server.Shutdown(ctx)
}

func (d *dash) handler() (http.Handler, error) {
	handler, err := d.uiHandler()
	if err != nil {
		return nil, err
	}

	router := mux.NewRouter()
	router.PathPrefix(apiPathPrefix).Handler(d.apiHandler.Handler())
	router.PathPrefix("/").Handler(handler)

	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Language", "Origin", "Content-Type"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	return handlers.CORS(allowedOrigins, allowedHeaders, allowedMethods)(router), nil
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

	d.logger.Infof("Proxying dashboard UI to %s", u.String())

	proxy := httputil.NewSingleHostReverseProxy(u)
	return proxy, nil
}
