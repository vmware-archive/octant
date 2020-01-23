/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
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

	"contrib.go.opencensus.io/exporter/jaeger"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/viper"
	"go.opencensus.io/trace"

	"github.com/vmware-tanzu/octant/internal/api"
	"github.com/vmware-tanzu/octant/internal/cluster"
	"github.com/vmware-tanzu/octant/internal/config"
	"github.com/vmware-tanzu/octant/internal/describer"
	internalErr "github.com/vmware-tanzu/octant/internal/errors"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/module"
	"github.com/vmware-tanzu/octant/internal/modules/applications"
	"github.com/vmware-tanzu/octant/internal/modules/clusteroverview"
	"github.com/vmware-tanzu/octant/internal/modules/configuration"
	"github.com/vmware-tanzu/octant/internal/modules/localcontent"
	"github.com/vmware-tanzu/octant/internal/modules/overview"
	"github.com/vmware-tanzu/octant/internal/modules/workloads"
	"github.com/vmware-tanzu/octant/internal/objectstore"
	"github.com/vmware-tanzu/octant/internal/portforward"
	"github.com/vmware-tanzu/octant/internal/terminal"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/plugin"
	pluginAPI "github.com/vmware-tanzu/octant/pkg/plugin/api"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/web"
)

type Options struct {
	EnableOpenCensus       bool
	DisableClusterOverview bool
	KubeConfig             string
	Namespace              string
	FrontendURL            string
	Context                string
	ClientQPS              float32
	ClientBurst            int
	UserAgent              string
}

// Run runs the dashboard.
func Run(ctx context.Context, logger log.Logger, shutdownCh chan bool, options Options) error {
	ctx = log.WithLoggerContext(ctx, logger)

	if options.Context != "" {
		logger.With("initial-context", options.Context).Infof("Setting initial context from user flags")
	}

	logger.Debugf("Loading configuration: %v", options.KubeConfig)
	restConfigOptions := cluster.RESTConfigOptions{
		QPS:       options.ClientQPS,
		Burst:     options.ClientBurst,
		UserAgent: options.UserAgent,
	}
	clusterClient, err := cluster.FromKubeConfig(ctx, options.KubeConfig, options.Context, options.Namespace, restConfigOptions)
	if err != nil {
		return fmt.Errorf("failed to init cluster client: %w", err)
	}

	if options.EnableOpenCensus {
		if err := enableOpenCensus(); err != nil {
			logger.Infof("Enabling OpenCensus")
			return fmt.Errorf("enabling open census: %w", err)
		}
	}

	nsClient, err := clusterClient.NamespaceClient()
	if err != nil {
		return fmt.Errorf("failed to create namespace client: %w", err)
	}

	// If not overridden, use initial namespace from current context in KUBECONFIG
	if options.Namespace == "" {
		options.Namespace = nsClient.InitialNamespace()
	}

	logger.Debugf("initial namespace for dashboard is %s", options.Namespace)

	appObjectStore, err := initObjectStore(ctx, clusterClient)
	if err != nil {
		return fmt.Errorf("initializing store: %w", err)
	}

	errorStore, err := internalErr.NewErrorStore()
	if err != nil {
		return fmt.Errorf("initializing error store: %w", err)
	}

	crdWatcher, err := describer.NewDefaultCRDWatcher(ctx, clusterClient, appObjectStore, errorStore)
	if err != nil {
		return fmt.Errorf("initializing CRD watcher: %w", err)
	}

	if err := crdWatcher.Watch(ctx); err != nil {
		return fmt.Errorf("unable to start CRD watcher: %w", err)
	}

	portForwarder, err := initPortForwarder(ctx, clusterClient, appObjectStore)
	if err != nil {
		return fmt.Errorf("initializing port forwarder: %w", err)
	}

	terminalManager, err := terminal.NewTerminalManager(ctx, clusterClient, appObjectStore)
	if err != nil {
		return fmt.Errorf("initializing terminal manager: %w", err)
	}

	actionManger := action.NewManager(logger)

	mo := &moduleOptions{
		clusterClient: clusterClient,
		namespace:     options.Namespace,
		logger:        logger,
		actionManager: actionManger,
	}
	moduleManager, err := initModuleManager(mo)
	if err != nil {
		return fmt.Errorf("init module manager: %w", err)
	}

	frontendProxy := pluginAPI.FrontendProxy{}

	pluginDashboardService := &pluginAPI.GRPCService{
		ObjectStore:   appObjectStore,
		PortForwarder: portForwarder,
		FrontendProxy: frontendProxy,
	}

	pluginManager, err := initPlugin(moduleManager, actionManger, pluginDashboardService)
	if err != nil {
		return fmt.Errorf("initializing plugin manager: %w", err)
	}

	dashConfig := config.NewLiveConfig(
		clusterClient,
		crdWatcher,
		options.KubeConfig,
		logger,
		moduleManager,
		appObjectStore,
		errorStore,
		pluginManager,
		portForwarder,
		terminalManager,
		options.Context,
		restConfigOptions)

	moduleList, err := initModules(ctx, dashConfig, options.Namespace, options)
	if err != nil {
		return fmt.Errorf("initializing modules: %w", err)
	}

	for _, mod := range moduleList {
		if err := moduleManager.Register(mod); err != nil {
			return fmt.Errorf("loading module %s: %w", mod.Name(), err)
		}
	}

	if err := pluginManager.Start(ctx); err != nil {
		return fmt.Errorf("start plugin manager: %w", err)
	}

	listener, err := buildListener()
	if err != nil {
		err = fmt.Errorf("failed to create net listener: %w", err)
		return fmt.Errorf("use OCTANT_LISTENER_ADDR to set host:port: %w", err)
	}

	// Initialize the API
	apiService := api.New(ctx, api.PathPrefix, actionManger, dashConfig)
	frontendProxy.FrontendUpdateController = apiService

	d, err := newDash(listener, options.Namespace, options.FrontendURL, apiService, logger)
	if err != nil {
		return fmt.Errorf("failed to create dash instance: %w", err)
	}

	if viper.GetBool("disable-open-browser") {
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
		return nil, fmt.Errorf("nil cluster client")
	}

	resourceAccess := objectstore.NewResourceAccess(client)
	appObjectStore, err := objectstore.NewDynamicCache(ctx, client, objectstore.Access(resourceAccess))

	if err != nil {
		return nil, fmt.Errorf("creating object store for app: %w", err)
	}

	return appObjectStore, nil
}

