/*
Copyright (c) 2021 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package config

import (
	"context"

	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware-tanzu/octant/internal/kubeconfig"
	"github.com/vmware-tanzu/octant/internal/module"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/internal/portforward"
	"github.com/vmware-tanzu/octant/pkg/cluster"
	"github.com/vmware-tanzu/octant/pkg/errors"
	"github.com/vmware-tanzu/octant/pkg/log"
	"github.com/vmware-tanzu/octant/pkg/plugin"
)

// ObjectHandler is a function that is run when a new object is available.
type ObjectHandler func(ctx context.Context, object *unstructured.Unstructured)

// CRDWatcher watches for CRDs.
type CRDWatcher interface {
	Watch(ctx context.Context) error
	AddConfig(config *CRDWatchConfig) error
}

// CRDWatchConfig is configuration for CRDWatcher.
type CRDWatchConfig struct {
	Add          ObjectHandler
	Delete       ObjectHandler
	IsNamespaced bool
}

type BuildInfo struct {
	Version string
	Commit  string
	Time    string
}

type Context struct {
	Name             string
	DefaultNamespace string
}

// CanPerform returns true if config can perform actions on an object.
func (c *CRDWatchConfig) CanPerform(u *unstructured.Unstructured) bool {
	spec, ok := u.Object["spec"].(map[string]interface{})
	if !ok {
		return false
	}

	scope, ok := spec["scope"].(string)
	if !ok {
		return false
	}

	if c.IsNamespaced && scope != string(apiextv1.NamespaceScoped) {
		return false
	}

	if !c.IsNamespaced && scope != string(apiextv1.ClusterScoped) {
		return false
	}

	return true
}

// Config is configuration for dash. It has knowledge of the all the major sections of
// dash.
type Dash interface {
	octant.LinkGenerator
	octant.Storage

	ClusterClient() cluster.ClientInterface

	CRDWatcher() CRDWatcher

	ErrorStore() errors.ErrorStore

	Logger() log.Logger

	PluginManager() plugin.ManagerInterface

	PortForwarder() portforward.PortForwarder

	SetContextChosenInUI(contextChosen bool)

	UseFSContext(ctx context.Context) error

	UseContext(ctx context.Context, contextName string) error

	CurrentContext() string

	Contexts() []kubeconfig.Context

	DefaultNamespace() string

	Validate() error

	ModuleManager() module.ManagerInterface

	BuildInfo() (string, string, string)

	KubeConfigPath() string
}
