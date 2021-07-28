/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package config_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	. "github.com/vmware-tanzu/octant/internal/config"

	clusterFake "github.com/vmware-tanzu/octant/internal/cluster/fake"
	configFake "github.com/vmware-tanzu/octant/internal/config/fake"
	internalErr "github.com/vmware-tanzu/octant/internal/errors"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/module"
	moduleFake "github.com/vmware-tanzu/octant/internal/module/fake"
	portForwardFake "github.com/vmware-tanzu/octant/internal/portforward/fake"
	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/cluster"
	"github.com/vmware-tanzu/octant/pkg/config"
	pluginFake "github.com/vmware-tanzu/octant/pkg/plugin/fake"
	objectStoreFake "github.com/vmware-tanzu/octant/pkg/store/fake"
)

func TestCRDWatchConfig_CanPerform(t *testing.T) {
	tests := []struct {
		name         string
		isNamespaced bool
		namespace    string
		expected     bool
	}{
		{
			name:         "is namespaced / populated namespace",
			isNamespaced: true,
			namespace:    "default",
			expected:     true,
		},
		{
			name:         "is not namespaced / blank namespace",
			isNamespaced: false,
			namespace:    "",
			expected:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := config.CRDWatchConfig{
				IsNamespaced: test.isNamespaced,
			}

			crd := testutil.CreateCRD("my-crd")
			if test.isNamespaced {
				crd.Spec.Scope = apiextv1.NamespaceScoped
			} else {
				crd.Spec.Scope = apiextv1.ClusterScoped
			}

			got := config.CanPerform(testutil.ToUnstructured(t, crd))

			assert.Equal(t, test.expected, got)
		})
	}
}

