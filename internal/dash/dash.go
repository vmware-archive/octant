package dash

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/heptio/developer-dash/internal/config"
	"github.com/heptio/developer-dash/internal/modules/clusteroverview"
	"github.com/heptio/developer-dash/internal/modules/configuration"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/heptio/developer-dash/internal/api"
	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/module"
	"github.com/heptio/developer-dash/internal/modules/localcontent"
	"github.com/heptio/developer-dash/internal/modules/overview"
	"github.com/heptio/developer-dash/internal/objectstore"
	"github.com/heptio/developer-dash/internal/portforward"
	"github.com/heptio/developer-dash/pkg/plugin"
	"github.com/heptio/developer-dash/web"

	"github.com/pkg/errors"
	"github.com/skratchdot/open-golang/open"
	"go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/trace"
)

const (
	apiPathPrefix       = "/api/v1"
	defaultListenerAddr = "127.0.0.1:0"
)

type Options struct {
	EnableOpenCensus bool
	KubeConfig       string
	Namespace        string
	FrontendURL      string
}

// Run runs the dashboard.
func Run(ctx context.Context, logger log.Logger, shutdownCh chan bool, options Options) error {
	ctx = log.WithLoggerContext(ctx, logger)

	logger.Debugf("Loading configuration: %v", options.KubeConfig)
	clusterClient, err := cluster.FromKubeconfig(ctx, options.KubeConfig)
	if err != nil {
		return errors.Wrap(err, "failed to init cluster client")
	}

	if options.EnableOpenCensus {
		if err := enableOpenCensus(); err != nil {
			return errors.Wrap(err, "enabling open census")
		}
	}

	nsClient, err := clusterClient.NamespaceClient()
	if err != nil {
		return errors.Wrap(err, "failed to create namespace client")
	}

	// If not overridden, use initial namespace from current context in KUBECONFIG
	if options.Namespace == "" {
		options.Namespace = nsClient.InitialNamespace()
	}

	logger.Debugf("initial namespace for dashboard is %s", options.Namespace)

	infoClient, err := clusterClient.InfoClient()
	if err != nil {
		return errors.Wrap(err, "failed to create info client")
	}

	appObjectStore, err := initObjectStore(ctx.Done(), clusterClient)
	if err != nil {
		return errors.Wrap(err, "initializing cache")
	}

	portForwarder, err := initPortForwarder(ctx, clusterClient, appObjectStore)
	if err != nil {
		return errors.Wrap(err, "initializing port forwarder")
	}

	pluginManager, err := initPlugin(ctx, portForwarder, appObjectStore)
	if err != nil {
		return errors.Wrap(err, "initializing plugin manager")
	}

	mo := moduleOptions{
		clusterClient: clusterClient,
		objectStore:   appObjectStore,
		namespace:     options.Namespace,
		logger:        logger,
		pluginManager: pluginManager,
		portForwarder: portForwarder,
	}
	moduleManager, err := initModuleManager(ctx, mo)
	if err != nil {
		return errors.Wrap(err, "init module manager")
	}

	listener, err := buildListener()
	if err != nil {
		return errors.Wrap(err, "failed to create net listener")
	}

	// Initialize the API
	ah := api.New(ctx, apiPathPrefix, nsClient, infoClient, moduleManager, logger)
	for _, m := range moduleManager.Modules() {
		if err := ah.RegisterModule(m); err != nil {
			return errors.Wrapf(err, "registering module: %v", m.Name())
		}
	}

	d, err := newDash(listener, options.Namespace, options.FrontendURL, ah, logger)
	if err != nil {
		return errors.Wrap(err, "failed to create dash instance")
	}

	if os.Getenv("CLUSTEREYE_DISABLE_OPEN_BROWSER") != "" {
		d.willOpenBrowser = false
	}

	go func() {
		if err := d.Run(ctx); err != nil {
			logger.Debugf("running dashboard service: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx := log.WithLoggerContext(context.Background(), logger)

	moduleManager.Unload()
	pluginManager.Stop(shutdownCtx)

	shutdownCh <- true

	return nil
}

// initObjectStore initializes the cluster object store interface
func initObjectStore(stopCh <-chan struct{}, client cluster.ClientInterface) (objectstore.ObjectStore, error) {
	if client == nil {
		return nil, errors.New("nil cluster client")
	}

	appObjectStore, err := objectstore.NewWatch(client, stopCh)

	if err != nil {
		return nil, errors.Wrapf(err, "creating object store for app")
	}

	return appObjectStore, nil
}

func initPortForwarder(ctx context.Context, client cluster.ClientInterface, appObjectStore objectstore.ObjectStore) (portforward.PortForwarder, error) {
	return portforward.Default(ctx, client, appObjectStore)
}

type moduleOptions struct {
	clusterClient *cluster.Cluster
	objectStore   objectstore.ObjectStore
	namespace     string
	logger        log.Logger
	pluginManager *plugin.Manager
	portForwarder portforward.PortForwarder
}

// initModuleManager initializes the moduleManager (and currently the modules themselves)
func initModuleManager(ctx context.Context, options moduleOptions) (*module.Manager, error) {
	moduleManager, err := module.NewManager(options.clusterClient, options.namespace, options.logger)
	if err != nil {
		return nil, errors.Wrap(err, "create module manager")
	}

	c := config.NewLiveConfig(
		options.clusterClient,
		options.logger,
		moduleManager,
		options.objectStore,
		options.pluginManager,
		options.portForwarder,
	)

	overviewOptions := overview.Options{
		Namespace:     options.namespace,
		DashConfig:    c,
	}
	overviewModule, err := overview.New(ctx, overviewOptions)
	if err != nil {
		return nil, errors.Wrap(err, "create overview module")
	}

	moduleManager.Register(overviewModule)

	clusterOverviewOptions := clusteroverview.Options{
		DashConfig:    c,
	}
	clusterOverviewModule := clusteroverview.New(ctx, clusterOverviewOptions)
	moduleManager.Register(clusterOverviewModule)

	configurationOptions := configuration.Options{
		DashConfig: c,
	}
	configurationModule := configuration.New(ctx, configurationOptions)
	moduleManager.Register(configurationModule)

	localContentPath := os.Getenv("CLUSTEREYE_LOCAL_CONTENT")
	if localContentPath != "" {
		localContentModule := localcontent.New(localContentPath)
		moduleManager.Register(localContentModule)
	}

	if err = moduleManager.Load(); err != nil {
		return nil, errors.Wrap(err, "modules load")
	}

	return moduleManager, nil
}

func buildListener() (net.Listener, error) {
	listenerAddr := defaultListenerAddr
	if customListenerAddr := os.Getenv("CLUSTEREYE_LISTENER_ADDR"); customListenerAddr != "" {
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
}

func newDash(listener net.Listener, namespace, uiURL string, apiHandler api.Service, logger log.Logger) (*dash, error) {
	return &dash{
		listener:        listener,
		namespace:       namespace,
		uiURL:           uiURL,
		defaultHandler:  web.Handler,
		willOpenBrowser: true,
		apiHandler:      apiHandler,
		logger:          logger,
	}, nil
}

func (d *dash) Run(ctx context.Context) error {
	handler, err := d.handler(ctx)
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
		if err = open.Run(dashboardURL); err != nil {
			d.logger.Warnf("unable to open browser: %v", err)
		}
	}

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	return server.Shutdown(shutdownCtx)
}

// handler configures primary http routes
func (d *dash) handler(ctx context.Context) (http.Handler, error) {
	handler, err := d.uiHandler()
	if err != nil {
		return nil, err
	}

	router := mux.NewRouter()
	router.PathPrefix(apiPathPrefix).Handler(d.apiHandler.Handler(ctx))
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

func enableOpenCensus() error {
	agentEndpointURI := "localhost:6831"
	collectorEndpointURI := "http://localhost:14268/api/traces"

	je, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint:     agentEndpointURI,
		CollectorEndpoint: collectorEndpointURI,
		Process: jaeger.Process{
			ServiceName: "clustereye",
		},
	})

	if err != nil {
		return errors.Wrap(err, "failed to create Jaeger exporter")
	}

	trace.RegisterExporter(je)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	return nil
}
