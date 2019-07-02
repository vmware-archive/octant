/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

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

	"github.com/vmware/octant/internal/action"
	"github.com/vmware/octant/internal/componentcache"
	"github.com/vmware/octant/internal/objectstore"

	"github.com/vmware/octant/internal/config"
	"github.com/vmware/octant/internal/describer"
	"github.com/vmware/octant/internal/modules/clusteroverview"
	"github.com/vmware/octant/internal/modules/configuration"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/vmware/octant/internal/api"
	"github.com/vmware/octant/internal/cluster"
	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/internal/module"
	"github.com/vmware/octant/internal/modules/localcontent"
	"github.com/vmware/octant/internal/modules/overview"
	"github.com/vmware/octant/internal/portforward"
	"github.com/vmware/octant/pkg/plugin"
	"github.com/vmware/octant/pkg/store"
	"github.com/vmware/octant/web"

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
	Context          string
}

// Run runs the dashboard.
func Run(ctx context.Context, logger log.Logger, shutdownCh chan bool, options Options) error {
	ctx = log.WithLoggerContext(ctx, logger)

	if options.Context != "" {
		logger.With("initial-context", options.Context).Infof("Setting initial context from user flags")
	}

	logger.Debugf("Loading configuration: %v", options.KubeConfig)
	clusterClient, err := cluster.FromKubeConfig(ctx, options.KubeConfig, options.Context)
	if err != nil {
		return errors.Wrap(err, "failed to init cluster client")
	}

	if options.EnableOpenCensus {
		if err := enableOpenCensus(); err != nil {
			logger.Infof("Enabling OpenCensus")
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

	appObjectStore, err := initObjectStore(ctx, clusterClient)
	if err != nil {
		return errors.Wrap(err, "initializing store")
	}

	componentCache, err := componentcache.NewComponentCache(ctx)
	if err != nil {
		return errors.Wrap(err, "initializing component cache")
	}

	crdWatcher, err := describer.NewDefaultCRDWatcher(ctx, appObjectStore)
	if err != nil {
		return errors.Wrap(err, "initializing CRD watcher")
	}

	portForwarder, err := initPortForwarder(ctx, clusterClient, appObjectStore)
	if err != nil {
		return errors.Wrap(err, "initializing port forwarder")
	}

	actionManger := action.NewManager(logger)

	mo := moduleOptions{
		clusterClient: clusterClient,
		namespace:     options.Namespace,
		logger:        logger,
		actionManager: actionManger,
	}
	moduleManager, err := initModuleManager(mo)
	if err != nil {
		return errors.Wrap(err, "init module manager")
	}

	pluginManager, err := initPlugin(portForwarder, appObjectStore, moduleManager)
	if err != nil {
		return errors.Wrap(err, "initializing plugin manager")
	}

	dashConfig := config.NewLiveConfig(
		clusterClient,
		crdWatcher,
		options.KubeConfig,
		logger,
		moduleManager,
		appObjectStore,
		componentCache,
		pluginManager,
		portForwarder,
		options.Context)

	moduleList, err := initModules(ctx, dashConfig, options.Namespace)
	if err != nil {
		return errors.Wrap(err, "initializing modules")
	}

	for _, mod := range moduleList {
		if err := moduleManager.Register(mod); err != nil {
			return errors.Wrapf(err, "loading module %s", mod.Name())
		}
	}

	if err := pluginManager.Start(ctx); err != nil {
		return errors.Wrapf(err, "start plugin manager")
	}

	listener, err := buildListener()
	if err != nil {
		return errors.Wrap(err, "failed to create net listener")
	}

	// Initialize the API
	ah := api.New(ctx, apiPathPrefix, clusterClient, moduleManager, actionManger, logger)
	for _, m := range moduleManager.Modules() {
		if err := ah.RegisterModule(m); err != nil {
			return errors.Wrapf(err, "registering module: %v", m.Name())
		}
	}

	d, err := newDash(listener, options.Namespace, options.FrontendURL, ah, logger)
	if err != nil {
		return errors.Wrap(err, "failed to create dash instance")
	}

	if os.Getenv("OCTANT_DISABLE_OPEN_BROWSER") != "" {
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
func initObjectStore(ctx context.Context, client cluster.ClientInterface) (store.Store, error) {
	if client == nil {
		return nil, errors.New("nil cluster client")
	}

	appObjectStore, err := objectstore.NewWatch(ctx, client)

	if err != nil {
		return nil, errors.Wrapf(err, "creating object store for app")
	}

	return appObjectStore, nil
}

func initPortForwarder(ctx context.Context, client cluster.ClientInterface, appObjectStore store.Store) (portforward.PortForwarder, error) {
	return portforward.Default(ctx, client, appObjectStore)
}

type moduleOptions struct {
	clusterClient  *cluster.Cluster
	crdWatcher     config.CRDWatcher
	objectStore    store.Store
	componentCache componentcache.ComponentCache
	namespace      string
	logger         log.Logger
	pluginManager  *plugin.Manager
	portForwarder  portforward.PortForwarder
	kubeConfigPath string
	initialContext string
	actionManager  *action.Manager
}

func initModules(ctx context.Context, dashConfig config.Dash, namespace string) ([]module.Module, error) {
	var list []module.Module

	overviewOptions := overview.Options{
		Namespace:  namespace,
		DashConfig: dashConfig,
	}
	overviewModule, err := overview.New(ctx, overviewOptions)
	if err != nil {
		return nil, errors.Wrap(err, "create overview module")
	}

	list = append(list, overviewModule)

	clusterOverviewOptions := clusteroverview.Options{
		DashConfig: dashConfig,
	}
	clusterOverviewModule, err := clusteroverview.New(ctx, clusterOverviewOptions)
	if err != nil {
		return nil, errors.Wrap(err, "create cluster overview module")
	}

	list = append(list, clusterOverviewModule)

	configurationOptions := configuration.Options{
		DashConfig:     dashConfig,
		KubeConfigPath: dashConfig.KubeConfigPath(),
	}
	configurationModule := configuration.New(ctx, configurationOptions)

	list = append(list, configurationModule)

	localContentPath := os.Getenv("OCTANT_LOCAL_CONTENT")
	if localContentPath != "" {
		localContentModule := localcontent.New(localContentPath)
		list = append(list, localContentModule)
	}

	return list, nil
}

// initModuleManager initializes the moduleManager (and currently the modules themselves)
func initModuleManager(options moduleOptions) (*module.Manager, error) {
	moduleManager, err := module.NewManager(options.clusterClient, options.namespace, options.actionManager, options.logger)
	if err != nil {
		return nil, errors.Wrap(err, "create module manager")
	}

	return moduleManager, nil
}

func buildListener() (net.Listener, error) {
	listenerAddr := defaultListenerAddr
	if customListenerAddr := os.Getenv("OCTANT_LISTENER_ADDR"); customListenerAddr != "" {
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
	apiHandler, err := d.apiHandler.Handler(ctx)
	if err != nil {
		return nil, err
	}

	router.PathPrefix(apiPathPrefix).Handler(apiHandler)
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
			ServiceName: "octant",
		},
	})

	if err != nil {
		return errors.Wrap(err, "failed to create Jaeger exporter")
	}

	trace.RegisterExporter(je)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	return nil
}