func TestLiveConfig(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	logger := log.NopLogger()

	clusterClient := clusterFake.NewMockClientInterface(controller)
	crdWatcher := stubCRDWatcher{}

	moduleManager := moduleFake.NewMockManagerInterface(controller)
	moduleManager.EXPECT().
		ObjectPath(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return("/pod", nil)

	objectStore := objectStoreFake.NewMockStore(controller)
	errorStore, err := internalErr.NewErrorStore()
	assert.NoError(t, err)
	pluginManager := pluginFake.NewMockManagerInterface(controller)
	portForwarder := portForwardFake.NewMockPortForwarder(controller)
	buildInfo := config.BuildInfo{}

	restConfigOptions := cluster.RESTConfigOptions{}

	config := NewLiveConfig(
		StaticClusterClient(clusterClient),
		crdWatcher,
		logger,
		moduleManager,
		objectStore,
		errorStore,
		pluginManager,
		portForwarder,
		restConfigOptions,
		buildInfo,
		"",
		false,
	)

	assert.NoError(t, config.Validate())
	assert.Equal(t, clusterClient, config.ClusterClient())
	assert.Equal(t, crdWatcher, config.CRDWatcher())
	assert.Equal(t, logger, config.Logger())
	assert.Equal(t, objectStore, config.ObjectStore())
	assert.Equal(t, pluginManager, config.PluginManager())
	assert.Equal(t, portForwarder, config.PortForwarder())

	objectPath, err := config.ObjectPath("", "", "", "")
	require.NoError(t, err)
	assert.Equal(t, "/pod", objectPath)
}

func TestLiveConfig_UseContext_WithContextChosenByUISetToTrue(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	logger := log.NopLogger()

	crdWatcher := stubCRDWatcher{}

	moduleManager := moduleFake.NewMockManagerInterface(controller)

	objectStore := objectStoreFake.NewMockStore(controller)
	errorStore, err := internalErr.NewErrorStore()
	assert.NoError(t, err)
	pluginManager := pluginFake.NewMockManagerInterface(controller)
	portForwarder := portForwardFake.NewMockPortForwarder(controller)
	buildInfo := config.BuildInfo{}

	restConfigOptions := cluster.RESTConfigOptions{}

	contextDecorator := configFake.NewMockKubeContextDecorator(controller)

	config := NewLiveConfig(
		contextDecorator,
		crdWatcher,
		logger,
		moduleManager,
		objectStore,
		errorStore,
		pluginManager,
		portForwarder,
		restConfigOptions,
		buildInfo,
		"",
		true, // contextChosenInUI
	)

	objectStore.EXPECT().UpdateClusterClient(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	moduleManager.EXPECT().Modules().Return(make([]module.Module, 0)).AnyTimes()
	pluginManager.EXPECT().SetOctantClient(gomock.Eq(config)).AnyTimes()
	contextDecorator.EXPECT().ClusterClient().AnyTimes()

	newContext := ""
	currentContext := "socketContext"
	// Since UseContext is called with "" with contextChosenInUI = true, our context should stay as currentContext
	moduleManager.EXPECT().UpdateContext(gomock.Any(), currentContext).Return(nil)
	contextDecorator.EXPECT().CurrentContext().Return(currentContext).AnyTimes()
	contextDecorator.EXPECT().SwitchContext(gomock.Any(), gomock.Eq(currentContext)).Return(nil)

	config.UseContext(context.TODO(), newContext)

	newContext = "newContext"
	// Since UseContext is called with "newContext" with contextChosenInUI = true, our context should change to newContext
	moduleManager.EXPECT().UpdateContext(gomock.Any(), newContext).Return(nil)
	contextDecorator.EXPECT().SwitchContext(gomock.Any(), gomock.Eq(newContext)).Return(nil)

	config.UseContext(context.TODO(), newContext)
}

func TestLiveConfig_UseContext_WithContextChosenByUISetToFalse(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	logger := log.NopLogger()

	crdWatcher := stubCRDWatcher{}

	moduleManager := moduleFake.NewMockManagerInterface(controller)

	objectStore := objectStoreFake.NewMockStore(controller)
	errorStore, err := internalErr.NewErrorStore()
	assert.NoError(t, err)
	pluginManager := pluginFake.NewMockManagerInterface(controller)
	portForwarder := portForwardFake.NewMockPortForwarder(controller)
	buildInfo := config.BuildInfo{}

	restConfigOptions := cluster.RESTConfigOptions{}

	contextDecorator := configFake.NewMockKubeContextDecorator(controller)

	config := NewLiveConfig(
		contextDecorator,
		crdWatcher,
		logger,
		moduleManager,
		objectStore,
		errorStore,
		pluginManager,
		portForwarder,
		restConfigOptions,
		buildInfo,
		"",
		false, // contextChosenInUI
	)

	objectStore.EXPECT().UpdateClusterClient(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	moduleManager.EXPECT().Modules().Return(make([]module.Module, 0)).AnyTimes()
	pluginManager.EXPECT().SetOctantClient(gomock.Eq(config)).AnyTimes()
	contextDecorator.EXPECT().ClusterClient().AnyTimes()

	newContext := ""
	currentContext := "socketContext"
	// Since UseContext is called with "" with contextChosenInUI = false, our context should change to newContext
	moduleManager.EXPECT().UpdateContext(gomock.Any(), newContext).Return(nil)
	contextDecorator.EXPECT().CurrentContext().Return(currentContext).AnyTimes()
	contextDecorator.EXPECT().SwitchContext(gomock.Any(), gomock.Eq(newContext)).Return(nil)

	config.UseContext(context.TODO(), newContext)

	newContext = "newContext"
	// Since UseContext is called with "newContext" with contextChosenInUI = false, our context should change to newContext
	moduleManager.EXPECT().UpdateContext(gomock.Any(), newContext).Return(nil)
	contextDecorator.EXPECT().SwitchContext(gomock.Any(), gomock.Eq(newContext)).Return(nil)

	config.UseContext(context.TODO(), newContext)
}

type stubCRDWatcher struct{}

var _ config.CRDWatcher = (*stubCRDWatcher)(nil)

func (w stubCRDWatcher) AddConfig(config *config.CRDWatchConfig) error {
	return nil
}

func (stubCRDWatcher) Watch(_ context.Context) error {
	return nil
}
