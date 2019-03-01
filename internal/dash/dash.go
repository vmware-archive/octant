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

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"k8s.io/client-go/restmapper"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/heptio/developer-dash/internal/api"
	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/localcontent"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/module"
	"github.com/heptio/developer-dash/internal/overview"
	"github.com/heptio/developer-dash/web"
	"github.com/skratchdot/open-golang/open"
)

const (
	apiPathPrefix       = "/api/v1"
	defaultListenerAddr = "127.0.0.1:0"
)

// Run runs the dashboard.
func Run(ctx context.Context, namespace, uiURL, kubeconfig string, logger log.Logger) error {
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

	infoClient, err := clusterClient.InfoClient()
	if err != nil {
		return errors.Wrap(err, "failed to create info client")
	}

	cache, err := initCache(ctx.Done(), clusterClient, logger)
	if err != nil {
		return errors.Wrap(err, "initializing cache")
	}

	moduleManager, err := initModuleManager(ctx, clusterClient, cache, namespace, logger)
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

	d, err := newDash(listener, namespace, uiURL, ah, logger)
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

// initCache initializes the cluster cache interface
func initCache(stopCh <-chan struct{}, client cluster.ClientInterface, logger log.Logger) (cache.Cache, error) {
	var opts []cache.InformerCacheOpt

	if os.Getenv("DASH_VERBOSE_CACHE") != "" {
		ch := make(chan cache.Notification)

		go func() {
			for notif := range ch {
				spew.Dump(notif)
			}
		}()

		opts = append(opts, cache.InformerCacheNotificationOpt(ch, stopCh))
	}

	if client == nil {
		return nil, errors.New("nil cluster client")
	}

	dynamicClient, err := client.DynamicClient()
	if err != nil {
		return nil, errors.Wrapf(err, "creating DynamicClient")
	}
	di, err := client.DiscoveryClient()
	if err != nil {
		return nil, errors.Wrapf(err, "creating DiscoveryClient")
	}

	groupResources, err := restmapper.GetAPIGroupResources(di)
	if err != nil {
		logger.Errorf("discovering APIGroupResources: %v", err)
		return nil, errors.Wrapf(err, "mapping APIGroupResources")
	}
	rm := restmapper.NewDiscoveryRESTMapper(groupResources)

	opts = append(opts, cache.InformerCacheLoggerOpt(logger))
	informerCache := cache.NewInformerCache(stopCh, dynamicClient, rm, opts...)

	return informerCache, nil
}

// initModuleManager initializes the moduleManager (and currently the modules themselves)
func initModuleManager(ctx context.Context, clusterClient *cluster.Cluster, cache cache.Cache, namespace string, logger log.Logger) (*module.Manager, error) {
	moduleManager, err := module.NewManager(clusterClient, namespace, logger)
	if err != nil {
		return nil, errors.Wrap(err, "create module manager")
	}

	overviewModule, err := overview.NewClusterOverview(ctx, clusterClient, cache, namespace, logger)
	if err != nil {
		return nil, errors.Wrap(err, "create overview module")
	}

	moduleManager.Register(overviewModule)

	localContentPath := os.Getenv("DASH_LOCAL_CONTENT")
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
		if err = open.Run(dashboardURL); err != nil {
			d.logger.Warnf("unable to open browser: %v", err)
		}
	}

	<-ctx.Done()

	// TODO context is already done - pass a different one too allow time for graceful shutdown
	return server.Shutdown(ctx)
}

// handler configures primary http routes
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