func initPortForwarder(ctx context.Context, client cluster.ClientInterface, appObjectStore store.Store) (portforward.PortForwarder, error) {
	return portforward.Default(ctx, client, appObjectStore)
}

type moduleOptions struct {
	clusterClient  *cluster.Cluster
	crdWatcher     config.CRDWatcher
	namespace      string
	logger         log.Logger
	pluginManager  *plugin.Manager
	portForwarder  portforward.PortForwarder
	kubeConfigPath string
	actionManager  *action.Manager
}

func initModules(ctx context.Context, dashConfig config.Dash, namespace string, options Options) ([]module.Module, error) {
	var list []module.Module

	podViewOptions := workloads.Options{
		DashConfig: dashConfig,
	}
	workloadModule, err := workloads.New(ctx, podViewOptions)
	if err != nil {
		return nil, fmt.Errorf("initialize workload module: %w", err)
	}

	list = append(list, workloadModule)

	if viper.GetBool("enable-feature-applications") {
		applicationsOptions := applications.Options{
			DashConfig: dashConfig,
		}
		applicationsModule := applications.New(ctx, applicationsOptions)
		list = append(list, applicationsModule)
	}

	overviewOptions := overview.Options{
		Namespace:  namespace,
		DashConfig: dashConfig,
	}
	overviewModule, err := overview.New(ctx, overviewOptions)
	if err != nil {
		return nil, fmt.Errorf("create overview module: %w", err)
	}

	list = append(list, overviewModule)

	if !options.DisableClusterOverview {
		clusterOverviewOptions := clusteroverview.Options{
			DashConfig: dashConfig,
		}
		clusterOverviewModule, err := clusteroverview.New(ctx, clusterOverviewOptions)
		if err != nil {
			return nil, fmt.Errorf("create cluster overview module: %w", err)
		}

		list = append(list, clusterOverviewModule)
	}

	configurationOptions := configuration.Options{
		DashConfig:     dashConfig,
		KubeConfigPath: dashConfig.KubeConfigPath(),
	}
	configurationModule := configuration.New(ctx, configurationOptions)

	list = append(list, configurationModule)

	localContentPath := viper.GetString("local-content")
	if localContentPath != "" {
		localContentModule := localcontent.New(localContentPath)
		list = append(list, localContentModule)
	}

	return list, nil
}

// initModuleManager initializes the moduleManager (and currently the modules themselves)
func initModuleManager(options *moduleOptions) (*module.Manager, error) {
	moduleManager, err := module.NewManager(options.clusterClient, options.namespace, options.actionManager, options.logger)
	if err != nil {
		return nil, fmt.Errorf("create module manager: %w", err)
	}

	return moduleManager, nil
}

func buildListener() (net.Listener, error) {
	listenerAddr := api.ListenerAddr()
	conn, err := net.DialTimeout("tcp", listenerAddr, time.Millisecond*500)
	if err != nil {
		return net.Listen("tcp", listenerAddr)
	}
	_ = conn.Close()
	return nil, fmt.Errorf("tcp %s: dial: already in use", listenerAddr)
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
			os.Exit(1) // TODO graceful shutdown for other goroutines (GH#494)
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
	var frontendHandler http.Handler
	frontendPath := viper.GetString("proxy-frontend")
	if frontendPath == "" {
		d.logger.Infof("Using embedded Octant frontend")
		// use embedded assets
		handler, err := d.uiHandler()
		if err != nil {
			return nil, err
		}
		frontendHandler = handler
	} else {
		d.logger.With("proxy-path", frontendPath).Infof("Creating reverse proxy to Octant frontend")
		// use reverse proxy
		proxyURL, err := url.Parse(frontendPath)
		if err != nil {
			return nil, err
		}
		frontendHandler = httputil.NewSingleHostReverseProxy(proxyURL)
	}

	router := mux.NewRouter()
	apiHandler, err := d.apiHandler.Handler(ctx)
	if err != nil {
		return nil, err
	}

	router.PathPrefix(api.PathPrefix).Handler(apiHandler)

	router.PathPrefix("/").Handler(frontendHandler)

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

	je, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint: agentEndpointURI,
		Process: jaeger.Process{
			ServiceName: "octant",
		},
	})

	if err != nil {
		return fmt.Errorf("failed to create Jaeger exporter: %w", err)
	}

	trace.RegisterExporter(je)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	return nil
}
