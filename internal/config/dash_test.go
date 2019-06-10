package config

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"


	clusterFake "github.com/heptio/developer-dash/internal/cluster/fake"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/module"
	moduleFake "github.com/heptio/developer-dash/internal/module/fake"
	objectstoreFake "github.com/heptio/developer-dash/internal/objectstore/fake"
	portForwardFake "github.com/heptio/developer-dash/internal/portforward/fake"
	"github.com/heptio/developer-dash/internal/testutil"
	pluginFake "github.com/heptio/developer-dash/pkg/plugin/fake"
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
	m := moduleFake.NewModule("module", logger)

	clusterClient := clusterFake.NewMockClientInterface(controller)
	crdWatcher := stubCRDWatcher{}
	moduleManager := moduleFake.NewStubManager("", []module.Module{m})
	objectStore := objectstoreFake.NewMockObjectStore(controller)
	pluginManager := pluginFake.NewMockManagerInterface(controller)
	portForwarder := portForwardFake.NewMockPortForwarder(controller)

	config := NewLiveConfig(clusterClient, crdWatcher, logger, moduleManager, objectStore, pluginManager, portForwarder)

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
