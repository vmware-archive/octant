/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package dash

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/vmware-tanzu/octant/internal/util/path_util"

	"contrib.go.opencensus.io/exporter/jaeger"
	"github.com/skratchdot/open-golang/open"
	"github.com/soheilhy/cmux"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"go.opencensus.io/trace"

	internalAPI "github.com/vmware-tanzu/octant/internal/api"
	internalConfig "github.com/vmware-tanzu/octant/internal/config"
	ocontext "github.com/vmware-tanzu/octant/internal/context"
	"github.com/vmware-tanzu/octant/internal/describer"
	oerrors "github.com/vmware-tanzu/octant/internal/errors"
	"github.com/vmware-tanzu/octant/internal/kubeconfig"
	internalLog "github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/module"
	"github.com/vmware-tanzu/octant/internal/modules/applications"
	"github.com/vmware-tanzu/octant/internal/modules/clusteroverview"
	"github.com/vmware-tanzu/octant/internal/modules/configuration"
	"github.com/vmware-tanzu/octant/internal/modules/localcontent"
	"github.com/vmware-tanzu/octant/internal/modules/overview"
	"github.com/vmware-tanzu/octant/internal/modules/workloads"
	"github.com/vmware-tanzu/octant/internal/objectstore"
	"github.com/vmware-tanzu/octant/internal/portforward"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/api"
	"github.com/vmware-tanzu/octant/pkg/api/websockets"
	"github.com/vmware-tanzu/octant/pkg/cluster"
	"github.com/vmware-tanzu/octant/pkg/config"
	"github.com/vmware-tanzu/octant/pkg/log"
	"github.com/vmware-tanzu/octant/pkg/octant"
	"github.com/vmware-tanzu/octant/pkg/plugin"
	pluginAPI "github.com/vmware-tanzu/octant/pkg/plugin/api"
	"github.com/vmware-tanzu/octant/pkg/store"
	_ "github.com/vmware-tanzu/octant/web"
)

type Runner struct {
	ctx                        context.Context
	dash                       *dash
	pluginManager              *plugin.Manager
	moduleManager              *module.Manager
	actionManager              *action.Manager
	streamingConnectionManager *api.StreamingConnectionManager
	apiCreated                 bool
	fs                         afero.Fs
}

func NewRunner(ctx context.Context, logger log.Logger, opts ...RunnerOption) (*Runner, error) {
	options := Options{}
	for _, opt := range opts {
		opt.nonClusterOption(&options)
	}

	r := Runner{}
	ctx = internalLog.WithLoggerContext(ctx, logger)
	ctx = ocontext.WithKubeConfigCh(ctx)
	r.ctx = ctx

	if options.Context != "" {
		logger.With("initial-context", options.Context).Infof("Setting initial context from user flags")
	}

	actionManger := action.NewManager(logger)
	r.actionManager = actionManger

	var streamingConnectionManager *api.StreamingConnectionManager
	if options.streamingClientFactory != nil {
		streamingConnectionManager = api.NewStreamingConnectionManager(ctx, r.actionManager, options.streamingClientFactory)
	} else {
		streamingConnectionManager = api.NewStreamingConnectionManager(ctx, r.actionManager, websockets.NewWebsocketConnectionFactory())
	}
	r.streamingConnectionManager = streamingConnectionManager
	go streamingConnectionManager.Run(ctx)

	var err error

	var pluginService *pluginAPI.GRPCService
	var apiService internalAPI.Service
	var apiErr error

	r.fs = afero.NewOsFs()

	if options.clusterClient != nil {
		apiService, pluginService, apiErr = r.initAPI(ctx, logger, opts...)
	} else {
		apiService, pluginService, apiErr = r.apiFromKubeConfig(options.KubeConfig, opts...)
	}
	if apiErr != nil {
		return nil, fmt.Errorf("failed to start service api: %w", apiErr)
	}

	d, err := newDash(options.Listener, options.Namespace, options.FrontendURL, options.BrowserPath, apiService, pluginService, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create dash instance: %w", err)
	}

	if viper.GetBool("disable-open-browser") {
		d.willOpenBrowser = false
	}

	r.dash = d

	return &r, nil
}

func (r *Runner) apiFromKubeConfig(kubeConfig string, opts ...RunnerOption) (internalAPI.Service, *pluginAPI.GRPCService, error) {
	logger := internalLog.From(r.ctx)
	validKubeConfig, err := ValidateKubeConfig(logger, kubeConfig, r.fs)
	if err == nil {
		opts = append(opts, WithKubeConfig(validKubeConfig))
		return r.initAPI(r.ctx, logger, opts...)
	} else {
		logger.Infof("no valid kube config found, initializing loading API")
		return internalAPI.NewLoadingAPI(r.ctx, internalAPI.PathPrefix, r.actionManager, r.streamingConnectionManager, logger), nil, nil
	}
}

