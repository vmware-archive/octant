/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package config

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"

	"github.com/vmware-tanzu/octant/internal/cluster"
	clusterFake "github.com/vmware-tanzu/octant/internal/cluster/fake"
	internalErr "github.com/vmware-tanzu/octant/internal/errors"
	"github.com/vmware-tanzu/octant/internal/log"
	moduleFake "github.com/vmware-tanzu/octant/internal/module/fake"
	portForwardFake "github.com/vmware-tanzu/octant/internal/portforward/fake"
	terminalFake "github.com/vmware-tanzu/octant/internal/terminal/fake"
	"github.com/vmware-tanzu/octant/internal/testutil"
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
			config := CRDWatchConfig{
				IsNamespaced: test.isNamespaced,
			}

			crd := testutil.CreateCRD("my-crd")
			if test.isNamespaced {
				crd.Spec.Scope = apiextv1beta1.NamespaceScoped
			} else {
				crd.Spec.Scope = apiextv1beta1.ClusterScoped
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
	terminalManager := terminalFake.NewMockManager(controller)
	kubeConfigPath := "/path"

	objectStore.EXPECT().
		RegisterOnUpdate(gomock.Any())

	contextName := "context-name"
	restConfigOptions := cluster.RESTConfigOptions{}

	config := NewLiveConfig(clusterClient, crdWatcher, kubeConfigPath, logger, moduleManager, objectStore,
		errorStore, pluginManager, portForwarder,
		terminalManager, contextName, restConfigOptions)

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

type stubCRDWatcher struct{}

var _ CRDWatcher = (*stubCRDWatcher)(nil)

func (stubCRDWatcher) Watch(_ context.Context, config *CRDWatchConfig) error {
	return nil
}
