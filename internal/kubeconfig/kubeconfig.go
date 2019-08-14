/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package kubeconfig

import (
	"path/filepath"
	"sort"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/vmware/octant/internal/util/strings"
)

//go:generate mockgen -destination=./fake/mock_loader.go -package=fake github.com/vmware/octant/internal/kubeconfig Loader

// KubeConfig describes a kube config for dash.
type KubeConfig struct {
	Contexts       []Context
	CurrentContext string
}

// Context describes a kube config context.
type Context struct {
	Name string `json:"name"`
}

// Loader is an interface for loading kube config.
type Loader interface {
	Load(filename string) (*KubeConfig, error)
}

// FSLoaderOpt is an option for configuring FSLoader.
type FSLoaderOpt func(loader *FSLoader)

// FSLoader loads kube configs from the file system.
type FSLoader struct {
}

var _ Loader = (*FSLoader)(nil)

// NewFSLoader creates an instance of FSLoader.
func NewFSLoader(options ...FSLoaderOpt) *FSLoader {
	l := &FSLoader{}

	for _, option := range options {
		option(l)
	}

	return l
}

// Load loads a kube config contexts from a list of files.
func (l *FSLoader) Load(fileList string) (*KubeConfig, error) {
	chain := strings.Deduplicate(filepath.SplitList(fileList))

	loadingRules := &clientcmd.ClientConfigLoadingRules{
		Precedence: chain,
	}

	config, err := loadingRules.Load()
	if err != nil {
		return nil, err
	}

	var list []Context

	for name := range config.Contexts {
		list = append(list, Context{Name: name})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

	return &KubeConfig{
		Contexts:       list,
		CurrentContext: config.CurrentContext,
	}, nil
}