func (r *Runner) Start(startupCh, shutdownCh chan bool, opts ...RunnerOption) error {
	options := Options{}
	for _, opt := range opts {
		opt.nonClusterOption(&options)
	}

	logger := internalLog.From(r.ctx)
	go func() {
		if err := r.dash.Run(r.ctx, startupCh); err != nil {
			logger.Debugf("running dashboard service: %v", err)
		}
	}()

	if !r.apiCreated {
		go func() {
			logger.Infof("waiting for kube config ...")
			options.KubeConfig = <-ocontext.KubeConfigChFrom(r.ctx)
			opts = append(opts, WithKubeConfig(options.KubeConfig))

			if options.KubeConfig == "" {
				logger.Errorf("unexpected empty kube config")
				return
			}
			logger.Debugf("Loading configuration: %v", options.KubeConfig)
			apiService, pluginService, err := r.initAPI(r.ctx, logger, opts...)
			if err != nil {
				logger.Errorf("cannot create api: %v", err)
			}
			r.dash.apiHandler = apiService
			r.dash.pluginService = pluginService
			hf := octant.NewHandlerFactory(
				octant.BackendHandler(r.dash.apiHandler.Handler),
				octant.FrontendURL(viper.GetString("proxy-frontend")))

			r.dash.server.Handler, err = hf.Handler(r.ctx)
			if err != nil {
				logger.Errorf("cannot create handler: %v", err)
			}

			logger.Infof("using api service")
		}()
	}

	<-r.ctx.Done()

	shutdownCtx := internalLog.WithLoggerContext(context.Background(), logger)

	if r.apiCreated {
		r.moduleManager.Unload()
		r.pluginManager.Stop(shutdownCtx)
	}

	shutdownCh <- true
	return nil
}

