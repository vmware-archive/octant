/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package plugin_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware/octant/internal/gvk"
	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/plugin"
	"github.com/vmware/octant/pkg/plugin/fake"
	"github.com/vmware/octant/pkg/view/component"
)

func TestDefaultRunner(t *testing.T) {
	counter := 0

	pr := plugin.DefaultRunner{
		RunFunc: func(name string, gvk schema.GroupVersionKind, object runtime.Object) error {
			counter++
			return nil
		},
	}

	object := testutil.CreateDeployment("deployment")

	clientNames := []string{"plugin1", "plugin2"}

	err := pr.Run(object, clientNames)
	require.NoError(t, err)

	assert.Equal(t, 2, counter)
}

func TestDefaultRunner_object_is_nil(t *testing.T) {
	pr := plugin.DefaultRunner{
		RunFunc: func(name string, gvk schema.GroupVersionKind, object runtime.Object) error {
			return nil
		},
	}

	clientNames := []string{"plugin1", "plugin2"}

	err := pr.Run(nil, clientNames)
	require.Error(t, err)
}

func TestDefaultRunner_run_func_is_nil(t *testing.T) {
	pr := plugin.DefaultRunner{}

	object := testutil.CreateDeployment("deployment")
	clientNames := []string{"plugin1", "plugin2"}

	err := pr.Run(object, clientNames)
	require.Error(t, err)
}

func TestDefaultRunner_run_func_returns_error(t *testing.T) {
	pr := plugin.DefaultRunner{
		RunFunc: func(name string, gvk schema.GroupVersionKind, object runtime.Object) error {
			return errors.Errorf("error")
		},
	}

	object := testutil.CreateDeployment("deployment")
	clientNames := []string{"plugin1", "plugin2"}

	err := pr.Run(object, clientNames)
	require.Error(t, err)
}

func Test_PrintRunner(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	store := fake.NewMockManagerStore(controller)
	service := fake.NewMockService(controller)

	object := testutil.CreateDeployment("deployment")
	clientNames := []string{"plugin1", "plugin2"}

	plugin1Metadata := plugin.Metadata{
		Capabilities: plugin.Capabilities{
			SupportsPrinterConfig: []schema.GroupVersionKind{gvk.DeploymentGVK},
		},
	}
	store.EXPECT().
		GetMetadata(gomock.Eq("plugin1")).Return(plugin1Metadata, nil)

	plugin2Metadata := plugin.Metadata{}
	store.EXPECT().
		GetMetadata(gomock.Eq("plugin2")).Return(plugin2Metadata, nil)

	store.EXPECT().
		GetService(gomock.Eq("plugin1")).Return(service, nil)

	pr := plugin.PrintResponse{}

	service.EXPECT().
		Print(gomock.Eq(object)).Return(pr, nil)

	ch := make(chan plugin.PrintResponse)
	defer close(ch)

	runner := plugin.PrintRunner(store, ch)

	done := make(chan bool)
	go func() {
		resp := <-ch
		assert.Equal(t, pr, resp)
		done <- true
	}()

	defer func() {
		<-done
	}()

	require.NoError(t, runner.Run(object, clientNames))
}

func Test_TabRunner(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	store := fake.NewMockManagerStore(controller)
	service := fake.NewMockService(controller)

	object := testutil.CreateDeployment("deployment")
	clientNames := []string{"plugin1", "plugin2"}

	plugin1Metadata := plugin.Metadata{
		Capabilities: plugin.Capabilities{
			SupportsTab: []schema.GroupVersionKind{gvk.DeploymentGVK},
		},
	}
	store.EXPECT().
		GetMetadata(gomock.Eq("plugin1")).Return(plugin1Metadata, nil)

	plugin2Metadata := plugin.Metadata{}
	store.EXPECT().
		GetMetadata(gomock.Eq("plugin2")).Return(plugin2Metadata, nil)

	store.EXPECT().
		GetService(gomock.Eq("plugin1")).Return(service, nil)

	tab := component.Tab{}

	service.EXPECT().
		PrintTab(gomock.Eq(object)).Return(&tab, nil)

	ch := make(chan component.Tab)
	defer close(ch)

	runner := plugin.TabRunner(store, ch)

	done := make(chan bool)
	go func() {
		resp := <-ch
		assert.Equal(t, tab, resp)
		done <- true
	}()

	defer func() {
		<-done
	}()

	require.NoError(t, runner.Run(object, clientNames))
}
