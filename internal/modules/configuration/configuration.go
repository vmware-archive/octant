/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package configuration

import (
	"context"
	"path"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware-tanzu/octant/internal/api"
	"github.com/vmware-tanzu/octant/internal/describer"
	"github.com/vmware-tanzu/octant/internal/event"
	"github.com/vmware-tanzu/octant/internal/module"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/config"
	"github.com/vmware-tanzu/octant/pkg/icon"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

type Options struct {
	DashConfig     config.Dash
	KubeConfigPath string
}

type Configuration struct {
	Options

	pathMatcher          *describer.PathMatcher
	kubeContextGenerator *event.ContextsGenerator
}

var _ module.Module = (*Configuration)(nil)
var _ module.ActionReceiver = (*Configuration)(nil)

func New(ctx context.Context, options Options) *Configuration {
	pm := describer.NewPathMatcher("configuration")
	for _, pf := range rootDescriber.PathFilters() {
		pm.Register(ctx, pf)
	}

	return &Configuration{
		Options:              options,
		pathMatcher:          pm,
		kubeContextGenerator: event.NewContextsGenerator(options.DashConfig),
	}
}

func (Configuration) Name() string {
	return "configuration"
}

func (Configuration) Description() string {
	return `Plugins module displays all registered plugins and their properties.
		To find list of known plugins go to https://github.com/topics/octant-plugin`
}

func (c Configuration) ClientRequestHandlers() []octant.ClientRequestHandler {
	return nil
}

func (c *Configuration) SetContext(ctx context.Context, contextName string) error {
	return nil
}

func (c *Configuration) Content(ctx context.Context, contentPath string, opts module.ContentOptions) (component.ContentResponse, error) {
	pf, err := c.pathMatcher.Find(contentPath)
	if err != nil {
		if err == describer.ErrPathNotFound {
			return component.EmptyContentResponse, api.NewNotFoundError(contentPath)
		}
		return component.EmptyContentResponse, err
	}

	options := describer.Options{
		Fields:   pf.Fields(contentPath),
		LabelSet: opts.LabelSet,
		Dash:     c.DashConfig,
	}

	cResponse, err := pf.Describer.Describe(ctx, "", options)
	if err != nil {
		return component.EmptyContentResponse, err
	}

	return cResponse, nil
}

func (c *Configuration) ContentPath() string {
	return c.Name()
}

func (c *Configuration) Navigation(ctx context.Context, namespace, root string) ([]navigation.Navigation, error) {
	return []navigation.Navigation{
		{
			Module:   "Configuration",
			Title:    "Plugins",
			Path:     path.Join(c.ContentPath(), "plugins"),
			IconName: icon.ConfigurationPlugin,
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

func (c Configuration) ResetCRDs(ctx context.Context) error {
	return nil
}

// Generators allow modules to send events to the frontend.
func (c Configuration) Generators() []octant.Generator {
	return []octant.Generator{
		c.kubeContextGenerator,
	}
}

func (c *Configuration) ActionPaths() map[string]action.DispatcherFunc {
	objectDeleter := NewObjectDeleter(c.DashConfig.Logger(), c.DashConfig.ObjectStore())

	return map[string]action.DispatcherFunc{
		objectDeleter.ActionName(): objectDeleter.Handle,
	}
}

func (c *Configuration) GvkFromPath(contentPath, namespace string) (schema.GroupVersionKind, error) {
	return schema.GroupVersionKind{}, errors.Errorf("not supported")
}
