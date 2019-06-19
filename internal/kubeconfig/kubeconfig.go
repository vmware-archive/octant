/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package kubeconfig

import (
	"github.com/spf13/afero"
	"k8s.io/apimachinery/pkg/util/yaml"
	clientcmdapiv1 "k8s.io/client-go/tools/clientcmd/api/v1"
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
	AppFS afero.Fs
}

var _ Loader = (*FSLoader)(nil)

// NewFSLoader creates an instance of FSLoader.
func NewFSLoader(options ...FSLoaderOpt) *FSLoader {
	l := &FSLoader{
		AppFS: afero.NewOsFs(),
	}

	for _, option := range options {
		option(l)
	}

	return l
}

// Load loads a kube config contexts from a file.
func (l *FSLoader) Load(filename string) (*KubeConfig, error) {
	var rawConfig *clientcmdapiv1.Config

	f, err := l.AppFS.Open(filename)
	if err != nil {
		return nil, err
	}

	defer func() {
		if cErr := f.Close(); cErr != nil && err == nil {
			err = cErr
		}
	}()

	if err := yaml.NewYAMLToJSONDecoder(f).Decode(&rawConfig); err != nil {
		return nil, err
	}

	var list []Context

	for _, kubeContext := range rawConfig.Contexts {
		list = append(list, Context{Name: kubeContext.Name})
	}

	return &KubeConfig{
		Contexts:       list,
		CurrentContext: rawConfig.CurrentContext,
	}, nil
}
