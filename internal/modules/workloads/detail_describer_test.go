package workloads

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	clusterFake "github.com/vmware-tanzu/octant/internal/cluster/fake"
	configFake "github.com/vmware-tanzu/octant/internal/config/fake"
	"github.com/vmware-tanzu/octant/internal/describer"
	"github.com/vmware-tanzu/octant/internal/queryer"
	queryerFake "github.com/vmware-tanzu/octant/internal/queryer/fake"
	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/plugin"
	pluginFake "github.com/vmware-tanzu/octant/pkg/plugin/fake"
	"github.com/vmware-tanzu/octant/pkg/store"
	objectStoreFake "github.com/vmware-tanzu/octant/pkg/store/fake"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func TestDetailDescriber_Describe(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	ctx := context.Background()

	dd, err := NewDetailDescriber()
	require.NoError(t, err)

	tdo := newTestDescriberOptions(t, controller)
	describerOptions := tdo.ToOptions()

	result, err := dd.Describe(ctx, "namespace", describerOptions)
	require.NoError(t, err)

	for _, c := range result.Components {
		f, ok := c.(*component.FlexLayout)
		require.Equal(t, true, ok)
		require.Equal(t, component.TitleFromString("Workload layout"), f.Title)
		require.Equal(t, 2, len(f.Config.Sections))
	}
}

type testDescriberOptions struct {
	dashConfig *configFake.MockDash
	queryer    queryer.Queryer
}

func newTestDescriberOptions(t *testing.T, controller *gomock.Controller) *testDescriberOptions {
	dashConfig := configFake.NewMockDash(controller)

	clusterClient := clusterFake.NewMockClientInterface(controller)
	objectStore := objectStoreFake.NewMockStore(controller)
	discoveryInterface := clusterFake.NewMockDiscoveryInterface(controller)
	pluginManager := pluginFake.NewMockManagerInterface(controller)

	podKey := store.Key{
		Namespace:  "namespace",
		APIVersion: "v1",
		Kind:       "Pod",
	}

	pod := testutil.CreatePod("pod")
	u := testutil.ToUnstructured(t, pod)
	objectStore.EXPECT().List(gomock.Any(), podKey).Return(testutil.ToUnstructuredList(t, pod), false, nil).AnyTimes()

	clusterClient.EXPECT().DiscoveryClient().Return(discoveryInterface, nil).AnyTimes()
	discoveryInterface.EXPECT().ServerPreferredNamespacedResources().AnyTimes()
	dashConfig.EXPECT().ClusterClient().Return(clusterClient).AnyTimes()
	dashConfig.EXPECT().ObjectStore().Return(objectStore).AnyTimes()
	dashConfig.EXPECT().PluginManager().Return(pluginManager).AnyTimes()
	dashConfig.EXPECT().TerminateThreshold().Return(int64(5)).AnyTimes()

	queryer := queryerFake.NewMockQueryer(controller)
	queryer.EXPECT().PersistentVolumeClaimsForPod(gomock.Any(), pod)
	queryer.EXPECT().ConfigMapsForPod(gomock.Any(), pod)
	queryer.EXPECT().SecretsForPod(gomock.Any(), pod)
	queryer.EXPECT().ServicesForPod(gomock.Any(), pod)
	queryer.EXPECT().OwnerReference(gomock.Any(), u)
	queryer.EXPECT().Children(gomock.Any(), u).Return(&unstructured.UnstructuredList{}, nil)

	dashConfig.EXPECT().ObjectPath(pod.Namespace, pod.APIVersion, pod.Kind, pod.Name)
	dashConfig.EXPECT().ObjectPath(pod.Namespace, "v1", "ServiceAccount", "")
	eventKey := store.Key{
		Namespace:  "namespace",
		APIVersion: "v1",
		Kind:       "Event",
	}
	objectStore.EXPECT().List(gomock.Any(), eventKey).Return(&unstructured.UnstructuredList{}, false, nil)
	pluginManager.EXPECT().ObjectStatus(gomock.Any(), u).Return(&plugin.ObjectStatusResponse{}, nil)

	tdo := &testDescriberOptions{
		dashConfig: dashConfig,
		queryer:    queryer,
	}
	return tdo
}

func (o *testDescriberOptions) ToOptions() describer.Options {
	return describer.Options{
		Dash:    o.dashConfig,
		Fields:  map[string]string{"name": "pod"},
		Queryer: o.queryer,
	}
}