func (r *Runner) initAPI(ctx context.Context, logger log.Logger, opts ...RunnerOption) (*internalAPI.API, *pluginAPI.GRPCService, error) {
	kubeConfigOptions := []kubeconfig.KubeConfigOption{}
	options := Options{}
	for _, opt := range opts {
		kubeConfigOptions = append(kubeConfigOptions, opt.kubeConfigOption)
		opt.nonClusterOption(&options)
	}
	frontendProxy := pluginAPI.FrontendProxy{}

	var kubeContextDecorator internalConfig.KubeContextDecorator
	if options.clusterClient != nil {
		kubeContextDecorator = internalConfig.StaticClusterClient(options.clusterClient)
	} else {
		var err error
		kubeContextDecorator, err = kubeconfig.NewKubeConfigContextManager(ctx, kubeConfigOptions...)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to init cluster client, does your kube config have a current-context set?: %w", err)
		}
	}
	clusterClient := kubeContextDecorator.ClusterClient()

	if options.EnableOpenCensus {
		if err := enableOpenCensus(); err != nil {
			logger.Infof("Enabling OpenCensus")
			return nil, nil, fmt.Errorf("enabling open census: %w", err)
		}
	}

	if options.EnableMemStats {
		if err := memStats(); err != nil {
			logger.Infof("Enable MemStat")
			return nil, nil, fmt.Errorf("enabling memstat: %w", err)
		}
	}

	nsClient, err := clusterClient.NamespaceClient()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create namespace client: %w", err)
	}

	// If not overridden, use initial namespace from current context in KUBECONFIG
	if options.Namespace == "" {
		options.Namespace = nsClient.InitialNamespace()
	}

	logger.Debugf("initial namespace for dashboard is %s", options.Namespace)

	var appObjectStore store.Store
	if options.objectStore != nil {
		appObjectStore = options.objectStore
	} else {
		factoryOption := objectstore.WithDynamicSharedInformerFactory(options.factory)
		appObjectStore, err = initObjectStore(ctx, clusterClient, factoryOption)
	}

	if err != nil {
		return nil, nil, fmt.Errorf("initializing store: %w", err)
	}

	errorStore, err := oerrors.NewErrorStore()
	if err != nil {
		return nil, nil, fmt.Errorf("initializing error store: %w", err)
	}

	crdWatcher, err := describer.NewDefaultCRDWatcher(ctx, clusterClient, appObjectStore, errorStore)
	if err != nil {
		var ae *oerrors.AccessError
		if errors.As(err, &ae) {
			if ae.Name() == oerrors.OctantAccessError {
				logger.Warnf("skipping CRD watcher due to access denied error starting watcher")
			}
		} else {
			return nil, nil, fmt.Errorf("initializing CRD watcher: %w", err)
		}
	}

	portForwarder, err := initPortForwarder(ctx, clusterClient, appObjectStore)
	if err != nil {
		return nil, nil, fmt.Errorf("initializing port forwarder: %w", err)
	}

	mo := &moduleOptions{
		clusterClient: clusterClient,
		namespace:     options.Namespace,
		logger:        logger,
		actionManager: r.actionManager,
	}
	moduleManager, err := initModuleManager(mo)
	if err != nil {
		return nil, nil, fmt.Errorf("init module manager: %w", err)
	}

	r.moduleManager = moduleManager
	pluginDashboardService := &pluginAPI.GRPCService{
		LinkGenerator:          moduleManager,
		ObjectStore:            appObjectStore,
		PortForwarder:          portForwarder,
		NamespaceInterface:     nsClient,
		FrontendProxy:          frontendProxy,
		WebsocketClientManager: r.streamingConnectionManager,
	}

	pluginManager, err := initPlugin(moduleManager, r.actionManager, r.streamingConnectionManager, pluginDashboardService)
	if err != nil {
		return nil, nil, fmt.Errorf("initializing plugin manager: %w", err)
	}

	r.pluginManager = pluginManager

	buildInfo := config.BuildInfo{
		Version: options.BuildInfo.Version,
		Commit:  options.BuildInfo.Commit,
		Time:    options.BuildInfo.Time,
	}

	restConfigOptions := cluster.RESTConfigOptions{
		QPS:       options.ClientQPS,
		Burst:     options.ClientBurst,
		UserAgent: options.UserAgent,
	}
	dashConfig := internalConfig.NewLiveConfig(
		kubeContextDecorator,
		crdWatcher,
		logger,
		moduleManager,
		appObjectStore,
		errorStore,
		pluginManager,
		portForwarder,
		restConfigOptions,
		buildInfo,
		options.KubeConfig,
		false,
	)

	pluginManager.SetOctantClient(dashConfig)

	if err := watchConfigs(ctx, dashConfig, options.KubeConfig); err != nil {
		return nil, nil, fmt.Errorf("set up config watcher: %w", err)
	}

	moduleList, err := initModules(ctx, dashConfig, options.Namespace, options)
	if err != nil {
		return nil, nil, fmt.Errorf("initializing modules: %w", err)
	}

	for _, mod := range moduleList {
		if err := moduleManager.Register(mod); err != nil {
			return nil, nil, fmt.Errorf("loading module %s: %w", mod.Name(), err)
		}
	}

	if err := pluginManager.Start(ctx); err != nil {
		return nil, nil, fmt.Errorf("start plugin manager: %w", err)
	}

	// Watch for CRDs after modules initialized
	if err := crdWatcher.Watch(ctx); err != nil {
		return nil, nil, fmt.Errorf("unable to start CRD watcher: %w", err)
	}

	apiService := internalAPI.New(ctx, internalAPI.PathPrefix, r.actionManager, r.streamingConnectionManager, dashConfig)
	frontendProxy.FrontendUpdateController = apiService

	r.apiCreated = true
	return apiService, pluginDashboardService, nil
}

// initObjectStore initializes the cluster object store interface
func initObjectStore(ctx context.Context, client cluster.ClientInterface, opts ...objectstore.Option) (store.Store, error) {
	if client == nil {
		return nil, fmt.Errorf("nil cluster client")
	}

	appObjectStore, err := objectstore.NewDynamicCache(ctx, client, opts...)

	if err != nil {
		return nil, fmt.Errorf("creating object store for app: %w", err)
	}

	return appObjectStore, nil
}

func initPortForwarder(ctx context.Context, client cluster.ClientInterface, appObjectStore store.Store) (portforward.PortForwarder, error) {
	return portforward.Default(ctx, client, appObjectStore)
}

