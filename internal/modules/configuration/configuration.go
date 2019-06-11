package configuration

import (
	"context"
	"net/http"
	"path"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/heptio/developer-dash/internal/api"
	"github.com/heptio/developer-dash/internal/clustereye"
	"github.com/heptio/developer-dash/internal/config"
	"github.com/heptio/developer-dash/internal/describer"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/module"
	"github.com/heptio/developer-dash/pkg/view/component"
)

type Options struct {
	DashConfig     config.Dash
	KubeConfigPath string
}

type Configuration struct {
	Options

	pathMatcher *describer.PathMatcher
}

var _ module.Module = (*Configuration)(nil)

func New(ctx context.Context, options Options) *Configuration {
	pm := describer.NewPathMatcher("configuration")
	for _, pf := range rootDescriber.PathFilters() {
		pm.Register(ctx, pf)
	}

	return &Configuration{
		Options:     options,
		pathMatcher: pm,
	}
}

func (Configuration) Name() string {
	return "configuration"
}

func (c *Configuration) Handlers(ctx context.Context) map[string]http.Handler {
	logger := log.From(ctx)

	update := &updateCurrentContextHandler{
		logger: logger,
		contextUpdateFunc: func(name string) error {
			logger.With("new-kube-context", name).Debugf("update context stub")
			return nil
		},
	}

	return map[string]http.Handler{
		"/kube-contexts": update,
	}
}

func (c *Configuration) Content(ctx context.Context, contentPath, prefix, namespace string, opts module.ContentOptions) (component.ContentResponse, error) {
	pf, err := c.pathMatcher.Find(contentPath)
	if err != nil {
		if err == describer.ErrPathNotFound {
			return describer.EmptyContentResponse, api.NewNotFoundError(contentPath)
		}
		return describer.EmptyContentResponse, err
	}

	options := describer.Options{
		Fields:   pf.Fields(contentPath),
		LabelSet: opts.LabelSet,
		Dash:     c.DashConfig,
	}

	cResponse, err := pf.Describer.Describe(ctx, prefix, namespace, options)
	if err != nil {
		return describer.EmptyContentResponse, err
	}

	return cResponse, nil
}

func (c *Configuration) ContentPath() string {
	return c.Name()
}

func (c *Configuration) Navigation(ctx context.Context, namespace, root string) ([]clustereye.Navigation, error) {
	return []clustereye.Navigation{
		{
			Title: "Configuration",
			Path:  path.Join("/content", c.ContentPath(), "/"),
			Children: []clustereye.Navigation{
				{
					Title: "Plugins",
					Path:  path.Join("/content", c.ContentPath(), "plugins"),
				},
			},
		},
	}, nil
}

func (Configuration) SetNamespace(namespace string) error {
	return nil
}

func (Configuration) Start() error {
	return nil
}

func (Configuration) Stop() {
}

func (c Configuration) SupportedGroupVersionKind() []schema.GroupVersionKind {
	return []schema.GroupVersionKind{}
}

func (c Configuration) GroupVersionKindPath(namespace, apiVersion, kind, name string) (string, error) {
	return "", errors.Errorf("configuration can't create paths for %s %s", apiVersion, kind)
}

func (c Configuration) AddCRD(ctx context.Context, crd *unstructured.Unstructured) error {
	return nil
}

func (c Configuration) RemoveCRD(ctx context.Context, crd *unstructured.Unstructured) error {
	return nil
}

// Generators allow modules to send events to the frontend.
func (c Configuration) Generators() []clustereye.Generator {
	kcg := newKubeContextGenerator(c.DashConfig.KubeConfigsPaths()[0])
	return []clustereye.Generator{
		kcg,
	}
}
