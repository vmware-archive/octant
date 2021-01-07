/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package config

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"

	"github.com/vmware-tanzu/octant/internal/cluster"
	clusterFake "github.com/vmware-tanzu/octant/internal/cluster/fake"
	internalErr "github.com/vmware-tanzu/octant/internal/errors"
	"github.com/vmware-tanzu/octant/internal/gvk"
	"github.com/vmware-tanzu/octant/internal/log"
	moduleFake "github.com/vmware-tanzu/octant/internal/module/fake"
	portForwardFake "github.com/vmware-tanzu/octant/internal/portforward/fake"
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
	kubeConfigPath := "/path"
	buildInfo := BuildInfo{}

	objectStore.EXPECT().
		RegisterOnUpdate(gomock.Any())

	contextName := "context-name"
	restConfigOptions := cluster.RESTConfigOptions{}

	config := NewLiveConfig(clusterClient, crdWatcher, kubeConfigPath, logger, moduleManager, objectStore,
		errorStore, pluginManager, portForwarder,
		contextName, restConfigOptions, buildInfo)

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

func TestServerPreferredResources(t *testing.T) {
	tests := []struct {
		name         string
		resourceList []*metav1.APIResourceList
		errGroup     error
		returnErr    bool
	}{
		{
			name: "all groups discovered, no error is returned",
			resourceList: []*metav1.APIResourceList{
				{GroupVersion: "groupB/v1"},
				{GroupVersion: "apps/v1beta1"},
				{GroupVersion: "extensions/v1beta1"},
			},
		},
		{
			name: "failed to discover some groups, no error is returned",
			resourceList: []*metav1.APIResourceList{
				{GroupVersion: "groupB/v1"},
				{GroupVersion: "apps/v1beta1"},
				{GroupVersion: "extensions/v1beta1"},
			},
			errGroup: &discovery.ErrGroupDiscoveryFailed{
				Groups: map[schema.GroupVersion]error{
					gvk.PodMetrics.GroupVersion(): errors.New("server currently unavailable"),
				},
			},
		},
		{
			name:      "non ErrGroupDiscoveryFailed error, returns error",
			errGroup:  fmt.Errorf("Generic error"),
			returnErr: true,
		},
	}

	logger := log.NopLogger()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			discoveryClient := clusterFake.NewMockDiscoveryInterface(controller)
			discoveryClient.EXPECT().ServerPreferredResources().Return(test.resourceList, test.errGroup)

			resources, err := ServerPreferredResources(discoveryClient, logger)
			if test.returnErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, test.resourceList, resources)
		})
	}
}

type stubCRDWatcher struct{}

var _ CRDWatcher = (*stubCRDWatcher)(nil)

func (w stubCRDWatcher) AddConfig(config *CRDWatchConfig) error {
	return nil
}

func (stubCRDWatcher) Watch(_ context.Context) error {
	return nil
}