type moduleOptions struct {
	clusterClient  cluster.ClientInterface
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
		DashConfig: dashConfig,
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

type dash struct {
	mux             cmux.CMux
	listener        net.Listener
	uiURL           string
	browserPath     string
	namespace       string
	defaultHandler  func() (http.Handler, error)
	apiHandler      internalAPI.Service
	willOpenBrowser bool
	logger          log.Logger
	handlerFactory  *octant.HandlerFactory
	server          http.Server
	pluginService   pluginAPI.Service
}

func newDash(listener net.Listener, namespace, uiURL string, browserPath string, apiHandler internalAPI.Service, pluginHandler pluginAPI.Service, logger log.Logger) (*dash, error) {
	hf := octant.NewHandlerFactory(
		octant.BackendHandler(apiHandler.Handler),
		octant.FrontendURL(viper.GetString("proxy-frontend")))

	return &dash{
		mux:             cmux.New(listener),
		handlerFactory:  hf,
		listener:        listener,
		namespace:       namespace,
		uiURL:           uiURL,
		browserPath:     browserPath,
		defaultHandler:  octant.GetFrontendHandler,
		willOpenBrowser: true,
		apiHandler:      apiHandler,
		pluginService:   pluginHandler,
		logger:          logger,
	}, nil
}

func (d *dash) SetAPIService(ctx context.Context, apiService internalAPI.Service) error {
	d.apiHandler = apiService
	hf := octant.NewHandlerFactory(
		octant.BackendHandler(d.apiHandler.Handler),
		octant.FrontendURL(viper.GetString("proxy-frontend")))
	var err error
	d.server.Handler, err = hf.Handler(ctx)
	return err
}

func (d *dash) SetFrontendHandler(fn octant.HandlerFactoryFunc) {
	d.handlerFactory.SetFrontend(fn)
}

func (d *dash) Run(ctx context.Context, startupCh chan bool) error {
	handler, err := d.handlerFactory.Handler(ctx)
	if err != nil {
		return err
	}

	d.server = http.Server{Handler: handler}

	// Enable serving the plugin API on the same endpoint as the Octant streaming API.
	// This enables remote gRPC plugins.
	// if d.pluginService != nil {
	//     grpcl := d.mux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	//     go serveGRPC(grpcl, d.pluginService)
	// }

	http1 := d.mux.Match(cmux.Any())
	go func() {
		if err = d.server.Serve(http1); err != nil && err != http.ErrServerClosed {
			d.logger.Errorf("http server: %v", err)
			os.Exit(1) // TODO graceful shutdown for other goroutines (GH#494)
		}
	}()

	go func() {
		if err := d.mux.Serve(); err != nil {
			errMessage := err.Error()
			if !strings.Contains(errMessage, "use of closed network connection") {
				panic(err)
			}
		}
	}()

	dashboardURL := fmt.Sprintf("http://%s", d.listener.Addr())

	d.logger.Infof("Dashboard is available at %s\n", dashboardURL)

	if startupCh != nil {
		startupCh <- true
	}

	if d.willOpenBrowser {
		runURL := dashboardURL
		if d.browserPath != "" {
			runURL += path_util.PrefixedPath(d.browserPath)
		}
		if err = open.Run(runURL); err != nil {
			d.logger.Warnf("unable to open browser: %v", err)
		}
	}

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	return d.server.Shutdown(shutdownCtx)
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

func memStats() error {
	memstats := viper.GetString("memstats")
	if memstats != "" {
		interval, err := time.ParseDuration(viper.GetString("meminterval"))
		if err != nil {
			interval, _ = time.ParseDuration("100ms")
		}

		fileMemStats, err := os.Create(memstats)
		if err != nil {
			return err
		}

		fileMemStats.WriteString("# Time\tHeapSys\tHeapAlloc\tHeapIdle\tHeapReleased\n")
		go func() {
			var stats runtime.MemStats
			start := time.Now().UnixNano()
			for {
				runtime.ReadMemStats(&stats)
				if fileMemStats != nil {
					fileMemStats.WriteString(fmt.Sprintf("%d\t%d\t%d\t%d\t%d\n",
						(time.Now().UnixNano()-start)/1000000, stats.HeapSys, stats.HeapAlloc, stats.HeapIdle, stats.HeapReleased))
					time.Sleep(interval)
				} else {
					break
				}
			}
		}()
	}
	return nil
}

// ValidateKubeConfig returns a valid file list of kube config(s)
func ValidateKubeConfig(logger log.Logger, kubeConfig string, fs afero.Fs) (string, error) {
	fileList := []string{}
	paths := filepath.SplitList(kubeConfig)

	for _, path := range paths {
		exists, err := afero.Exists(fs, path)
		if err != nil {
			logger.Errorf("check path exists: %v", err)
		}

		if exists {
			fileList = append(fileList, path)
			continue
		}
		logger.Infof("cannot find kube config: %v", path)
	}

	if len(fileList) > 0 {
		return strings.Join(fileList, string(filepath.ListSeparator)), nil
	}
	return "", fmt.Errorf("no kubeconfig found")
}
